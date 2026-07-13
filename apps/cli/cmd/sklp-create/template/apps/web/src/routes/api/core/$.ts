import { createFileRoute } from "@tanstack/react-router";

// Same-origin proxy to the core API. The browser calls /api/core/* (relative,
// served by this TanStack Start server), and this handler forwards to core.
// Because every browser request is same-origin, there is NO CORS: no preflight,
// no ALLOWED_ORIGINS to keep in sync. CORE_API_URL points at the core service.
const CORE_URL = (process.env.CORE_API_URL ?? "http://localhost:4100").replace(/\/+$/, "");

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

async function proxy({ request, params }: { request: Request; params: { _splat?: string } }) {
  const incoming = new URL(request.url);
  const path = params._splat ?? "";
  const target = `${CORE_URL}/${path}${incoming.search}`;

  const headers = new Headers();
  for (const [key, value] of request.headers) {
    if (!STRIP_REQUEST_HEADERS.has(key.toLowerCase())) headers.set(key, value);
  }

  const method = request.method;
  const hasBody = method !== "GET" && method !== "HEAD";

  const upstream = await fetch(target, {
    method,
    headers,
    body: hasBody ? await request.arrayBuffer() : undefined,
    redirect: "manual",
  });

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

export const Route = createFileRoute("/api/core/$")({
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
