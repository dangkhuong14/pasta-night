# 001 — Backend stack: stdlib-first Go, two `x/` modules, JSON-file cache

## Status

Accepted (2026-09-27)

## Context

The backend serves 3 viewing options × ≤ 40 movies to anonymous customers. The data changes weekly and is under 1 MB. PROJECT-RULES requires an ADR for every dependency and storage decision; this one records the choices made when the backend was first built (ARCHITECTURE.md §1 Tech stack).

## Decision

- Go standard library for HTTP (`net/http` with Go 1.22 `ServeMux` patterns), logging (`log/slog`), JSON, and crypto.
- `golang.org/x/sync` **v0.11.0**: `errgroup` bounds the parallel TMDB detail fetches; `singleflight` coalesces refreshes of the same option. v0.11.0 is the newest release that supports Go 1.22 (v0.12+ requires Go 1.23).
- `gopkg.in/yaml.v3` **v3.0.1**: parses `configs/options.yaml` and rejects unknown fields.
- Storage: JSON files under `CACHE_DIR` plus an immutable in-memory snapshot swapped with `atomic.Pointer`. No database, no Redis.

## Consequences

- Two third-party modules to keep updated.
- The cache survives restarts through the files and can always be rebuilt from TMDB (~126 calls).
- Each instance keeps its own cache and calls TMDB itself. Running several replicas multiplies TMDB calls; revisit this decision before scaling out.
- Moving to Go 1.23+ unlocks newer `x/sync` releases.

## Alternatives considered

- **chi / gin / echo:** unnecessary for 5 routes; `ServeMux` patterns cover methods and path params.
- **SQLite / Postgres / Redis:** more to run, back up, and secure for < 1 MB of read-mostly data.
- **Options as JSON or env vars:** YAML is easier for the business to review and tune through PRs.
- **`sync.WaitGroup` + a hand-rolled semaphore:** `errgroup.SetLimit` gives bounded concurrency and error propagation with less code.
