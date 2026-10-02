# Project Status — Frontend

_Last updated: 2026-10-02_

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
- **Deployed to Vercel** (2026-10-02): `https://pasta-night.vercel.app`, functions in Singapore (`sin1`, pinned by `vercel.json`), against the Fly backend. Checked in headless Chrome at 390 px: the detail sheet loads its 9-slide carousel with the trailer playing and all cast photos, no console errors
- Real shop details (2026-10-02): address on Đường Tạ Quang Bửu (Quận 8), hotline 0939 465 840, opening hours Mon–Fri 19:00–24:00 and Sunday 12:00–24:00, in the footer and the about sheet
- Footer redesign + "Về Pásta Night" about sheet (2026-09-30): brand story, address, opening hours, hotline, social links and the shop CTA, opened from every page via `?about=1`. Content lives in `src/config/brand.ts`

- Detail hero carousel (2026-09-30, revised the same day after a device check): the trailer loops so YouTube's suggestion end screen never appears, a transparent layer over the player restores swiping (a cross-origin iframe keeps the touch), the position dots became buttons, and a mouse can drag the strip (`overflow-x-auto` offers a mouse nothing once the scrollbar is hidden). `MediaCarousel` replaces the static hero — the trailer in a muted-autoplay `youtube-nocookie` embed, then up to 8 stills, on a snap scroller with position dots. Images auto-advance every 5 s; the video slide, a touch, a hidden tab and `prefers-reduced-motion` each stop it.

## In Progress
- **Custom domain `pastanight.io.vn`**: the frontend is live on Vercel, but the domain's DNS at Tino Host still points to Tino's server (`103.142.24.254`). Waiting on the domain being added in Vercel and the DNS records changed at Tino (`docs/DEPLOYMENT.md` §3.3 step 6); Vercel issues the TLS certificate once DNS resolves

## Next
- **The rest of `src/config/brand.ts`**: the Google Maps URL (the address has no house number, so the "CHỈ ĐƯỜNG" button stays hidden), **Saturday opening hours** (the business gave Mon–Fri and Sunday only), and the Facebook / Instagram / TikTok URLs. Empty links simply hide their icon, so nothing breaks until they land
- Brand assets: `public/images/options/{netflix-chill,solo,friends}.webp`, dish photos (drop in; no code change). Sizes are in DESIGN-SYSTEM §6 "Brand assets still missing"
- Pasta pairings per option in `src/config/brand.ts` (card hidden until then)
- `NEXT_PUBLIC_SHOP_URL` for "MUA NGAY" (PromoCard hidden until then)

## Known Issues
- Vietnamese genre names start with "Phim …", so the second chip on movie cards is usually truncated.
- The trailer embed pulls ~1 MB of YouTube player JS on a page customers open over mobile data. If that hurts, switch to a still plus a play button and load the iframe on tap (SCREENS §3).
- The looping trailer keeps streaming while the sheet stays open on the video slide. Acceptable for a short trailer; revisit if it shows up in data usage.
- ESLint 9 (pinned by create-next-app) is marked unsupported by npm; upgrade when `eslint-config-next` supports ESLint 10.

## Open Questions
- "CÓ MẶT TRÊN" lists only services available in Vietnam, but TMDB cannot confirm a title is licensed here. Should the card copy say so?
- Brand name confirmation: "Pásta Night" vs "Lusso Pasta" (SCREENS §0).
- Pasta pairing copy for `netflix-chill`, `solo`, `friends` (business).
- Should `netflix-chill` drop the Netflix provider filter, or the landing page hide options without a list? Needs a product decision (backend `configs/options.yaml`).
