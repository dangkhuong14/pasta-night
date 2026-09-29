# Frontend — CLAUDE.md

Next.js website for Pásta Night. Before writing code, read `docs/PROJECT-RULES.md`; its TL;DR is mandatory.

## Quick Facts

- Next.js App Router, TypeScript strict, Tailwind v4, shadcn/ui. Server Components by default.
- Only `src/lib/api.ts` calls the backend. Contract: `../docs/API_SPEC.md`.
- Dark-only, mobile-first, black & gold. Tokens only, never raw hex.
- UI copy in Vietnamese. Fonts must load the `vietnamese` subset.

## Where to Look

| Read | When |
|---|---|
| `docs/SCREENS.md` | building or changing any screen: layout, data mapping, states, hide rules |
| `docs/DESIGN-SYSTEM.md` | any styling or new component |
| `design/*.png` | visual reference for the screen you're building |
| `docs/ARCHITECTURE.md` | new route, feature folder, env var, or data flow |
| `src/features/<feature>/context.md` | before editing that feature |

- Design images are a **reference**, not a spec. When an image and SCREENS.md disagree, SCREENS.md wins.
- SCREENS.md §0 lists design elements the API doesn't provide. Don't invent API fields for them.

## Commands

```bash
cp .env.example .env.local     # first time; leave NEXT_PUBLIC_API_BASE_URL empty to use mock data
npm run dev                    # http://localhost:3000
npm run lint
npx tsc --noEmit               # type check
npm run build
npx shadcn@latest add drawer   # add a shadcn component
```

## Build Order (first implementation)

1. Tokens + fonts (`globals.css`, `layout.tsx`) and `components/common`.
2. Landing with mock data.
3. Recommendations grid + skeleton + shuffle.
4. Detail sheet.
5. Error states, then switch to the real API.

## Feature Context Files

- Read `src/features/<feature>/context.md` before editing a feature. Create it the first time you touch a feature that has none. Keep it under ~40 lines.

```md
# <feature>

## Purpose
One or two sentences.

## Key Files
- `RecommendationGrid.tsx` — client; owns shuffle state

## Invariants
- Never show the same movie twice before the list is exhausted.

## Gotchas
- `?movie=` must survive shuffling; don't reset it on "Gợi ý khác".

## Related Docs
- `../../../docs/SCREENS.md` §2
```