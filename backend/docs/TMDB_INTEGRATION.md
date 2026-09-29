# TMDB Integration

> How the backend calls TMDB and maps responses to our domain types.
> This file is the **source of truth for mapping rules**. Field types/constraints live in [DATABASE.md](../../docs/DATABASE.md).
> Related: [ARCHITECTURE.md](./ARCHITECTURE.md) §5 Refresh · [PROJECT-RULES.md](./PROJECT-RULES.md)

## 1. Scope

- Only two packages touch TMDB:
  - `internal/platform/tmdb` — HTTP client + DTOs. Knows TMDB, knows nothing about our domain.
  - `internal/refresh` — calls the client, maps DTOs → `domain` types (`mapper.go`).
- The request path (handlers/services) **never** calls TMDB.
- We use exactly **two endpoints**: `GET /discover/movie` and `GET /movie/{id}`.

---

## 2. Setup

| Item | Value |
|---|---|
| API base URL | `https://api.themoviedb.org/3` |
| Image base URL | `https://image.tmdb.org/t/p/` |
| Credential | API Read Access Token (TMDB account → Settings → API) |
| Env var | `TMDB_READ_TOKEN` |
| License | free key = non-commercial only; **commercial agreement required before launch** |

Smoke test (requires `jq`):

```bash
curl -s "https://api.themoviedb.org/3/movie/603?language=vi-VN" \
  -H "Authorization: Bearer $TMDB_READ_TOKEN" | jq '{id, title, runtime}'
```

---

## 3. Client Rules (`internal/platform/tmdb`)

- Headers on every request:
  ```text
  Authorization: Bearer <TMDB_READ_TOKEN>
  Accept: application/json
  ```
- Timeout: `http.Client{Timeout: 10 * time.Second}` + the caller's `ctx`.
- Build queries with `url.Values`. Never concatenate strings into URLs.
- Always send `language=$TMDB_LANGUAGE`.
- Concurrency is limited by the **caller** (`refresh`, `errgroup.SetLimit(TMDB_CONCURRENCY)`). TMDB rejects bursts over ~40 requests/s per IP with `429`.
- DTOs decode **only the fields we use**. Fewer fields = fewer breakages when TMDB changes.
- DTOs use plain Go types. JSON `null` and missing fields both decode to zero values (`""`, `0`); the mapper turns zero values into `null` where the domain allows it.
- DTOs live in `platform/tmdb/dto.go`, named after TMDB (`MovieDetails`, not `Movie`). They never leave `platform/tmdb` and `refresh`.

---

## 4. Endpoints

### 4.1 Phase 1 — `GET /discover/movie`

**Request params**

| Param | Example | Source |
|---|---|---|
| `language` | `vi-VN` | `TMDB_LANGUAGE` |
| `include_adult` | `false` | constant |
| `include_video` | `false` | constant |
| `with_genres` | `35\|28\|27` | option `genre_ids`, joined by `\|` (or) / `,` (and) per `genre_mode` |
| `vote_average.gte` | `6.5` | option `min_vote_average` |
| `vote_count.gte` | `200` | option `min_vote_count` |
| `sort_by` | `popularity.desc` | option `sort_by` |
| `with_watch_providers` | `8` | option `watch_provider_ids` joined by `\|` = any of them (only if set) |
| `with_watch_monetization_types` | `flatrate` | constant (only if `with_watch_providers` is set) |

- **`watch_region` is never sent.** TMDB matches zero movies for a region it has no provider data for: `with_watch_providers=8&watch_region=VN` returns `total_results: 0` for every query, which silently emptied the `netflix-chill` option. Without the parameter the provider filter still works.
| `page` | `1` | loop `1..DISCOVER_PAGES`; stop early when `page > total_pages` |

```bash
curl -s -G "https://api.themoviedb.org/3/discover/movie" \
  -H "Authorization: Bearer $TMDB_READ_TOKEN" \
  --data-urlencode "with_genres=35|28|27" \
  -d language=vi-VN -d include_adult=false -d vote_count.gte=200 \
  -d sort_by=popularity.desc -d page=1 | jq '[.results[] | {id, title}]'
```

**Response DTO** — phase 1 only needs IDs:

```go
type DiscoverResponse struct {
	Page       int             `json:"page"`
	TotalPages int             `json:"total_pages"`
	Results    []DiscoverMovie `json:"results"`
}

type DiscoverMovie struct {
	ID    int  `json:"id"`
	Adult bool `json:"adult"`
}
```

**Phase 1 rules**

- Keep TMDB's rank order across pages (page 1 first).
- De-duplicate IDs: rankings can shift between page calls, so the same movie can appear on two pages.
- Skip `id == 0` and `adult == true`.
- Stop at 40 IDs (`domain.MaxListSize`); pages after that are not fetched.

### 4.2 Phase 2 — `GET /movie/{id}` with `append_to_response`

One call per movie returns details + credits + videos + watch providers.

```text
GET /movie/603?language=vi-VN
              &append_to_response=credits,videos,watch/providers
              &include_video_language=vi,en
```

- `include_video_language=vi,en`: with `language=vi-VN` alone, most movies return no trailers.
- TMDB allows at most **20** appended objects per call; we use 3.

**Response DTO:**

```go
type MovieDetails struct {
	ID             int            `json:"id"`
	Adult          bool           `json:"adult"`
	Title          string         `json:"title"`
	OriginalTitle  string         `json:"original_title"`
	Overview       string         `json:"overview"`
	Tagline        string         `json:"tagline"`
	PosterPath     string         `json:"poster_path"`
	BackdropPath   string         `json:"backdrop_path"`
	ReleaseDate    string         `json:"release_date"` // "YYYY-MM-DD" or ""
	VoteAverage    float64        `json:"vote_average"`
	VoteCount      int            `json:"vote_count"`
	Runtime        int            `json:"runtime"`
	Genres         []Genre        `json:"genres"`
	Credits        Credits        `json:"credits"`
	Videos         Videos         `json:"videos"`
	WatchProviders WatchProviders `json:"watch/providers"` // ⚠️ the JSON key contains a slash
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Credits struct {
	Cast []CastMember `json:"cast"`
	Crew []CrewMember `json:"crew"`
}

type CastMember struct {
	Name        string `json:"name"`
	Order       int    `json:"order"`        // billing order, 0 = top
	ProfilePath string `json:"profile_path"` // "" when TMDB has no photo
}

type CrewMember struct {
	Name string `json:"name"`
	Job  string `json:"job"`
}

type Videos struct {
	Results []Video `json:"results"`
}

type Video struct {
	Key      string `json:"key"`
	Site     string `json:"site"`      // "YouTube", "Vimeo"
	Type     string `json:"type"`      // "Trailer", "Teaser", "Clip", ...
	Official bool   `json:"official"`
	Language string `json:"iso_639_1"` // "vi", "en"
}

type WatchProviders struct {
	Results map[string]RegionProviders `json:"results"` // key = region code, e.g. "VN"
}

type RegionProviders struct {
	Flatrate []Provider `json:"flatrate"`
	Free     []Provider `json:"free"`
	Ads      []Provider `json:"ads"`
	Rent     []Provider `json:"rent"`
	Buy      []Provider `json:"buy"`
}

type Provider struct {
	ProviderID      int    `json:"provider_id"`
	ProviderName    string `json:"provider_name"`
	LogoPath        string `json:"logo_path"`
	DisplayPriority int    `json:"display_priority"`
}
```

---

## 5. Mapping Rules (`refresh/mapper.go`)

`tmdb.MovieDetails` → `domain.MovieDetail`. Rules are applied in this order; a movie is **dropped** only where stated.

| Domain field | Rule | Example |
|---|---|---|
| — | drop if `id == 0` or `adult == true` | |
| `id` | `id` | `603` |
| `title` | `title`; if `""` → `original_title`; if both `""` → **drop** | `"Ma Trận"` |
| `original_title` | `original_title` | `"The Matrix"` |
| `overview` | as-is; may stay `""` (no vi-VN translation) | |
| `tagline` | as-is | |
| `poster_url` | `""` → `null`; else `image base + "w500" + poster_path` | `…/t/p/w500/abc.jpg` |
| `backdrop_url` | `""` → `null`; else `image base + "w1280" + backdrop_path` | |
| `release_year` | first 4 chars of `release_date` as int; `""` or invalid → `null` | `"1999-03-30"` → `1999` |
| `rating` | `math.Round(vote_average*10) / 10` | `8.216` → `8.2` |
| `vote_count` | as-is | |
| `runtime_minutes` | `0` → `null` | |
| `genres` | `genres[].name`, TMDB order | `["Phim Hành Động"]` |
| `directors` | `credits.crew` where `job == "Director"`; names de-duplicated, TMDB order | `["Lana Wachowski", "Lilly Wachowski"]` |
| `cast` | `credits.cast` sorted by `order` asc; first 5, each mapped to `{name, profile_url}`. `profile_url`: `""` → `null`; else `image base + "w185" + profile_path` (avatars render at 48 px, so w185 covers 3x screens) | `{"name": "Keanu Reeves", "profile_url": "…/t/p/w185/abc.jpg"}` |
| `trailer_url` | see trailer selection below; none → `null` | `https://www.youtube.com/watch?v=<key>` |
| `providers` | see provider rules below; region missing → `[]` | |
| `fetched_at` | `time.Now().UTC()` at fetch time | |

### Trailer selection

1. Keep videos with `site == "YouTube"` and `type == "Trailer"`.
2. Prefer language `vi`, then `en`, then others.
3. Within the same language, prefer `official == true`.
4. Take the first → `https://www.youtube.com/watch?v=` + `key`.

### Provider rules

**No region filtering, but an allowlist instead.** TMDB has no watch-provider data at all for some regions, `VN` among them, so a region-scoped list would always be empty here. Every region is merged instead, and only services a customer in Vietnam can actually use are kept — otherwise a movie lists Rakuten TV, Viaplay or Sky Store, which is worse than showing nothing.

The allowlist lives in `internal/refresh/providers.go` (`allowedProviderIDs`), keyed by TMDB provider ID:

| Service | TMDB ID |
|---|---|
| Netflix | 8 |
| Amazon Prime Video | 9, 119 (TMDB uses two IDs for it) |
| Apple TV Store (rent/buy) | 2 |
| Apple TV (subscription, formerly Apple TV+) | 350 |
| YouTube | 192 |
| Crunchyroll | 283 |
| Rakuten Viki | 344 |
| iQIYI | 581 |
| WeTV | 623 |

These run in Vietnam but TMDB has no entry for them, so they can never appear: **FPT Play, VieON, Galaxy Play, TV360, K+, Danet, Bilibili**. Re-check the list with `GET /watch/providers/movie` when TMDB adds region `VN`.

1. Read every region in `watch/providers.results`. No regions → `[]`.
2. Per region, walk types in order `flatrate`, `free`, `ads`, `rent`, `buy`; map each entry to `{id, name, logo_url, type}`. Skip entries with `provider_id <= 0`, an empty `provider_name`, or an ID outside `allowedProviderIDs`.
3. `logo_url`: `""` → `null`; else `image base + "w92" + logo_path`.
4. De-duplicate by provider `id` across all regions. The kept `type` is the best one seen anywhere (flatrate beats rent), and the kept `display_priority` is the lowest seen.
5. Count how many regions offer each provider; a provider listed under several types in one region still counts once for it.
6. Sort by region count desc, then type order, then `display_priority` asc, then `id` asc. Regions are walked in sorted order and `id` breaks the last tie, so the same input always maps to the same list (Go iterates maps in random order).
7. Drop an entry whose name (trimmed, case-insensitive) is already in the list: TMDB gives one service several IDs — `9` and `119` are both "Amazon Prime Video" — and the ranking already put the better one first.
8. Keep the first **8** (`maxProviders`). Popular movies are offered by 40+ services worldwide, most of them in only one or two countries.

> The list answers "which services available in Vietnam carry this movie somewhere in the world". TMDB has no `VN` data, so it still cannot prove the title is licensed in Vietnam today — but every service named is one a customer here can subscribe to.

---

## 6. Error Handling

TMDB error body:

```json
{ "success": false, "status_code": 7, "status_message": "Invalid API key: You must be granted a valid key." }
```

- The client parses it into `*tmdb.APIError{HTTPStatus, Code, Message}` and wraps a sentinel so callers use `errors.Is`.

| HTTP (TMDB `status_code`) | Sentinel | Cause | Phase 1 (discover) | Phase 2 (details) |
|---|---|---|---|---|
| 401 (3, 7, 10) | `tmdb.ErrUnauthorized` | bad/suspended token | abort the **whole run**, log `Error` | abort the **whole run**, log `Error` |
| 404 (6) | `tmdb.ErrNotFound` | movie removed from TMDB | — | drop the movie, continue |
| 429 (25) | `tmdb.ErrRateLimited` | too many requests | retry | retry |
| 400 / 422 (5, 18, 22, 27), any other 4xx | `tmdb.ErrBadRequest` | bug in our params | fail this option, log `Error` | fail this option, log `Error` |
| 5xx / 504 (9, 11, 24), timeout, network, unreadable body | `tmdb.ErrUnavailable` | TMDB or network down | fail this option, keep old list | drop the movie, continue |

**Retry policy**

- Only `429` is retried: wait `Retry-After` seconds (default `5s` if absent), max **2** retries, respecting `ctx`.
- Still `429` after the retries → `tmdb.ErrRateLimited`: phase 1 fails the option (old list kept); phase 2 drops the movie.
- Everything else: no per-call retry. The next job tick (`REFRESH_CHECK_INTERVAL`) is the retry.

**Publish guards** (per option, in `refresher.go`)

- Phase 1 returns 0 IDs → don't publish; keep the old list; log `Warn` (the option's filters may be too strict to match anything).
- Phase 2 keeps < 50% of phase-1 IDs → don't publish; keep the old list; log `Warn`. This prevents a TMDB outage from shrinking a 40-movie list to 3.

---

## 7. Testing

- Fixtures in `internal/platform/tmdb/testdata/`. They are hand-written in TMDB's response shape (tests and scripts never call TMDB); replace them with real captures when someone runs the manual checks:
  ```text
  discover_page1.json
  movie_603_full.json          # all appends, providers in two regions, one actor without a photo
  movie_other_region.json      # watch/providers with a single non-VN region
  movie_missing_fields.json    # no poster, runtime 0, release_date ""
  error_401.json
  error_429.json
  ```
- Client tests: `httptest.Server` serves fixtures; assert headers, query params, error sentinels, 429 retry.
- Mapper tests (table-driven), at minimum:
  - title fallback to `original_title`; drop when both empty
  - `""` paths → `null` URLs; `runtime 0` → `null`; `release_date ""` → `null`
  - cast: billing order, capped at 5, and an actor without `profile_path` → `profile_url: null`
  - trailer preference (vi > en; official first)
  - provider merging across regions: ranking by region count, best type wins, cap at 8, deterministic order
- Refresher tests: fake client (consumer-side interface) for partial failures, 401 abort, publish guards.
- Manual check of VN provider coverage before promising the feature:
  ```bash
  curl -s "https://api.themoviedb.org/3/movie/603/watch/providers" \
    -H "Authorization: Bearer $TMDB_READ_TOKEN" | jq '.results.VN'
  ```

---

## 8. Attribution Checklist (frontend)

- TMDB logo + notice that the product uses the TMDB API but is not endorsed or certified by TMDB.
- JustWatch attribution wherever `providers` are shown.
- Commercial agreement with TMDB signed before production launch.