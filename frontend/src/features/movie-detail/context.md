# movie-detail

## Purpose
Bottom sheet over the recommendations page for `?movie={id}`: summary from the page's list instantly, backdrop and cast after `GET /movies/{id}`.

## Key Files
- `MovieDetailSheet.tsx` — client; `?movie=` ↔ open state, close behavior, section composition
- `useMovieDetail.ts` — fetches the detail on open; in-memory cache per session
- `ProvidersCard.tsx`, `PastaPairingCard.tsx`, `CastList.tsx`, `PromoCard.tsx` — sections

## Invariants
- Never show an empty section: follow the hide rules in SCREENS.md §3.
- Close: `history.back()` if `?movie=` appeared after mount (opened in-app), else `replaceState` (deep link).
- `?movie=` must parse as a positive integer, else it is ignored.
- `MOVIE_NOT_FOUND` → close and `router.refresh()`; other detail errors keep the summary and hide cast.
- A `?movie=` that is not in the list and fails to load closes silently.

## Gotchas
- The page renders the sheet next to the grid; features never import each other.
- The pasta pairing is resolved on the server (`page.tsx`) because it checks `public/` for the dish image.

## Related Docs
- `../../../docs/SCREENS.md` §3
- `../../../docs/DESIGN-SYSTEM.md` §5
