# Backend Project Rules

> Conventions every contributor (human or AI) MUST follow in `backend/`.
> Related: [ARCHITECTURE.md](./ARCHITECTURE.md) · [DATABASE.md](../../docs/DATABASE.md) · [API_SPEC.md](../../docs/API_SPEC.md)

## TL;DR — read this first

- Go, stdlib-first: `net/http`, `log/slog`. No web framework, no ORM, no database.
- The request path reads **only** from the in-memory cache snapshot. Never call TMDB from a handler or service.
- Layers: `handler → service → repository → cache`. Never skip a layer.
- `ctx context.Context` is the first parameter of every function that does I/O or waits.
- Wrap errors with context (`%w`). Map them to API error codes **only** in `platform/httpx`.
- Cache files are written atomically with `fsutil.WriteFileAtomic`. Never write in place.
- Make targeted changes. Don't reformat, rename, or restructure unrelated code.
- New dependency, storage type, or architectural pattern → write an ADR in `docs/decisions/` first.

---

## 1. Naming Conventions

### Files & folders

- Packages: lowercase, one word, singular, no underscores.
  - ✅ `movie`, `refresh`, `httpx` · ❌ `movies`, `movie_service`, `movieService`
- Files: lowercase, `snake_case` when multi-word. Tests end in `_test.go`.
  - ✅ `handler.go`, `movie_list.go`, `service_test.go` · ❌ `MovieHandler.go`, `movie-list.go`
- Standard files in an HTTP feature package (`option/`, `movie/`):

  | File | Contains |
  |---|---|
  | `handler.go` | HTTP handlers + route registration |
  | `service.go` | business logic |
  | `repository.go` | data access (only if the feature reads the cache) |
  | `context.md` | feature notes for humans and AI |

- Config & data: kebab-case → `configs/options.yaml`, option IDs like `netflix-chill`.

### Variables, functions, types

| Kind | Convention | ✅ Do | ❌ Don't |
|---|---|---|---|
| Exported | PascalCase | `MovieDetail`, `NewService` | |
| Unexported | camelCase | `fetchDiscoverPage` | `fetch_discover_page` |
| Initialisms | consistent caps | `MovieID`, `posterURL`, `TMDBClient` | `MovieId`, `posterUrl`, `TmdbClient` |
| Constructors | `New` + type | `movie.NewService(repo, log)` | `movie.CreateService(...)` |
| Interfaces | behavior name, declared by the **consumer** | `MovieReader` in `movie/` | `IMovieRepository` |
| Sentinel errors | `Err` prefix | `ErrOptionNotFound` | `OptionNotFound` |
| Constants | Go casing | `DefaultListTTL` | `DEFAULT_LIST_TTL` |
| Booleans | `is/has/can` | `isStale`, `hasPoster` | `stale`, `poster` |
| Receivers | 1–2 letters, same in every method | `func (s *Service)` | `func (this *Service)` |

### Names outside Go code

| Where | Convention | Example |
|---|---|---|
| JSON fields | snake_case | `poster_url`, `fetched_at` |
| Env vars | UPPER_SNAKE | `TMDB_READ_TOKEN`, `CACHE_DIR` |
| Log keys | snake_case | `option_id`, `duration_ms` |
| API error codes | UPPER_SNAKE | `OPTION_NOT_FOUND` |
| Route paths | lowercase, plural nouns | `/api/v1/options/{option_id}/recommendations` |

---

## 2. Code Style

### Formatting & linting

- `gofmt` + `goimports` are required. CI rejects unformatted code.
- `golangci-lint run` must pass with: `govet`, `errcheck`, `staticcheck`, `revive`, `gosec`, `errorlint`, `bodyclose`, `contextcheck`, `unused`.
- Soft limits: ~120 chars per line, ~50 lines per function. Split functions that do more than one thing.

### Import order

Three groups separated by a blank line. `goimports -local pasta_night/be` enforces this.

```go
import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"pasta_night/be/internal/cache"
	"pasta_night/be/internal/domain"
)
```

### Comments

- Every exported identifier has a doc comment that starts with its name.
- Explain **why**, not what the code already says.
- TODO format: `// TODO(<name>): <action> (<issue link>)`

```go
// ✅ Explains intent and a non-obvious constraint.
// Publish writes all detail files before the list file so readers
// never see a list that references a missing detail.
func (r *Refresher) Publish(ctx context.Context, list domain.MovieList, details []domain.MovieDetail) error

// ❌ Restates the code.
// loop over ids and save
```

---

## 3. Mandatory Patterns

### Error handling

- Wrap with context. Messages are lowercase, have no trailing punctuation, and name the operation.

  ```go
  if err != nil {
  	return fmt.Errorf("fetch discover page %d for option %s: %w", page, opt.ID, err)
  }
  ```

- Domain errors are sentinels in `domain/errors.go`. Check them with `errors.Is` / `errors.As`, never by comparing strings.
- Handlers never build error JSON by hand. Always call `httpx.WriteError(w, r, err)`.
- The domain error → HTTP status + code mapping lives in **one** table: `platform/httpx/errors.go`.

  ```go
  var errorMap = map[error]apiError{
  	domain.ErrOptionNotFound: {http.StatusNotFound, "OPTION_NOT_FOUND"},
  	domain.ErrMovieNotFound:  {http.StatusNotFound, "MOVIE_NOT_FOUND"},
  	domain.ErrCacheNotReady:  {http.StatusServiceUnavailable, "CACHE_NOT_READY"},
  }
  ```

- Unmapped errors → `500 INTERNAL_ERROR`. The real cause is logged, never returned to the client.
- No `panic` or `log.Fatal` outside `cmd/server/main.go`. The recovery middleware is a last-resort net, not control flow.
- Ignoring an error needs a reason: `_ = f.Close() // read-only file; close error is irrelevant`.

### Logging

- `log/slog` with the JSON handler. Inject `*slog.Logger` through constructors. No package-level logger.
- Log once, at the boundary: request middleware for HTTP, the refresh job for background work. Never log **and** return the same error.
- Levels:

  | Level | Use for | Example |
  |---|---|---|
  | Debug | local diagnostics | TMDB response sizes |
  | Info | lifecycle events | startup, refresh completed |
  | Warn | degraded but still serving | refresh failed, serving stale list |
  | Error | needs a human | cache dir not writable |

- Required keys:
  - Requests: `request_id`, `method`, `path`, `status`, `duration_ms`
  - Refresh: `option_id`, `phase`, `fetched`, `skipped`, `failed`

  ```go
  logger.Warn("refresh failed, serving stale list",
  	"option_id", opt.ID, "age_h", int(age.Hours()), "err", err)
  ```

- Never log secrets (`TMDB_READ_TOKEN`, `X-Admin-Token`) or raw upstream request headers.

### Validation

- Validate in the handler. Services assume their input is valid.

  | Input | Rule | On failure |
  |---|---|---|
  | `option_id` (path) | exists in loaded options (allowlist) | `404 OPTION_NOT_FOUND` |
  | `movie_id` (path) | integer > 0 | `400 VALIDATION_ERROR` |
  | `option_id` (admin query) | empty, or exists in options | `404 OPTION_NOT_FOUND` |

- Config (env + `options.yaml`) is validated at startup. Invalid config → exit from `main`.
- Upstream data is validated while mapping in `refresh/mapper.go`, following [TMDB_INTEGRATION.md](./TMDB_INTEGRATION.md) §5 (drop rules, `null` conversions).

### Context & concurrency

- Bounded fan-out for upstream calls: `errgroup` + `g.SetLimit(cfg.TMDBConcurrency)`.
- Every goroutine exits on `ctx.Done()`. No fire-and-forget goroutines.
- Shared read state is an **immutable snapshot** swapped with `atomic.Pointer`. Readers never lock and never see partial updates.
- Concurrent refreshes of the same option are coalesced with `singleflight` (key = `option_id`).

### File writes

```go
// ✅ temp file + rename: atomic on the same filesystem
if err := fsutil.WriteFileAtomic(path, data, 0o644); err != nil {
	return fmt.Errorf("write list %s: %w", opt.ID, err)
}

// ❌ a crash or a concurrent reader can observe a half-written file
os.WriteFile(path, data, 0o644)
```

### Testing

- Table-driven tests. `httptest` for handlers.
- TMDB is faked through the consumer-side interface in `refresh`. Tests never hit the network.
- Every change adds tests for new service logic and for every new error-code mapping.

---

## 4. Don'ts

### Anti-patterns

- ❌ Calling TMDB from handlers or services. Only `refresh` imports `platform/tmdb`.
- ❌ Returning TMDB DTOs from the API. Map to `domain` types in `refresh/mapper.go`.
- ❌ Business logic in handlers. A handler = parse → validate → call service → write response.
- ❌ Importing another feature package. Depend on `domain` + a consumer-side interface instead.
- ❌ `init()` for wiring, or package-level mutable state. Wire everything in `cmd/server/main.go`.
- ❌ `time.NewTicker(7 * 24 * time.Hour)` for refresh — it resets on every restart. Tick hourly and compare `fetched_at` with the TTL.
- ❌ Publishing a list before its details are written, or deleting details before all lists are published.
- ❌ Adding a DB, Redis, or a new third-party module without an ADR.

### Deprecated approaches

| Don't | Use instead |
|---|---|
| `io/ioutil` | `os`, `io` |
| stdlib `log` | `log/slog` |
| `interface{}` | `any` |
| gorilla/mux, chi, gin | `http.ServeMux` patterns: `mux.HandleFunc("GET /api/v1/movies/{movie_id}", h)` |
| `http.ListenAndServe` with defaults | `http.Server` with explicit timeouts |

### Security

- Secrets come from env only. `.env` is gitignored; keep `.env.example` up to date.
- Never return `err.Error()` or stack traces to clients.
- Build upstream URLs with `url.Values`. User input never reaches upstream requests.
- CORS: allowlist from `CORS_ALLOWED_ORIGINS`. No `*` in production.
- Compare admin tokens with `subtle.ConstantTimeCompare`.
- Set `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, and `IdleTimeout` on `http.Server`.
- Don't download or re-host posters/logos. Store TMDB image CDN URLs only (licensing).