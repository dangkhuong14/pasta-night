# Data Model & Storage

> **There is no database.** Data is precomputed from TMDB, persisted as JSON files under `$CACHE_DIR`, and served from an in-memory snapshot. This document is the data model for that store.
> Adding a DB or Redis requires an ADR in `docs/decisions/`.
> Related: [API_SPEC.md](./API_SPEC.md) · [backend/ARCHITECTURE.md](../backend/docs/ARCHITECTURE.md)
> The "TMDB source" columns below are a summary; full mapping and edge-case rules: [backend/TMDB_INTEGRATION.md](../backend/docs/TMDB_INTEGRATION.md) §5.

## 1. Storage Overview

| Entity | Source of truth | Stored at | Written by |
|---|---|---|---|
| ViewingOption | Git | `backend/configs/options.yaml` | developers (via PR) |
| MovieList | TMDB `/discover/movie` | `$CACHE_DIR/lists/{option_id}.json` | refresh job |
| MovieDetail | TMDB `/movie/{id}` | `$CACHE_DIR/details/{movie_id}.json` | refresh job |

- The cache is **disposable**: deleting `$CACHE_DIR` loses nothing. It is rebuilt on the next refresh (`503 CACHE_NOT_READY` until then).
- Timestamps: RFC 3339, UTC.

---

## 2. Entities

### ViewingOption — a "who are you watching with" choice

| Field | Type | Constraints | TMDB discover param |
|---|---|---|---|
| `id` | string | PK; kebab-case; unique; immutable (used in FE routes) | — |
| `label` | string | required; Vietnamese display text | — |
| `description` | string | required | — |
| `icon` | string | required; icon name used by FE | — |
| `discover.genre_ids` | int[] | ≥ 1 TMDB genre ID | `with_genres` |
| `discover.genre_mode` | `and` \| `or` | default `or` | joins IDs with `,` (and) or `\|` (or) |
| `discover.min_vote_average` | float | 0–10; default `6.5` | `vote_average.gte` |
| `discover.min_vote_count` | int | ≥ 0; default `200` (filters obscure titles) | `vote_count.gte` |
| `discover.watch_provider_ids` | int[] | optional | `with_watch_providers` + `watch_region` |
| `discover.sort_by` | string | default `popularity.desc` | `sort_by` |

```yaml
# backend/configs/options.yaml — list order = display order.
# Genre/provider values are starting points; tune them with the business.
options:
  - id: netflix-chill
    label: "Netflix & Chill"
    description: "Cho hai người"
    icon: heart
    discover:
      genre_ids: [10749, 35]      # Romance, Comedy
      genre_mode: or
      min_vote_average: 7.0
      watch_provider_ids: [8]     # Netflix (region = WATCH_REGION)

  - id: solo
    label: "Một mình"
    description: "Thời gian cho riêng bạn"
    icon: user
    discover:
      genre_ids: [18, 53, 9648]   # Drama, Thriller, Mystery
      genre_mode: or

  - id: friends
    label: "Hội bạn"
    description: "Xem cùng nhóm bạn"
    icon: users
    discover:
      genre_ids: [35, 28, 27]     # Comedy, Action, Horror
      genre_mode: or
```

### MovieList — ranked movie IDs for one option

| Field | Type | Constraints |
|---|---|---|
| `schema_version` | int | = `cache.SchemaVersion` |
| `option_id` | string | FK → `ViewingOption.id` |
| `fetched_at` | timestamp | when phase 1 completed |
| `movie_ids` | int[] | rank order; unique; 1–40 items; every ID has a MovieDetail |

```json
{
  "schema_version": 1,
  "option_id": "friends",
  "fetched_at": "2026-09-27T10:00:00Z",
  "movie_ids": [603, 27205, 680]
}
```

### MovieDetail — everything the API can return about one movie

| Field | Type | Constraints | TMDB source |
|---|---|---|---|
| `schema_version` | int | = `cache.SchemaVersion` | — |
| `id` | int | PK; > 0 | `id` |
| `title` | string | required; localized, falls back to original | `title` |
| `original_title` | string | required | `original_title` |
| `overview` | string | may be `""` (no vi-VN translation) | `overview` |
| `tagline` | string | may be `""` | `tagline` |
| `poster_url` | string \| null | absolute URL | `https://image.tmdb.org/t/p/w500` + `poster_path` |
| `backdrop_url` | string \| null | absolute URL | `https://image.tmdb.org/t/p/w1280` + `backdrop_path` |
| `release_year` | int \| null | | year of `release_date` |
| `rating` | float | 0–10, 1 decimal | `vote_average` (rounded) |
| `vote_count` | int | ≥ 0 | `vote_count` |
| `runtime_minutes` | int \| null | > 0 when set | `runtime` (`0` → `null`) |
| `genres` | string[] | localized names | `genres[].name` |
| `directors` | string[] | may be empty | `credits.crew` where `job == "Director"` |
| `cast` | string[] | ≤ 5, billing order | `credits.cast` sorted by `order` |
| `trailer_url` | string \| null | YouTube URL | first `videos.results` with `site == "YouTube"` and `type == "Trailer"` |
| `providers` | Provider[] | may be empty | `watch/providers.results[WATCH_REGION]` |
| `fetched_at` | timestamp | when phase 2 fetched this movie | — |

**Provider** (embedded in MovieDetail)

| Field | Type | Constraints |
|---|---|---|
| `id` | int | TMDB provider ID (e.g. `8` = Netflix) |
| `name` | string | required |
| `logo_url` | string \| null | `https://image.tmdb.org/t/p/w92` + `logo_path` |
| `type` | enum | `flatrate` \| `free` \| `ads` \| `rent` \| `buy` |

```json
{
  "schema_version": 1,
  "id": 603,
  "title": "Ma Trận",
  "original_title": "The Matrix",
  "overview": "Một hacker phát hiện thế giới anh đang sống chỉ là một mô phỏng.",
  "tagline": "",
  "poster_url": "https://image.tmdb.org/t/p/w500/<poster_path>.jpg",
  "backdrop_url": null,
  "release_year": 1999,
  "rating": 8.2,
  "vote_count": 26000,
  "runtime_minutes": 136,
  "genres": ["Phim Hành Động", "Phim Khoa Học Viễn Tưởng"],
  "directors": ["Lana Wachowski", "Lilly Wachowski"],
  "cast": ["Keanu Reeves", "Laurence Fishburne", "Carrie-Anne Moss", "Hugo Weaving", "Gloria Foster"],
  "trailer_url": "https://www.youtube.com/watch?v=<video_key>",
  "providers": [
    { "id": 8, "name": "Netflix", "logo_url": "https://image.tmdb.org/t/p/w92/<logo_path>.jpg", "type": "flatrate" }
  ],
  "fetched_at": "2026-09-27T10:00:05Z"
}
```

---

## 3. Relationships

```text
ViewingOption 1 ──── 1 MovieList N ──────── N MovieDetail 1 ──── N Provider (embedded)
   id  ◀──────────── option_id
                     movie_ids[] ──────────▶ id
```

- **ViewingOption 1–1 MovieList:** `MovieList.option_id` → `ViewingOption.id`.
- **MovieList N–N MovieDetail:** `movie_ids[]` → `details/{id}.json`. One movie can appear in several lists but is stored once.
- **MovieDetail 1–N Provider:** embedded array, no separate file.

Referential integrity is enforced by **write order**, not by an engine:

- A list is published only after every referenced detail file exists.
- Details referenced by no list are deleted by GC. GC runs only after all lists are published.
- A list whose `option_id` is no longer in `options.yaml` is ignored on load and deleted by GC.
- A list that references a missing or invalid detail file is ignored on load (treated as missing) and refetched by the check that runs at startup.

---

## 4. Indexes (in-memory snapshot)

There are no DB indexes; the snapshot's maps play that role. Total size ≈ 3 options × 40 movies × ~3 KB < 1 MB, so everything lives in RAM.

| Index | Go type | Serves | Why |
|---|---|---|---|
| options by ID | `map[string]domain.Option` | option validation on every option route | O(1) allowlist check |
| options in order | `[]domain.Option` | `GET /options` | keep YAML display order |
| lists by option ID | `map[string]domain.MovieList` | `GET /options/{id}/recommendations` | O(1) lookup |
| details by movie ID | `map[int]domain.MovieDetail` | recommendations, `GET /movies/{id}` | O(1); shared across lists |

- The snapshot is immutable. Every publish builds a new one and swaps it via `atomic.Pointer`, so reads are lock-free and never see partial updates.
- The two option indexes are built once at startup from `options.yaml` (options never change at runtime) and passed to the packages that need them. The swapped snapshot holds the list and detail maps; its detail map holds exactly the movies the lists reference.

---

## 5. Freshness (TTL)

| Data | Env var | Default | Why |
|---|---|---|---|
| MovieList | `LIST_TTL` | `168h` (7 d) | suggestions feel new every week |
| MovieDetail | `DETAIL_TTL` | `336h` (14 d) | providers change roughly monthly; metadata almost never |

- Age = now − `fetched_at` **stored in the file**. Never use file mtime (it changes on copy/restore).
- Expired data is still served (stale-while-revalidate).
- `stale = true` when list age > `LIST_TTL + REFRESH_CHECK_INTERVAL`, i.e. a refresh was due but has not succeeded.

---

## 6. Migration Rules

- Every cache file has `schema_version`. The current value is the const `cache.SchemaVersion` in `internal/cache`.
- Any change to a persisted struct (add, rename, remove, retype, or change the meaning of a field):
  1. Bump `cache.SchemaVersion`.
  2. Update the field tables and JSON examples in this file.
  3. Write `cache schema vN → vN+1` in the PR description.
- **No in-place migrations.** On load, files with a different version are treated as missing and refetched on the next tick. Expect `503 CACHE_NOT_READY` for the first refresh after deploy (~126 TMDB calls).
- File naming: `lists/{option_id}.json`, `details/{movie_id}.json`. Temp files `*.tmp` are ignored on load and removed at startup.
- Never hand-edit cache files. To reset: stop the service → delete `$CACHE_DIR/*` → start.
- Changing a ViewingOption `id` breaks FE routes → treat it as an API breaking change (see API_SPEC.md §1).