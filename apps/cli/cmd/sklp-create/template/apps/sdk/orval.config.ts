import { defineConfig } from "orval";

export default defineConfig({
  core: {
    input: "../core/docs/swagger.yaml",
    output: {
      mode: "split",
      target: "./src/generated/core.ts",
      schemas: "./src/generated/model",
      client: "axios",
      clean: true,
      prettier: false,
      override: {
        mutator: {
          path: "./src/http-client.ts",
          name: "coreHttp",
        },
      },
    },
  },
});
