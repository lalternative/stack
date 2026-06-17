import { createPlatformAuth } from "@lalternative/auth/server";
import { tanstackStartCookies } from "better-auth/tanstack-start";
import { pool } from "./db";

const authSecret = process.env.BETTER_AUTH_SECRET;
if (!authSecret) {
  throw new Error("BETTER_AUTH_SECRET environment variable is required");
}

export const auth = createPlatformAuth({
  database: pool,
  baseURL: process.env.BETTER_AUTH_URL ?? "http://localhost:5273",
  secret: authSecret,
  appName: process.env.APP_NAME ?? "app",
  // No transactional mailer wired by default: OTP codes are logged to stdout.
  // Wire a real mailer here once the app has an email transport.
  google: process.env.GOOGLE_CLIENT_ID
    ? {
        clientId: process.env.GOOGLE_CLIENT_ID,
        clientSecret: process.env.GOOGLE_CLIENT_SECRET!,
      }
    : undefined,
  github: process.env.GITHUB_CLIENT_ID
    ? {
        clientId: process.env.GITHUB_CLIENT_ID,
        clientSecret: process.env.GITHUB_CLIENT_SECRET!,
      }
    : undefined,
  plugins: [tanstackStartCookies()],
});

export type Auth = typeof auth;
