# CLAUDE.md

Guide for AI coding agents in this repo. This file stays short; details live in the docs it points to.

## Project

Web app for Pásta Night, a pasta startup. Customers scan the QR code on their order card, pick a viewing mood (`netflix-chill`, `solo`, `friends`), and get movie recommendations.

- `frontend/` — Next.js (App Router, TypeScript, Tailwind, shadcn/ui). Mobile-first, black & gold theme, UI copy in Vietnamese.
- `backend/` — Go API. A background job precomputes recommendations from TMDB; requests are served from memory. No database.

## Where to look

| Read | When |
|---|---|
| `backend/CLAUDE.md` | before any backend work |
| `frontend/CLAUDE.md` | before any frontend work |
| `frontend/docs/PROJECT-RULES.md` | before writing frontend code |
| `frontend/docs/SCREENS.md` | building or changing a screen |
| `frontend/docs/DESIGN-SYSTEM.md` | any styling or new UI component |
| `backend/docs/PROJECT-RULES.md` | before writing backend code |
| `backend/docs/ARCHITECTURE.md` | adding packages or env vars, changing flows |
| `backend/docs/TMDB_INTEGRATION.md` | touching `internal/platform/tmdb` or `internal/refresh` |
| `docs/API_SPEC.md` | any endpoint, request/response field, or error code (FE or BE) |
| `docs/DATABASE.md` | cached data, cache files, `configs/options.yaml` |
| `backend/internal/<feature>/context.md` | before editing that feature |
| `docs/decisions/`, `backend/docs/decisions/` | before proposing a dependency, storage, or pattern |
| `*/docs/PROJECT-STATUS.md` | start of a session, if it exists |

- Read only what the task needs.

**When sources disagree**

- Contract: `API_SPEC.md` wins over code. Data shape: `DATABASE.md`. TMDB mapping: `TMDB_INTEGRATION.md`.
- Screens: `frontend/docs/SCREENS.md` wins over `frontend/design/*.png` and code.
- Doc and code mismatch → stop and report it. Don't silently change either one.

## Workflow

1. Scope the task: FE, BE, or shared? Which features?
2. Read the docs listed above for that scope.
3. Non-trivial change (contract change, new package or dependency, more than 3 files) → state the plan and the files first, then wait for approval.
4. Make targeted changes. No drive-by refactors, renames, or reformatting of unrelated code.
5. Update docs in the same change (see Doc Sync).
6. Update <side>/docs/PROJECT-STATUS.md if the task changed what's done, in progress, or broken. End with a short summary: what changed, which docs were updated, open questions.

**Definition of done**

- Backend: `goimports` clean; `golangci-lint run` and `go test ./...` pass.
- Frontend: lint and build pass.
- Docs describe the code as it now is.

## Doc Sync — change X, update Y

| Change | Update |
|---|---|
| endpoint, field, status or error code | `docs/API_SPEC.md` (+ Changelog), `frontend/src/lib/api-types.ts`, `api.ts`, `mock-data.ts` |
| persisted struct / cache file format | `docs/DATABASE.md` + bump `cache.SchemaVersion` |
| `configs/options.yaml` schema | `docs/DATABASE.md` §2 ViewingOption |
| TMDB params, DTOs, mapping, error handling | `backend/docs/TMDB_INTEGRATION.md` |
| package, folder, env var, flow | `backend/docs/ARCHITECTURE.md` (+ `backend/.env.example`) |
| new or changed convention | the relevant `PROJECT-RULES.md` |
| dependency, storage, or pattern decision | new ADR |
| feature behavior or invariants | `internal/<feature>/context.md` (BE) · `src/features/<feature>/context.md` (FE) |
| screen layout, data mapping, states, hide rules | `frontend/docs/SCREENS.md` |
| design token or shared UI component | `frontend/docs/DESIGN-SYSTEM.md` |
| task finished, blocked, or a known issue found | <side>/docs/PROJECT-STATUS.md |

## General Conventions

- **Language:** code, comments, docs, commits, logs → English. User-facing UI copy → Vietnamese.
- **Commits:** Conventional Commits; scope = area or package.
  - `feat(movie): add stale flag to recommendations meta`
  - `fix(refresh): drop duplicate IDs across discover pages`
  - `docs(api): add MOVIE_NOT_FOUND error code`
- **Branches:** `<type>/<short-kebab>` → `feat/admin-refresh`, `fix/tmdb-429-retry`.
- **ADRs:**
  - System-wide: `docs/decisions/SYS-NNN-<kebab-title>.md`. One side only: `<side>/docs/decisions/NNN-<kebab-title>.md`.
  - Sections: Status (Proposed · Accepted · Superseded by …) · Context · Decision · Consequences · Alternatives considered.
  - Never edit an Accepted ADR; supersede it with a new one.
- **Secrets:** env only. Never print, log, or commit `.env`, `TMDB_READ_TOKEN`, or `ADMIN_TOKEN`.

## Don't

- Don't commit, push, or open PRs unless asked.
- Don't call the real TMDB API from tests or scripts. Use fixtures in `backend/internal/platform/tmdb/testdata/`. Run manual `curl` checks only when asked.
- Don't change the API contract, add a dependency, or add storage (DB, Redis) without the Doc Sync / ADR steps.
- Don't edit runtime or generated files: `backend/data/cache/`, `frontend/.next/`, lockfiles (unless the task is dependency work).
- Don't edit `frontend/design/`: UX Pilot exports are reference only.
- Don't guess on decisions the docs don't cover when they affect the contract or data. Ask.