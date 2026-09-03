# AgentGavel site (GitHub Pages)

Vanilla HTML/CSS/JS under `dashboard/` — marketing home plus the Opt-in /
Unratified leaderboard. No build step and no npm dependency.

## Serve locally

From the repository root:

```bash
python3 -m http.server -d dashboard 8000
```

Then open:

- Marketing home: <http://127.0.0.1:8000/>
- Leaderboard: <http://127.0.0.1:8000/leaderboard/>

## Layout

| Path | Role |
| ---- | ---- |
| `index.html` / `site.css` / `site.js` | Marketing home (cinematic hero + scroll sections) |
| `brand/` | Seal mark assets served with Pages |
| `leaderboard/` | Opt-in + Unratified tabs (ADR 006) |
| `data/index.json` | JSON array of entry filenames under `data/` |
| `data/<run_id>.json` | One leaderboard entry (see `data/schema.json`) |
| `data/schema.json` | JSON Schema draft-07 for entries |

`AgentGavel report --publish` is the only writer of `data/`. Do not hand-edit
`index.json` or entry files for real runs; sample rows for the Opt-in tab are
an exception documented in the ADR 006 addendum and are labeled
`sample: true`.

## Rendering

`leaderboard/app.js` loads `../data/index.json`, fetches each listed entry,
and fills one table per tab (`#opt-in`, `#unratified`). An empty index
renders empty tables.
