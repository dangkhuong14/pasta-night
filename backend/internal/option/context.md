# option

## Purpose
Serves `GET /api/v1/options`: the viewing options from `configs/options.yaml`, in display order.

## Key Files
- `service.go` — holds the options loaded at startup
- `handler.go` — route and response shape (`id`, `label`, `description`, `icon`)

## Invariants
- Options are config: loaded and validated once at startup (`platform/config`), never changed at runtime.
- Discover params stay internal; the API exposes only `id`, `label`, `description`, `icon`.

## Gotchas
- Option IDs are frontend routes: renaming one is an API breaking change.

## Related Docs
- `../../../docs/API_SPEC.md` §5.1
- `../../../docs/DATABASE.md` §2 ViewingOption
