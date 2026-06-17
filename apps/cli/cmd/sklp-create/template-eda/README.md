# app

Skalpai-style application generated from the `bootstrap` template. Ships
with the `sklp` dev/CI/publish workflow, OTEL observability via skalpai
sdk-go, a Go DDD backend (`apps/core`), a React + TanStack web app
(`apps/web`).

## Quick start

```bash
cp .env.example .env                # fill SKALPAI_PROJECT_ID + SKALPAI_API_KEY
git config core.hooksPath .githooks
sklp dev                            # boots core (4100) + web (5273)
```

## Layout

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
| `sklp run build` | Build `bin/core` + web bundle |
| `sklp run secops` | Security scan: gitleaks + semgrep + govulncheck + trivy |
| `sklp run publish` | Build, scan (Trivy), push impacted images (main only) |

**No GitHub Actions** drive `sklp` pipelines. Local pre-push gates use
`.githooks/pre-commit`. Remote dispatch is the skalpai UI.

## Observability

`apps/core/observability/` wraps `github.com/digstack/skalpai/packages/sdk-go`
and configures OTEL traces + metrics + logs at boot. Traffic flows to the URL
in `SKALPAI_ENDPOINT` (using `SKALPAI_API_KEY`). No-op when either env var is
empty.

## Image policy

`rg.fr-par.scw.cloud/skaplai/<service>:YYYYMMDD-HHMM-<short-sha>` — immutable.
No `latest`. Trivy gates HIGH/CRITICAL fixable findings.
