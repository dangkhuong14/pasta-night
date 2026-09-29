# recommendations

## Purpose
Grid of 6 movies for one option, with "GỢI Ý KHÁC" to page through the whole list (≤ 40) without repeats.

## Key Files
- `RecommendationGrid.tsx` — client; owns the shuffled order and the current page
- `useShuffle.ts` — paging over a shuffled list; reshuffles when fewer than `DISPLAY_COUNT` remain
- `MovieCard.tsx` — client; opens the detail sheet with `history.pushState(?movie=ID)`
- `CacheNotReady.tsx` — client; waits `Retry-After` seconds, then `router.refresh()`
- `constants.ts` — `DISPLAY_COUNT = 6`

## Invariants
- The first shuffle happens on the server per request (`page.tsx`); the grid initializes its state once from props, so re-renders never reorder it.
- Never show the same movie twice before the list is exhausted.
- Fewer than `DISPLAY_COUNT` movies in total → hide the button.

## Gotchas
- Opening the sheet must not trigger a server round trip: use `pushState`, not `router.push`.
- `?movie=` must survive "Gợi ý khác"; paging never touches the URL.

## Related Docs
- `../../../docs/SCREENS.md` §2
- `../../../docs/ARCHITECTURE.md` §4, §5
