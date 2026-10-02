# Design System — Pásta Night

> Visual source of truth for the frontend. Reference screens: `frontend/design/*.png` (UX Pilot export).
> Use **tokens**, never raw hex values, in components.
> Related: [SCREENS.md](./SCREENS.md) · [PROJECT-RULES.md](./PROJECT-RULES.md)

## 1. Principles

- **Cinematic & premium:** near-black canvas, gold used sparingly for brand, titles, and the one primary action per screen.
- **Mobile-first:** designed at ~375 px. On larger screens, content stays in a centered column (`max-w-md`, 28 rem).
- **Dark only:** there is no light theme.
- **Vietnamese copy:** fonts must include the `vietnamese` subset, or diacritics fall back to a system font.

---

## 2. Color Tokens

`✓` = specified in the UX Pilot prompt. `≈` = measured from the export; verify in UX Pilot before launch.

| Token (CSS var) | Value | | Usage |
|---|---|---|---|
| `--background` | `#0A0A0A` | ✓ | page background |
| `--foreground` | `#F5F5F5` | ≈ | primary text |
| `--card` | `#141414` | ≈ | cards, sheet surface |
| `--muted` | `#1C1C1C` | ≈ | skeletons, neutral chips, rating badge |
| `--muted-foreground` | `#A3A3A3` | ≈ | secondary text, section labels |
| `--border` | `rgb(255 255 255 / 0.08)` | ≈ | card borders, dividers |
| `--primary` | `#D4AF37` | ✓ | gold: wordmark, titles, CTA, active chip |
| `--primary-foreground` | `#0A0A0A` | ≈ | text on gold buttons |
| `--ring` | `#D4AF37` | | focus ring |
| `--gold-bright` | `#F5D77A` | ≈ | hover/pressed gold, star fill |
| `--gold-soft` | `rgb(212 175 55 / 0.12)` | ≈ | gold chip background, promo card tint |
| `--gold-line` | `rgb(212 175 55 / 0.35)` | ≈ | gold borders (chip, pairing card, outline button) |

### Implementation (Tailwind v4 + shadcn/ui)

```css
/* src/app/globals.css (excerpt). Dark-only: :root holds the dark palette. */
:root {
  --background: #0a0a0a;
  --foreground: #f5f5f5;
  --card: #141414;
  --card-foreground: #f5f5f5;
  --muted: #1c1c1c;
  --muted-foreground: #a3a3a3;
  --primary: #d4af37;
  --primary-foreground: #0a0a0a;
  --border: rgb(255 255 255 / 0.08);
  --ring: #d4af37;
  --radius: 0.75rem;

  /* brand extras (not part of shadcn) */
  --gold-bright: #f5d77a;
  --gold-soft: rgb(212 175 55 / 0.12);
  --gold-line: rgb(212 175 55 / 0.35);
}

@theme inline {
  /* shadcn's generated block already maps --background, --primary, ... */
  --color-gold-bright: var(--gold-bright);
  --color-gold-soft: var(--gold-soft);
  --color-gold-line: var(--gold-line);
  --font-serif: var(--font-playfair);
  --font-sans: var(--font-inter);
  --shadow-gold: 0 0 24px rgb(212 175 55 / 0.35);
}
```

Resulting classes: `bg-background`, `bg-card`, `text-primary`, `text-muted-foreground`, `border-border`, `bg-gold-soft`, `border-gold-line`, `shadow-gold`, `font-serif`.

---

## 3. Typography

```tsx
// src/app/layout.tsx
const playfair = Playfair_Display({ subsets: ["latin", "vietnamese"], variable: "--font-playfair" });
const inter = Inter({ subsets: ["latin", "vietnamese"], variable: "--font-inter" });
```

| Role | Font | Size / weight | Style | Example |
|---|---|---|---|---|
| Wordmark | Playfair | 20 px / 400 | uppercase, `tracking-[0.3em]`, gold | PÁSTA NIGHT |
| Eyebrow / section label | Inter | 11 px / 500 | uppercase, `tracking-[0.2em]`, muted | GỢI Ý PHIM TỐI NAY · NỘI DUNG |
| Display (landing H1) | Playfair | 28 px / 400 | centered, foreground, `leading-tight` | Tối nay bạn xem phim cùng ai? |
| Page title (H2) | Playfair | 24 px / 400 | gold | Dành riêng cho bạn |
| Detail title | Playfair | 26 px / 400 | gold | Hương Vị Tình Yêu |
| Card title (option) | Playfair | 20 px / 400 | gold | Netflix & Chill |
| Movie card title | Inter | 14 px / 600 | foreground, 1 line, truncate | |
| Body | Inter | 14 px / 400 | `leading-relaxed`, foreground/90 | synopsis |
| Caption / meta | Inter | 12 px / 400 | muted | 1h 45p · 2023 |
| Button | Inter | 13 px / 600 | uppercase, `tracking-[0.15em]` | GỢI Ý KHÁC |

---

## 4. Spacing, Radius, Elevation, Motion

- **Spacing scale:** 4 · 8 · 12 · 16 · 20 · 24 · 32. Page padding `px-5`; grid gap `gap-4`; section gap `space-y-6` to `space-y-8`.
- **Radius:** cards `rounded-xl` (12 px) · posters `rounded-lg` · chips, badges, buttons `rounded-full` · sheet top `rounded-t-3xl`.
- **Elevation:** no heavy shadows. Use `border-border`; the primary CTA gets `shadow-gold`.
- **Motion:**
  - Card enter: fade + 8 px slide-up, 200 ms, stagger 40 ms.
  - Skeleton shimmer: 1.5 s linear infinite.
  - Press: `active:scale-[0.98]`.
  - Respect `prefers-reduced-motion`: disable shimmer, slide, and stagger (`motion-reduce:` variants).

---

## 5. Components

| Component | Spec | Used on |
|---|---|---|
| `BrandWordmark` | "PÁSTA NIGHT" per Wordmark style | all headers, footer |
| `SectionLabel` | Eyebrow style; optional gold variant (`text-primary`) | detail sections, landing eyebrow |
| `OptionCard` | full width, ~176 px tall, `rounded-xl`, 1 px `border-border`; background photo + bottom gradient `from-black/90 via-black/40 to-transparent`; title (Card title) + description (Body, 13 px) bottom-left, `p-4`. Whole card is a link | landing |
| `MovieCard` | poster 2:3 `rounded-lg` with `RatingBadge` top-right; title below (1 line); up to 2 `GenreChip`s (first `gold`, second `neutral`) | recommendations |
| `RatingBadge` | pill, `bg-black/60 backdrop-blur`, gold star icon + rating (1 decimal), 11 px | poster, detail title row |
| `GenreChip` | pill, 10–11 px, `px-2 py-0.5`. `gold`: `bg-gold-soft border-gold-line text-primary`. `neutral`: `bg-muted border-border text-muted-foreground` | cards, detail meta |
| `PrimaryButton` | pill, `bg-primary text-primary-foreground shadow-gold`, h-12, Button text style | "GỢI Ý KHÁC", "THỬ LẠI" |
| `OutlineButton` | pill, `border-gold-line text-primary`, h-9, small | "MUA NGAY" |
| `StickyActionBar` | fixed bottom, full width inside the column, `p-4`, gradient fade above; respects `env(safe-area-inset-bottom)` | recommendations |
| `PageHeader` | sticky, h-14, back chevron (gold) left + centered `BrandWordmark`, `border-b border-border`, `bg-background/80 backdrop-blur` | recommendations |
| `MovieDetailSheet` | shadcn `Drawer`; full height on mobile; hero `MediaCarousel` with a bottom gradient, round close button top-left; sheet body `rounded-t-3xl bg-card` with drag handle | detail |
| `MediaCarousel` | h-65 `snap-x snap-mandatory` scroller, one full-width slide per `media[]` entry, scrollbar hidden, `cursor-grab` with mouse drag-to-scroll; video slide = `youtube-nocookie` iframe (muted, looping autoplay) under a transparent tap-to-control layer; dot buttons at `bottom-5`, clear of the sheet body that overlaps the hero — 6 px dot in a 28 × 32 px target, active one `w-4 bg-primary`, the rest `w-1.5 bg-foreground/40`. Images auto-advance every 5 s; stops on the video slide, on touch, and under `prefers-reduced-motion` (SCREENS §3) | detail hero |
| `InfoCard` | `bg-card rounded-xl border-border p-4`, `SectionLabel` on top | "CÓ MẶT TRÊN" |
| `ProviderLogo` | 32 px logo `rounded-md` + 10 px name below | providers |
| `PastaPairingCard` | like `InfoCard` but `border-gold-line`; gold `SectionLabel`; 56 px dish image `rounded-lg` + dish name (Playfair 16 px) + note (Caption) | detail |
| `CastAvatar` | 48 px circle, 1 px `border-gold-line`, name (Caption) below. Photo from `cast[].profile_url`; `null` → initials on `bg-muted` | detail |
| `PromoCard` | `rounded-xl`, gradient `from-gold-soft to-transparent`, `border-gold-line`; Playfair 18 px gold title + Caption + `OutlineButton` right | detail bottom |
| `SheetActions` | closing note (Playfair italic 14 px, muted, centered) + two equal pill buttons, h-10, `border-border bg-muted/40`, 12 px muted text with a 16 px icon; hover turns the border gold | detail bottom |
| `Skeleton` | `bg-muted rounded-*` + shimmer; mirrors the real layout's sizes | all loading states |
| `ErrorState` | centered: gold line icon (lucide `film` / `clapperboard`), Playfair 20 px title, Body muted text, `PrimaryButton` | errors |
| `Footer` | centered stack behind a 1 px `bg-gold-line/50` divider: `BrandWordmark` (16 px) · contact and social icon row (44 px targets, gold) · "Về Pásta Night" link · TMDB logo + attribution (10 px muted). Empty destinations drop their icon | all pages |
| `SocialIcons` | inline SVG on lucide's 24 px grid: `FacebookIcon`, `InstagramIcon`, `TiktokIcon`. lucide v1 removed brand icons | footer, about sheet |
| `AboutSheet` | shadcn `Drawer`, `rounded-t-3xl`, max height 85 dvh, scrollable; drag handle + close ×; wordmark, serif-italic tagline, story, info rows (gold circle icon + label + value), social circles, gold CTA | `?about=1`, every page |

Icons: `lucide-react`, stroke 1.5, gold for actions, muted for metadata (clock, calendar).

---

## 6. Imagery

| Image | Source | Fallback |
|---|---|---|
| Option card background | `public/images/options/{option_id}.webp` (brand asset) | gold radial gradient + the option's lucide `icon` |
| Cast avatar | `cast[].profile_url` (TMDB CDN, `w185`) | initials on `bg-muted` |
| Movie poster | `poster_url` (TMDB CDN) | `bg-muted` + film icon + title |
| Detail hero | `backdrop_url` → `poster_url` | gradient `from-card to-background` |
| Provider logo | `logo_url` (TMDB CDN) | provider name only |
| Dish image | `public/images/dishes/{dish_id}.webp` (brand asset) | hide the image slot |

- Use `next/image` with `sizes`. `priority` only for above-the-fold images (the first 2 posters, the detail hero).
- The UX Pilot photos are AI-generated mockups: replace them with licensed brand assets before launch.

### Brand assets still missing

The app renders its fallbacks until these land; dropping the files in is enough, no code change (`src/lib/assets.ts` checks whether each file exists).

| Asset | Path | Size | Notes |
|---|---|---|---|
| Option card background | `public/images/options/{option_id}.webp` — `netflix-chill`, `solo`, `friends` | ≥ 896×352, ideally 1344×528 (≈ 2.55:1) | Rendered `object-cover` at up to 448×176 CSS px. Keep the lower half dark: a black gradient and gold title sit on top |
| Dish photo | `public/images/dishes/{dish_id}.webp` | ≥ 168×168 square (56 px at 3x) | Only needed once `pastaPairings` in `src/config/brand.ts` has entries |