#!/usr/bin/env bash
# Verify fixture paths referenced by the suite catalog and fixtures/manifest.json
# exist on disk. Exits non-zero on the first missing path.
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

missing=0

fail() {
  echo "missing fixture: $1" >&2
  missing=1
}

# Paths listed in fixtures/manifest.json must exist.
python3 - <<'PY'
import json
import os
import sys

with open("fixtures/manifest.json", encoding="utf-8") as f:
    manifest = json.load(f)

missing = []
for scenario_id, paths in manifest.get("scenarios", {}).items():
    if not isinstance(paths, list):
        missing.append(f"{scenario_id}: expected list of paths, got {type(paths).__name__}")
        continue
    for path in paths:
        if not os.path.isfile(path):
            missing.append(f"{scenario_id}: {path}")

if missing:
    for item in missing:
        print(f"missing fixture: {item}", file=sys.stderr)
    sys.exit(1)
print(f"manifest ok: {sum(len(v) for v in manifest.get('scenarios', {}).values())} path(s)")
PY

# Paths under fixtures/ referenced from suites/ source must exist.
# Skip the manifest itself when it appears only as a documentation mention.
while IFS= read -r path; do
  case "$path" in
    fixtures/manifest.json) continue ;;
  esac
  if [[ ! -f "$path" ]]; then
    fail "$path (referenced from suites/)"
  else
    echo "suite ref ok: $path"
  fi
done < <(grep -RhoE 'fixtures/[A-Za-z0-9_./-]+' suites || true)

if [[ "$missing" -ne 0 ]]; then
  exit 1
fi

echo "fixture path checks passed"
