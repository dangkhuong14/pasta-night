# Backend — CLAUDE.md

Go API. Before writing code, read `docs/PROJECT-RULES.md`; its TL;DR is mandatory.

## Quick Facts

- Go, stdlib-first (`net/http`, `log/slog`). No framework, no ORM, no database.
- The request path reads only the in-memory snapshot. Only `internal/refresh` calls TMDB.
- Layers: handler → service → repository → cache. Shared types and sentinel errors live in `internal/domain`.

## Package Map

| Package | Responsibility | Read |
|---|---|---|
| `internal/option` | viewing options endpoint | `internal/option/context.md` |
| `internal/movie` | recommendations + movie detail endpoints | `internal/movie/context.md`, `../docs/API_SPEC.md` §5 |
| `internal/refresh` | two-phase background refresh + admin refresh endpoint | `internal/refresh/context.md`, `docs/TMDB_INTEGRATION.md`, `docs/ARCHITECTURE.md` §5 |
| `internal/cache` | snapshot + JSON persistence + GC | `internal/cache/context.md`, `../docs/DATABASE.md` |
| `internal/domain` | shared types + sentinel errors | `../docs/DATABASE.md` §2 |
| `internal/platform/tmdb` | TMDB client + DTOs | `docs/TMDB_INTEGRATION.md` §3, §4, §6 |
| `internal/platform/httpx` | envelope, error mapping, middleware | `../docs/API_SPEC.md` §4, §6 |
| `internal/platform/config` | env + `options.yaml` loading | `docs/ARCHITECTURE.md` §6 |

## Commands

```bash
cp .env.example .env                                  # first time; fill in TMDB_READ_TOKEN
set -a && . ./.env && set +a && go run ./cmd/server   # run on :8080
go test ./...                                         # all tests, no network
go test ./internal/refresh/... -run TestMapper        # one package / test
golangci-lint run
goimports -local pasta_night/be -w .
```

## Feature Context Files

- Read `internal/<feature>/context.md` before editing a feature.
- Create it the first time you touch a feature that has none.
- Keep it under ~40 lines. Update it when invariants or gotchas change.

Template:

```md
# <feature>

## Purpose
One or two sentences on what this package does and for whom.

## Key Files
- `service.go` — builds MovieSummary list, computes `stale`

## Invariants
- A list is published only after all its details are written.

## Gotchas
- TMDB can return the same movie ID on two discover pages.

## Related Docs
- `../../docs/TMDB_INTEGRATION.md` §5
```