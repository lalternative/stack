# stack

Skalpai-style application stack. A ready-to-run starter that gives a new
project everything skalpai has on day one — observability wired in, the `sklp`
dev/CI/publish workflow, a Go DDD backend, a TanStack web app, and a project
CLI.

This repo doubles as **scaffolder**: install the CLI once, then use it to
spin up new projects from the embedded template.

## Quick start

```bash
# Install the scaffolder once
go install github.com/lalternative/stack/apps/cli/cmd/sklp-create@latest

# Create a new project
sklp-create start myapp
cd myapp
git config core.hooksPath .githooks
sklp dev                            # boots core (4100) + web (5273)
```

`sklp-create start <name>` writes the embedded template to `./<name>/`,
rewrites the Go module paths (`app/` → `<name>/`) and npm scopes
(`@app/` → `@<name>/`), seeds `.env` from `.env.example`, and `git init`s.

## Layout (generated project)

```
.sklp/                 dev.yaml + tasks/{ci,build,publish,secops}.yaml
apps/
  core/                Go 1.25 + Echo v4 + DuckDB (DDD per bounded context)
  web/                 React 19 + TanStack Router + Vite
  cli/                 Project CLI binary (cobra)
infra/nginx/           Static bundle nginx config
Dockerfile*            One per shipped image (core, web)
```

## Workflow

Source of truth: `.sklp/dev.yaml` + `.sklp/tasks/*.yaml`.

| Command | What it runs |
|---|---|
| `sklp dev` | Local supervisor — core + web |
| `sklp run ci` | Lint + test impacted services |
| `sklp run build` | Build `bin/app` + web bundle |
| `sklp run secops` | Security scan: gitleaks + semgrep + govulncheck + trivy |
| `sklp run publish` | Build, scan (Trivy), push impacted images (main only) |

**No GitHub Actions** drive `sklp` pipelines. Local pre-push gates use
`.githooks/pre-commit`. Remote dispatch is the skalpai UI.

## Observability

`apps/core/observability/` wraps `github.com/digstack/skalpai/packages/sdk-go`
and configures OTEL traces + metrics + logs at boot. Traffic flows to the URL
in `SKALPAI_ENDPOINT` (using `SKALPAI_API_KEY`). The wrapper is a no-op when
either env var is empty, so `go test ./...` and offline dev work without
telemetry.

## Image policy

`rg.fr-par.scw.cloud/skaplai/<service>:YYYYMMDD-HHMM-<short-sha>` — immutable.
No `latest`. Trivy gates HIGH/CRITICAL fixable findings.

## Repo layout (this template)

This repo holds the source AND the embedded template:

```
apps/
  cli/cmd/sklp-create/   scaffolder CLI (this binary)
    main.go              root + version
    start.go             `sklp-create start <name>` (NOT in generated projects)
    template/            checked-in copy of the project template
  core/, web/, ...       the actual source code
```

Touching anything under `apps/core/`, `apps/web/`, or `.sklp/` requires the
same change in `apps/cli/cmd/sklp-create/template/` for the scaffolder to
ship the update. (Drift is gated by visual diff on PRs — there is no
auto-sync.)
