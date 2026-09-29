# Project Status — Backend

_Last updated: 2026-09-29 (watch-provider region filtering removed; cache schema v2)_

## Done
- All 5 endpoints of API_SPEC v1: `GET /api/v1/options`, `GET /api/v1/options/{option_id}/recommendations`, `GET /api/v1/movies/{movie_id}`, `POST /api/v1/admin/refresh`, `GET /healthz`
- Two-phase TMDB refresh (discover → details) with publish guards, GC, singleflight, and admin-forced runs
- JSON file cache (`schema_version` 1) + in-memory snapshot; startup load skips invalid files
- Config from env (`ADMIN_TOKEN` empty = admin disabled) and validated `configs/options.yaml`
- Middleware: request ID, request log, panic recovery, CORS allowlist, admin token (constant-time)
- Dockerfile (distroless, nonroot, cache volume)
- Tests: unit, handler, e2e against a fake TMDB; race detector passes; golangci-lint clean
- ADR `docs/decisions/001-backend-stack.md`
- First real run against TMDB: all 3 options cached in `data/cache/`
- Dropped watch-provider region filtering (2026-09-29), because TMDB has no data for region `VN` at all. `WATCH_REGION` is gone; discover no longer sends `watch_region`; `providers` merges every region, ranked by how many offer the service, de-duplicated by name (TMDB gives Amazon Prime Video two IDs) and capped at 8. Cache schema v1 → v2. Result: `netflix-chill` went from 0 to 40 movies, and 115 of 120 movies now list providers (was 0).
- Doc/code audit against API_SPEC, DATABASE, TMDB_INTEGRATION: no contract mismatches found. Two doc bugs fixed:
  - API_SPEC §3, §5.5: documented that `/healthz` sends `Cache-Control: no-store` (it reports live state, unlike other GETs)
  - DATABASE §2: `original_title` changed from "required" to "may be `\"\"`\"" (rare TMDB data gap; mapped as-is, no fallback, matches TMDB_INTEGRATION §5)

## In Progress
- None

## Next
- Replace hand-written TMDB fixtures with real captures
- Set up a deploy environment (host, DNS, cache volume)
- Decide whether to document the remaining stricter-than-spec behaviors found in the audit (fetched_at truncated to the second, mapper drops blank names/invalid providers, GET routes also match HEAD, options.yaml validates sort_by/IDs more tightly, admin `?option_id=` treats blank as "all" and ignores malformed values, movie_id accepts `+5`/`007`) or leave them as implementation details

## Known Issues
- Vietnamese ISP DNS blocks `api.themoviedb.org` (NXDOMAIN). Use 1.1.1.1 / 8.8.8.8 or a VPN, locally and in production.
- `providers` is a worldwide list, not a "watchable in Vietnam" list (see Done). TMDB simply has no `VN` watch-provider data to offer.
- `revive` flags `movie.MovieReader` as stuttering, but PROJECT-RULES §1 uses that name as its example; suppressed with `//nolint:revive`.
- No gcc or golangci-lint on the dev machine: run lint and `go test -race` through Docker.

## Open Questions
- TMDB commercial agreement before production launch.
- Several replicas would each call TMDB (ADR 001): revisit before scaling out.
