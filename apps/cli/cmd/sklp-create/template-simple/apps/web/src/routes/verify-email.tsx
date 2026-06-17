import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { VerifyEmailForm } from "@lalternative/auth";
import { authClient } from "@/lib/auth-client";
import { mintCoreToken } from "@/lib/mint-core-token";

export const Route = createFileRoute("/verify-email")({
  validateSearch: (search: Record<string, unknown>) => ({
    email: typeof search.email === "string" ? search.email : "",
  }),
  component: VerifyEmail,
});

function VerifyEmail() {
  const { email } = Route.useSearch();
  const navigate = useNavigate();

  return (
    <main className="mx-auto flex min-h-screen max-w-sm flex-col justify-center gap-6 p-6 font-sans">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Verify your email</h1>
        <p className="mt-1 text-sm text-gray-600">
          {email
            ? `We sent a 6-digit code to ${email}.`
            : "Enter the code we emailed you."}
        </p>
      </div>
      <VerifyEmailForm
        email={email}
        authClient={authClient}
        onSuccess={async () => {
          await mintCoreToken();
          await navigate({ to: "/" });
        }}
      />
    </main>
  );
}
