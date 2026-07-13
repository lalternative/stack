# stack

Skalpai-style application stack. Generated from this template, each new
app starts with the full skalpai dev workflow + observability already wired.

## Stack
- **Backend**: Go 1.25 + Echo v4 + DuckDB (embedded) — port 4100
- **Frontend**: React 19 + TanStack Router + Tailwind v4 + shadcn/ui (`new-york`, zinc) — port 5273
- **SDK**: `apps/sdk` — typed TS client generated from the core OpenAPI (swag → orval → tsup)
- **CLI**: cobra-based project CLI (`apps/cli`)
- **Pipeline**: host-installed `sklp` CLI (NOT vendored), config in `.sklp/`

## Architecture
- DDD per bounded context under `apps/core/<context>/`
- Pattern: `domain/ → application/<usecase>/ → infrastructure/ → api.go + dto.go + bootstrap.go`
- Database abstraction via `pkg/db.Executor` interface
- File-based SQL migrations in `apps/core/migrations/duckdb/`
- Observability wired at boot through `apps/core/observability` using `@digstack/sdk-go`

## Conventions
- All code, comments, commit messages in English
- Commits must include `Co-Authored-By: codesyl <codesyl@pm.me>`
- Do NOT include `Co-Authored-By: Claude` or any Anthropic co-author
- Use `github.com/google/uuid` for IDs
- JWT middleware extracts user via `middleware.GetUser(c)`
- Frontend: import UI from the `@/` alias (`@/components`, `@/lib/utils`); compose classes with `cn()`; add shadcn components via `pnpm dlx shadcn@latest add <name>`
- The TS API client is generated, never hand-written: annotate Echo handlers with swaggo `// @...`, then `sklp run generate` (swag → orval → tsup). `apps/sdk/src/generated` and `apps/core/docs` are checked in; the CI `sdk-up-to-date` step fails on drift.

## Feature / PR / Publish flow

Use the `sklp flow` helpers — they wrap this flow and, by creating a
dedicated worktree per feature, let several features be worked on in
parallel without branch-switching in a shared checkout. Prefer them over
raw `git checkout -b`.

1. `sklp flow start <name>` — opens `feat/<name>` in a dedicated worktree
   (default; `--branch` checks out in place instead). One worktree per
   feature keeps parallel work isolated.
2. `sklp run ci` before pushing.
3. `sklp flow end --pr` — diffs vs `--onto`, runs CI on the impacted scope,
   and opens the PR targeting `main`. CI runs again on PR.
4. After merge, update local `main`.
5. `sklp run publish` from clean local `main`.
6. Publish only impacted services. Tags `YYYYMMDD-HHMM-<short-sha>`.
7. ArgoCD / Image Updater handles rollout.

Rules:
- Never publish from a feature branch.
- Never claim deployment success when only the image was pushed.
- Treat `.sklp/tasks/ci.yaml` and `.sklp/tasks/publish.yaml` as the source of truth.

## Development

```bash
cp .env.example .env
sklp dev                   # boots core → web
sklp dev --validate        # parse-only sanity check
```

Ports: core 4100, web 5273.

Other entry points:
- `sklp run build` — build `bin/app` and the web bundle
- `sklp run ci` — lint + test impacted services
- `sklp run secops` — gitleaks, semgrep, govulncheck, pnpm audit, trivy
- `sklp run publish` — build, scan, push impacted Docker images (main only)

## CI / triggers — NO GitHub Actions

Do NOT add `.github/workflows/*.yml` that drive `sklp run`. The skalpai UI is
the CI control plane (remote dispatch). Local pre-push gates use git hooks:

```bash
git config core.hooksPath .githooks
```

## Container registry

`rg.fr-par.scw.cloud/skaplai/<service>:<date>-<short-sha>`. Immutable tags
only. No `latest`. Trivy gates HIGH/CRITICAL fixable findings.
