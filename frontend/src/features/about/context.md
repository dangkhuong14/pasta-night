# about

## Purpose
"Về Pásta Night": the brand story and shop details, as a sheet over whatever page the customer is on. Opened from the footer link on every page.

## Key Files
- `AboutSheet.tsx` — client; `?about=1` ↔ open state, all sections
- `AboutLink.tsx` — client; the footer link, opens the sheet with `history.pushState`
- `constants.ts` — `ABOUT_PARAM`

## Invariants
- Content comes from `config/brand.ts`, never from the API.
- Never show an empty row: address, map button, hours, hotline, each social icon and the shop button each hide when their value is empty.
- Close: `history.back()` if the sheet was opened in-app, else `replaceState` (someone shared the `?about=1` link).
- Rendered once in `app/layout.tsx`, so it works from every page.

## Gotchas
- It reads `useSearchParams`, so the layout wraps it in `<Suspense>`; without that, statically prerendered routes such as `/_not-found` fail to build.
- The footer lives in `components/common` and may not import features, so the layout passes `<AboutLink />` in as a prop.

## Related Docs
- `../../../docs/SCREENS.md` §6
- `../../../docs/DESIGN-SYSTEM.md` §5
