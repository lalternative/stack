# stack

Skalpai-style application stack. Generated from this template, each new
app starts with the full skalpai dev workflow + observability already wired.

This repo is **backend + scaffolder**, not a full app: it holds the Go core
and the `sklp-create` CLI. It has no frontend or TS SDK of its own — a
generated project gets its web app scaffolded on the fly (see below).

## Stack
- **Backend**: Go 1.25 + Echo v4 + DuckDB (embedded) — port 4100
- **CLI**: cobra-based scaffolder + project CLI (`apps/cli`)
- **Pipeline**: host-installed `sklp` CLI (NOT vendored), config in `.sklp/`
- **Frontend (generated projects only)**: `sklp-create start` scaffolds
  `apps/web` on the fly with the official TanStack CLI
  (`pnpm dlx @tanstack/cli create` — React 19 + TanStack Start + Nitro SSR
  + tanstack-query, port 5273). Nothing frontend is embedded or checked in
  here, so it never drifts.

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
- The core OpenAPI spec is generated from swaggo annotations: annotate Echo handlers with `// @...`, then `sklp run generate` (swag → `apps/core/docs`). `apps/core/docs` is checked in; the CI `swagger-up-to-date` step fails on drift.
- When editing the scaffolder, mirror any template change under `apps/cli/cmd/sklp-create/template/` (drift is caught by PR diff, not auto-synced).

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
