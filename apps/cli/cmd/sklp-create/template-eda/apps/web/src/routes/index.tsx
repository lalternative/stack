import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";

export const Route = createFileRoute("/")({
  component: Home,
});

type Project = { id: string; name: string };

function Home() {
  const [projects, setProjects] = useState<Project[]>([]);

  useEffect(() => {
    fetch("/api/v1/projects")
      .then((r) => (r.ok ? r.json() : []))
      .then(setProjects)
      .catch(() => {});
  }, []);

  return (
    <main className="mx-auto max-w-2xl p-6 font-sans">
      <h1 className="text-2xl font-semibold tracking-tight">app</h1>
      <p className="mt-2 text-sm text-gray-600">
        Skalpai-style stack on TanStack Start. Edit{" "}
        <code className="rounded bg-gray-100 px-1">src/routes/index.tsx</code> to
        begin.
      </p>
      <ul className="mt-4 list-disc pl-5">
        {projects.map((p) => (
          <li key={p.id}>{p.name}</li>
        ))}
      </ul>
    </main>
  );
}
