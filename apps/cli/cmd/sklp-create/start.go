package main

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed all:template-simple
var templateSimpleFS embed.FS

//go:embed all:template-eda
var templateEDAFS embed.FS

func templateFor(mode string) (embed.FS, string, error) {
	switch mode {
	case "simple":
		return templateSimpleFS, "template-simple", nil
	case "eda":
		return templateEDAFS, "template-eda", nil
	default:
		return embed.FS{}, "", fmt.Errorf("unknown --mode %q (want: simple, eda)", mode)
	}
}

var nameRe = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

const (
	digstackPrivatePattern = "github.com/digstack/*"
	digstackHTTPSPrefix    = "https://github.com/digstack/"
	digstackSSHPrefix      = "git@github.com:digstack/"
)

func newStartCmd() *cobra.Command {
	var (
		assumeYes bool
		mode      string
		render    string
	)
	cmd := &cobra.Command{
		Use:   "start <name>",
		Short: "Scaffold a new skalpai-style project in ./<name>",
		Long: `Create a new project from an embedded bootstrap template.

Two modes are available:
  --mode simple   Postgres + Echo REST (DDD CRUD)            [default]
  --mode eda      Postgres + NATS JetStream + go-eda (event-driven)

The web app is built on TanStack Start. Choose how it renders:
  --render ssr    Server-rendered via Nitro                  [default]
  --render spa    Client-only single-page app (static shell)

On first run, prompts to set GOPRIVATE and a git insteadOf rule so the
private skalpai sdk-go module can be fetched over SSH. Pass --yes to
apply both without confirmation.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStart(args[0], mode, render, assumeYes)
		},
	}
	cmd.Flags().BoolVarP(&assumeYes, "yes", "y", false, "apply machine config changes (GOPRIVATE, git insteadOf) without confirmation")
	cmd.Flags().StringVarP(&mode, "mode", "m", "simple", "template profile: simple | eda")
	cmd.Flags().StringVarP(&render, "render", "r", "ssr", "web rendering mode: ssr | spa")
	return cmd
}

func runStart(name, mode, render string, assumeYes bool) error {
	if !nameRe.MatchString(name) {
		return fmt.Errorf("name must be kebab-case [a-z][a-z0-9-]*, got %q", name)
	}
	if render != "ssr" && render != "spa" {
		return fmt.Errorf("unknown --render %q (want: ssr, spa)", render)
	}
	if _, err := os.Stat(name); err == nil {
		return fmt.Errorf("./%s already exists", name)
	}

	fsys, root, err := templateFor(mode)
	if err != nil {
		return err
	}

	if err := ensureMachineConfig(assumeYes); err != nil {
		return err
	}

	if err := extractTemplate(fsys, root, name); err != nil {
		return fmt.Errorf("extract template: %w", err)
	}
	if err := rename(name); err != nil {
		return fmt.Errorf("rename: %w", err)
	}
	if err := applyRenderMode(name, render); err != nil {
		return fmt.Errorf("set render mode: %w", err)
	}
	if err := copyEnvExample(name); err != nil {
		return fmt.Errorf("seed .env: %w", err)
	}
	if err := gitInit(name); err != nil {
		fmt.Fprintf(os.Stderr, "warn: git init failed (%v) — continuing\n", err)
	}

	fmt.Printf("✓ created ./%s (mode=%s, render=%s)\n\n", name, mode, render)
	fmt.Printf("next:\n")
	fmt.Printf("  cd %s\n", name)
	fmt.Printf("  git config core.hooksPath .githooks\n")
	fmt.Printf("  sklp dev\n")
	return nil
}

// applyRenderMode shapes the generated web app for the chosen rendering mode.
//
// The template ships the SSR variant (TanStack Start + Nitro). For "ssr" we
// just bake the env-driven flag to a literal. For "spa" we convert the app to
// a true client-only Vite SPA — TanStack Router (no Start server, no shell
// prerender), with a classic index.html + client entry. This sidesteps
// TanStack Start's SPA shell-prerender, which is not viable here.
func applyRenderMode(name, render string) error {
	web := filepath.Join(name, "apps", "web")
	if _, err := os.Stat(web); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if render != "spa" {
		// SSR is the template's native shape — nothing to transform.
		return nil
	}

	// SPA: replace the Start/Nitro setup with a client-only Vite app.
	writes := map[string]string{
		"vite.config.ts":           spaViteConfig,
		"index.html":               spaIndexHTML,
		"src/main.tsx":             spaMainTSX,
		"src/routes/__root.tsx":    spaRootTSX,
		"src/routes/index.tsx":     spaIndexRoute,
		"Caddyfile":                caddyfile,
	}
	for rel, content := range writes {
		target := filepath.Join(web, rel)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
			return err
		}
	}
	// Better Auth needs the TanStack Start server (better-auth runs server-side).
	// A SPA has no server, so drop the auth wiring entirely: the generated SPA
	// ships without auth, and the Caddyfile above serves the static build and
	// reverse-proxies /api to the Go core. Re-add auth by switching to SSR.
	authPaths := []string{
		"src/lib/auth.ts",
		"src/lib/auth-client.ts",
		"src/lib/db.ts",
		"src/lib/get-server-session.ts",
		"src/lib/mint-core-token.ts",
		"src/routes/api",            // api/auth/$ + api/core/$ server routes
		"src/routes/_protected.tsx", // pathless auth guard
		"src/routes/_protected",     // protected routes (index lives here in SSR)
		"src/routes/login.tsx",
		"src/routes/signup.tsx",
		"src/routes/verify-email.tsx",
	}
	for _, rel := range authPaths {
		if err := os.RemoveAll(filepath.Join(web, rel)); err != nil {
			return err
		}
	}
	// The SSR-only server entry script is meaningless for a SPA; drop it so the
	// generated package.json does not advertise a `node .output/server` start.
	if err := stripStartScript(filepath.Join(web, "package.json")); err != nil {
		return err
	}
	return nil
}

// stripStartScript removes the Nitro-only "start" script from package.json.
func stripStartScript(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	s := strings.Replace(string(b),
		"\n    \"start\": \"node .output/server/index.mjs\",", "", 1)
	if s == string(b) {
		return nil
	}
	return os.WriteFile(path, []byte(s), 0o644)
}

// --- SPA variant files (client-only Vite + TanStack Router) ---

const spaViteConfig = `import { defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";
import viteReact from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";

// Client-only single-page app: TanStack Router (no SSR), mounted from
// index.html via src/main.tsx. The API is reached through the dev proxy.
export default defineConfig({
  server: {
    port: 5273,
    proxy: { "/api": "http://localhost:4100" },
  },
  plugins: [
    tsconfigPaths({ projects: ["./tsconfig.json"] }),
    tailwindcss(),
    tanstackRouter({ target: "react", autoCodeSplitting: true }),
    viteReact(),
  ],
});
`

const spaIndexHTML = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>app</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
`

const spaMainTSX = `import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider } from "@tanstack/react-router";
import { getRouter } from "./router";
import "./styles.css";

const router = getRouter();

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>,
);
`

// SPA root: a plain layout route (no shellComponent / HeadContent / Scripts,
// which are TanStack Start server primitives). The document lives in
// index.html; this just renders the matched route into #root.
const spaRootTSX = `import { Outlet, createRootRoute } from "@tanstack/react-router";

export const Route = createRootRoute({
  notFoundComponent: NotFound,
  component: RootLayout,
});

function RootLayout() {
  return <Outlet />;
}

function NotFound() {
  return (
    <main className="mx-auto max-w-2xl p-6 font-sans">
      <h1 className="text-2xl font-semibold tracking-tight">404</h1>
      <p className="mt-2 text-sm text-gray-600">Page not found.</p>
    </main>
  );
}
`

// spaIndexRoute is the public home route for the SPA. The SSR template's index
// lives under the auth guard (_protected/index.tsx), which is removed for SPA;
// this restores a plain public landing at "/".
const spaIndexRoute = `import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";

export const Route = createFileRoute("/")({
  component: Home,
});

type Project = { id: string; name: string };

function Home() {
  const [projects, setProjects] = useState<Project[]>([]);

  useEffect(() => {
    // /api is reverse-proxied to the Go core (Caddyfile in prod, Vite proxy in
    // dev). This SPA ships without auth — switch to SSR for better-auth.
    fetch("/api/v1/projects")
      .then((r) => (r.ok ? r.json() : []))
      .then(setProjects)
      .catch(() => {});
  }, []);

  return (
    <main className="mx-auto max-w-2xl p-6 font-sans">
      <h1 className="text-2xl font-semibold tracking-tight">app</h1>
      <p className="mt-2 text-sm text-gray-600">
        Skalpai-style bootstrap (SPA). Edit{" "}
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
`

// caddyfile serves the SPA static build and reverse-proxies /api to the Go
// core. A SPA has no server runtime, so Caddy is the production front:
// SPA fallback (try_files → index.html) + same-origin /api proxy (no CORS).
const caddyfile = `# Production front for the SPA build.
# Run: caddy run --config ./Caddyfile
# CORE_UPSTREAM defaults to the in-cluster core service.
{
	auto_https off
}

:8080 {
	encode gzip

	# API requests are reverse-proxied to the Go core (same-origin, no CORS).
	handle /api/* {
		reverse_proxy {$CORE_UPSTREAM:http://core:4100}
	}

	# Everything else is the static SPA build with client-side routing fallback.
	handle {
		root * /srv
		try_files {path} /index.html
		file_server
	}
}
`

// ensureMachineConfig checks the two prerequisites for fetching the
// private skalpai sdk-go module and offers to set them. Without these,
// `go run` from the generated project fails with a HTTPS prompt.
func ensureMachineConfig(assumeYes bool) error {
	if err := ensureGoPrivate(assumeYes); err != nil {
		return err
	}
	return ensureGitInsteadOf(assumeYes)
}

func ensureGoPrivate(assumeYes bool) error {
	out, err := exec.Command("go", "env", "GOPRIVATE").Output()
	if err != nil {
		return fmt.Errorf("read GOPRIVATE: %w", err)
	}
	current := strings.TrimSpace(string(out))
	if strings.Contains(current, "github.com/digstack/") || strings.Contains(current, digstackPrivatePattern) {
		return nil
	}

	next := digstackPrivatePattern
	if current != "" {
		next = current + "," + digstackPrivatePattern
	}
	prompt := fmt.Sprintf("GOPRIVATE does not include %q.\n  current:  %q\n  proposed: %q\nApply with `go env -w GOPRIVATE=...`?", digstackPrivatePattern, current, next)
	if !assumeYes && !confirm(prompt) {
		fmt.Fprintln(os.Stderr, "skipping — `sklp dev` will fail until GOPRIVATE is set.")
		return nil
	}
	if err := exec.Command("go", "env", "-w", "GOPRIVATE="+next).Run(); err != nil {
		return fmt.Errorf("set GOPRIVATE: %w", err)
	}
	fmt.Printf("✓ GOPRIVATE=%s\n", next)
	return nil
}

func ensureGitInsteadOf(assumeYes bool) error {
	key := "url." + digstackSSHPrefix + ".insteadOf"
	out, _ := exec.Command("git", "config", "--global", "--get", key).Output()
	if strings.TrimSpace(string(out)) == digstackHTTPSPrefix {
		return nil
	}
	prompt := fmt.Sprintf("git is not configured to fetch %s over SSH.\nApply `git config --global %s %q`?", digstackHTTPSPrefix, key, digstackHTTPSPrefix)
	if !assumeYes && !confirm(prompt) {
		fmt.Fprintln(os.Stderr, "skipping — `sklp dev` will fail until the SSH rewrite is set.")
		return nil
	}
	if err := exec.Command("git", "config", "--global", key, digstackHTTPSPrefix).Run(); err != nil {
		return fmt.Errorf("set git insteadOf: %w", err)
	}
	fmt.Printf("✓ git rewrite %s → %s\n", digstackHTTPSPrefix, digstackSSHPrefix)
	return nil
}

func confirm(prompt string) bool {
	fmt.Fprintf(os.Stderr, "\n%s [y/N] ", prompt)
	r := bufio.NewReader(os.Stdin)
	line, err := r.ReadString('\n')
	if err != nil {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes":
		return true
	}
	return false
}

func extractTemplate(fsys embed.FS, root, dst string) error {
	return fs.WalkDir(fsys, root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(p, root)
		rel = strings.TrimPrefix(rel, "/")
		if rel == "" {
			return os.MkdirAll(dst, 0o755)
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		b, err := fsys.ReadFile(p)
		if err != nil {
			return err
		}
		// go:embed refuses to walk into subtrees that contain a go.mod, so
		// template/.../go.mod is shipped as go.mod.tmpl and restored here.
		target = strings.TrimSuffix(target, ".tmpl")
		mode := fs.FileMode(0o644)
		if strings.HasSuffix(rel, ".sh") || strings.HasPrefix(rel, ".githooks/") {
			mode = 0o755
		}
		return os.WriteFile(target, b, mode)
	})
}

// rename rewrites the template placeholders to the chosen project name.
// Mirrors scripts/rename.sh but runs in-process from the embedded copy.
func rename(name string) error {
	return filepath.WalkDir(name, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		s := string(b)
		s = strings.ReplaceAll(s, "@app/", "@"+name+"/")
		s = strings.ReplaceAll(s, `"app/`, `"`+name+`/`)
		s = strings.ReplaceAll(s, "module app/", "module "+name+"/")
		s = strings.ReplaceAll(s, "skaplai/", name+"/")
		if s == string(b) {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		return os.WriteFile(p, []byte(s), info.Mode())
	})
}

func copyEnvExample(dst string) error {
	src := filepath.Join(dst, ".env.example")
	b, err := os.ReadFile(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return os.WriteFile(filepath.Join(dst, ".env"), b, 0o600)
}

func gitInit(dir string) error {
	if _, err := exec.LookPath("git"); err != nil {
		return err
	}
	cmd := exec.Command("git", "init", "--quiet")
	cmd.Dir = dir
	return cmd.Run()
}
