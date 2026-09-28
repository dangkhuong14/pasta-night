# refresh

## Purpose
Keeps the cache fresh from TMDB. A job checks every option on a ticker and refreshes lists older than `LIST_TTL`; `POST /api/v1/admin/refresh` forces a refresh. The only package that calls TMDB.

## Key Files
- `job.go` — check at startup and on every tick, TTL check, `Trigger`, GC after each run, result logging
- `refresher.go` — phase 1 (discover), phase 2 (details), publish guards, publish
- `mapper.go` — TMDB DTO → domain (TMDB_INTEGRATION.md §5)
- `handler.go` — admin endpoint; `main` registers it only when `ADMIN_TOKEN` is set

## Invariants
- Option refreshes run one at a time (`Refresher.mu`), so at most `TMDB_CONCURRENCY` calls are in flight.
- A refresh of an option that is already refreshing joins it (`singleflight`, key = option ID).
- On any failure the old list stays. Guards: phase 1 found 0 IDs, or phase 2 kept < 50%.
- A TMDB 401 stops the whole run: every other call would fail the same way.
- List age comes from `fetched_at`, so a restart does not reset the schedule.

## Gotchas
- Discover rankings shift between page calls: IDs are de-duplicated across pages.
- `Trigger` never blocks. The 202 response lists accepted options; the refresh runs afterwards.

## Related Docs
- `../../docs/TMDB_INTEGRATION.md`
- `../../docs/ARCHITECTURE.md` §5
