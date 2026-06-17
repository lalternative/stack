# Architecture

This bootstrap mirrors the skalpai layout so a new project lands on the same
mental model from day one.

## Layered backend (DDD)

```
apps/core/<context>/
├── api.go              ── HTTP handlers (Echo). Thin: bind → call handler → DTO
├── dto.go              ── transport types + mapping from domain
├── bootstrap.go        ── wires handlers + repo, exposes RegisterRoutes
├── domain/             ── entities, value objects, repository interfaces
├── application/<usecase>/
│                       ── one folder per command/query. Holds Command|Query +
│                          Handler. No HTTP, no SQL.
└── infrastructure/     ── concrete repositories (duckdb_repository.go, ...)
```

Rationale: same shape as `apps/core/project/` in skalpai. New contexts copy
the folder, no debate.

## Data access seam

`apps/core/pkg/db.Executor` is the only thing repositories see. Backed by
DuckDB out of the box; swap by implementing the three methods on a new
driver. Migrations are file-based SQL under `migrations/duckdb/`, applied in
lexical order at boot.

## Observability seam

`apps/core/observability` wraps `github.com/digstack/skalpai/packages/sdk-go`
and configures OTEL traces, metrics, and logs exporters at boot. Disabled
(no-op) when `SKALPAI_ENDPOINT` or `SKALPAI_API_KEY` are empty, so
`go test ./...` and offline dev don't fight the SDK. The web app should
ship `@digstack/sdk-browser` for RUM (left for the consuming project to
add — it's not in the bootstrap to avoid a forced dep).

## sklp seam

The bootstrap **uses** the host `sklp` CLI. `.sklp/dev.yaml` is the
declarative dev topology; `.sklp/tasks/*.yaml` are the CI/build/publish
pipelines. New behavior goes into helpers/templates/codegen — not into
new top-level YAML fields (skalpai convention).

## Renaming

Run `./scripts/rename.sh <newname>` after cloning to:
- replace `app/` Go module paths with `<newname>/`
- replace `@app/` npm scopes with `@<newname>/`
- replace `rg.fr-par.scw.cloud/skaplai/` images with your registry/team

## Container registry policy

`rg.fr-par.scw.cloud/skaplai/<service>:YYYYMMDD-HHMM-<short-sha>`. Set by
`tag.strategy: date-sha` in `.sklp/tasks/publish.yaml`. Never `latest`,
never overwritten, Trivy gates HIGH/CRITICAL fixable findings.
