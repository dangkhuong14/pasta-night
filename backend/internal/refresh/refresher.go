// Package refresh keeps the cache fresh from TMDB: a job checks every option
// on a ticker and runs the two-phase refresh when a list is due, and an admin
// endpoint forces one. It is the only package that calls TMDB.
package refresh

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	"pasta_night/be/internal/cache"
	"pasta_night/be/internal/domain"
	"pasta_night/be/internal/platform/tmdb"
)

// TMDB is the part of the TMDB client the refresher needs.
// platform/tmdb.Client satisfies it; tests use a fake.
type TMDB interface {
	Discover(ctx context.Context, q tmdb.DiscoverQuery) (tmdb.DiscoverResponse, error)
	MovieDetails(ctx context.Context, id int) (tmdb.MovieDetails, error)
}

// Settings are the refresh knobs from the environment.
type Settings struct {
	WatchRegion   string
	DiscoverPages int
	Concurrency   int
	DetailTTL     time.Duration
}

// Phases, as reported in the "phase" log key.
const (
	PhaseDiscover = "discover"
	PhaseDetails  = "details"
	PhasePublish  = "publish"
)

// Publish guards (TMDB_INTEGRATION.md §6): the old list stays in place.
var (
	// errNoMovies means phase 1 found no movies, e.g. filters too strict for the region.
	errNoMovies = errors.New("discover returned no movies")
	// errTooFewDetails means phase 2 kept under half of the phase-1 IDs, e.g. during a TMDB outage.
	errTooFewDetails = errors.New("too few movie details fetched")
)

// Result summarizes one option refresh for the job's log line.
type Result struct {
	OptionID    string
	Phase       string // last phase reached
	Fetched     int    // details fetched from TMDB
	Skipped     int    // details reused because they are younger than DETAIL_TTL
	Failed      int    // movies dropped by a TMDB error or a mapping rule
	IsPublished bool
}

// Refresher runs the two-phase refresh (ARCHITECTURE.md §5).
type Refresher struct {
	client   TMDB
	store    *cache.Store
	settings Settings
	log      *slog.Logger
	now      func() time.Time
	mu       sync.Mutex // one option at a time keeps TMDB calls within Settings.Concurrency
}

// NewRefresher returns a Refresher that fetches through client and publishes to store.
func NewRefresher(client TMDB, store *cache.Store, settings Settings, log *slog.Logger) *Refresher {
	return &Refresher{client: client, store: store, settings: settings, log: log, now: time.Now}
}

// RefreshOption fetches a new list for opt and publishes it with its
// details. On any error the existing list stays in place. An error wrapping
// tmdb.ErrUnauthorized means the whole run must stop.
func (r *Refresher) RefreshOption(ctx context.Context, opt domain.Option) (Result, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	res := Result{OptionID: opt.ID, Phase: PhaseDiscover}
	ids, err := r.discover(ctx, opt)
	if err != nil {
		return res, err
	}
	if len(ids) == 0 {
		return res, errNoMovies
	}
	listFetchedAt := r.timestamp()

	res.Phase = PhaseDetails
	details, stats, err := r.details(ctx, ids)
	res.Fetched, res.Skipped, res.Failed = stats.fetched, stats.reused, stats.failed
	if err != nil {
		return res, err
	}
	if len(details)*2 < len(ids) {
		return res, fmt.Errorf("%w: kept %d of %d", errTooFewDetails, len(details), len(ids))
	}

	res.Phase = PhasePublish
	movieIDs := make([]int, len(details))
	for i, d := range details {
		movieIDs[i] = d.ID
	}
	list := domain.MovieList{OptionID: opt.ID, FetchedAt: listFetchedAt, MovieIDs: movieIDs}
	if err := r.store.PublishList(ctx, list, details); err != nil {
		return res, err
	}
	res.IsPublished = true
	return res, nil
}

// discover runs phase 1: TMDB rank order across pages, de-duplicated, at
// most domain.MaxListSize IDs.
func (r *Refresher) discover(ctx context.Context, opt domain.Option) ([]int, error) {
	ids := make([]int, 0, domain.MaxListSize)
	seen := make(map[int]bool, domain.MaxListSize)
	for page := 1; page <= r.settings.DiscoverPages && len(ids) < domain.MaxListSize; page++ {
		resp, err := r.client.Discover(ctx, discoverQuery(opt, r.settings.WatchRegion, page))
		if err != nil {
			return nil, fmt.Errorf("fetch discover page %d for option %s: %w", page, opt.ID, err)
		}
		for _, m := range resp.Results {
			// Rankings can shift between page calls, so a movie can repeat.
			if m.ID <= 0 || m.Adult || seen[m.ID] {
				continue
			}
			seen[m.ID] = true
			ids = append(ids, m.ID)
			if len(ids) == domain.MaxListSize {
				break
			}
		}
		if page >= resp.TotalPages {
			break
		}
	}
	return ids, nil
}

// discoverQuery maps an option's discover params to TMDB query params
// (TMDB_INTEGRATION.md §4.1).
func discoverQuery(opt domain.Option, region string, page int) tmdb.DiscoverQuery {
	d := opt.Discover
	genreSep := "|"
	if d.GenreMode == domain.GenreModeAnd {
		genreSep = ","
	}
	q := tmdb.DiscoverQuery{
		Page:           page,
		WithGenres:     joinInts(d.GenreIDs, genreSep),
		VoteAverageGTE: d.MinVoteAverage,
		VoteCountGTE:   d.MinVoteCount,
		SortBy:         d.SortBy,
	}
	if len(d.WatchProviderIDs) > 0 {
		q.WithWatchProviders = joinInts(d.WatchProviderIDs, "|")
		q.WatchRegion = region
	}
	return q
}

func joinInts(ns []int, sep string) string {
	parts := make([]string, len(ns))
	for i, n := range ns {
		parts[i] = strconv.Itoa(n)
	}
	return strings.Join(parts, sep)
}

type detailStats struct {
	fetched, reused, failed int
}

// details runs phase 2: reuse details younger than DETAIL_TTL, fetch the
// rest in parallel, and drop the movies that fail. The result keeps rank order.
func (r *Refresher) details(ctx context.Context, ids []int) ([]domain.MovieDetail, detailStats, error) {
	snap := r.store.Snapshot()
	now := r.now()
	var stats detailStats
	results := make([]domain.MovieDetail, len(ids))
	isFilled := make([]bool, len(ids))
	var toFetch []int // indexes into ids
	for i, id := range ids {
		if d, ok := snap.Details[id]; ok && now.Sub(d.FetchedAt) < r.settings.DetailTTL {
			results[i], isFilled[i] = d, true
			stats.reused++
			continue
		}
		toFetch = append(toFetch, i)
	}

	var fetched, failed atomic.Int64
	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(r.settings.Concurrency)
	for _, i := range toFetch {
		g.Go(func() error {
			d, ok, err := r.fetchDetail(gctx, ids[i])
			if err != nil {
				return err
			}
			if !ok {
				failed.Add(1)
				return nil
			}
			results[i], isFilled[i] = d, true // each goroutine owns index i
			fetched.Add(1)
			return nil
		})
	}
	err := g.Wait()
	stats.fetched, stats.failed = int(fetched.Load()), int(failed.Load())
	if err != nil {
		return nil, stats, err
	}

	details := make([]domain.MovieDetail, 0, len(ids))
	for i, d := range results {
		if isFilled[i] {
			details = append(details, d)
		}
	}
	return details, stats, nil
}

// fetchDetail fetches and maps one movie. ok is false when the movie is
// dropped (TMDB_INTEGRATION.md §6); an error fails the whole option.
func (r *Refresher) fetchDetail(ctx context.Context, id int) (domain.MovieDetail, bool, error) {
	m, err := r.client.MovieDetails(ctx, id)
	if err != nil {
		switch {
		case ctx.Err() != nil:
			return domain.MovieDetail{}, false, ctx.Err()
		case errors.Is(err, tmdb.ErrUnauthorized), errors.Is(err, tmdb.ErrBadRequest):
			return domain.MovieDetail{}, false, fmt.Errorf("fetch details for movie %d: %w", id, err)
		default: // not found, unavailable, still rate limited after retries
			r.log.Debug("dropping movie", "movie_id", id, "err", err)
			return domain.MovieDetail{}, false, nil
		}
	}
	d, ok := mapMovieDetails(m, r.settings.WatchRegion, r.timestamp())
	if !ok || d.ID != id {
		r.log.Debug("dropping movie that fails the mapping rules", "movie_id", id)
		return domain.MovieDetail{}, false, nil
	}
	return d, true, nil
}

// timestamp is the stored fetched_at: UTC, whole seconds.
func (r *Refresher) timestamp() time.Time {
	return r.now().UTC().Truncate(time.Second)
}
