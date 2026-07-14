import { configureCoreClient } from "@app/sdk";

// The scaffold runs SSR (Nitro) + client. Generated routes are relative to the
// core basePath (/api/v1), so we only pick the origin here:
//
//  - Server (SSR/Nitro): relative URLs don't resolve — a fetch needs an
//    absolute origin. Use CORE_API_URL (e.g. http://localhost:4100), falling
//    back to the dev core port so SSR loaders work out of the box.
//  - Browser: same-origin; Vite proxies /api -> core in dev, and in prod the
//    web host proxies it. VITE_CORE_API_URL can override the origin.
function resolveOrigin(): string {
  if (import.meta.env.SSR) {
    return (
      process.env["CORE_API_URL"]?.trim().replace(/\/+$/, "") ||
      "http://localhost:4100"
    );
  }
  return (
    import.meta.env["VITE_CORE_API_URL"]?.trim().replace(/\/+$/, "") || ""
  );
}

configureCoreClient({ baseURL: `${resolveOrigin()}/api/v1` });
