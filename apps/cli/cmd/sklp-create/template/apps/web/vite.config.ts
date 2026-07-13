import { defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import viteReact from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import { nitro } from "nitro/vite";

// Server-rendered (SSR) via TanStack Start + Nitro — the default.
// `sklp create --render spa` rewrites this app to a client-only Vite SPA
// at generation time (no Start server), so this file stays pure SSR.
export default defineConfig({
  server: {
    port: 5273,
    proxy: {
      // Backend API (Go core) for non-auth routes. better-auth handles
      // /api/auth/* locally via the Start server route, so scope the dev
      // proxy to /api/v1 only — never the broad /api.
      "/api/v1": {
        target: "http://localhost:4100",
        changeOrigin: true,
      },
    },
  },
  plugins: [
    nitro(),
    tsconfigPaths({ projects: ["./tsconfig.json"] }),
    tailwindcss(),
    tanstackStart(),
    viteReact(),
  ],
});
