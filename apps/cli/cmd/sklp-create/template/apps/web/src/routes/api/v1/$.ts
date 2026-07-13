import { createFileRoute } from "@tanstack/react-router";

// Same-origin reverse proxy: forwards /api/v1/* to the Go core so the browser
// never talks to the core cross-origin. Keeping it same-origin means the
// `token` cookie stays SameSite=Lax and the core needs no CORS entry.
//
// The incoming path already starts with /api/v1, which the core mounts
// (apps/core/main.go), so we forward the path verbatim. Headers are passed
// through unchanged (Cookie, Authorization, Content-Type, ...) — so a caller
// may authenticate with the session cookie OR an `Authorization: Bearer`.
// CORE_API_URL points at the core host (e.g. http://core:4100 in-cluster).
const CORE_URL = (process.env.CORE_API_URL ?? "http://localhost:4100").replace(
  /\/+$/,
  "",
);

// Hop-by-hop headers must not be forwarded (RFC 7230 §6.1); `host` is dropped so
// fetch sets it for the upstream. `content-length`/`content-encoding` are
// dropped because we re-read the body, so the runtime recomputes the length.
const STRIP_REQUEST_HEADERS = new Set([
  "host",
  "connection",
  "keep-alive",
  "transfer-encoding",
  "content-length",
  "content-encoding",
  "accept-encoding",
  "upgrade",
  "proxy-authorization",
  "proxy-authenticate",
  "te",
  "trailer",
]);

async function proxy({ request }: { request: Request }) {
  const incoming = new URL(request.url);
  // Forward the path verbatim — the core mounts /api/v1.
  const target = `${CORE_URL}${incoming.pathname}${incoming.search}`;

  const headers = new Headers();
  for (const [key, value] of request.headers) {
    if (!STRIP_REQUEST_HEADERS.has(key.toLowerCase())) headers.set(key, value);
  }

  const method = request.method;
  const hasBody = method !== "GET" && method !== "HEAD";

  let upstream: Response;
  try {
    upstream = await fetch(target, {
      method,
      headers,
      body: hasBody ? await request.arrayBuffer() : undefined,
      redirect: "manual",
    });
  } catch {
    // Core unreachable (down, DNS, refused) — surface a gateway error, not 500.
    return new Response(JSON.stringify({ message: "core unavailable" }), {
      status: 502,
      headers: { "Content-Type": "application/json" },
    });
  }

  const respHeaders = new Headers(upstream.headers);
  respHeaders.delete("content-length");
  respHeaders.delete("content-encoding");
  respHeaders.delete("transfer-encoding");

  // SSE: opt out of any buffering layer so events flush as they arrive.
  if (respHeaders.get("content-type")?.includes("text/event-stream")) {
    respHeaders.set("Cache-Control", "no-cache, no-transform");
    respHeaders.set("X-Accel-Buffering", "no");
  }

  return new Response(upstream.body, {
    status: upstream.status,
    statusText: upstream.statusText,
    headers: respHeaders,
  });
}

export const Route = createFileRoute("/api/v1/$")({
  server: {
    handlers: {
      GET: proxy,
      POST: proxy,
      PUT: proxy,
      PATCH: proxy,
      DELETE: proxy,
    },
  },
});
