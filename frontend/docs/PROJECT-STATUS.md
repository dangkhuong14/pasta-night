# Project Status — Frontend

_Last updated: 2026-09-29_

## Done
- Next.js 16.3 app scaffolded per ARCHITECTURE §2 (TypeScript strict, Tailwind v4 tokens, shadcn `base-nova`, Playfair + Inter with the `vietnamese` subset); ADR `docs/decisions/001-frontend-stack.md`
- Landing `/`: options from `GET /options`, option cards with gradient + icon fallback (no brand photos yet), loading skeleton, error state
- Recommendations `/recommendations/[option]`: server-side shuffle per request, 6 per page, "GỢI Ý KHÁC" without repeats, sticky action bar, poster fallback, loading skeleton
- Detail sheet `?movie={id}`: opens via `history.pushState` (no round trip), full-height Base UI drawer, backdrop + cast after `GET /movies/{id}`, all SCREENS §3 hide rules, deep-link and back-gesture handling
- States: `OPTION_NOT_FOUND` → redirect, `CACHE_NOT_READY` → auto-retry on `Retry-After`, empty list → "Chưa có gợi ý", other errors → `ErrorState` + retry; `MOVIE_NOT_FOUND` → close sheet + refresh
- Footer with TMDB logo + notice (API Terms wording); "Nguồn: JustWatch" on the providers card
- Mock mode (no `NEXT_PUBLIC_API_BASE_URL`) with real TMDB sample data and edge cases
- Verified: lint (0 warnings), `tsc`, Prettier, `next build` (mock and API mode, build needs no backend); 390 px screenshots against `design/*.png` in both modes; error states against a stub API

- Cast photos in the detail sheet (2026-09-29), from `cast[].profile_url` (API v1.1); initials remain the fallback when TMDB has no photo
- Detail sheet closing actions (2026-09-30): "Chúc bạn ngon miệng." plus share (native share sheet, clipboard fallback) and "Chọn phim khác"
- Footer redesign + "Về Pásta Night" about sheet (2026-09-30): brand story, address, opening hours, hotline, social links and the shop CTA, opened from every page via `?about=1`. Content lives in `src/config/brand.ts`; every value still holds the mockup's placeholder

## In Progress
- None

## Next
- **Real shop details in `src/config/brand.ts`**: address, Google Maps URL, opening hours, hotline, and the Facebook / Instagram / TikTok URLs. Everything there now is the UX Pilot placeholder ("123 Đường Pasta", "0123 456 789"); empty links simply hide their icon, so nothing breaks until they land
- Brand assets: `public/images/options/{netflix-chill,solo,friends}.webp`, dish photos (drop in; no code change). Sizes are in DESIGN-SYSTEM §6 "Brand assets still missing"
- Pasta pairings per option in `src/config/brand.ts` (card hidden until then)
- `NEXT_PUBLIC_SHOP_URL` for "MUA NGAY" (PromoCard hidden until then)
- Deploy setup (host, env vars, backend `CORS_ALLOWED_ORIGINS` with the production origin)

## Known Issues
- Vietnamese genre names start with "Phim …", so the second chip on movie cards is usually truncated.
- ESLint 9 (pinned by create-next-app) is marked unsupported by npm; upgrade when `eslint-config-next` supports ESLint 10.

## Open Questions
- "CÓ MẶT TRÊN" lists only services available in Vietnam, but TMDB cannot confirm a title is licensed here. Should the card copy say so?
- Brand name confirmation: "Pásta Night" vs "Lusso Pasta" (SCREENS §0).
- Pasta pairing copy for `netflix-chill`, `solo`, `friends` (business).
- Should `netflix-chill` drop the Netflix provider filter, or the landing page hide options without a list? Needs a product decision (backend `configs/options.yaml`).
