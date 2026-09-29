# Frontend Architecture

> How the frontend is structured and why.
> Related: [PROJECT-RULES.md](./PROJECT-RULES.md) · [SCREENS.md](./SCREENS.md) · [DESIGN-SYSTEM.md](./DESIGN-SYSTEM.md) · [API_SPEC.md](../../docs/API_SPEC.md)

## 1. Overview

- **What:** a mobile-first website opened from the QR code on a Pásta Night order card. Customer picks a viewing option → sees movie recommendations → opens a movie's detail.
- **Three screens:** landing (`/`), recommendations (`/recommendations/[option]`), movie detail (sheet, `?movie={id}`).
- **Constraint:** users arrive on phones, often on mobile data → server-render the first paint, ship little JavaScript.

### System diagram

```text
QR card ─scan─▶ Phone browser
                   │  HTML (server-rendered)          │  JSON (client, sheet only)
                   ▼                                  ▼
             Next.js server ───── REST /api/v1 ────▶ Go API ◀── images: TMDB CDN
             (Server Components)                     (see backend/docs)
```

### Tech stack

| Choice | Why |
|---|---|
| Next.js (App Router) | server rendering for fast first paint after a QR scan; file-based routes match our 2 pages |
| TypeScript `strict` | API types checked against API_SPEC |
| Tailwind v4 | CSS-variable tokens (`@theme`) match the design system |
| shadcn/ui (`base-nova` style, Base UI primitives) | owned source components; `Drawer` (bottom sheet), `Button`, `Skeleton`; `cn` package for class merging |
| lucide-react | shadcn's default icons; option `icon` values from the API are lucide names |
| `next/font` | self-hosted Playfair Display + Inter with the `vietnamese` subset |
| `next/image` | responsive, lazy-loaded TMDB posters |

---

## 2. Folder Structure

```text
frontend/
├── src/
│   ├── app/
│   │   ├── layout.tsx                 # <html lang="vi">, fonts, column wrapper, Footer
│   │   ├── globals.css                # Tailwind + tokens (DESIGN-SYSTEM §2)
│   │   ├── page.tsx                   # Landing (server): getOptions → OptionList
│   │   ├── loading.tsx
│   │   ├── error.tsx
│   │   └── recommendations/[option]/
│   │       ├── page.tsx               # server: getRecommendations → grid + detail sheet
│   │       ├── loading.tsx            # skeleton (design/04-loading-state.png)
│   │       └── error.tsx
│   ├── features/
│   │   ├── options/
│   │   │   ├── OptionCard.tsx
│   │   │   ├── OptionList.tsx
│   │   │   ├── OptionIcon.tsx         # API `icon` name → lucide icon
│   │   │   └── context.md
│   │   ├── recommendations/
│   │   │   ├── RecommendationGrid.tsx # client: shuffle + paging
│   │   │   ├── MovieCard.tsx
│   │   │   ├── CacheNotReady.tsx      # client: auto-retry
│   │   │   ├── useShuffle.ts
│   │   │   ├── constants.ts           # DISPLAY_COUNT
│   │   │   └── context.md
│   │   └── movie-detail/
│   │       ├── MovieDetailSheet.tsx   # client: ?movie= ↔ open state
│   │       ├── ProvidersCard.tsx
│   │       ├── PastaPairingCard.tsx
│   │       ├── CastList.tsx
│   │       ├── PromoCard.tsx
│   │       ├── useMovieDetail.ts
│   │       └── context.md
│   ├── components/
│   │   ├── ui/                        # shadcn-generated (don't edit heavily)
│   │   └── common/                    # BrandWordmark, PageHeader, GenreChip, RatingBadge,
│   │                                  # PrimaryButton, OutlineButton, StickyActionBar,
│   │                                  # InfoCard, ErrorState, Footer, SectionLabel
│   ├── lib/
│   │   ├── api.ts                     # ONLY backend caller
│   │   ├── api-types.ts               # mirrors API_SPEC
│   │   ├── mock-data.ts               # used when NEXT_PUBLIC_API_BASE_URL is unset
│   │   ├── error-copy.ts              # API error code → Vietnamese copy
│   │   ├── format.ts                  # formatRuntime, formatRating, initials
│   │   ├── shuffle.ts                 # Fisher–Yates, shared by page.tsx and useShuffle
│   │   ├── assets.ts                  # server-only: optional brand photo in public/?
│   │   └── utils.ts                   # cn()
│   └── config/
│       └── brand.ts                   # pasta pairings by option_id, brand strings
├── public/images/
│   ├── options/{option_id}.webp       # optional; missing → gradient + icon
│   ├── dishes/{dish_id}.webp          # optional; missing → no image slot
│   └── brand/tmdb-logo.svg            # official TMDB logo (attribution)
├── design/                            # UX Pilot exports — reference only, not shipped
├── docs/
│   └── decisions/                     # frontend ADRs
├── next.config.ts                     # images.remotePatterns: image.tmdb.org
└── .env.example
```

- `features/*` = one screen area each; `components/common` = shared presentational pieces; `lib` = non-UI logic.

---

## 3. Layer Architecture

```text
app/ route (Server Component)      fetch page data, handle API errors, compose features
   ↓
features/*/components              render UI from props
   ↓
features/*/hooks (client only)     local interaction state: useShuffle, useMovieDetail
   ↓
lib/api.ts                         HTTP, envelope unwrap, ApiError, mock fallback
   ↓
Go API (/api/v1)
```

Dependency rules:

- `app/` → `features/`, `components/`, `lib/`, `config/`.
- `features/X` → `components/`, `lib/`, `config/`. **Never** `features/Y`.
- `components/common` → `lib/utils` only (pure presentational).
- `lib/` → nothing in `app/`, `features/`, or `components/`.

Example: the grid doesn't import the sheet. Cards link to `?movie=603`; `recommendations/[option]/page.tsx` renders both `<RecommendationGrid>` and `<MovieDetailSheet>`.

---

## 4. Rendering & Data Flow per Route

| Route | Rendering | Data | Client JS |
|---|---|---|---|
| `/` | Server Component | `getOptions()`, `revalidate: 300` | none (cards are `<Link>`) |
| `/recommendations/[option]` | Server Component | `getRecommendations(option)`, `revalidate: 300`; list shuffled per request | grid paging, sheet |
| sheet `?movie={id}` | Client Component | summary from the page's list + `getMovie(id)` on open | yes |

- Both pages call `connection()` and render per request: `next build` never needs the API, and every visit gets its own shuffle. API responses still come from the fetch cache (`revalidate: 300`). Only `200` responses are cached, so `CACHE_NOT_READY` is retried for real.
- The sheet opens without a server round trip: `MovieCard` sets `?movie=` with `history.pushState`, which Next.js syncs into `useSearchParams`.

Error handling in `recommendations/[option]/page.tsx`:

```tsx
const { option } = await params; // params is a Promise in Next.js 15+
try {
  const { data, meta } = await getRecommendations(option);
  // render grid + sheet
} catch (e) {
  if (e instanceof ApiError && e.code === "OPTION_NOT_FOUND") redirect("/");
  if (e instanceof ApiError && e.code === "CACHE_NOT_READY") return <CacheNotReady retryAfter={e.retryAfter} />;
  throw e; // → error.tsx
}
```

---

## 5. State Management

No global store. State lives where it belongs:

| State | Where | Why |
|---|---|---|
| selected option | URL path `/recommendations/{id}` | shareable, back button works |
| open movie | URL query `?movie={id}` (set with `history.pushState`) | back gesture closes the sheet; deep-linkable |
| initial order | shuffled in `recommendations/[option]/page.tsx` per request | server and client render the same cards (no hydration mismatch) |
| current page + later reshuffles | `useState` in `useShuffle` (`RecommendationGrid`) | throwaway; resets on reload (fine) |
| movie detail response | `useMovieDetail` (fetch on open, in-memory per session) | only needed while the sheet is open |

---

## 6. Communication

| Between | How | Notes |
|---|---|---|
| Next.js server → Go API | `fetch` in Server Components | no CORS involved |
| Browser → Go API | `fetch` from `useMovieDetail` | FE origin must be in backend `CORS_ALLOWED_ORIGINS` |
| Browser → TMDB CDN | `next/image` | `image.tmdb.org` in `images.remotePatterns` |
| Features | props + URL | no context providers, no event bus |

---

## 7. Configuration

| Env var | Example | Purpose |
|---|---|---|
| `NEXT_PUBLIC_API_BASE_URL` | `http://localhost:8080/api/v1` | backend base URL; unset → mock data |
| `NEXT_PUBLIC_SHOP_URL` | `https://…` | "Mua ngay" link; unset → `PromoCard` hidden |

- `NEXT_PUBLIC_*` values are inlined at build time → changing them requires a rebuild.
- Keep `.env.example` in sync.