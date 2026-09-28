package refresh

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"

	"pasta_night/be/internal/cache"
	"pasta_night/be/internal/domain"
	"pasta_night/be/internal/platform/tmdb"
)

var testNow = time.Date(2026, 9, 27, 10, 0, 0, 0, time.UTC)

// anyCount skips a count check where parallel fetches race an abort.
const anyCount = -1

var (
	errUnauthorized = &tmdb.APIError{HTTPStatus: http.StatusUnauthorized, Code: 7}
	errBadRequest   = &tmdb.APIError{HTTPStatus: http.StatusBadRequest, Code: 22}
	errUnavailable  = fmt.Errorf("%w: connection reset", tmdb.ErrUnavailable)
)

// fakeTMDB implements the TMDB interface from in-memory pages and movies.
// Movies without an entry answer 404.
type fakeTMDB struct {
	mu          sync.Mutex
	pages       map[int]tmdb.DiscoverResponse
	discoverErr error
	movies      map[int]tmdb.MovieDetails
	movieErrs   map[int]error
	delay       time.Duration

	discoverCalls []tmdb.DiscoverQuery
	detailCalls   []int
	inFlight      int
	maxInFlight   int
}

func newFakeTMDB() *fakeTMDB {
	return &fakeTMDB{pages: map[int]tmdb.DiscoverResponse{}, movies: map[int]tmdb.MovieDetails{}, movieErrs: map[int]error{}}
}

// withPage adds a discover page of total pages with the given movie IDs.
func (f *fakeTMDB) withPage(page, total int, ids ...int) *fakeTMDB {
	results := make([]tmdb.DiscoverMovie, len(ids))
	for i, id := range ids {
		results[i] = tmdb.DiscoverMovie{ID: id}
	}
	f.pages[page] = tmdb.DiscoverResponse{Page: page, TotalPages: total, Results: results}
	return f
}

// withMovies adds valid movie details for ids.
func (f *fakeTMDB) withMovies(ids ...int) *fakeTMDB {
	for _, id := range ids {
		f.movies[id] = tmdb.MovieDetails{ID: id, Title: fmt.Sprintf("Phim %d", id), OriginalTitle: fmt.Sprintf("Movie %d", id), VoteAverage: 7.5}
	}
	return f
}

func (f *fakeTMDB) Discover(_ context.Context, q tmdb.DiscoverQuery) (tmdb.DiscoverResponse, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.discoverCalls = append(f.discoverCalls, q)
	if f.discoverErr != nil {
		return tmdb.DiscoverResponse{}, f.discoverErr
	}
	return f.pages[q.Page], nil
}

func (f *fakeTMDB) MovieDetails(ctx context.Context, id int) (tmdb.MovieDetails, error) {
	f.mu.Lock()
	f.detailCalls = append(f.detailCalls, id)
	f.inFlight++
	f.maxInFlight = max(f.maxInFlight, f.inFlight)
	f.mu.Unlock()
	defer func() {
		f.mu.Lock()
		f.inFlight--
		f.mu.Unlock()
	}()
	if f.delay > 0 {
		select {
		case <-time.After(f.delay):
		case <-ctx.Done():
			return tmdb.MovieDetails{}, ctx.Err()
		}
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	if err, ok := f.movieErrs[id]; ok {
		return tmdb.MovieDetails{}, err
	}
	if m, ok := f.movies[id]; ok {
		return m, nil
	}
	return tmdb.MovieDetails{}, &tmdb.APIError{HTTPStatus: http.StatusNotFound, Code: 34}
}

func (f *fakeTMDB) calls() (discover int, details []int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.discoverCalls), append([]int(nil), f.detailCalls...)
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

var friends = domain.Option{
	ID: "friends",
	Discover: domain.DiscoverParams{
		GenreIDs: []int{35, 28, 27}, GenreMode: domain.GenreModeOr,
		MinVoteAverage: 6.5, MinVoteCount: 200, SortBy: "popularity.desc",
	},
}

func newTestRefresher(t *testing.T, client TMDB) (*Refresher, *cache.Store) {
	t.Helper()
	store := cache.NewStore(t.TempDir(), []string{"friends", "solo"}, discardLogger())
	r := NewRefresher(client, store, Settings{WatchRegion: "VN", DiscoverPages: 2, Concurrency: 3, DetailTTL: 336 * time.Hour}, discardLogger())
	r.now = func() time.Time { return testNow }
	return r, store
}

func seedList(t *testing.T, store *cache.Store, fetchedAt time.Time, ids ...int) {
	t.Helper()
	details := make([]domain.MovieDetail, len(ids))
	for i, id := range ids {
		details[i] = domain.MovieDetail{ID: id, Title: fmt.Sprintf("Cũ %d", id), FetchedAt: fetchedAt}
	}
	list := domain.MovieList{OptionID: "friends", FetchedAt: fetchedAt, MovieIDs: ids}
	if err := store.PublishList(context.Background(), list, details); err != nil {
		t.Fatal(err)
	}
}

func TestRefreshOption(t *testing.T) {
	tests := []struct {
		name          string
		fake          *fakeTMDB
		wantErr       error // nil = published
		wantIDs       []int // list after the refresh
		wantFetched   int
		wantFailed    int
		wantPhase     string
		isOldListKept bool
	}{
		{
			name:        "dedupes across pages and skips adult and id 0",
			fake:        newFakeTMDB().withPage(1, 2, 1, 2, 0, 3).withPage(2, 2, 3, 4).withMovies(1, 2, 3, 4),
			wantIDs:     []int{1, 2, 3, 4},
			wantFetched: 4,
			wantPhase:   PhasePublish,
		},
		{
			name: "drops movies that 404 or fail, keeps at least half",
			fake: func() *fakeTMDB {
				f := newFakeTMDB().withPage(1, 1, 1, 2, 3, 4).withMovies(1, 3, 4)
				f.movieErrs[3] = errUnavailable
				return f
			}(),
			wantIDs:     []int{1, 4},
			wantFetched: 2,
			wantFailed:  2,
			wantPhase:   PhasePublish,
		},
		{
			name:          "no movies discovered keeps the old list",
			fake:          newFakeTMDB().withPage(1, 1),
			wantErr:       errNoMovies,
			wantIDs:       []int{9},
			wantPhase:     PhaseDiscover,
			isOldListKept: true,
		},
		{
			name:          "under half of the details keeps the old list",
			fake:          newFakeTMDB().withPage(1, 1, 1, 2, 3, 4).withMovies(1),
			wantErr:       errTooFewDetails,
			wantIDs:       []int{9},
			wantFetched:   1,
			wantFailed:    3,
			wantPhase:     PhaseDetails,
			isOldListKept: true,
		},
		{
			name: "discover 401 aborts",
			fake: func() *fakeTMDB {
				f := newFakeTMDB()
				f.discoverErr = errUnauthorized
				return f
			}(),
			wantErr:       tmdb.ErrUnauthorized,
			wantIDs:       []int{9},
			wantPhase:     PhaseDiscover,
			isOldListKept: true,
		},
		{
			name: "discover outage keeps the old list",
			fake: func() *fakeTMDB {
				f := newFakeTMDB()
				f.discoverErr = errUnavailable
				return f
			}(),
			wantErr:       tmdb.ErrUnavailable,
			wantIDs:       []int{9},
			wantPhase:     PhaseDiscover,
			isOldListKept: true,
		},
		{
			name: "details 401 aborts",
			fake: func() *fakeTMDB {
				f := newFakeTMDB().withPage(1, 1, 1, 2).withMovies(1)
				f.movieErrs[2] = errUnauthorized
				return f
			}(),
			wantErr:       tmdb.ErrUnauthorized,
			wantIDs:       []int{9},
			wantFetched:   anyCount,
			wantFailed:    anyCount,
			wantPhase:     PhaseDetails,
			isOldListKept: true,
		},
		{
			name: "details 400 fails the option",
			fake: func() *fakeTMDB {
				f := newFakeTMDB().withPage(1, 1, 1, 2).withMovies(1)
				f.movieErrs[2] = errBadRequest
				return f
			}(),
			wantErr:       tmdb.ErrBadRequest,
			wantIDs:       []int{9},
			wantFetched:   anyCount,
			wantFailed:    anyCount,
			wantPhase:     PhaseDetails,
			isOldListKept: true,
		},
		{
			name: "mapping drops count as failed",
			fake: func() *fakeTMDB {
				f := newFakeTMDB().withPage(1, 1, 1, 2, 3).withMovies(1, 2)
				f.movies[3] = tmdb.MovieDetails{ID: 3, Adult: true, Title: "x"}
				return f
			}(),
			wantIDs:     []int{1, 2},
			wantFetched: 2,
			wantFailed:  1,
			wantPhase:   PhasePublish,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, store := newTestRefresher(t, tt.fake)
			seedList(t, store, testNow.Add(-400*time.Hour), 9)

			res, err := r.RefreshOption(context.Background(), friends)
			if tt.wantErr == nil && err != nil {
				t.Fatalf("RefreshOption: %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			isFetchedOK := tt.wantFetched == anyCount || res.Fetched == tt.wantFetched
			isFailedOK := tt.wantFailed == anyCount || res.Failed == tt.wantFailed
			if res.Phase != tt.wantPhase || !isFetchedOK || !isFailedOK {
				t.Errorf("result = %+v, want phase %s, fetched %d, failed %d", res, tt.wantPhase, tt.wantFetched, tt.wantFailed)
			}
			if res.IsPublished == tt.isOldListKept {
				t.Errorf("IsPublished = %v", res.IsPublished)
			}
			list := store.Snapshot().Lists["friends"]
			if !reflect.DeepEqual(list.MovieIDs, tt.wantIDs) {
				t.Errorf("list = %v, want %v", list.MovieIDs, tt.wantIDs)
			}
			if !tt.isOldListKept && !list.FetchedAt.Equal(testNow) {
				t.Errorf("fetched_at = %v, want %v", list.FetchedAt, testNow)
			}
		})
	}
}

func TestRefreshReusesFreshDetails(t *testing.T) {
	fake := newFakeTMDB().withPage(1, 1, 1, 2, 3).withMovies(1, 2, 3)
	r, store := newTestRefresher(t, fake)
	fresh := testNow.Add(-time.Hour)
	expired := testNow.Add(-400 * time.Hour) // older than DETAIL_TTL (336h)
	list := domain.MovieList{OptionID: "friends", FetchedAt: fresh, MovieIDs: []int{1, 2, 3}}
	if err := store.PublishList(context.Background(), list, []domain.MovieDetail{
		{ID: 1, Title: "Cũ 1", FetchedAt: fresh},
		{ID: 2, Title: "Cũ 2", FetchedAt: fresh},
		{ID: 3, Title: "Cũ 3", FetchedAt: expired},
	}); err != nil {
		t.Fatal(err)
	}

	res, err := r.RefreshOption(context.Background(), friends)
	if err != nil {
		t.Fatalf("RefreshOption: %v", err)
	}
	if _, details := fake.calls(); !reflect.DeepEqual(details, []int{3}) {
		t.Errorf("fetched details %v, want only the expired movie 3", details)
	}
	if res.Skipped != 2 || res.Fetched != 1 {
		t.Errorf("result = %+v, want 2 skipped and 1 fetched", res)
	}
	snap := store.Snapshot()
	if snap.Details[1].Title != "Cũ 1" || snap.Details[3].Title != "Phim 3" {
		t.Errorf("details = %v", snap.Details)
	}
}

func TestDiscoverPaging(t *testing.T) {
	ids := func(from, n int) []int {
		out := make([]int, n)
		for i := range out {
			out[i] = from + i
		}
		return out
	}
	tests := []struct {
		name         string
		pages        int
		fake         *fakeTMDB
		wantCalls    int
		wantListSize int
	}{
		{"stops at total_pages", 2, newFakeTMDB().withPage(1, 1, 1, 2), 1, 2},
		{"stops at DISCOVER_PAGES", 2, newFakeTMDB().withPage(1, 5, ids(1, 5)...).withPage(2, 5, ids(6, 5)...), 2, 10},
		{"stops at 40 movies", 3, newFakeTMDB().withPage(1, 3, ids(1, 20)...).withPage(2, 3, ids(21, 25)...).withPage(3, 3, ids(46, 20)...), 2, domain.MaxListSize},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := newTestRefresher(t, tt.fake)
			r.settings.DiscoverPages = tt.pages
			got, err := r.discover(context.Background(), friends)
			if err != nil {
				t.Fatal(err)
			}
			if calls, _ := tt.fake.calls(); calls != tt.wantCalls {
				t.Errorf("discover calls = %d, want %d", calls, tt.wantCalls)
			}
			if len(got) != tt.wantListSize {
				t.Errorf("list size = %d, want %d", len(got), tt.wantListSize)
			}
		})
	}
}

func TestDetailsRespectConcurrency(t *testing.T) {
	ids := make([]int, 20)
	for i := range ids {
		ids[i] = i + 1
	}
	fake := newFakeTMDB().withPage(1, 1, ids...).withMovies(ids...)
	fake.delay = 5 * time.Millisecond
	r, _ := newTestRefresher(t, fake)

	if _, err := r.RefreshOption(context.Background(), friends); err != nil {
		t.Fatal(err)
	}
	if fake.maxInFlight > r.settings.Concurrency {
		t.Errorf("max in-flight TMDB calls = %d, want <= %d", fake.maxInFlight, r.settings.Concurrency)
	}
}

func TestDiscoverQuery(t *testing.T) {
	tests := []struct {
		name string
		opt  domain.DiscoverParams
		want tmdb.DiscoverQuery
	}{
		{
			name: "or genres, no providers",
			opt:  domain.DiscoverParams{GenreIDs: []int{35, 28, 27}, GenreMode: domain.GenreModeOr, MinVoteAverage: 6.5, MinVoteCount: 200, SortBy: "popularity.desc"},
			want: tmdb.DiscoverQuery{Page: 1, WithGenres: "35|28|27", VoteAverageGTE: 6.5, VoteCountGTE: 200, SortBy: "popularity.desc"},
		},
		{
			name: "and genres with providers",
			opt: domain.DiscoverParams{GenreIDs: []int{10749, 35}, GenreMode: domain.GenreModeAnd, MinVoteAverage: 7,
				MinVoteCount: 200, WatchProviderIDs: []int{8, 119}, SortBy: "popularity.desc"},
			want: tmdb.DiscoverQuery{Page: 1, WithGenres: "10749,35", VoteAverageGTE: 7, VoteCountGTE: 200,
				SortBy: "popularity.desc", WithWatchProviders: "8|119", WatchRegion: "VN"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := discoverQuery(domain.Option{ID: "x", Discover: tt.opt}, "VN", 1)
			if got != tt.want {
				t.Errorf("query =\n%+v\nwant\n%+v", got, tt.want)
			}
		})
	}
}
