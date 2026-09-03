# Leaderboard GitHub Pages

How the static Opt-in / Unratified leaderboard in `dashboard/` reaches
the public web. Policy: [ADR 006](../adr/006-leaderboard-policy.md)
(including the v0.3 addendum). Provenance badges:
[ADR 007](../adr/007-adapter-ratification.md). Use case: UC-022.

## What gets published

The site is vanilla HTML/CSS/JS under `dashboard/` (no build step). GitHub
Pages serves that directory as static files.

| Path | Role |
| --- | --- |
| `dashboard/index.html` | Two tabs: `#opt-in` and `#unratified` |
| `dashboard/app.js` / `style.css` | Load `data/index.json` and render tables |
| `dashboard/data/index.json` | JSON array of entry filenames |
| `dashboard/data/<run_id>.json` | One scorecard entry (`data/schema.json`) |

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

Expected pattern after Pages is enabled (founder task T14.20):

```text
https://agentgavel.github.io/agentgavel/
```

Entry list:

```text
https://agentgavel.github.io/agentgavel/data/index.json
```

Until T14.20 enables the repository Pages setting (`build_type=workflow`)
and a successful `pages.yml` run, that URL may not resolve. Do not treat
it as live until the founder smoke check (`curl` for `id="opt-in"` /
`id="unratified"` and a JSON `index.json`) passes. If the live URL
differs from the pattern above, record the actual URL here after T14.20.

## Opt-in vs Unratified (ADR 006)

Two tabs stay separate on purpose:

| Tab | Who appears | Provenance rule |
| --- | --- | --- |
| **Opt-in** | Maintainer-signed submissions (v1.0+) | Bound to registered maintainer keys |
| **Unratified** | AgentGavel-operated / unsolicited runs against public releases | Never on the primary tab; show honest provenance |

Unsolicited runs never appear on Opt-in. Malicious auto-submissions are
rejected once signatures land (v1.0). Until then, the ADR 006 **addendum**
applies (see below).

## Provenance badges (ADR 007)

Every entry carries one of three labels:

| `provenance` | Meaning |
| --- | --- |
| `ratified` | Framework maintainers reviewed or contributed the adapter |
| `provisional` | Core maintainers granted provisional ratification after outreach + review; expires unless renewed or upgraded |
| `unofficial` | Not ratified; default for unsolicited / Unratified publish path |

Author-affiliated adapters cannot skip to `ratified` via the provisional
path. The dashboard UI must keep the three-way distinction obvious.

## `report --publish` in v0.3 (until v1.0)

Until the v1.0 signed submission process exists
([ADR 006](../adr/006-leaderboard-policy.md) addendum):

- `AgentGavel report --publish` writes **`tab: "unratified"` only**.
- `--tab opt-in` is rejected (exit `2`, cites ADR 006).
- Published entries copy `provenance` from the run Handshake (typically
  `unofficial` for unsolicited FakeAdapter / public-release runs).
- Real Opt-in rows wait for E15 signatures; the CI rule then flips to
  "opt-in requires a verified signature".

Example (after a security suite run has produced a scorecard under
`results/<run-id>/`):

```bash
./AgentGavel report --publish --dashboard dashboard <run-id>
```

Stdout is the path of the written `dashboard/data/<run-id>.json`. The
command also updates `dashboard/data/index.json`. Merge that change to
`main` before it can appear on Pages.

## Samples are labeled

Committed demo rows under `dashboard/data/` carry `sample: true` and a
framework name that includes `(sample)`:

- Opt-in demo: `sample-example-opt-in.json` — shows the Opt-in tab shape;
  not a real ratification.
- Unratified demo: `sample-fakeadapter-unratified.json` — derived from a
  FakeAdapter run; `provenance: unofficial`.

The dashboard tags sample entries visibly as **sample**. CI
(`scripts/check-dashboard.sh`) enforces `tab=opt-in ⇒ sample=true` so a
real run cannot be promoted onto Opt-in by editing JSON in v0.3.

## Unmerged drafts do not appear

| State | On public Pages? |
| --- | --- |
| Entry only on a feature branch / open PR | No |
| Entry merged to `main` under `dashboard/data/` | Yes, after `pages.yml` deploys |
| Local `report --publish` into a checkout | Local only until you commit + merge |

Because `pages.yml` uploads `dashboard/` from `main` only, an unmerged
draft scorecard cannot leak onto
`https://agentgavel.github.io/agentgavel/`.

## Serve locally

From the repository root (no Pages account needed):

```bash
python3 -m http.server -d dashboard 8000
```

Open <http://127.0.0.1:8000/>. You should see both Opt-in and Unratified
sections; with the committed samples, both tables have at least one
labeled sample row.
