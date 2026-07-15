# stack

Scaffolder for Skalpai-style application stacks. `sklp-create` writes a
ready-to-run project — observability wired in, the `sklp` dev/CI/publish
workflow, a Go DDD backend, and a TanStack web app — from an embedded template.

## Quick start

```bash
# Install the scaffolder once
go install github.com/lalternative/stack/cmd/sklp-create@latest

# Create a new project
sklp-create start myapp
cd myapp
sklp dev stack                      # boots postgres + nats + core (4100) + web (5273)
```

`sklp-create start <name>` writes the embedded template to `./<name>/`,
rewrites the Go module paths (`app/` → `<name>/`) and npm scopes
(`@app/` → `@<name>/`), seeds `.env` from `.env.example`, and `git init`s.

The web app is **not** embedded: `sklp-create start` scaffolds it on the fly
with the official TanStack CLI (`pnpm dlx @tanstack/cli create`), so every new
project gets an up-to-date TanStack Start app rather than a frozen, drifting
copy. `pnpm` must be on PATH.

## This repo

A single Go module — the scaffolder binary and the template it ships. Nothing
else runs here; the backend lives inside the template.

```
cmd/sklp-create/
  main.go              root + version
  start.go             `sklp-create start <name>` (NOT copied into projects)
  template/            checked-in project template (the backend + .sklp lives here)
go.mod                 module github.com/lalternative/stack
```

Edit the generated stack by editing `cmd/sklp-create/template/`. Its `.sklp/`
(dev.yaml + pipelines) and `Dockerfile` are what every new project receives.
