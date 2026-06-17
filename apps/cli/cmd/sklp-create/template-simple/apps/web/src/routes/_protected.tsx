import { Outlet, createFileRoute, redirect } from "@tanstack/react-router";
import { getServerSession } from "@/lib/get-server-session";
import { mintCoreToken } from "@/lib/mint-core-token";

// Pathless layout guarding every route nested under it. Gates on the
// better-auth session (source of truth) and refreshes the short-lived JWT
// cookie the Go core verifies on each entry.
export const Route = createFileRoute("/_protected")({
  beforeLoad: async () => {
    const session = await getServerSession();
    if (!session) {
      throw redirect({ to: "/login" });
    }
    await mintCoreToken();
    return { session };
  },
  component: () => <Outlet />,
});
