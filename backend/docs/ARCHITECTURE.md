# Backend Architecture

> How the backend is structured and why.
> Related: [PROJECT-RULES.md](./PROJECT-RULES.md) · [DATABASE.md](../../docs/DATABASE.md) · [API_SPEC.md](../../docs/API_SPEC.md)

## 1. Overview

- **What:** a Go API that returns movie recommendations for a viewing option (`netflix-chill`, `solo`, `friends`).
- **Who calls it:** the Next.js frontend, opened by customers who scan the QR code on a pasta order card.
- **Shape of the problem:** read-only for customers, tiny dataset (3 options × ≤ 40 movies), data changes slowly.
- **Core design:** precompute recommendations from TMDB in the background; serve every request from memory.

### System diagram

```text
QR card ─scan─▶ Browser ─HTTPS─▶ Next.js frontend
                                       │ REST/JSON (/api/v1)
                                       ▼
┌──────────────────────────── Go API ────────────────────────────┐
│                                                                │
│  READ PATH (per request)          WRITE PATH (background)      │
│                                                                │
│  handler                          refresh job (hourly tick)    │
│     ↓                                ↓                         │
│  service                          refresher ───────▶ TMDB API  │
│     ↓                                ↓ 1. write JSON files     │
│  repository                          │    ($CACHE_DIR, volume) │
│     ↓                                ↓ 2. swap snapshot        │
│  in-memory snapshot ◀────────────────┘                         │
│                                                                │
└────────────────────────────────────────────────────────────────┘
```

**Invariant:** the read path never touches TMDB or the disk.

### Tech stack

| Choice | Why |
|---|---|
| Go (version pinned in `go.mod`) | team expertise; single static binary; goroutines fit the background job |
| `net/http` (Go 1.22+ `ServeMux`) | method + path patterns cover ~5 routes; no framework to learn or upgrade |
| `log/slog` | structured JSON logs in the stdlib |
| `golang.org/x/sync` | `errgroup` (bounded parallel fetches), `singleflight` (dedupe refreshes) |
| `gopkg.in/yaml.v3` | human-editable option config |
| JSON files + in-memory snapshot | < 1 MB of data; survives restarts; no DB to run, back up, or secure |
| TMDB API | discover by genre, details, credits, trailers, and watch providers from one source |
| Docker | reproducible deploys; `$CACHE_DIR` mounted as a volume |

### External constraints

- TMDB's free tier is **non-commercial only** → a commercial agreement is required before production launch.
- Watch-provider data comes from JustWatch → the frontend must show attribution.
- TMDB has no watch-provider data for region `VN` at all. `providers` therefore merges every region and keeps only services that operate in Vietnam (allowlist in `internal/refresh/providers.go`, TMDB_INTEGRATION.md §5). Local services such as FPT Play or VieON are missing from TMDB entirely, so they can never appear.

---

## 2. Folder Structure

```text
backend/
├── cmd/
│   └── server/
│       └── main.go          # entry point: config, wiring, HTTP server, refresh job, shutdown
├── internal/
│   ├── domain/              # pure types + sentinel errors; imports nothing internal
│   │   ├── option.go
│   │   ├── movie.go
│   │   └── errors.go
│   ├── option/              # feature: viewing options
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── context.md
│   ├── movie/               # feature: recommendations + movie detail
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── response.go      # API response shapes
│   │   └── context.md
│   ├── refresh/             # feature: background two-phase refresh + admin trigger
│   │   ├── job.go           # hourly ticker + TTL checks, Trigger, GC after each run
│   │   ├── refresher.go     # phase 1, phase 2, publish
│   │   ├── mapper.go        # TMDB DTO → domain types
│   │   ├── handler.go       # POST /api/v1/admin/refresh
│   │   └── context.md
│   ├── cache/               # snapshot (atomic.Pointer) + JSON load/save + GC
│   │   ├── store.go
│   │   ├── persist.go
│   │   └── context.md
│   └── platform/            # infrastructure only, no business rules
│       ├── config/          # env + options.yaml parsing & validation
│       ├── httpx/           # response envelope, error mapping, middleware
│       ├── fsutil/          # WriteFileAtomic
│       └── tmdb/            # TMDB HTTP client + DTOs
├── configs/
│   └── options.yaml         # option → TMDB discover params (changed via PR)
├── data/cache/              # runtime cache (gitignored; Docker volume in prod)
├── docs/
│   └── decisions/           # backend ADRs
├── .env.example
├── .golangci.yml
├── Dockerfile
└── go.mod
```

- The team template's `src/features/*/context.md` maps to `internal/<feature>/context.md` (Go convention: `internal/`, not `src/`).
- `internal/` also stops other Go modules from importing these packages.

---

## 3. Layer Architecture

### Read path (per request)

| Layer | Responsibility | Example |
|---|---|---|
| Handler | parse + validate input, write the response envelope | `movie.Handler.ListRecommendations` |
| Service | business rules | compute `stale`, build summaries from details |
| Repository | feature-shaped reads over the cache | `GetList(optionID)`, `GetDetails(ids)` |
| Cache | current snapshot, lock-free | `cache.Store.Snapshot()` |

### Write path (background)

```text
refresh.Job ─▶ refresh.Refresher ─▶ platform/tmdb.Client     (fetch)
                       │
                       ├─▶ refresh.mapper                     (TMDB DTO → domain)
                       └─▶ cache.Store                        (write files, swap snapshot)
```

### Dependency rules

```text
cmd/server   → everything (wiring only)
option       → domain, platform/httpx
movie        → domain, cache, platform/httpx
refresh      → domain, cache, platform/tmdb, platform/httpx (admin handler)
cache        → domain, platform/fsutil (cache owns every cache-file write)
platform/*   → domain at most (httpx: error mapping, config: option types)
domain       → nothing
```

- Feature packages never import each other. `movie` doesn't know `refresh` exists; they meet only in `cache`.
- Interfaces are declared by the consumer. Example: `refresh` declares the `TMDB` interface it needs; `platform/tmdb.Client` satisfies it.

---

## 4. Communication

| Between | How | Notes |
|---|---|---|
| FE ↔ BE | REST, JSON over HTTPS | contract in API_SPEC.md; CORS allowlist |
| Modules (in-process) | direct calls through consumer-side interfaces | no event bus, no channels between features |
| Refresh → read path | cache snapshot | single writer (refresher), many readers (handlers) |
| BE → TMDB | HTTPS, `Authorization: Bearer $TMDB_READ_TOKEN` | 10 s timeout per call; max `TMDB_CONCURRENCY` in flight; retry only on `429` — see [TMDB_INTEGRATION.md](./TMDB_INTEGRATION.md) |

---

## 5. Key Flows

### Startup

1. Load and validate env + `configs/options.yaml`. Exit on error.
2. Load cache files into a snapshot. Skip `*.tmp` files, unknown options, and wrong `schema_version`.
3. Start the HTTP server. Options without a cached list return `503 CACHE_NOT_READY`.
4. Start the refresh job. It runs one check immediately.

### Request: `GET /api/v1/options/{option_id}/recommendations`

1. Handler validates `option_id` against the loaded options.
2. Service reads the list + details through the repository.
3. Service builds `MovieSummary[]` in list order and sets `meta.stale` (see DATABASE.md §5).
4. Handler writes `{ data, meta }`.

### Refresh (per option)

1. Every `REFRESH_CHECK_INTERVAL`: skip the option if list age (from `fetched_at`) < `LIST_TTL`.
2. **Phase 1 — list:** `GET /discover/movie` for pages `1..DISCOVER_PAGES` → ordered, de-duplicated IDs.
   ```text
   /discover/movie?with_genres=35|28|27&vote_average.gte=6.5&vote_count.gte=200
                  &sort_by=popularity.desc&language=vi-VN&page=1
   ```
3. **Phase 2 — details:** only for IDs that are missing or older than `DETAIL_TTL`. Failed IDs are dropped from the list.
   ```text
   /movie/{id}?language=vi-VN&append_to_response=watch/providers,credits,videos
              &include_video_language=vi,en
   ```
4. **Publish:** write detail files → write list file → swap snapshot. Details go first so a list never references a missing detail. Details that no list references leave the in-memory snapshot here.
5. **GC** (after all options of a run): delete detail files not referenced by any list.
6. **On failure:** keep existing data, log `Warn`, retry on the next tick.

- Option refreshes run one at a time, so at most `TMDB_CONCURRENCY` TMDB calls are in flight.
- `POST /api/v1/admin/refresh` queues a forced run that skips step 1. A request for an option that is already refreshing joins it (`singleflight`).

Call budget: 3 options × (2 discover + ≤ 40 details) ≈ **126 calls** on a cold start. Later runs are much cheaper because fresh details are reused.

### Shutdown

`SIGTERM` → cancel root context (refresh job stops) → `server.Shutdown` with a 10 s timeout → exit.

---

## 6. Configuration

| Env var | Default | Purpose |
|---|---|---|
| `PORT` | `8080` | HTTP port |
| `TMDB_READ_TOKEN` | — (required) | TMDB read access token (Bearer) |
| `TMDB_LANGUAGE` | `vi-VN` | localized titles and overviews |
| `CACHE_DIR` | `./data/cache` | JSON cache location |
| `LIST_TTL` | `168h` | list freshness (7 days) |
| `DETAIL_TTL` | `336h` | detail freshness (14 days) |
| `REFRESH_CHECK_INTERVAL` | `1h` | job tick |
| `DISCOVER_PAGES` | `2` | pages per option (20 movies/page) |
| `TMDB_CONCURRENCY` | `5` | max parallel TMDB calls |
| `CORS_ALLOWED_ORIGINS` | — | comma-separated frontend origins |
| `ADMIN_TOKEN` | empty | empty = admin endpoints disabled (`404 NOT_FOUND`); set = required `X-Admin-Token` |