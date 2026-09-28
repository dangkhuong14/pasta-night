# cache

## Purpose
The in-memory snapshot the read path serves from, persisted as JSON files under `CACHE_DIR`. The only package that writes cache files.

## Key Files
- `store.go` — `Snapshot`, `Store`: `Snapshot()`, `IsReady()`, `PublishList`, `GC`
- `persist.go` — file format (`schema_version`), `Load`, read/write helpers

## Invariants
- A Snapshot is immutable; every change swaps in a new one (`atomic.Pointer`). Readers never lock.
- `Snapshot.Details` holds exactly the movies that `Snapshot.Lists` reference.
- `PublishList` writes detail files, then the list file, then swaps the snapshot.
- Writers (`PublishList`, `GC`) hold `Store.mu`, so GC never runs between a publish's file writes.
- `Load` treats a list as missing if any detail it references is missing or invalid.

## Gotchas
- A detail reused from memory may have been deleted from disk by GC. `PublishList` rewrites every detail that is not in the current snapshot with the same `fetched_at`.
- Bump `SchemaVersion` on any persisted struct change (DATABASE.md §6); old files then load as missing.
- Writes recreate missing directories, so wiping `CACHE_DIR` at runtime heals on the next publish.

## Related Docs
- `../../../docs/DATABASE.md`
