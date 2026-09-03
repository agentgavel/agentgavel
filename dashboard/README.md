# AgentGavel static leaderboard

Vanilla HTML/CSS/JS site for the Opt-in and Unratified tabs defined in
[ADR 006](../docs/adr/006-leaderboard-policy.md). Provenance badges use the
three-way labels from [ADR 007](../docs/adr/007-adapter-ratification.md)
(`ratified` / `provisional` / `unofficial`).

There is no build step and no npm dependency. Serve the directory as static
files.

## Serve locally

From the repository root:

```bash
python3 -m http.server -d dashboard 8000
```

Then open <http://127.0.0.1:8000/>.

## Data layout

| Path | Role |
| ---- | ---- |
| `data/index.json` | JSON array of entry filenames under `data/` |
| `data/<run_id>.json` | One leaderboard entry (see `data/schema.json`) |
| `data/schema.json` | JSON Schema draft-07 for entries |

`AgentGavel report --publish` is the only writer of `data/`. Do not hand-edit
`index.json` or entry files for real runs; sample rows for the Opt-in tab are
an exception documented in the ADR 006 addendum and are labeled
`sample: true`.

## Rendering

`app.js` loads `data/index.json`, fetches each listed entry, and fills one
table per tab (`#opt-in`, `#unratified`). An empty index renders empty tables.
