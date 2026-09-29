# options

## Purpose
Landing screen list: one card per viewing option from `GET /options`, linking to `/recommendations/{id}`.

## Key Files
- `OptionList.tsx` — renders cards in API order (display order)
- `OptionCard.tsx` — server component; whole card is a `<Link>`; background image or fallback
- `OptionIcon.tsx` — API `icon` name → lucide icon (static switch, no client JS)

## Invariants
- Labels, descriptions, and the option list always come from the API; never hardcode them.
- Card background: `public/images/options/{id}.webp` if the file exists, else gold gradient + icon (DESIGN-SYSTEM §6).
- An empty `data` array is a misconfiguration: the page throws so `error.tsx` shows `ErrorState`.

## Gotchas
- Unknown `icon` names fall back to `Film`; add a case to `OptionIcon.tsx` when options.yaml gains one.
- Dropping a new `.webp` into `public/images/options/` is enough; no code change.

## Related Docs
- `../../../docs/SCREENS.md` §1
- `../../../docs/DESIGN-SYSTEM.md` §5 OptionCard, §6
