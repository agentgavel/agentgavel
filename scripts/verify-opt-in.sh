#!/usr/bin/env bash
# Verify non-sample Opt-in dashboard entries against dashboard/keys/registry.json (ADR 013).
# Usage: bash scripts/verify-opt-in.sh [dashboard-dir]
# Default dashboard-dir: dashboard
# Exits 0 when every tab=opt-in sample!=true entry verifies; 1 on failure.
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

dash="${1:-dashboard}"
data_dir="$dash/data"
registry="$dash/keys/registry.json"

if [[ ! -d "$data_dir" ]]; then
  echo "verify-opt-in: missing data dir: $data_dir" >&2
  exit 1
fi
if [[ ! -f "$registry" ]]; then
  echo "verify-opt-in: missing registry: $registry" >&2
  exit 1
fi

# Prefer a prebuilt binary (CI), else go run.
run_verify() {
  local entry_path="$1"
  if [[ -n "${AGENTGAVEL_BIN:-}" ]]; then
    "$AGENTGAVEL_BIN" verify-entry --registry "$registry" "$entry_path"
  elif [[ -x "$root/AgentGavel" ]]; then
    "$root/AgentGavel" verify-entry --registry "$registry" "$entry_path"
  else
    GOWORK=off go run ./cmd/AgentGavel verify-entry --registry "$registry" "$entry_path"
  fi
}

# List entry files that need signature verification (tab=opt-in, sample!=true).
to_verify="$(python3 - "$data_dir" <<'PY'
import json
import os
import sys

data_dir = sys.argv[1]
for name in sorted(os.listdir(data_dir)):
    if not name.endswith(".json") or name in ("schema.json", "index.json"):
        continue
    path = os.path.join(data_dir, name)
    try:
        with open(path, encoding="utf-8") as f:
            entry = json.load(f)
    except (OSError, json.JSONDecodeError) as e:
        print(f"ERROR:{name}:{e}", file=sys.stderr)
        sys.exit(2)
    if not isinstance(entry, dict):
        print(f"ERROR:{name}:not an object", file=sys.stderr)
        sys.exit(2)
    if entry.get("tab") != "opt-in":
        continue
    if entry.get("sample") is True:
        continue
    print(path)
PY
)"

checked=0
while IFS= read -r path; do
  [[ -z "$path" ]] && continue
  checked=$((checked + 1))
  name="$(basename "$path")"
  if ! run_verify "$path"; then
    echo "verify-opt-in: bad signature or verify failed: $name" >&2
    exit 1
  fi
done <<EOF
$to_verify
EOF

echo "verify-opt-in ok: $checked non-sample opt-in entr(y/ies) verified (registry=$registry)"
