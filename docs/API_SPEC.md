# API Specification — v1

> Contract between the Next.js frontend and the Go backend. Any change here needs review from both FE and BE.
> Related: [DATABASE.md](./DATABASE.md) · [backend/ARCHITECTURE.md](../backend/docs/ARCHITECTURE.md)

## 1. Base URL & Versioning

| Environment | Base URL |
|---|---|
| Local | `http://localhost:8080/api/v1` |
| Production | `https://<api-domain>/api/v1` |

- The major version lives in the path (`/api/v1`).
- **Non-breaking** (stay on v1): new endpoints, new optional response fields, new error codes.
- **Breaking** (→ `/api/v2`, keep v1 until FE migrates): removing or renaming fields, changing types or meaning, changing status codes.
- The frontend must ignore unknown response fields.

---

## 2. Authentication

- **Public endpoints: no auth.** Customers are anonymous (they arrive by QR scan), so there is no login and **no token refresh flow**.
- **Admin endpoints:** header `X-Admin-Token: <ADMIN_TOKEN>`.
  - Missing or wrong token → `401 UNAUTHORIZED`.
  - `ADMIN_TOKEN` unset on the server → admin routes are disabled (`404 NOT_FOUND`).
  - Never call admin endpoints from the frontend.

---

## 3. Conventions

- JSON, UTF-8, `snake_case` field names.
- Timestamps: RFC 3339 UTC, e.g. `2026-09-27T10:00:00Z`.
- Nullable fields are always present as `null`, never omitted. Arrays are `[]`, never `null`.
- Text is Vietnamese (`vi-VN`). `title` may fall back to the original title; `overview` may be `""` when TMDB has no translation.
- Every response carries an `X-Request-ID` header. Include it in bug reports.
- Successful `GET` responses under `/api/v1` send `Cache-Control: public, max-age=300`. Exception: `/healthz` (§5.5) always sends `no-store`, since it reports live state.

---

## 4. Response Format

### Success

```json
{
  "data": {},
  "meta": {}
}
```

- `data`: object or array.
- `meta`: always an object (`{}` when empty).

### Error

```json
{
  "code": "OPTION_NOT_FOUND",
  "message": "option \"family\" does not exist",
  "details": { "option_id": "family" }
}
```

- `code`: stable and machine-readable. The frontend branches on it.
- `message`: English, for developers. The frontend shows its own Vietnamese copy.
- `details`: object with extra context, or `null`.

---

## 5. Endpoints

| Method | Path | Auth | Purpose |
|---|---|---|---|
| GET | `/options` | public | list viewing options |
| GET | `/options/{option_id}/recommendations` | public | ranked movies for one option |
| GET | `/movies/{movie_id}` | public | full detail of one movie |
| POST | `/admin/refresh` | admin | force a refresh |
| GET | `/healthz` (no `/api/v1` prefix) | public | liveness for Docker/orchestrator |

Values in the examples below are illustrative.

### 5.1 `GET /options`

- **Params:** none.
- **200:**

  ```json
  {
    "data": [
      { "id": "netflix-chill", "label": "Netflix & Chill", "description": "Cho hai người", "icon": "heart" },
      { "id": "solo", "label": "Một mình", "description": "Thời gian cho riêng bạn", "icon": "user" },
      { "id": "friends", "label": "Hội bạn", "description": "Xem cùng nhóm bạn", "icon": "users" }
    ],
    "meta": {}
  }
  ```

- Order of `data` = display order.
- **Errors:** `500 INTERNAL_ERROR`.

### 5.2 `GET /options/{option_id}/recommendations`

- **Path params:**

  | Name | Type | Rules |
  |---|---|---|
  | `option_id` | string | required; one of the IDs from `GET /options` |

- **200:** `data` is `MovieSummary[]` in rank order.

  ```json
  {
    "data": [
      {
        "id": 603,
        "title": "Ma Trận",
        "overview": "Một hacker phát hiện thế giới anh đang sống chỉ là một mô phỏng.",
        "poster_url": "https://image.tmdb.org/t/p/w500/<poster_path>.jpg",
        "release_year": 1999,
        "rating": 8.2,
        "runtime_minutes": 136,
        "genres": ["Phim Hành Động", "Phim Khoa Học Viễn Tưởng"],
        "providers": [
          { "id": 8, "name": "Netflix", "logo_url": "https://image.tmdb.org/t/p/w92/<logo_path>.jpg", "type": "flatrate" }
        ]
      }
    ],
    "meta": {
      "option_id": "friends",
      "total": 40,
      "fetched_at": "2026-09-27T10:00:00Z",
      "stale": false
    }
  }
  ```

- **MovieSummary** fields: `id`, `title`, `overview`, `poster_url`, `release_year`, `rating`, `runtime_minutes`, `genres`, `providers`. Types and nullability: DATABASE.md §2 MovieDetail.
- The full list (≤ 40) is returned. Shuffling and "Gợi ý khác" are done client-side.
- `meta.stale: true` → a refresh is overdue (DATABASE.md §5). The data is still safe to display.
- **Errors:** `404 OPTION_NOT_FOUND`, `503 CACHE_NOT_READY`, `500 INTERNAL_ERROR`.

### 5.3 `GET /movies/{movie_id}`

- **Path params:**

  | Name | Type | Rules |
  |---|---|---|
  | `movie_id` | int | required; > 0 |

- Only movies in a current recommendation list are available. There is no live TMDB lookup.
- **200:** `data` is `MovieDetail` = all MovieSummary fields + `original_title`, `tagline`, `backdrop_url`, `vote_count`, `directors`, `cast`, `trailer_url`.

  ```json
  {
    "data": {
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
      ]
    },
    "meta": { "fetched_at": "2026-09-27T10:00:05Z" }
  }
  ```

- **Errors:** `400 VALIDATION_ERROR`, `404 MOVIE_NOT_FOUND`, `500 INTERNAL_ERROR`.

### 5.4 `POST /admin/refresh`

- **Headers:** `X-Admin-Token: <ADMIN_TOKEN>` (required).
- **Query params:**

  | Name | Type | Rules |
  |---|---|---|
  | `option_id` | string | optional; omitted = all options |

- **Body:** none.
- **Behavior:**
  - Ignores `LIST_TTL` (always re-runs phase 1). Phase 2 still reuses details younger than `DETAIL_TTL`.
  - Runs asynchronously. A request for an option that is already refreshing joins the running refresh.
- **202:**

  ```json
  { "data": { "accepted": ["friends"] }, "meta": {} }
  ```

- **Example:**

  ```bash
  curl -X POST "http://localhost:8080/api/v1/admin/refresh?option_id=friends" \
       -H "X-Admin-Token: $ADMIN_TOKEN"
  ```

- **Errors:** `401 UNAUTHORIZED`, `404 OPTION_NOT_FOUND`, `404 NOT_FOUND` (admin disabled).

### 5.5 `GET /healthz`

- Infrastructure endpoint: **not** wrapped in the `{ data, meta }` envelope.
- Sends `Cache-Control: no-store` (§3 exception): `cache_ready` changes over time and must never be served stale by a cache or proxy sitting in front of the orchestrator's probe.
- **200** while the process is alive:

  ```json
  { "status": "ok", "cache_ready": true }
  ```

- `cache_ready`: `true` when every option has a cached list.

---

## 6. Error Codes

| HTTP | `code` | When | Frontend handling |
|---|---|---|---|
| 400 | `VALIDATION_ERROR` | malformed path or query param | generic error screen |
| 401 | `UNAUTHORIZED` | admin token missing or invalid | — (admin only) |
| 404 | `NOT_FOUND` | unknown route, or admin routes disabled | generic error screen |
| 404 | `OPTION_NOT_FOUND` | `option_id` not configured | redirect to `/` |
| 404 | `MOVIE_NOT_FOUND` | movie not in any current list | close detail sheet, refetch the list |
| 405 | `METHOD_NOT_ALLOWED` | wrong HTTP method | — (client bug) |
| 503 | `CACHE_NOT_READY` | first refresh not finished (new deploy, empty volume, schema bump); sends `Retry-After: 30` | "Đang chuẩn bị gợi ý…" + auto-retry |
| 500 | `INTERNAL_ERROR` | unexpected server error | error screen with retry button |

- TMDB outages never show up here: the request path doesn't call TMDB, and stale data is served with `meta.stale: true`.

---

## 7. Changelog

| Version | Date | Change |
|---|---|---|
| v1.0 | 2026-09-27 | initial contract |