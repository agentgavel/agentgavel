# Signed Opt-in submission

How framework maintainers publish a real Opt-in leaderboard row through
GitHub. Policy and trust root:
[ADR 012](../adr/012-github-native-signed-submissions.md) (GitHub-native
PR path) and [ADR 013](../adr/013-opt-in-signature-format.md) (Ed25519
signature format + key registry). Use case: UC-032.

Pages publish loop after merge:
[Leaderboard GitHub Pages](leaderboard-pages.md). Key registration:
[dashboard/keys/README.md](../../dashboard/keys/README.md).

## What you submit

v1.0 Opt-in rows are **GitHub-native** (ADR 012): you produce a signed
scorecard entry, open a pull request that adds
`dashboard/data/<run_id>.json` and updates `dashboard/data/index.json`,
CI verifies the signature, and merge to `main` deploys via
`.github/workflows/pages.yml`.

There is no separate hosted submission API. Trust is the Ed25519
signature bound to an **active** key in
`dashboard/keys/registry.json`, not an identity-provider login.

## Prerequisites

1. An **active** key for your framework in
   `dashboard/keys/registry.json` (register via PR first; see
   [dashboard/keys/README.md](../../dashboard/keys/README.md)).
2. The matching private key as a base64 file on disk (keep it out of git).
3. A built `AgentGavel` binary (or `GOWORK=off go run ./cmd/AgentGavel`).

```bash
export GOWORK=off
make build
# or: go build -o AgentGavel ./cmd/AgentGavel
./AgentGavel version
```

The registry `framework` string must match the entry's `framework`
field exactly (no `(sample)` suffix on real Opt-in rows).

## Workflow

### 1. Generate a scorecard

Run the security suite (or your ratified adapter path) so
`results/<run-id>/` contains a scorecard or summary the report command
can load:

```bash
./AgentGavel run --adapter <your-adapter> --suite security --out . --run-id <run-id>
./AgentGavel report <run-id>
```

Confirm the text (or `--json`) scorecard looks right before you sign.

### 2. Sign with `report --sign`

Sign either from a completed run or from an existing entry JSON.
Requires `--key` (path to base64 Ed25519 private key) and `--key-id`
(must match an **active** registry row).

**From a run** (builds an Opt-in entry with `sample=false`, then signs):

```bash
./AgentGavel report --sign \
  --key /path/to/maintainer.priv.b64 \
  --key-id your-org-2026-09 \
  --framework "Exact Framework Display Name" \
  --adapter-name your.adapter.module \
  --out dashboard/data/<run-id>.json \
  <run-id>
```

**From an existing entry file:**

```bash
./AgentGavel report --sign \
  --key /path/to/maintainer.priv.b64 \
  --key-id your-org-2026-09 \
  --entry /path/to/unsigned-entry.json \
  --out dashboard/data/<run-id>.json
```

Without `--out`, the signed JSON goes to stdout.

**Sign and publish in one step** (writes the entry and updates
`index.json`):

```bash
./AgentGavel report --publish --sign --tab opt-in \
  --key /path/to/maintainer.priv.b64 \
  --key-id your-org-2026-09 \
  --framework "Exact Framework Display Name" \
  --adapter-name your.adapter.module \
  --dashboard dashboard \
  <run-id>
```

Stdout is the path of the written `dashboard/data/<run-id>.json`.
Unsigned `--tab opt-in` is rejected (exit `2`, cites ADR 013).

Signed fields follow ADR 013: Ed25519 over canonical JSON of the entry
**without** `signature`, `key_id`, or `sample`. The on-disk file carries
`key_id` and `signature` (base64 raw 64-byte signature).

### 3. Verify locally (optional)

Before you open a PR, check the entry against the registry:

```bash
./AgentGavel verify-entry \
  --registry dashboard/keys/registry.json \
  dashboard/data/<run-id>.json
```

Exit `0` on success, `1` on verification failure, `2` on usage errors.

Or run the same script CI uses:

```bash
bash scripts/verify-opt-in.sh
```

That script verifies every `tab=opt-in` entry with `sample` not `true`.

### 4. Open a PR that adds the data entry and index

Commit both:

| Path | Change |
| --- | --- |
| `dashboard/data/<run_id>.json` | Signed Opt-in entry (`tab: "opt-in"`, `sample: false`, `key_id`, `signature`) |
| `dashboard/data/index.json` | Filename listed in the index array |

If you used `report --publish --sign --tab opt-in`, both files are
already updated in your checkout. Open a PR against `main` with those
changes (and nothing else unless maintainers ask).

Do not commit private keys. Do not reuse the fixture key
`example-framework-test-1` for production Opt-in rows.

### 5. CI verifies

On `pull_request` and `push` to `main`, the **`dashboard`** job in
[`.github/workflows/ci.yml`](../../.github/workflows/ci.yml) runs:

1. `bash scripts/check-dashboard.sh` — schema and ADR 006/007 rules,
   including Opt-in = sample **or** verified signature (ADR 013).
2. `bash scripts/verify-opt-in.sh` — `AgentGavel verify-entry` for every
   non-sample Opt-in file against `dashboard/keys/registry.json`.

A bad or missing signature fails the job. Fix the entry or key
registration, then push again.

### 6. Merge → Pages

After CI is green and the PR merges to `main`,
[`.github/workflows/pages.yml`](../../.github/workflows/pages.yml)
uploads `dashboard/` and deploys GitHub Pages. Only content on `main`
appears on the public site.

```text
https://agentgavel.dev/leaderboard/
```

Unmerged PR drafts never appear on Pages. See
[Leaderboard GitHub Pages](leaderboard-pages.md).

## Samples need no signature

Committed demo rows under `dashboard/data/` with `sample: true` (for
example `sample-example-opt-in.json`) demonstrate the Opt-in tab shape.
They **do not** need `key_id` or `signature`.

CI rule (ADR 013): `tab=opt-in` ⇒ (`sample=true` **or** signature
verifies against an **active** registry key whose `framework` matches
the entry). Samples stay unsigned; real Opt-in rows must verify.

## Quick reference

| Step | Command / action |
| --- | --- |
| Generate scorecard | `AgentGavel run` … then `AgentGavel report <run-id>` |
| Sign | `AgentGavel report --sign --key … --key-id …` |
| Sign + write data/index | `AgentGavel report --publish --sign --tab opt-in …` |
| Local verify | `AgentGavel verify-entry --registry dashboard/keys/registry.json <entry>` |
| CI verify | `dashboard` job → `scripts/check-dashboard.sh` + `scripts/verify-opt-in.sh` |
| Publish | Merge to `main` → `pages.yml` → agentgavel.dev |

## Related

- [ADR 012: GitHub-native signed Opt-in submissions](../adr/012-github-native-signed-submissions.md)
- [ADR 013: Opt-in signature format and key registry](../adr/013-opt-in-signature-format.md)
- [Leaderboard GitHub Pages](leaderboard-pages.md)
- [dashboard/keys/README.md](../../dashboard/keys/README.md)
- [dashboard/README.md](../../dashboard/README.md)
