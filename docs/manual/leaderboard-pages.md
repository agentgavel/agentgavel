# Leaderboard GitHub Pages

How the static Opt-in / Unratified leaderboard in `dashboard/` reaches
the public web. Policy: [ADR 006](../adr/006-leaderboard-policy.md).
The ADR 006 **v0.3 addendum** (unsigned Opt-in blocked; samples only)
**expired** when E15 landed signed Opt-in — CI rule is now ADR 013
(`tab=opt-in` ⇒ sample **or** verified signature). Signature format and
key registry: [ADR 013](../adr/013-opt-in-signature-format.md). Maintainer
submission walkthrough: [opt-in-submission.md](opt-in-submission.md).
Provenance badges: [ADR 007](../adr/007-adapter-ratification.md). Use
cases: UC-022 (Pages publish), UC-032 (signed Opt-in).

## What gets published

The site is vanilla HTML/CSS/JS under `dashboard/` (no build step). GitHub
Pages serves that directory as static files: marketing home at `/`,
leaderboard at `/leaderboard/`.

| Path | Role |
| --- | --- |
| `dashboard/index.html` | Marketing home (inline CSS/JS) |
| `dashboard/leaderboard/index.html` | Two tabs: `#opt-in` and `#unratified` |
| `dashboard/leaderboard/app.js` / `style.css` | Load `../data/index.json` and render tables |
| `dashboard/data/index.json` | JSON array of entry filenames |
| `dashboard/data/<run_id>.json` | One scorecard entry (`data/schema.json`) |
| `dashboard/CNAME` | Custom domain `agentgavel.dev` |

`AgentGavel report --publish` is the only writer of real entries under
`data/`. See [dashboard/README.md](../../dashboard/README.md).

## Pages source: `pages.yml` from `main`

Workflow: [`.github/workflows/pages.yml`](../../.github/workflows/pages.yml).

- **Trigger:** push to `main` that touches `dashboard/**` or the workflow
  file itself, plus `workflow_dispatch`.
- **Artifact:** `actions/upload-pages-artifact` with `path: dashboard`.
- **Deploy:** `actions/deploy-pages@v4` under the `github-pages`
  environment.

Only content that is **merged on `main`** is uploaded. Branches, PRs, and
unmerged drafts never appear on the public site — there is no preview
Pages deploy for pull requests.

## Public URL

Canonical domain: `agentgavel.dev` (custom domain, `dashboard/CNAME`).
The default `agentgavel.github.io/agentgavel/` URL still resolves and
redirects to the custom domain once DNS + HTTPS enforcement finish.

```text
https://agentgavel.dev/
https://agentgavel.dev/leaderboard/
```

Entry list:

```text
https://agentgavel.dev/data/index.json
```

**Live (T14.20 verified 2026-09-05):** `build_type=workflow`, CNAME
`agentgavel.dev`. Smoke:

```bash
curl -sfL https://agentgavel.dev/ | grep -q 'Governance'
curl -sfL https://agentgavel.dev/leaderboard/ | grep -E 'id="opt-in"|id="unratified"'
curl -sfL https://agentgavel.dev/data/index.json   # JSON array of sample entries
```

HTTPS certificate provisioning / `https_enforced` may still catch up
asynchronously in GitHub Pages settings; HTTP and TLS via the shared
edge already serve the site. Flip **Enforce HTTPS** when the custom
cert shows ready.

## Opt-in vs Unratified (ADR 006)

Two tabs stay separate on purpose:

| Tab | Who appears | Provenance rule |
| --- | --- | --- |
| **Opt-in** | Maintainer-signed submissions | Bound to registered maintainer keys ([ADR 013](../adr/013-opt-in-signature-format.md)) |
| **Unratified** | AgentGavel-operated / unsolicited runs against public releases | Never on the primary tab; show honest provenance |

Unsolicited runs never appear on Opt-in. Unsigned or forged Opt-in
entries fail CI (`scripts/check-dashboard.sh` +
`scripts/verify-opt-in.sh`). How maintainers sign and open a PR:
[opt-in-submission.md](opt-in-submission.md).

**ADR 006 addendum status:** the v0.3 interim (reject `--tab opt-in`;
Opt-in = samples only) is **expired**. Real Opt-in rows require an
Ed25519 signature that verifies against an **active** key in
`dashboard/keys/registry.json` (ADR 013).

## Provenance badges (ADR 007)

Every entry carries one of three labels:

| `provenance` | Meaning |
| --- | --- |
| `ratified` | Framework maintainers reviewed or contributed the adapter |
| `provisional` | Core maintainers granted provisional ratification after outreach + review; expires unless renewed or upgraded |
| `unofficial` | Not ratified; default for unsolicited / Unratified publish path |

Author-affiliated adapters cannot skip to `ratified` via the provisional
path. The dashboard UI must keep the three-way distinction obvious.

## `report --publish` and signed Opt-in

**Unratified (default publish path):**

- `AgentGavel report --publish` (without `--sign`) writes
  **`tab: "unratified"`**.
- Entries copy `provenance` from the run Handshake (typically
  `unofficial` for unsolicited FakeAdapter / public-release runs).

```bash
./AgentGavel report --publish --dashboard dashboard <run-id>
```

Stdout is the path of the written `dashboard/data/<run-id>.json`. The
command also updates `dashboard/data/index.json`. Merge that change to
`main` before it can appear on Pages.

**Opt-in (signed, ADR 013):**

- `--tab opt-in` requires `--sign` with `--key` / `--key-id` (exit `2`
  without a signature; cites ADR 013).
- Full maintainer flow (sign → verify → PR → CI → Pages):
  [opt-in-submission.md](opt-in-submission.md).

```bash
./AgentGavel report --publish --sign --tab opt-in \
  --key /path/to/maintainer.priv.b64 \
  --key-id your-org-2026-09 \
  --framework "Exact Framework Display Name" \
  --adapter-name your.adapter.module \
  --dashboard dashboard \
  <run-id>
./AgentGavel verify-entry \
  --registry dashboard/keys/registry.json \
  dashboard/data/<run-id>.json
```

## Samples are labeled

Committed demo rows under `dashboard/data/` carry `sample: true` and a
framework name that includes `(sample)`:

- Opt-in demo: `sample-example-opt-in.json` — shows the Opt-in tab shape;
  not a real ratification (no signature required).
- Unratified demo: `sample-fakeadapter-unratified.json` — derived from a
  FakeAdapter run; `provenance: unofficial`.

The dashboard tags sample entries visibly as **sample**. CI
(`scripts/check-dashboard.sh`) enforces ADR 013:
`tab=opt-in` ⇒ (`sample=true` **or** signature verifies against an
**active** registry key). Samples stay unsigned; real Opt-in rows must
verify.

## Unmerged drafts do not appear

| State | On public Pages? |
| --- | --- |
| Entry only on a feature branch / open PR | No |
| Entry merged to `main` under `dashboard/data/` | Yes, after `pages.yml` deploys |
| Local `report --publish` into a checkout | Local only until you commit + merge |

Because `pages.yml` uploads `dashboard/` from `main` only, an unmerged
draft scorecard cannot leak onto `https://agentgavel.dev/`.

## Serve locally

From the repository root (no Pages account needed):

```bash
python3 -m http.server -d dashboard 8000
```

Open <http://127.0.0.1:8000/> for the marketing home and
<http://127.0.0.1:8000/leaderboard/> for Opt-in and Unratified. With the
committed samples, both tables have at least one labeled sample row.
