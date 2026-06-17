import { createServerFn } from "@tanstack/react-start";
import { getRequestHeaders } from "@tanstack/react-start/server";
import { auth } from "@/lib/auth";

// Resolve the better-auth session on the server side, forwarding the incoming
// request cookies so the in-app beforeLoad hook can keep the user signed in
// across hard refreshes.
export const getServerSession = createServerFn({ method: "GET" }).handler(
  async () => {
    const headers = getRequestHeaders() as unknown as HeadersInit;
    return auth.api.getSession({ headers: new Headers(headers) });
  },
);
