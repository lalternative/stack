#!/usr/bin/env bash
# Rename the bootstrap placeholders to a real project name.
#
# Replaces:
#   app/           → <name>/        (Go module paths)
#   @app/          → @<name>/       (npm scopes)
#   skaplai/       → <name>/        (Scaleway registry path component)
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "usage: scripts/rename.sh <name>" >&2
  exit 1
fi

NAME=$1
if ! [[ "$NAME" =~ ^[a-z][a-z0-9-]*$ ]]; then
  echo "name must be kebab-case [a-z0-9-]" >&2
  exit 1
fi

ROOT=$(cd "$(dirname "$0")/.." && pwd)
cd "$ROOT"

# Skip generated + binary dirs.
mapfile -t FILES < <(git ls-files 2>/dev/null || find . \
  -not -path './node_modules/*' \
  -not -path './.git/*' \
  -not -path './bin/*' \
  -not -path './dist/*' \
  -not -path './.sklp/cache/*' \
  -not -path './.sklp/logs/*' \
  -type f)

for f in "${FILES[@]}"; do
  # macOS / BSD sed compat
  sed -i.bak \
    -e "s|@app/|@${NAME}/|g" \
    -e "s|\"app/|\"${NAME}/|g" \
    -e "s|module app/|module ${NAME}/|g" \
    -e "s|skaplai/|${NAME}/|g" \
    "$f" && rm -f "${f}.bak"
done

echo "renamed → ${NAME}"
echo "next: cp .env.example .env && sklp dev"
