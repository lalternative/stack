import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { authClient } from "@/lib/auth-client";

export const Route = createFileRoute("/_protected/")({
  component: Home,
});

type Project = { id: string; name: string };

function Home() {
  const navigate = useNavigate();
  const [projects, setProjects] = useState<Project[]>([]);

  useEffect(() => {
    // Same-origin proxy to the Go core (see routes/api/core/$.ts) — the auth
    // cookie travels automatically, no CORS. /api/core/* maps 1:1 to core /*.
    fetch("/api/core/v1/projects")
      .then((r) => (r.ok ? r.json() : []))
      .then(setProjects)
      .catch(() => {});
  }, []);

  async function logout() {
    await authClient.signOut();
    await navigate({ to: "/login" });
  }

  return (
    <main className="mx-auto max-w-2xl p-6 font-sans">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold tracking-tight">app</h1>
        <button
          onClick={logout}
          className="text-sm text-gray-600 underline"
        >
          Sign out
        </button>
      </div>
      <p className="mt-2 text-sm text-gray-600">
        Skalpai-style bootstrap on TanStack Start with better-auth. Edit{" "}
        <code className="rounded bg-gray-100 px-1">
          src/routes/_protected/index.tsx
        </code>{" "}
        to begin.
      </p>
      <ul className="mt-4 list-disc pl-5">
        {projects.map((p) => (
          <li key={p.id}>{p.name}</li>
        ))}
      </ul>
    </main>
  );
}
