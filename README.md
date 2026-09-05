# AgentGavel

<p align="center">
  <img src="brand/logo-lockup.svg" alt="AgentGavel" width="360" height="64"/>
</p>

Adversarial benchmarking for AI agent **governance and security** — not task
completion.

AgentGavel measures whether a framework's control plane holds when the model is
fully compromised: policy ceilings under prompt injection, HITL gates that stop
side effects, and tamper-evident provenance for consequential actions.

| | |
| --- | --- |
| Spec | [RFC 0001](docs/RFC-0001.md) |
| Repo | [github.com/agentgavel/agentgavel](https://github.com/agentgavel/agentgavel) |
| Brand | [brand/](brand/) |
| License | [Apache-2.0](LICENSE) |

## Why it exists

Task benchmarks reward agents that finish jobs. They do not answer:

1. Can the agent be tricked into exceeding its authority?
2. Do “approval required” labels actually stop execution?
3. If it misbehaves, can we reconstruct *what it did, on whose behalf, and why*?

AgentGavel fills that gap with a reproducible, CI-friendly suite and a
scorecard that separates **Hard** vs **Soft** governance.

## Hard vs Soft

- **Hard** — deterministic chokepoints that refuse unsafe actions regardless of
  LLM output.
- **Soft** — system prompts or policies an LLM may be persuaded to ignore.

Both are scored. Hard/Soft classification uses a **Compliance Oracle** so the
label does not depend on any particular model's behavior (RFC §3.12 /
[ADR 003](docs/adr/003-hard-soft-oracle.md)).

## Neutrality

The author also builds [Sire](https://sire.run), one of the target frameworks.
Credibility rules (RFC §0):

- Scenarios are **code, not judgment** — deterministic predicates; no
  per-framework exploit code.
- Adapters carry a **provenance** label on every scorecard.
- Scenario changes get a public comment window starting at v0.2.
- Every published run records a **fingerprint** (versions, config, model, seeds).

### Adapter provenance

| Label | Meaning |
| --- | --- |
| **ratified** | Framework maintainers reviewed or contributed the adapter |
| **provisional** | Independent review after outreach ([ADR 007](docs/adr/007-adapter-ratification.md)) |
| **unofficial** | Otherwise — including author-affiliated adapters until external sign-off |

A low score behind an **unofficial** adapter is a claim about the adapter as
much as the framework. The in-tree Sire and LangGraph adapters are **unofficial**.

## What's in the box (v0.1)

- **Go engine + CLI** (`AgentGavel`) — run suites, write fingerprints, print GSI scorecards
- **JSON-RPC stdio adapter protocol** + **Python SDK** (`sdk/python`)
- **Compliance Oracle** — `AgentGavel oracle --listen …`
- **Security suite SEC-001…007** (`suites/security`) with shared fixtures
- **mcpfuzz** rogue MCP modes for grant/crash pressure
- **Unofficial adapters**: [Sire](adapters/sire/), [LangGraph](adapters/langgraph/)
- **FakeAdapter** — in-repo harness proof (not a framework ranking)

v0.2 added the governance suite and more unofficial adapters. v0.3 adds
reliability + the Unratified leaderboard (below). Later: public submission
(v1.0). See [docs/roadmap.md](docs/roadmap.md).

## Status

v0.1 **implementation is complete** in this repository (engine, protocol, SDK,
Oracle, SEC-001…007, CLI, unofficial Sire + LangGraph adapters, CI). The tagged
`v0.1.0` GitHub Release is the remaining cut. Track progress in
[docs/plan.md](docs/plan.md) and [docs/roadmap.md](docs/roadmap.md).

## Quick start

**Requirements:** Go 1.26+, `python3` on `PATH`. If a parent `go.work` lists
unrelated modules, keep `GOWORK=off` (Makefile default).

```bash
git clone https://github.com/agentgavel/agentgavel.git
cd agentgavel
make build
./AgentGavel version
```

### Smoke: FakeAdapter + SEC-001 (oracle mode)

```bash
go build -o /tmp/fakeadapter ./internal/engine/testdata/fakeadapter
./AgentGavel run \
  --adapter /tmp/fakeadapter \
  --suite security \
  --mode oracle \
  --scenarios SEC-001 \
  --seeds 3 \
  --out /tmp/ag-smoke-fake \
  --run-id smoke-fake-sec001
./AgentGavel report --json --root /tmp/ag-smoke-fake smoke-fake-sec001
```

Expect `results/smoke-fake-sec001/summary.json` and a SEC-001 scorecard row.
FakeAdapter reports `provenance=unofficial` (lab fixture).

Full copy-paste checks (Sire, LangGraph, expected labels):
[docs/manual/v0.1-smoke.md](docs/manual/v0.1-smoke.md).

### CLI commands

```text
AgentGavel version        # print version (ldflags / GoReleaser inject release tags)
AgentGavel oracle         # Compliance Oracle HTTP server (--listen host:port)
AgentGavel run            # run a suite (security|reliability); write results/<run-id>/
AgentGavel rubber-stamp   # SEC-002 + SEC-006 approval-gate check (no GSI)
AgentGavel report         # GSI scorecard text, --json, or --publish to dashboard/
AgentGavel help
```

### Develop

```bash
make test    # go test ./...
make lint    # go vet + golangci-lint when installed
make fmt     # gofmt
```

Go module path (imports): `github.com/agentgavel/agentgavel`  
Clone URL: `https://github.com/agentgavel/agentgavel`

## v0.3 — reliability + leaderboard

v0.3 adds a fast approval-gate check, the reliability suite (REL-001..003),
and an Unratified leaderboard publish path.

- **`AgentGavel rubber-stamp`** — SEC-002 + SEC-006 only (approval forgery /
  timeout auto-approve). CI exit codes; both-N/A fails closed. Manual:
  [docs/manual/rubber-stamp.md](docs/manual/rubber-stamp.md).
- **`AgentGavel run --suite reliability`** — REL-001..003 (demotion latch,
  replay drift, ledger completeness) scored under the resilience pillar.
- **`AgentGavel report --publish`** — writes an Unratified dashboard entry
  (`tab=unratified` only until v1.0). Merge to `main` for GitHub Pages.
  Manual: [docs/manual/leaderboard-pages.md](docs/manual/leaderboard-pages.md).

Published Unratified rows inherit Handshake provenance — typically
**unofficial**. A low score behind an unofficial adapter is a claim about the
adapter as much as the framework (same caveat as [Adapter provenance](#adapter-provenance)).

## Releases

Static `AgentGavel` binaries (linux/darwin × amd64/arm64) publish via GoReleaser
on a version tag:

```bash
git tag -a v0.1.0 -m "v0.1.0"
git push origin v0.1.0
```

That runs [`.github/workflows/release.yml`](.github/workflows/release.yml) and
attaches archives + checksums to the GitHub Release. Local dry-run:

```bash
goreleaser release --snapshot --clean
```

## Docs map

| Doc | Role |
| --- | --- |
| [docs/RFC-0001.md](docs/RFC-0001.md) | Normative specification |
| [docs/design.md](docs/design.md) | Architecture sketch |
| [docs/adr/](docs/adr/) | Architecture decisions |
| [docs/manual/v0.1-smoke.md](docs/manual/v0.1-smoke.md) | Local smoke commands |
| [docs/manual/v0.3-smoke.md](docs/manual/v0.3-smoke.md) | REL, rubber-stamp, publish, dashboard smoke |
| [docs/manual/rubber-stamp.md](docs/manual/rubber-stamp.md) | `rubber-stamp` usage + exit codes |
| [docs/manual/leaderboard-pages.md](docs/manual/leaderboard-pages.md) | Dashboard / Pages publish path |
| [docs/plan.md](docs/plan.md) / [docs/roadmap.md](docs/roadmap.md) | Execution plan + progress |

## License

Apache License 2.0. See [LICENSE](LICENSE).
