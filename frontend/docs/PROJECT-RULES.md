# Frontend Project Rules

> Conventions every contributor (human or AI) MUST follow in `frontend/`.
> Related: [ARCHITECTURE.md](./ARCHITECTURE.md) · [DESIGN-SYSTEM.md](./DESIGN-SYSTEM.md) · [SCREENS.md](./SCREENS.md) · [API_SPEC.md](../../docs/API_SPEC.md)

## TL;DR — read this first

- Next.js App Router, TypeScript `strict`, Tailwind v4, shadcn/ui. Server Components by default.
- Only `src/lib/api.ts` calls the backend. Components never call `fetch`.
- Style with design tokens only (`bg-card`, `text-primary`, …). No raw hex in components.
- All UI copy in Vietnamese; code, comments, and commits in English.
- Never show empty sections: follow the hide rules in SCREENS.md.
- Make targeted changes. Don't reformat, rename, or restructure unrelated code.
- New dependency, state library, or pattern → ADR in `docs/decisions/` first.

---

## 1. Naming Conventions

### Files & folders

| Kind | Convention | ✅ Do | ❌ Don't |
|---|---|---|---|
| Folders | kebab-case | `movie-detail/` | `movieDetail/` |
| Components | PascalCase `.tsx`, one component per file | `MovieCard.tsx` | `movie-card.tsx` |
| Hooks | camelCase, `use` prefix | `useShuffle.ts` | `shuffle-hook.ts` |
| Other modules | kebab-case `.ts` | `api-types.ts`, `error-copy.ts` | `apiTypes.ts` |
| Next.js special files | as Next requires | `page.tsx`, `loading.tsx`, `error.tsx` | |
| Static assets | kebab-case, keyed by ID | `images/options/netflix-chill.webp` | `NetflixChill.png` |

### Variables, functions, components

| Kind | Convention | Example |
|---|---|---|
| Components | PascalCase | `PastaPairingCard` |
| Props type | `<Component>Props` | `MovieCardProps` |
| Types | PascalCase, no `I` prefix | `MovieSummary` ❌ `IMovieSummary` |
| Functions, variables | camelCase | `formatRuntime`, `visibleMovies` |
| Module constants | UPPER_SNAKE | `DISPLAY_COUNT = 6` |
| Booleans | `is/has/should` | `isSheetOpen`, `hasProviders` |
| Event props / handlers | `onX` / `handleX` | `onShuffle` → `handleShuffle` |

### API types mirror the API

- Types in `src/lib/api-types.ts` use the **exact** API field names (`snake_case`), so they can be diffed against API_SPEC.md.
- No camelCase mapping layer.

```ts
// ✅ mirrors API_SPEC §5.2
export type MovieSummary = {
  id: number;
  title: string;
  overview: string;
  poster_url: string | null;
  release_year: number | null;
  rating: number;
  runtime_minutes: number | null;
  genres: string[];
  providers: Provider[];
};
```

---

## 2. Code Style

### Formatting & linting

- Prettier (defaults) + `prettier-plugin-tailwindcss` (sorts classes). CI rejects unformatted code.
- ESLint: Next.js `core-web-vitals` + `typescript` configs. Zero warnings.
- TypeScript `strict: true`. No `any`, no `@ts-ignore` (use `@ts-expect-error` with a reason, sparingly).

### Import order

```tsx
import { useState } from "react";               // 1. react / next
import Link from "next/link";

import { Star } from "lucide-react";            // 2. third-party

import { getMovie } from "@/lib/api";          // 3. absolute (@/ → src/)
import type { MovieSummary } from "@/lib/api-types";

import { GenreChip } from "./GenreChip";        // 4. relative
```

- Use `import type` for type-only imports.
- Use the `@/` alias for anything outside the current feature folder.

### Comments

- JSDoc on exported helpers in `lib/`. Explain **why**, not what.
- `// TODO(<name>): <action> (<issue link>)`

---

## 3. Mandatory Patterns

### Data fetching

- `src/lib/api.ts` is the only module that calls the backend. It:
  - reads `NEXT_PUBLIC_API_BASE_URL` (unset → returns mock data from `mock-data.ts`);
  - unwraps `{ data, meta }`;
  - throws `ApiError` for non-2xx responses.

```ts
export class ApiError extends Error {
  constructor(
    public status: number,
    public code: string,          // API error code, e.g. "OPTION_NOT_FOUND"
    public details: unknown,
    public retryAfter?: number,   // seconds, from Retry-After
  ) {
    super(code);
  }
}

export function getRecommendations(optionId: string) {
  return request<MovieSummary[], RecommendationsMeta>(
    `/options/${encodeURIComponent(optionId)}/recommendations`,
    { next: { revalidate: 300 } },
  );
}
```

- Server Components fetch for page data. Client components fetch only for interaction-driven data (movie detail on sheet open).

### Error handling

- Branch on `ApiError.code`, never on `message` (it's English, for developers).
- User-facing copy per code lives in `src/lib/error-copy.ts`:

```ts
export const ERROR_COPY: Record<string, { title: string; body: string }> = {
  CACHE_NOT_READY: { title: "Đang chuẩn bị gợi ý…", body: "Vui lòng đợi trong giây lát." },
  DEFAULT:         { title: "Rất tiếc, đã có lỗi xảy ra", body: "Vui lòng thử lại sau ít phút." },
};
```

- Page-level errors → the route's `error.tsx` renders `ErrorState`. Per-code behavior (redirect, auto-retry, silent close) → SCREENS.md.

### Null handling & validation

- Every nullable field from the API has a rendered fallback or a hide rule (DESIGN-SYSTEM §6, SCREENS.md hide rules).
- Route param `option` is not validated client-side; the API answers `OPTION_NOT_FOUND` → `redirect("/")`.
- `?movie=` must parse as a positive integer, else ignore it.

### Server vs. Client Components

- Default: Server Component.
- `"use client"` only for interactivity: shuffle, sheet, auto-retry. Keep the client boundary as low in the tree as possible.
- Pass plain serializable props (API types) from server to client.

### Styling

- Tokens only (DESIGN-SYSTEM §2). Conditional classes via `cn()` from `@/lib/utils`.
- Variants via props, not duplicated components: `<GenreChip variant="gold" />`.
- No inline `style` except truly dynamic values (e.g. animation delay).

### Accessibility

- `<html lang="vi">`.
- Every image has `alt` (poster: movie title; decorative background: `alt=""`).
- Tap targets ≥ 44 px. Visible focus ring (`ring-ring`).
- Icon-only buttons have `aria-label` ("Quay lại", "Đóng").
- Motion respects `prefers-reduced-motion`.

### Logging

- No `console.log` in committed code. `console.error` only inside error boundaries.

---

## 4. Don'ts

### Anti-patterns

- ❌ `fetch` outside `src/lib/api.ts`.
- ❌ Hardcoding option labels, descriptions, or the option list — they come from `GET /options`.
- ❌ Raw hex, arbitrary color values, or fonts outside the tokens.
- ❌ Showing the API `message` or `err.message` to users.
- ❌ Rendering a section with empty data (empty providers card, blank synopsis).
- ❌ Features importing other features. Compose them in the route `page.tsx`; share through `components/common` or `lib`.
- ❌ `components/common` importing a feature. Take the feature-owned piece as a prop and let the route or layout pass it in (`Footer` does this with `aboutLink`).
- ❌ Global state libraries (Redux, Zustand) or data libraries (React Query, SWR) without an ADR. URL + local state are enough (ARCHITECTURE §5).
- ❌ Heavy edits to `components/ui/*` (shadcn). Wrap or compose instead.

### Deprecated approaches

| Don't | Use instead |
|---|---|
| Pages Router (`pages/`, `getServerSideProps`) | App Router, Server Components |
| `<img>` | `next/image` |
| `tailwind.config.js` for tokens | CSS variables + `@theme` in `globals.css` |
| Google Fonts `<link>` | `next/font` with the `vietnamese` subset |

### Security & compliance

- Only `NEXT_PUBLIC_*` vars reach the browser; never put secrets there. They are inlined **at build time** → changing one requires a rebuild.
- No `dangerouslySetInnerHTML`.
- External links (`NEXT_PUBLIC_SHOP_URL`) use `rel="noopener noreferrer"`.
- Never remove the TMDB or JustWatch attribution (licensing).
- Images come from the TMDB CDN (`images.remotePatterns`) or `public/`. Don't hotlink posters from other sites.