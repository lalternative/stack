import { defineConfig } from 'vite'
import { devtools } from '@tanstack/devtools-vite'

import { tanstackStart } from '@tanstack/react-start/plugin/vite'

import viteReact from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import { nitro } from 'nitro/vite'

// Origin of the core API, proxied under /api in dev. Overridable so the same
// config works when core runs elsewhere (CI, containers).
const CORE_API_ORIGIN =
  process.env['CORE_API_URL']?.replace(/\/+$/, '') || 'http://localhost:4100'

const config = defineConfig({
  resolve: { tsconfigPaths: true },
  server: { port: 5273 },
  plugins: [
    devtools(),
    nitro({
      rollupConfig: { external: [/^@sentry\//] },
      // Browser-side calls to /api are proxied to the core API through Nitro.
      // Vite's `server.proxy` does NOT work under TanStack Start (Nitro owns
      // the dev request pipeline), so proxying must live in routeRules.
      // Server-side (SSR) fetches bypass this and use an absolute URL — see
      // src/lib/api.ts.
      routeRules: {
        '/api/**': { proxy: `${CORE_API_ORIGIN}/api/**` },
      },
    }),
    tailwindcss(),
    tanstackStart(),
    viteReact(),
  ],
})

export default config
