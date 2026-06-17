import { useEffect, useState } from "react";
import type { ProjectProjectDTO } from "@app/sdk";

import { coreAPI } from "./lib/api";

export function App() {
  const [projects, setProjects] = useState<ProjectProjectDTO[]>([]);
  useEffect(() => {
    coreAPI
      .listProjects()
      .then(setProjects)
      .catch(() => {});
  }, []);
  return (
    <main className="mx-auto max-w-2xl p-6">
      <h1 className="text-2xl font-semibold tracking-tight">app</h1>
      <p className="mt-2 text-muted-foreground">
        Skalpai-style stack. Edit{" "}
        <code className="rounded bg-muted px-1.5 py-0.5 text-sm">src/App.tsx</code> to begin.
      </p>
      <ul className="mt-4 space-y-1">
        {projects.map((p) => (
          <li key={p.id} className="rounded-md border border-border px-3 py-2">
            {p.name}
          </li>
        ))}
      </ul>
    </main>
  );
}
