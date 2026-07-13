# app

Skalpai-style application generated from the `stack` template. Ships
with the `sklp` dev/CI/publish workflow, OTEL observability via skalpai
sdk-go, a Go DDD backend (`apps/core`), a React + TanStack web app
(`apps/web`).

## Quick start

```bash
cp .env.example .env                # fill SKALPAI_PROJECT_ID + SKALPAI_API_KEY
git config core.hooksPath .githooks
sklp dev stack                      # boots core (4100) + web (5273)
```

## Layout

```
.sklp/                 space.yaml + stack/dev.yaml + pipelines/{ci,build,generate,publish,secops}.yaml
apps/
  core/                Go 1.25 + Echo v4 + Postgres (DDD per bounded context)
  sdk/                 Typed TS client generated from the core OpenAPI (swag → orval → tsup)
  web/                 React 19 + TanStack Router + Vite
infra/nginx/           Static bundle nginx config
Dockerfile*            One per shipped image (core, web)
```

## API SDK

The core exposes an OpenAPI spec (swaggo annotations on the Echo handlers)
and `apps/sdk` is a typed TypeScript client generated from it. The web app
consumes it via the `@app/sdk` workspace package. Never hand-edit the
generated client — annotate the handlers, then regenerate:

```bash
sklp run generate    # swag → apps/core/docs, orval → apps/sdk/src/generated, tsup → dist
```

`apps/core/docs` and `apps/sdk/src/generated` are checked in.

## Workflow

Source of truth: `.sklp/space.yaml` (runner recipe), `.sklp/stack/dev.yaml`
(dev topology) + `.sklp/pipelines/*.yaml` (CI/build/publish).

| Command | What it runs |
|---|---|
| `sklp dev stack` | Local supervisor — core + web |
| `sklp run ci` | Lint + test impacted services |
| `sklp run generate` | Regenerate OpenAPI spec + typed SDK (swag → orval → tsup) |
| `sklp run build` | Build `bin/core` + web bundle |
| `sklp run secops` | Security scan: gitleaks + semgrep + govulncheck + trivy |
| `sklp run publish` | Build, scan (Trivy), push impacted images (main only) |

**No GitHub Actions** drive `sklp` pipelines. Local pre-push gates use
`.githooks/pre-commit`. Remote dispatch is the skalpai UI.

## Observability

`apps/core/observability/` wraps `github.com/lalternative/packages/skalpai/sdk-go`
and configures OTEL traces + metrics + logs at boot. Traffic flows to the URL
in `SKALPAI_ENDPOINT` (using `SKALPAI_API_KEY`). No-op when either env var is
empty.

## Image policy

`rg.fr-par.scw.cloud/skaplai/<service>:YYYYMMDD-HHMM-<short-sha>` — immutable.
No `latest`. Trivy gates HIGH/CRITICAL fixable findings.
