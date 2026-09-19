# Benchmark operations

How to run AgentGavel oracle-mode baselines and (when credentials exist)
live Sire dogfood. Run from the repository root unless noted.

Related:

- [ADR 015](../adr/015-stub-vs-live-runtime.md) — `runtime` stub | live | harness
- [v0.3-smoke.md](v0.3-smoke.md) — reliability / rubber-stamp / Unratified publish
- [leaderboard-pages.md](leaderboard-pages.md) — Pages publish path
- [adapters/sire/README.md](../../adapters/sire/README.md) — live env vars

## Runtime warning (read first)

| `runtime` | Meaning |
| --- | --- |
| `harness` | FakeAdapter / engine self-test. **Not a product ranking.** |
| `stub` | In-process fixture shaped like a framework. **Not a product ranking.** |
| `live` | Real framework package, API, or gateway. |

Published Unratified rows must keep honest `provenance` and `runtime`. Do not
market stub or harness scorecards as framework rankings.

## Prerequisites

```bash
export GOWORK=off   # if a parent go.work lists unrelated modules
make build
./AgentGavel version
# expect exit 0
```

## FakeAdapter oracle baseline (harness)

Build the in-tree FakeAdapter and run full security + reliability suites in
oracle mode (`runtime=harness`).

```bash
go build -o /tmp/fakeadapter ./internal/engine/testdata/fakeadapter

./AgentGavel run \
  --adapter /tmp/fakeadapter \
  --suite security \
  --mode oracle \
  --seeds 25 \
  --out .claude/scratch/benchmark-baselines \
  --run-id fake-sec-oracle-baseline
echo "exit=$?"

./AgentGavel run \
  --adapter /tmp/fakeadapter \
  --suite reliability \
  --mode oracle \
  --seeds 25 \
  --out .claude/scratch/benchmark-baselines \
  --run-id fake-rel-oracle-baseline
echo "exit=$?"

./AgentGavel report --json --root .claude/scratch/benchmark-baselines fake-sec-oracle-baseline
./AgentGavel report --json --root .claude/scratch/benchmark-baselines fake-rel-oracle-baseline
```

Expect:

- Exit `0` for both runs.
- `summary.json` under
  `.claude/scratch/benchmark-baselines/results/<run-id>/`
  (gitignored via `.claude/scratch/`).
- Security summary lists SEC-001..SEC-010; reliability lists REL-001..REL-003.
- `runtime` is `harness`; `provenance` is `unofficial`.

Automated mirrors:

```bash
go test ./cmd/AgentGavel -count=1 -run 'TestRunOracleFakeAllPassWritesSummary'
go test ./cmd/AgentGavel -count=1 -run 'TestRunReliabilityFakeAdapterWritesRELRows'
```

### rubber-stamp (harness)

```bash
./AgentGavel rubber-stamp \
  --adapter /tmp/fakeadapter \
  --seeds 25 \
  --out .claude/scratch/benchmark-baselines \
  --run-id fake-rubber-stamp
echo "exit=$?"
# expect 0 on FakeAdapter
```

### Dashboard sample (harness, not a framework row)

Committed sample `dashboard/data/sample-fakeadapter-unratified.json` carries
`runtime: "harness"` and framework label `FakeAdapter (sample)`. Validate:

```bash
bash scripts/check-dashboard.sh
# expect exit 0
```

## Sire live dogfood (when credentials exist)

Default `python -m adapters.sire` stays on `StubSireClient` (`runtime=stub`).
To switch to live HTTP without code changes:

| Env | Required | Purpose |
| --- | --- | --- |
| `AGENTGAVEL_SIRE_TOKEN` (or `SIRE_API_TOKEN`) | yes | Bearer token |
| `AGENTGAVEL_SIRE_WORKER_ID` (or `SIRE_WORKER_ID`) | yes | Worker id |
| `AGENTGAVEL_SIRE_API_BASE` (or `SIRE_API_BASE`) | no | Default `https://api.sire.run/api/v1` |

Never commit tokens. Store them in gitignored scratch or a password manager.

```bash
# Point the worker model base_url at a running Compliance Oracle first.
export AGENTGAVEL_SIRE_TOKEN=...   # from 1Password / scratch — do not echo
export AGENTGAVEL_SIRE_WORKER_ID=...

# From adapters/sire after PYTHONPATH setup (see adapter README):
PYTHONPATH=src:../../sdk/python/src python -m adapters.sire
# Handshake must report runtime=live, ledger=false, observability=false,
# provenance=unofficial until external review (ADR 007).
```

Full live SEC/REL + Unratified publish steps wait on founder task T16.10
(credentials). After that:

```bash
# Example shape (exact adapter launch flag may be a path to a wrapper script):
./AgentGavel run \
  --adapter <sire-stdio-launcher> \
  --suite security \
  --mode oracle \
  --seeds 25 \
  --out .claude/scratch/benchmark-baselines \
  --run-id sire-live-sec-oracle
./AgentGavel report --publish --dashboard dashboard <run-id>
# entry must show runtime=live provenance=unofficial
```

## LangGraph live (optional PyPI package)

Default `python -m adapters.langgraph` stays on the in-process stub
(`runtime=stub`). Product-comparable runs need the optional extra:

```bash
cd adapters/langgraph
pip install -e '.[live]'
export AGENTGAVEL_LANGGRAPH_RUNTIME=live
# Oracle must be reachable from this process (local Compliance Oracle is fine).
export PYTHONPATH=src:../../sdk/python/src
```

From repo root (after `make build`):

```bash
./AgentGavel run \
  --adapter "python3 -m adapters.langgraph" \
  --suite security \
  --mode oracle \
  --seeds 25 \
  --out .claude/scratch/benchmark-baselines \
  --run-id langgraph-live-sec-oracle
./AgentGavel report --json --root .claude/scratch/benchmark-baselines langgraph-live-sec-oracle
```

Expect Handshake / summary: `runtime=live`, `provenance=unofficial`,
`framework_name=langgraph`, `framework_version` matching
`importlib.metadata.version("langgraph")`. See
[adapters/langgraph/README.md](../../adapters/langgraph/README.md).

Without `[live]` installed, `AGENTGAVEL_LANGGRAPH_RUNTIME=live` fails closed.

## Publish reminder

- Unratified: `AgentGavel report --publish` (no `--sign`).
- Opt-in: requires `--sign` (ADR 013).
- Stub/harness publishes are evidence only — not product rankings (ADR 015).
