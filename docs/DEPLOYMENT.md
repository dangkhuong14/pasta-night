# Deployment

How Pásta Night runs in production, how to set it up the first time, and how to operate it afterwards.

- **Frontend** (`frontend/`): Vercel, region Singapore.
- **Backend** (`backend/`): Fly.io, region Singapore, one machine with one volume.
- **Domain**: the shop's own (`pastanight.io.vn`, registered with Tino Host) serves the frontend. The backend uses Fly's own hostname, `pasta-night-api.fly.dev`, whose TLS certificate Fly provides.

**Current production** (first deployed 2026-10-02):

| | URL |
|---|---|
| Frontend | `https://pastanight.io.vn` (DNS at Tino Host) · `https://pasta-night.vercel.app` (Vercel's own URL) |
| Backend | `https://pasta-night-api.fly.dev` (API base `…/api/v1`) |

In the steps below, `<api-host>` means `pasta-night-api.fly.dev`; `api.<domain>` appears only in the optional §3.2.

Why these hosts, and what was ruled out: [`decisions/SYS-001-deployment-targets.md`](decisions/SYS-001-deployment-targets.md).

Throughout this guide, `<domain>` is the shop's domain (for example `pastanight.vn`) and `pasta-night-api` is the Fly app name from `backend/fly.toml`.

---

## 1. Topology

```text
customer's phone (QR code → https://<domain>/)
   │
   ├── page HTML ───────────▶ Vercel · Next.js server (sin1)
   │                             │  GET /options, /options/{id}/recommendations
   │                             │  server-side, cached 5 min — no CORS involved
   │                             ▼
   └── movie detail ────────▶ Fly.io · Go API (sin) ──▶ TMDB API
       browser-side fetch          │   background refresh, 1 h ticker
       — needs CORS                ▼
                              volume /app/data/cache (JSON cache, no database)
```

Two things in this picture drive most of the setup:

- **The backend must keep its disk.** All state is a JSON cache on a volume. Losing it is not data loss — it is rebuilt from TMDB — but until then the API answers `503 CACHE_NOT_READY` for a few minutes (~126 TMDB calls).
- **Only the movie detail sheet calls the API from the browser.** The landing page and the grid are fetched by the Next.js server. So a wrong CORS setting does not break the site; it silently breaks only the detail sheet (no cast, no carousel). Check it explicitly (§4).

---

## 2. Before the first deploy

| | Item | Why |
|---|---|---|
| ☐ | **TMDB commercial terms** checked | Pásta Night is a business; TMDB's free API terms cover non-commercial use. Open question in `backend/docs/PROJECT-STATUS.md`. |
| ☐ | Repository **pushed to GitHub** | Vercel builds from the Git repository. (Fly deploys from your machine, no Git needed.) |
| ☐ | Backend run locally once on the current code | Cache schema v5 (`media` for the carousel) has not run against a real cache outside tests. Confirm `GET /api/v1/movies/{id}` returns `media`. |
| ☐ | Carousel tried on a **real phone** | The trailer slide's swipe handling cannot be tested in headless Chrome (`frontend/src/features/movie-detail/context.md`). |
| ☐ | DNS access for `<domain>` | You will add records for the apex and for `api`. |
| ☐ | Accounts: Fly.io (`flyctl` installed, `fly auth login`), Vercel | |
| ☐ | A password manager entry ready for `ADMIN_TOKEN` | Fly cannot show a secret back once set. |

Not blockers, but visible to customers: Saturday opening hours, the Google Maps link and the social URLs in `frontend/src/config/brand.ts` (all `TODO(business)`).

**Print the QR code last**, once `https://<domain>/` works end to end. It should encode the landing page URL.

---

## 3. First-time setup

The order matters. The frontend build needs the backend URL (it is inlined at build time), and the backend's CORS needs the frontend origin. Deploying in this order breaks the loop.

### 3.1 Backend on Fly.io

From `backend/`:

```bash
fly apps create pasta-night-api                       # name is global; if taken, change `app` in fly.toml
fly volumes create pasta_cache --region sin --size 1  # 1 GB is the minimum; the cache is under 1 MB
```

Set the secrets. They never go in `fly.toml` or Git:

```bash
# Reads the token from your local .env without printing it.
fly secrets set TMDB_READ_TOKEN="$(grep '^TMDB_READ_TOKEN=' .env | cut -d= -f2-)"

# Generate ADMIN_TOKEN (openssl rand -hex 32), save it in your password manager, then:
fly secrets set ADMIN_TOKEN="<the token you saved>"
```

Deploy and pin to one machine:

```bash
fly deploy
fly scale count 1
```

`CORS_ALLOWED_ORIGINS` is still `http://localhost:3000` at this point. That is expected; §3.4 fixes it.

Check: `curl https://pasta-night-api.fly.dev/healthz` → `{"status":"ok","cache_ready":false}`. `cache_ready` turns `true` once the first refresh finishes (§3.5).

### 3.2 (Optional) a custom hostname for the backend

Skipped in production: customers never see the API hostname, so `pasta-night-api.fly.dev` is used as is. To move it to `api.<domain>` later:

```bash
fly certs add api.<domain>
fly certs show api.<domain>   # prints the exact DNS records to add
```

At your DNS provider, add a `CNAME` record: `api` → `pasta-night-api.fly.dev`. Wait until `fly certs show` reports the certificate as issued.

### 3.3 Frontend on Vercel

1. **New Project** → import the GitHub repository.
2. **Root Directory**: `frontend`. Framework preset: Next.js (detected). Leave the build settings at their defaults.
3. **Environment Variables** (scope: Production):

   | Name | Value |
   |---|---|
   | `NEXT_PUBLIC_API_BASE_URL` | `https://pasta-night-api.fly.dev/api/v1` |
   | `NEXT_PUBLIC_SHOP_URL` | the online shop link, or leave empty to hide the "MUA NGAY" card |

4. Region: nothing to set. `frontend/vercel.json` pins functions to Singapore (`sin1`), next to the backend. The import screen has no region picker, so it lives in code.
5. Deploy.
6. **Settings → Domains**: add `<domain>`, and `www.<domain>` set to redirect to it. Add the DNS records Vercel displays.

`NEXT_PUBLIC_*` values are baked into the build. Changing one later requires a **Redeploy**, not just a save.

Preview deployments get their own `*.vercel.app` origin, which the backend's CORS will not allow. On previews the detail sheet fails; that is expected unless you add the preview origin too.

### 3.4 Close the CORS loop

In `backend/fly.toml`, set the real origin:

```toml
CORS_ALLOWED_ORIGINS = "https://<domain>"
```

Rules (enforced at startup by `validateOrigin` in `backend/internal/platform/config/config.go`): scheme and host only, lowercase, **no trailing slash**, no path, no `*`. Several origins are comma-separated. If `www` redirects to the apex, only the apex is needed — the redirect happens before the page loads.

Production allows three origins: `https://www.pastanight.io.vn` (the primary in Vercel; the apex redirects to it), `https://pastanight.io.vn`, and `https://pasta-night.vercel.app`, the project's own Vercel URL. **If the primary domain is switched in Vercel, check this list**: the browser's origin is whichever host the redirect lands on. Preview URLs stay excluded.

```bash
fly deploy
```

Commit the `fly.toml` change.

### 3.5 Warm the cache

The first boot starts with an empty volume, so the API answers `503 CACHE_NOT_READY` with `Retry-After: 30` while the first refresh runs. The frontend already handles this: it shows "Đang chuẩn bị gợi ý…" and retries on its own.

The refresh starts by itself at boot. To start it explicitly:

```bash
curl -X POST "https://<api-host>/api/v1/admin/refresh" -H "X-Admin-Token: <ADMIN_TOKEN>"
# → 202 Accepted
```

Poll `https://<api-host>/healthz` until `"cache_ready": true` (a few minutes).

---

## 4. Verify

Run all of these after the first deploy, and after any change to hosting, domain or CORS.

| # | Check | Expected |
|---|---|---|
| 1 | `curl https://<api-host>/healthz` | `{"status":"ok","cache_ready":true}` |
| 2 | `curl https://<api-host>/api/v1/options` | the 3 options |
| 3 | `curl https://<api-host>/api/v1/movies/<id>` | the body contains `"media"`, not `"trailer_url"` — the backend runs the current code |
| 4 | Open `https://<domain>/` on a **real phone**: option → grid → open a movie | the carousel plays the trailer muted and swipes, and the cast photos load. **This is the CORS check**: the grid works even with CORS wrong; the sheet does not |
| 5 | Browser devtools → Network, on the detail sheet | no CORS errors |
| 6 | `curl -sI -H "Origin: https://example.com" https://<api-host>/api/v1/options` | no `Access-Control-Allow-Origin` header |
| 7 | `fly machines list` | exactly **one** machine |
| 8 | `fly machine restart <id>`, then check #1 at once | `cache_ready: true` right away — the volume is mounted |
| 9 | `fly logs` | `"msg":"cache loaded"` with non-zero `lists` and `details` |

---

## 5. Day-to-day operations

| Task | How | Notes |
|---|---|---|
| Deploy backend code | `cd backend && fly deploy` | |
| Deploy frontend code | push to `main` | Vercel builds automatically |
| Change `NEXT_PUBLIC_*` | edit in Vercel, then **Redeploy** | build-time values |
| Change `configs/options.yaml` | `fly deploy`, then force a refresh of the changed option | The file is baked into the image. A changed option keeps its old list until `LIST_TTL` (7 days) unless forced: `POST /api/v1/admin/refresh?option_id=<id>`. A brand-new option refreshes on its own. |
| Bump `cache.SchemaVersion` | `fly deploy` at a quiet hour | Every cached file becomes invalid: `503` until the refresh finishes. The shop opens at 19:00 on weekdays, so deploy in the morning. |
| Force a refresh | `curl -X POST https://<api-host>/api/v1/admin/refresh -H "X-Admin-Token: …"` | add `?option_id=<id>` for one option |
| Rotate `ADMIN_TOKEN` or `TMDB_READ_TOKEN` | `fly secrets set NAME=…` | restarts the machine; the cache survives on the volume |
| Logs | `fly logs` · Vercel → project → Logs | backend logs are JSON (`log/slog`) |
| Roll back backend | `fly releases`, then `fly deploy --image <image of the good release>` | Rolling back across a schema bump means another `503` window: the older binary rejects the newer cache files |
| Roll back frontend | Vercel → Deployments → Instant Rollback | |

---

## 6. Troubleshooting

| Symptom | Likely cause | Fix |
|---|---|---|
| Landing page and grid work; the detail sheet shows no cast and no carousel; the console shows a CORS error | `CORS_ALLOWED_ORIGINS` does not match the page's origin: `www` vs apex, `http` vs `https`, trailing slash | correct it in `fly.toml`, `fly deploy` (§3.4) |
| Backend will not start: `invalid environment: CORS_ALLOWED_ORIGINS: invalid origin …` | path, query or trailing slash in an origin | `https://<domain>`, nothing after the host |
| Backend will not start: `TMDB_READ_TOKEN: is required` | secret missing | `fly secrets set TMDB_READ_TOKEN=…` |
| "Đang chuẩn bị gợi ý…" never goes away | the first refresh keeps failing | `fly logs`: a TMDB `401` means a wrong token; anything else, see `backend/docs/TMDB_INTEGRATION.md` |
| Detail sheet shows one still image, no trailer | the backend image predates the carousel and still serves `trailer_url` | `fly deploy` |
| Frontend still calls the old backend after changing the env var | `NEXT_PUBLIC_*` is inlined at build time | Vercel → Redeploy |
| Every restart brings minutes of `503` | the volume is not mounted, so the cache lives on the machine's temporary disk | `fly volumes list`; `[[mounts]] destination` must be `/app/data/cache` |
| Two machines in `fly machines list` | Fly's default for availability | `fly scale count 1` |
| `api.themoviedb.org` does not resolve | Vietnamese ISP DNS blocks it | Affects local development in Vietnam only, never Fly in Singapore. Locally, use DNS 1.1.1.1 or a VPN |

---

## 7. Production configuration

**Backend** (`backend/fly.toml` `[env]` and `fly secrets`; full list in `backend/docs/ARCHITECTURE.md` §6):

| Variable | Where | Production value |
|---|---|---|
| `TMDB_READ_TOKEN` | `fly secrets` | TMDB read access token |
| `ADMIN_TOKEN` | `fly secrets` | 64 hex characters (`openssl rand -hex 32`); empty disables the admin endpoint |
| `CORS_ALLOWED_ORIGINS` | `fly.toml` | `https://<domain>` |
| `PORT` | `fly.toml` | `8080` |
| `CACHE_DIR` | Dockerfile | `/app/data/cache`, the volume mount point |
| everything else | defaults | `LIST_TTL=168h`, `DETAIL_TTL=336h`, `REFRESH_CHECK_INTERVAL=1h`, `DISCOVER_PAGES=2`, `TMDB_CONCURRENCY=5` |

**Frontend** (Vercel → Environment Variables; reference in `frontend/docs/ARCHITECTURE.md` §7):

| Variable | Production value |
|---|---|
| `NEXT_PUBLIC_API_BASE_URL` | `https://<api-host>/api/v1` |
| `NEXT_PUBLIC_SHOP_URL` | shop link, or empty |

---

## 8. Limits to watch

- **Vercel image optimization.** TMDB posters, backdrops and cast photos go through Next.js image optimization (`next.config.ts`), and each optimized image counts toward the Vercel plan's quota. Watch it if traffic grows.
- **The trailer embed** loads ~1 MB of YouTube player script, and the trailer loops while the sheet stays open on it. A lighter option is in `frontend/docs/SCREENS.md` §3.
- **One backend machine** is a single point of failure by design (backend ADR 001: replicas would each call TMDB). If it is down, the pages fail; Fly restarts it when health checks fail.
- **No volume backup is needed.** The cache is disposable and rebuilt from TMDB.
