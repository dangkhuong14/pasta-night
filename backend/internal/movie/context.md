# movie

## Purpose
Serves ranked recommendations per option and movie details, read only from the cache snapshot.

## Key Files
- `handler.go` — routes; validates `option_id` (allowlist) and `movie_id` (integer > 0)
- `service.go` — builds recommendations, computes `stale`
- `repository.go` — reads one snapshot per call
- `response.go` — API shapes (no `schema_version`; `fetched_at` moves to `meta`)

## Invariants
- Never calls TMDB or reads the disk.
- `stale` = list age > `LIST_TTL + REFRESH_CHECK_INTERVAL`; stale data is still served.
- A movie is available only while a current list references it (the snapshot holds exactly those).
- Arrays are `[]`, never `null`; nullable fields are always present.

## Gotchas
- Read a list and its details from the same snapshot; a refresh between two reads can drop details.

## Related Docs
- `../../../docs/API_SPEC.md` §5.2, §5.3
- `../../../docs/DATABASE.md` §5
