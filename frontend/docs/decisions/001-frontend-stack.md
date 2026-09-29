# 001 — Frontend stack: Next.js 16, Tailwind v4, shadcn/ui on Base UI

## Status

Accepted (2026-09-28)

## Context

The frontend is a mobile-first site opened from a QR code, often on mobile data: 2 routes plus a detail sheet, dark black & gold theme, Vietnamese copy. PROJECT-RULES requires an ADR for every dependency; this one records the choices made when the frontend was first built (ARCHITECTURE.md §1 Tech stack).

## Decision

| Package | Version | Why |
|---|---|---|
| `next` | 16.3.6 | App Router, Server Components, `next/image` (TMDB posters proxied through our server), `next/font` |
| `react`, `react-dom` | 19.2.8 | required by Next 16 |
| `tailwindcss` + `@tailwindcss/postcss` | 4.x | CSS-variable tokens via `@theme` (DESIGN-SYSTEM §2) |
| `shadcn` (CLI + `shadcn/tailwind.css`) | 4.21 | owned-source components; style `base-nova` |
| `@base-ui/react` | 1.8 | primitives behind shadcn `base-nova`: `Drawer`, `Button` |
| `cn` | 0.4 | shadcn's class-merging helper (replaces `clsx` + `tailwind-merge`) |
| `class-variance-authority` | 0.7 | variants in shadcn components |
| `tw-animate-css` | 1.4 | shadcn animation utilities |
| `lucide-react` | 1.48 | icons; option `icon` values from the API are lucide names |
| `prettier` + `prettier-plugin-tailwindcss` (dev) | 3.9 / 0.8 | formatting + class sorting (PROJECT-RULES §2) |
| `eslint` + `eslint-config-next` (dev) | 9 / 16.3.6 | `core-web-vitals` + `typescript` configs |

- No global state or data-fetching library: URL + local state (ARCHITECTURE §5).
- shadcn's default primitive library is now Base UI (`base-nova`). The old Radix `Drawer` depends on `vaul`, which is no longer maintained; the Base UI drawer is maintained and supports swipe-to-dismiss.

## Consequences

- `components/ui/*` follow Base UI APIs (`render` prop, `onOpenChange(open, details)`), not Radix; examples written for Radix-based shadcn need adapting.
- ESLint 9 is flagged as unsupported by npm; move to ESLint 10 once `eslint-config-next` supports it.
- The images in `public/` are optional brand assets; the UI has fallbacks, so the site works before the business delivers them.

## Alternatives considered

- **Radix-based shadcn (`new-york`):** familiar, but its Drawer relies on unmaintained `vaul`.
- **React Query / SWR for the detail fetch:** one client fetch with an in-memory cache doesn't justify a library.
- **CSS Modules or a component library (MUI, Chakra):** Tailwind tokens + owned shadcn source keep the bundle small and the design exact.
