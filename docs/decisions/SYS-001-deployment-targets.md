# SYS-001 — Deployment targets: Vercel for the frontend, Fly.io for the backend

## Status

Accepted (2026-10-02)

## Context

Pásta Night has never been deployed. Customers reach it by scanning a QR code printed on the order card, so it needs a stable public domain before any card is printed, and it is opened on phones, mostly over mobile data, in Vietnam.

The two sides have different hosting needs:

- **Backend** (`backend/`): a single Go binary whose only state is a JSON cache on disk (backend ADR 001). It needs a **persistent volume**: without it, every restart serves `503 CACHE_NOT_READY` while ~126 TMDB calls rebuild the cache. It must run as **exactly one instance**, because replicas would each call TMDB, and it must stay **awake**, because the refresh job runs on a 1 h ticker inside the process.
- **Frontend** (`frontend/`): Next.js 16 with per-request rendering (`connection()`) and image optimization for TMDB artwork, so it needs a Node server, not a static export. `NEXT_PUBLIC_API_BASE_URL` is inlined at build time.

The shop has its own domain.

## Decision

- **Frontend on Vercel**, function region Singapore (`sin1`), serving `<domain>`.
- **Backend on Fly.io**, region Singapore (`sin`), serving `api.<domain>`: one `shared-cpu-1x` machine, always on (`auto_stop_machines = "off"`), with one 1 GB volume mounted at `/app/data/cache`. Configuration in `backend/fly.toml`; secrets through `fly secrets`.
- Singapore for both, as the region closest to Vietnam that both providers offer, which also keeps the Vercel → Fly hop short.
- First-time setup, operations and troubleshooting: `docs/DEPLOYMENT.md`.

## Consequences

- No container or server to manage for the frontend; previews and rollbacks come with Vercel.
- The existing `backend/Dockerfile` deploys unchanged.
- Two providers, so the first deploy has an ordering constraint: backend first (the frontend build needs its URL), then the frontend, then the backend's `CORS_ALLOWED_ORIGINS` (it needs the frontend origin). `docs/DEPLOYMENT.md` §3 spells it out.
- Only the browser-side detail sheet goes through CORS; the landing page and grid are fetched by the Next.js server. A CORS mistake therefore breaks the detail sheet alone, which is easy to miss.
- The backend is a single point of failure. Fly restarts it on failed health checks, and losing the volume costs minutes of `503`, not data — the cache is rebuilt from TMDB.
- An always-on machine is billed around the clock. Accepted: a customer who just scanned the QR code should not wait for a cold start.
- TMDB artwork is proxied by Vercel image optimization and counts toward its quota.
- Preview deployments' `*.vercel.app` origins are not in the backend's CORS list, so the detail sheet fails on previews unless one is added.

## Alternatives considered

- **Backend on Vercel serverless or another functions platform:** no persistent disk, and no long-running process for the refresh ticker. Would need a database or object store plus a scheduled job — the storage that backend ADR 001 deliberately avoided.
- **Fly.io with scale-to-zero:** cheaper, but the refresh job stops while asleep, and the first scan after idle waits for a cold start.
- **One VPS for both, with Docker Compose and Caddy:** cheapest and avoids the cross-provider ordering, but leaves OS updates, TLS, builds and monitoring to the shop, which has no one to run them.
- **Vercel for the frontend, a VPS for the backend:** full control of the disk, but the same operations burden for the backend alone, plus TLS to set up by hand.
- **Frontend on Fly.io alongside the backend:** one provider, but Next.js image optimization, previews and rollbacks would need to be built and run by hand.
