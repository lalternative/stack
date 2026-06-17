import { createServerFn } from "@tanstack/react-start";
import { getRequestHeaders, setResponseHeader } from "@tanstack/react-start/server";
import { SignJWT } from "jose";
import { auth } from "@/lib/auth";

const TOKEN_TTL_MS = 15 * 60 * 1000;

// Mint a short-lived HS256 JWT for the current better-auth session and set it
// as the `token` cookie that apps/core (Go) reads and verifies. Returns true
// when a token was issued, false when there is no active session.
//
// Called after sign-in (login, verify-email) and on every entry into the
// protected layout so the JWT cookie stays in sync with the longer-lived
// better-auth session cookie.
export const mintCoreToken = createServerFn({ method: "GET" }).handler(
  async (): Promise<boolean> => {
    const headers = getRequestHeaders() as unknown as HeadersInit;
    const session = await auth.api.getSession({ headers: new Headers(headers) });
    if (!session) return false;

    if (!process.env.JWT_SECRET) {
      throw new Error("JWT_SECRET environment variable is required");
    }
    const secret = new TextEncoder().encode(process.env.JWT_SECRET);
    const expiresAt = new Date(Date.now() + TOKEN_TTL_MS);

    const token = await new SignJWT({
      sub: session.user.id,
      email: session.user.email,
      name: session.user.name,
    })
      .setProtectedHeader({ alg: "HS256" })
      .setIssuedAt()
      .setExpirationTime(expiresAt)
      .sign(secret);

    const secure = process.env.NODE_ENV === "production" ? "; Secure" : "";
    const cookieDomain = process.env.COOKIE_DOMAIN
      ? `; Domain=${process.env.COOKIE_DOMAIN}`
      : "";
    setResponseHeader(
      "Set-Cookie",
      `token=${token}; HttpOnly; SameSite=Lax; Path=/; Expires=${expiresAt.toUTCString()}${secure}${cookieDomain}`,
    );
    return true;
  },
);
