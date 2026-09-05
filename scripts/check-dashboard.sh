#!/usr/bin/env bash
# Validate dashboard/data entries and leaderboard HTML (ADR 006 / ADR 007).
# Usage: bash scripts/check-dashboard.sh [dashboard-dir]
# Default dashboard-dir: dashboard
# Exits 0 on success, 1 on validation failure.
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

dash="${1:-dashboard}"
data_dir="$dash/data"
# Leaderboard tables live under leaderboard/ after the marketing home landed;
# fall back to dash/index.html for older layouts.
html=""
for candidate in "$dash/leaderboard/index.html" "$dash/index.html"; do
  if [[ -f "$candidate" ]]; then
    html="$candidate"
    break
  fi
done

if [[ ! -d "$data_dir" ]]; then
  echo "check-dashboard: missing data dir: $data_dir" >&2
  exit 1
fi
if [[ -z "$html" ]]; then
  echo "check-dashboard: missing leaderboard HTML under $dash" >&2
  exit 1
fi

python3 - "$data_dir" "$html" <<'PY'
import json
import os
import sys

data_dir, html_path = sys.argv[1], sys.argv[2]

REQUIRED = [
    "run_id",
    "framework",
    "adapter",
    "adapter_version",
    "provenance",
    "tab",
    "sample",
    "gsi",
    "grade",
    "pillars",
    "catastrophic",
    "na",
    "fingerprint",
    "generated_at",
]
PROVENANCE = {"ratified", "provisional", "unofficial"}
TABS = {"opt-in", "unratified"}

errors = []

index_path = os.path.join(data_dir, "index.json")
if not os.path.isfile(index_path):
    errors.append(f"missing {index_path}")
    print("\n".join(errors), file=sys.stderr)
    sys.exit(1)

try:
    with open(index_path, encoding="utf-8") as f:
        index = json.load(f)
except json.JSONDecodeError as e:
    print(f"index.json: invalid JSON: {e}", file=sys.stderr)
    sys.exit(1)

if not isinstance(index, list):
    print("index.json: must be a JSON array of entry filenames", file=sys.stderr)
    sys.exit(1)

entry_files = sorted(
    name
    for name in os.listdir(data_dir)
    if name.endswith(".json") and name not in ("schema.json", "index.json")
)

# index.json must list exactly the entry files (sorted).
sorted_index = sorted(index)
if sorted_index != entry_files:
    errors.append(
        "index.json must equal the sorted entry file list: "
        f"index={sorted_index!r} files={entry_files!r}"
    )

for name in index:
    if not isinstance(name, str) or not name:
        errors.append(f"index.json: invalid entry name {name!r}")
        continue
    path = os.path.join(data_dir, name)
    if not os.path.isfile(path):
        errors.append(f"index.json lists missing file: {name}")
        continue
    try:
        with open(path, encoding="utf-8") as f:
            entry = json.load(f)
    except json.JSONDecodeError as e:
        errors.append(f"{name}: invalid JSON: {e}")
        continue
    if not isinstance(entry, dict):
        errors.append(f"{name}: expected object")
        continue
    for key in REQUIRED:
        if key not in entry:
            errors.append(f"{name}: missing required key {key}")
    if "provenance" in entry and entry["provenance"] not in PROVENANCE:
        errors.append(f"{name}: provenance must be one of {sorted(PROVENANCE)}")
    if "tab" in entry and entry["tab"] not in TABS:
        errors.append(f"{name}: tab must be one of {sorted(TABS)}")
    if entry.get("tab") == "opt-in" and entry.get("sample") is not True:
        errors.append(f"{name}: tab=opt-in requires sample=true (v0.3 ADR 006 addendum)")
    if "sample" in entry and not isinstance(entry["sample"], bool):
        errors.append(f"{name}: sample must be boolean")
    if "gsi" in entry and not isinstance(entry["gsi"], (int, float)):
        errors.append(f"{name}: gsi must be a number")
    if "pillars" in entry and not isinstance(entry["pillars"], dict):
        errors.append(f"{name}: pillars must be an object")
    if "catastrophic" in entry and not isinstance(entry["catastrophic"], list):
        errors.append(f"{name}: catastrophic must be an array")
    if "na" in entry and not isinstance(entry["na"], list):
        errors.append(f"{name}: na must be an array")
    if "fingerprint" in entry and not isinstance(entry["fingerprint"], dict):
        errors.append(f"{name}: fingerprint must be an object")

# Also validate entry files not listed in index (so a stray bad file fails).
for name in entry_files:
    if name in index:
        continue
    errors.append(f"entry file not listed in index.json: {name}")

with open(html_path, encoding="utf-8") as f:
    html = f.read()
for section_id in ("opt-in", "unratified"):
    needle = f'id="{section_id}"'
    if needle not in html:
        errors.append(f"{html_path}: missing {needle}")

if errors:
    for err in errors:
        print(f"check-dashboard: {err}", file=sys.stderr)
    sys.exit(1)

print(f"dashboard ok: {len(entry_files)} entr(y/ies), html={html_path}")
PY
