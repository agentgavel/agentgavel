# AgentGavel site (GitHub Pages)

Vanilla HTML/CSS/JS under `dashboard/`. No build step and no npm dependency.

## Serve locally

```bash
python3 -m http.server -d dashboard 8000
```

- Marketing home: <http://127.0.0.1:8000/>
- Leaderboard: <http://127.0.0.1:8000/leaderboard/>

## Layout

| Path | Role |
| ---- | ---- |
| `index.html` | Marketing home (cinematic hero + editorial sections; inline CSS/JS) |
| `brand/` | Seal mark assets served with Pages |
| `leaderboard/` | Opt-in + Unratified tabs (ADR 006) |
| `data/` | Scorecard entries; written by `AgentGavel report --publish` |
| `CNAME` | Custom domain `agentgavel.dev` |

`leaderboard/app.js` loads `../data/index.json` and fills `#opt-in` /
`#unratified`.
