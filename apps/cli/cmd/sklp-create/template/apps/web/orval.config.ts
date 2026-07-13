import { defineConfig } from "orval";

// Generates the typed React Query client from the core OpenAPI spec
// (apps/core/docs/swagger.json, produced by `swag init`). Output lands in the
// shared front lib (@app/front) so any front-end can consume it. Committed;
// regenerate with `pnpm --filter @app/web generate:api` (or `sklp run generate`).
export default defineConfig({
  core: {
    input: {
      target: "../core/docs/swagger.json",
    },
    output: {
      mode: "tags-split",
      target: "../../lib/front/src/api",
      schemas: "../../lib/front/src/api/generated.schemas",
      client: "react-query",
      override: {
        mutator: {
          path: "../../lib/front/src/orval-fetcher.ts",
          name: "coreFetcher",
        },
      },
    },
  },
});
