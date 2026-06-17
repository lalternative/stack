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
    proxy: { "/api": "http://localhost:4100" },
  },
  plugins: [
    nitro(),
    tsconfigPaths({ projects: ["./tsconfig.json"] }),
    tailwindcss(),
    tanstackStart(),
    viteReact(),
  ],
});
