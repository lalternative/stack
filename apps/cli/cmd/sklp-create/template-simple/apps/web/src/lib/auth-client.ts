import { createPlatformAuthClient } from "@lalternative/auth/client";

// On the server, `window` is undefined and the package default fallback is
// `http://localhost:3000` — wrong for this app (web runs on :5273). Pass an
// explicit baseURL so SSR loaders and beforeLoad hooks reach the right origin.
// The browser keeps window.location.origin via the package's default branch.
const baseURL =
  typeof window === "undefined"
    ? process.env.BETTER_AUTH_URL ?? "http://localhost:5273"
    : window.location.origin;

export const authClient = createPlatformAuthClient({ baseURL });
