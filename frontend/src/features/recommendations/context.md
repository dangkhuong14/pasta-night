# recommendations

## Purpose
Grid of 8 movies for one option (2×4 on phones, 4×2 from `md`), with "GỢI Ý KHÁC" to page through the whole list (≤ 40) without repeats.

## Key Files
- `RecommendationGrid.tsx` — client; owns the shuffled order and the current page
- `useShuffle.ts` — paging over a shuffled list; reshuffles when fewer than `DISPLAY_COUNT` remain
- `MovieCard.tsx` — client; opens the detail sheet with `history.pushState(?movie=ID)`
- `CacheNotReady.tsx` — client; waits `Retry-After` seconds, then `router.refresh()`
- `constants.ts` — `DISPLAY_COUNT = 8`

## Invariants
- The first shuffle happens on the server per request (`page.tsx`); the grid initializes its state once from props, so re-renders never reorder it.
- Never show the same movie twice before the list is exhausted. A list shorter than 2 × `DISPLAY_COUNT` (16) cannot keep this across a reshuffle: there are not enough unseen movies for the next page (SCREENS §2).
- The page size is the same on every screen; only the column count is responsive. Keep it that way: a size that depends on the viewport would need JS (the server cannot know the width), and hiding cards with CSS would make "GỢI Ý KHÁC" skip movies nobody saw.
- Fewer than `DISPLAY_COUNT` movies in total → hide the button.

## Gotchas
- Opening the sheet must not trigger a server round trip: use `pushState`, not `router.push`.
- `?movie=` must survive "Gợi ý khác"; paging never touches the URL.
- No 3-column step: 8 cards in 3 columns leave a hole in the last row.
- The mock list has 14 movies, so in mock mode the second page always repeats 2 — expected, not a bug.

## Related Docs
- `../../../docs/SCREENS.md` §2
- `../../../docs/ARCHITECTURE.md` §4, §5
