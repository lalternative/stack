import { createFileRoute } from "@tanstack/react-router";
import { auth } from "@/lib/auth";

// Server route mounting the better-auth handler. `server.handlers` is codegen'd
// by the TanStack Start plugin.
export const Route = createFileRoute("/api/auth/$")({
  server: {
    handlers: {
      GET: async ({ request }: { request: Request }) => auth.handler(request),
      POST: async ({ request }: { request: Request }) => auth.handler(request),
    },
  },
});
