# Architecture

This stack mirrors the skalpai layout so a new project lands on the same
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

## Event-driven seam (durable consumers)

Integration events are consumed through `github.com/lalternative/packages/go/eda/pkg/consumer`,
**never** by hand-rolling a JetStream subscription. The brick owns every
redelivery concern once, so handlers can't get them wrong:

| Concern | Provided by |
| --- | --- |
| Permanent error → `Term()` | `consumer.Permanent(err)` / `ErrPermanent` |
| Bounded `MaxDeliver` | `EventHandler.MaxDeliver()` |
| Staged `BackOff` | `Config.BackOff` (default 30/60/120s) |
| Dead-letter stream (DLQ) | `MAX_DELIVERIES` advisory stream |
| Heartbeat anti-`AckWait` | `InProgress()` ticker (concurrent handlers) |
| Idempotency | optional `Config.Idempotency` (by `event_id`) |
| Reconnect / retry loop | `consumer.Run` |

A handler lives under `apps/core/<context>/events/`, implements
`consumer.EventHandler` (`Name/Subject/DurableName/MaxDeliver/Handle`) and
writes business logic **only** in `Handle`. `main.go` starts one
`consumer.Run` goroutine per handler. `apps/core/project/events/` ships a
working example — copy it, change the subject and the body of `Handle`.

> Anti-pattern this prevents: a service that re-implements a minimal JetStream
> consumer (Nak everywhere, no `Term`, no `MaxDeliver`/`BackOff`/DLQ) and
> silently loses or infinitely retries messages. If you find yourself calling
> `js.Subscribe`/`Consume` directly, stop — use the brick.

## Observability seam

`apps/core/observability` wraps `github.com/lalternative/packages/skalpai/sdk-go`
and configures OTEL traces, metrics, and logs exporters at boot. Disabled
(no-op) when `SKALPAI_ENDPOINT` or `SKALPAI_API_KEY` are empty, so
`go test ./...` and offline dev don't fight the SDK. The web app should
ship `@digstack/sdk-browser` for RUM (left for the consuming project to
add — it's not in the stack to avoid a forced dep).

## sklp seam

The stack **uses** the host `sklp` CLI. `.sklp/space.yaml` is the runner
recipe (toolchain + registry); `.sklp/stack/dev.yaml` is the declarative
dev topology; `.sklp/pipelines/*.yaml` are the CI/build/publish pipelines.
New behavior goes into helpers/templates/codegen — not into new top-level
YAML fields (skalpai convention).

## Renaming

`sklp start <name>` renames the template in-process while scaffolding:
- replace `app/` Go module paths with `<name>/`
- replace `@app/` npm scopes with `@<name>/`
- replace `rg.fr-par.scw.cloud/skaplai/` images with your registry/team

## Container registry policy

`rg.fr-par.scw.cloud/skaplai/<service>:YYYYMMDD-HHMM-<short-sha>`. Set by
`tag.strategy: date-sha` in `.sklp/pipelines/publish.yaml`. Never `latest`,
never overwritten, Trivy gates HIGH/CRITICAL fixable findings.
