# movie-detail

## Purpose
Bottom sheet over the recommendations page for `?movie={id}`: summary from the page's list instantly, hero carousel and cast after `GET /movies/{id}`.

## Key Files
- `MovieDetailSheet.tsx` — client; `?movie=` ↔ open state, close behavior, section composition
- `useMovieDetail.ts` — fetches the detail on open; in-memory cache per session
- `ProvidersCard.tsx`, `PastaPairingCard.tsx`, `CastList.tsx`, `PromoCard.tsx` — sections
- `MediaCarousel.tsx` — client; hero carousel of `media[]`: YouTube embed + stills, auto-advance, position dots
- `SheetActions.tsx` — client; closing note, share (Web Share API, clipboard fallback), close
- `CastList.tsx` — TMDB photo per actor, initials when `profile_url` is null

## Invariants
- Never show an empty section: follow the hide rules in SCREENS.md §3.
- Close: `history.back()` if `?movie=` appeared after mount (opened in-app), else `replaceState` (deep link).
- `?movie=` must parse as a positive integer, else it is ignored.
- `MOVIE_NOT_FOUND` → close and `router.refresh()`; other detail errors keep the summary and hide cast.
- A `?movie=` that is not in the list and fails to load closes silently.
- The carousel renders `media[]` in the order the API sends it and never reorders: the trailer is index 0 by contract (API_SPEC §5.3).
- Autoplay is always muted, and both autoplay and auto-advance are off under `prefers-reduced-motion`.
- The trailer loops so YouTube's end screen never renders; do not "fix" this back to `rel=0`, which has not suppressed suggestions since 2018.
- The dots are the one navigation that always works. Keep them buttons: over the video, a swipe may belong to the player.
- Three ways in, all landing on the same scroll position: touch (native), mouse drag (ours), dots. Do not add a fourth source of truth for the current slide.

## Gotchas
- The page renders the sheet next to the grid; features never import each other.
- From `lg` the body is a two-column grid with `items-start`, which is what lets the carousel column be `sticky`; stretch the columns and the sticky element has nowhere to move. The hero's `-mt-6` overlap, drag handle and bottom gradient are switched off there.
- The sheet covers the page, so it takes a `footer` prop and the route passes `<Footer aboutLink={<AboutLink />} />` in — the page's own footer is behind the sheet.
- The pasta pairing is resolved on the server (`page.tsx`) because it checks `public/` for the dish image.
- The carousel's only state is the scroll position; `current` mirrors it for the dots. Auto-advance calls `scrollTo` on the same element a swipe moves, so the two can never disagree.
- A `scrollTo` of ours fires scroll events too. `selfScrollUntilRef` marks that window, or every auto-advance would look like the customer taking over and pause the next one.
- `prefers-reduced-motion` is read with `useSyncExternalStore`, not an effect: it decides the iframe's `autoplay` during render, and `react-hooks/set-state-in-effect` rejects the effect form.
- Pointer capture for the mouse drag is taken only after the 4 px threshold. Taking it on `pointerdown` retargets the click to the scroller and the video slide silently loses its tap-to-control.
- `draggable={false}` on the stills: otherwise a mouse drag starts the browser's image drag instead of scrolling.
- A cross-origin iframe swallows touches, so the video slide carries a transparent button over it for swipes. Headless Chrome cannot reproduce this — the player never loads there, and the swipe chains through to the scroller — so test the trailer slide on a real device.

## Related Docs
- `../../../docs/SCREENS.md` §3
- `../../../docs/DESIGN-SYSTEM.md` §5
