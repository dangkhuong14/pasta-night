package refresh

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"

	"pasta_night/be/internal/cache"
	"pasta_night/be/internal/domain"
	"pasta_night/be/internal/platform/tmdb"
)

// Job decides when options refresh: on every tick for lists older than
// LIST_TTL, and on demand through Trigger. A restart does not reset the
// schedule because list age comes from fetched_at, not from the ticker.
type Job struct {
	refresher *Refresher
	store     *cache.Store
	options   []domain.Option
	byID      map[string]domain.Option
	interval  time.Duration
	listTTL   time.Duration
	log       *slog.Logger
	now       func() time.Time

	group singleflight.Group // key = option ID: concurrent refreshes of one option share a run

	mu      sync.Mutex
	pending map[string]bool // option IDs forced by Trigger and not started yet
	wake    chan struct{}   // capacity 1: tells Run that pending changed
}

// NewJob returns a Job that checks the options every interval.
func NewJob(refresher *Refresher, store *cache.Store, options []domain.Option,
	interval, listTTL time.Duration, log *slog.Logger,
) *Job {
	byID := make(map[string]domain.Option, len(options))
	for _, o := range options {
		byID[o.ID] = o
	}
	return &Job{
		refresher: refresher,
		store:     store,
		options:   options,
		byID:      byID,
		interval:  interval,
		listTTL:   listTTL,
		log:       log,
		now:       time.Now,
		pending:   make(map[string]bool),
		wake:      make(chan struct{}, 1),
	}
}

// Run checks every option immediately, then on every tick and on every
// Trigger, until ctx is done. It returns once all refreshes it started stop.
func (j *Job) Run(ctx context.Context) {
	var wg sync.WaitGroup
	defer wg.Wait()
	start := func(optionIDs []string, isForced bool) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			j.runOnce(ctx, optionIDs, isForced)
		}()
	}

	start(j.optionIDs(), false)
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			start(j.optionIDs(), false)
		case <-j.wake:
			if ids := j.takePending(); len(ids) > 0 {
				start(ids, true)
			}
		}
	}
}

// Trigger forces a refresh of the given options, ignoring LIST_TTL. It never
// blocks. An option that is already refreshing joins the running refresh.
func (j *Job) Trigger(optionIDs []string) {
	j.mu.Lock()
	for _, id := range optionIDs {
		j.pending[id] = true
	}
	j.mu.Unlock()
	select {
	case j.wake <- struct{}{}:
	default: // Run has a wake-up queued already; it will see the new IDs
	}
}

func (j *Job) optionIDs() []string {
	ids := make([]string, len(j.options))
	for i, o := range j.options {
		ids[i] = o.ID
	}
	return ids
}

// takePending returns the forced option IDs in display order and clears them.
func (j *Job) takePending() []string {
	j.mu.Lock()
	defer j.mu.Unlock()
	var ids []string
	for _, o := range j.options {
		if j.pending[o.ID] {
			ids = append(ids, o.ID)
		}
	}
	clear(j.pending)
	return ids
}

// runOnce refreshes the due (or forced) options one after another, then
// deletes cache files that are no longer referenced.
func (j *Job) runOnce(ctx context.Context, optionIDs []string, isForced bool) {
	for _, id := range optionIDs {
		if ctx.Err() != nil {
			return
		}
		opt, ok := j.byID[id]
		if !ok || (!isForced && !j.isDue(id)) {
			continue
		}
		if err := j.refresh(ctx, opt); errors.Is(err, tmdb.ErrUnauthorized) {
			return // every other call would fail the same way
		}
	}
	if ctx.Err() != nil {
		return
	}
	stats, err := j.store.GC(ctx)
	if err != nil {
		j.log.Error("cache gc failed", "err", err)
		return
	}
	if stats.Details > 0 || stats.Lists > 0 {
		j.log.Info("cache gc completed", "deleted_details", stats.Details, "deleted_lists", stats.Lists)
	}
}

// isDue reports whether an option has no list or its list is at least LIST_TTL old.
func (j *Job) isDue(optionID string) bool {
	list, ok := j.store.Snapshot().Lists[optionID]
	return !ok || j.now().Sub(list.FetchedAt) >= j.listTTL
}

// refresh runs one option refresh through singleflight and logs the outcome
// once, inside the shared call.
func (j *Job) refresh(ctx context.Context, opt domain.Option) error {
	_, err, _ := j.group.Do(opt.ID, func() (any, error) {
		start := time.Now()
		res, err := j.refresher.RefreshOption(ctx, opt)
		j.logResult(res, err, time.Since(start))
		return res, err
	})
	return err
}

func (j *Job) logResult(res Result, err error, took time.Duration) {
	attrs := []any{
		"option_id", res.OptionID,
		"phase", res.Phase,
		"fetched", res.Fetched,
		"skipped", res.Skipped,
		"failed", res.Failed,
		"duration_ms", took.Milliseconds(),
	}
	switch {
	case err == nil:
		j.log.Info("refresh completed", attrs...)
	case errors.Is(err, context.Canceled):
		j.log.Info("refresh cancelled", attrs...)
	case errors.Is(err, tmdb.ErrUnauthorized):
		j.log.Error("refresh aborted: TMDB rejected the token", append(attrs, "err", err)...)
	case errors.Is(err, tmdb.ErrBadRequest):
		j.log.Error("refresh failed: TMDB rejected our request", append(attrs, "err", err)...)
	default:
		if list, ok := j.store.Snapshot().Lists[res.OptionID]; ok {
			age := j.now().Sub(list.FetchedAt)
			j.log.Warn("refresh failed, serving stale list", append(attrs, "age_h", int(age.Hours()), "err", err)...)
			return
		}
		j.log.Warn("refresh failed, option has no list yet", append(attrs, "err", err)...)
	}
}
