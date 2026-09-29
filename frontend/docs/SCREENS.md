# Screens

> What each screen shows, where every piece of data comes from, and how every state looks.
> Visuals: `frontend/design/*.png` · Tokens & components: [DESIGN-SYSTEM.md](./DESIGN-SYSTEM.md) · API: [API_SPEC.md](../../docs/API_SPEC.md)

## 0. Design vs. API Gaps (read first)

The UX Pilot screens show things the v1 API does not provide. v1 decisions:

| Design element | In API? | v1 decision |
|---|---|---|
| Option card background photos | no | static brand assets `public/images/options/{option_id}.webp` |
| Option descriptions ("Phim lãng mạn, ấm áp cho hai người.") | yes (`description`) | update copy in `backend/configs/options.yaml`; FE never hardcodes it |
| "Gợi ý món mỳ hoàn hảo" pasta pairing | no | static `src/config/brand.ts`, keyed by `option_id`; no entry → hide card |
| "Mua ngay" shop link | no | env `NEXT_PUBLIC_SHOP_URL`; unset → hide `PromoCard` |
| Cast photos | yes, since API v1.1 | `cast[].profile_url` from TMDB. Initials avatars remain the fallback for the few actors TMDB has no photo of |
| TMDB + JustWatch attribution | n/a | **required by license, missing from the design** → add to `Footer` and the providers card |
| Brand name | design: "Pásta Night"; UX Pilot prompts: "Lusso Pasta" | use **Pásta Night** (confirm with the business) |
| Error state | not in export (only the "LỖI" tab) | spec in §4 below |

---

## 1. Landing — `/`

Reference: `design/01-landing.png`

**Layout (top → bottom)**

1. Background: `bg-background` + faint gold radial glow at the top.
2. `BrandWordmark` (centered).
3. Eyebrow: "GỢI Ý PHIM TỐI NAY".
4. Display H1: "Tối nay bạn xem phim cùng ai?"
5. `OptionCard` list, `space-y-4`, API order.
6. `Footer`.

**Data**

| UI | Source |
|---|---|
| option cards (order, count) | `GET /options` → `data[]` |
| card title | `label` |
| card description | `description` |
| card link | `/recommendations/{id}` |
| card background | `public/images/options/{id}.webp` (fallback: gradient + `icon`) |

**States**

| State | Behavior |
|---|---|
| loading | rarely visible (server-rendered); `loading.tsx`: wordmark + H1 + 3 card skeletons (176 px) |
| error | `ErrorState` with "Thử lại" → `retry()` |
| empty `data` | treat as error (misconfiguration) |

---

## 2. Recommendations — `/recommendations/[option]`

Reference: `design/02-recommendations.png`

**Layout**

1. `PageHeader`: back chevron → `/`, centered wordmark.
2. Page title "Dành riêng cho bạn" + eyebrow "NHỮNG BỘ PHIM HOÀN HẢO CHO TỐI NAY".
3. 2-column grid of `MovieCard`, `gap-4`, bottom padding so the last row clears the action bar.
4. `StickyActionBar` with `PrimaryButton` "GỢI Ý KHÁC".
5. `Footer`.

**Data**

| UI | Source |
|---|---|
| movie list | `GET /options/{option}/recommendations` → `data[]` (≤ 40) |
| poster | `poster_url` |
| rating badge | `rating` (1 decimal, e.g. `8.4`) |
| title | `title` |
| chips | first 2 of `genres` (1st `gold`, 2nd `neutral`) |
| card tap | set `?movie={id}` → opens detail sheet |

**"Gợi ý khác" behavior**

- Show **6** movies at a time (`DISPLAY_COUNT` in `features/recommendations/constants.ts`).
- On load: `page.tsx` shuffles the full list once on the server, per request, and the grid shows the first 6. (A first shuffle on the client would either break hydration or visibly swap the cards after JS loads.)
- On tap: show the next 6 from the shuffled order, no repeats. When fewer than 6 remain, reshuffle the full list and start over.
- Scroll to the top of the grid; cards re-animate.
- Fewer than 6 movies in total → hide the button.

**States**

| State | Behavior |
|---|---|
| loading | `loading.tsx` = `design/04-loading-state.png`: header, title bar, 2×2 card skeletons |
| `OPTION_NOT_FOUND` | `redirect("/")` |
| `CACHE_NOT_READY` | `ErrorState` title "Đang chuẩn bị gợi ý…", no button; auto-retry after `Retry-After` (default 30 s) |
| other errors | `ErrorState` (§4) |
| `meta.stale: true` | render normally; no UI change |
| `data: []` | `ErrorState` title "Chưa có gợi ý", button "Chọn lựa khác" → `/` |

---

## 3. Movie Detail — sheet over recommendations (`?movie={id}`)

Reference: `design/03-movie-detail.png`

**Open/close**

- Opens when `?movie=` is present. A card tap sets it with `history.pushState` (no server round trip). Close = remove the param: `history.back()` if the sheet was opened in-app, else `history.replaceState` (deep link).
- Mobile back gesture closes the sheet (param is in history).
- Direct link `/recommendations/friends?movie=603` opens the page with the sheet open.

**Layout**

1. Hero image (≈ 260 px) with bottom gradient; round close button (×) top-left.
2. Sheet body (`rounded-t-3xl`, drag handle):
   1. Title row: title (Detail title) + `RatingBadge`.
   2. Meta row: clock + runtime · calendar + year · first genre as `neutral` chip.
   3. `InfoCard` "CÓ MẶT TRÊN": `ProviderLogo` row (≤ 8) + "Nguồn: JustWatch" caption. The backend limits `providers` to services available in Vietnam (Netflix, Amazon Prime Video, Apple TV, YouTube, Crunchyroll, Viki, iQIYI, WeTV). TMDB has no `VN` licensing data, so this says the service carries the movie somewhere, not that it is streamable in Vietnam today (DATABASE.md §2).
   4. `PastaPairingCard` "GỢI Ý MÓN MỲ HOÀN HẢO".
   5. `SectionLabel` "NỘI DUNG" + overview.
   6. `SectionLabel` "DIỄN VIÊN CHÍNH" + `CastAvatar` row (max 5, horizontal scroll, photo or initials).
   7. `PromoCard` "Cần thêm gia vị?" · "Mua ngay Pásta Night để trải nghiệm phim thêm trọn vẹn." · "MUA NGAY".

**Data**

| UI | Source | Available |
|---|---|---|
| title, rating, runtime, year, genre, providers, overview | the `MovieSummary` already in the list | instantly |
| hero | `backdrop_url` → `poster_url` | after `GET /movies/{id}` (use `poster_url` meanwhile) |
| cast names + photos | `cast[].name`, `cast[].profile_url` | after `GET /movies/{id}` |
| pasta pairing | `brand.ts` → `pastaPairings[option_id]` | instantly |
| shop link | `NEXT_PUBLIC_SHOP_URL` | instantly |

**Formatting**

- Runtime: `105` → "1h 45p" · `60` → "1h" · `45` → "45p" (`lib/format.ts: formatRuntime`).
- Rating: always 1 decimal (`8` → "8.0").

**Hide rules** (never show empty sections)

| Condition | Hide |
|---|---|
| `runtime_minutes` / `release_year` null | that meta item |
| `providers` empty | whole "CÓ MẶT TRÊN" card |
| `overview` = `""` | "NỘI DUNG" section |
| `cast` empty | "DIỄN VIÊN CHÍNH" section |
| no pairing for option | `PastaPairingCard` |
| `NEXT_PUBLIC_SHOP_URL` unset | `PromoCard` |

**States**

| State | Behavior |
|---|---|
| detail loading | summary parts render; hero uses poster; cast row shows 4 circle skeletons |
| `MOVIE_NOT_FOUND` | close the sheet, `router.refresh()` the list |
| other detail errors | keep summary parts; hide cast; no error screen |
| `?movie=` not in current list | fetch detail; if it fails → close the sheet silently |

---

## 4. Error State (shared)

Not in the export; spec from the UX Pilot prompt.

- Centered in the column, vertically ~40% from top.
- Gold line icon (lucide `film`), 48 px.
- Title (Playfair 20 px): default "Rất tiếc, đã có lỗi xảy ra".
- Text (Body, muted): default "Vui lòng thử lại sau ít phút."
- `PrimaryButton` "THỬ LẠI" → `retry()` in `error.tsx` (Next.js 16.3: re-fetches and re-renders the segment).
- Copy per API error code lives in `src/lib/error-copy.ts` (never show the API `message`).

---

## 5. Footer (all pages)

- "— PÁSTA NIGHT —" (thin lines, muted, 10 px, tracking wide).
- Below, 10 px muted: TMDB logo + TMDB's required notice, kept in its **original English wording** (don't translate a legal notice). Current text, from TMDB's API Terms of Use: "This website uses TMDB and the TMDB APIs but is not endorsed, certified, or otherwise approved by TMDB." Logo: `public/images/brand/tmdb-logo.svg` (official "blue_short" SVG), smaller than the Pásta Night wordmark as the terms require.