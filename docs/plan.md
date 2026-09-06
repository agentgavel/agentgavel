# AgentGavel -- implement RFC 0001 in full

## 1. Context

Problem: Agent frameworks publish little comparable data about governance and
security under adversarial pressure. Task-completion benchmarks do not answer
whether policy ceilings, HITL gates, and audit provenance hold when the model
is fully compromised.

Objectives:
- Implement `docs/RFC-0001.md` end to end across releases v0.1, v0.2, v0.3, v1.0.
- v0.1 shipped (`v0.1.0`): Go engine, protocol, Python SDK, Oracle, SEC-001..007,
  unofficial Sire + LangGraph adapters.
- v0.2 shipped (`v0.2.0`): SEC-008..010, GOV-v0 scaffold, six RFC §8.1 unofficial
  adapters, `run --ci`, scenario governance.
- v0.3 shipped (`v0.3.0`): REL-v0, rubber-stamp (ADR 011), live leaderboard on
  `agentgavel.dev` (ADR 006), Pages + dashboard CI.
- Frontier is now v1.0 (E15 executable after T15.0): GitHub-native signed Opt-in
  (ADR 012 + ADR 013), harness red-team bounty, ratification ops (ADR 007),
  and unofficial gateway-style adapters for OpenClaw and Hermes Agent
  (ADR 014 / RFC §8.3).

Non goals (v1.0):
- Per-framework exploit code (forbidden by RFC section 0).
- Firebase / BaaS as Opt-in trust root (rejected; ADR 012).
- Hosted submission API beyond signed PRs into `dashboard/data/`.
- n8n / Dify gateway adapters (still deferred under RFC §8.2).
- Nested Codex / Claude Agent SDK / Copilot plugins as the OpenClaw SUT;
  non-reference Hermes terminal backends as the Hermes SUT.

Constraints and assumptions:
- Module/repo: `github.com/agentgavel/agentgavel`.
- Neutrality rules in RFC section 0 are binding.
- Open questions Q3-Q7 resolved in ADRs 002, 004, 005, 006, 007; SEC-008
  semantic judge in ADR 009; REL-v0 in ADR 010; rubber-stamp in ADR 011;
  GitHub-native Opt-in in ADR 012; Ed25519 signature format in ADR 013;
  gateway-style OpenClaw adapter in ADR 014.
- No hosted production service beyond GitHub Pages dashboard; release "done"
  means tagged GitHub release + CI green + documented local smoke + live
  leaderboard URL.

Success metrics:
- Signed Opt-in PR path verifies in CI; `sample` no longer required when
  signature verifies (ADR 013).
- `SECURITY.md` + harness bounty scope published.
- At least one non-author adapter reaches provisional or ratified (ADR 007).
- Unofficial OpenClaw and Hermes Agent adapters each have Handshake +
  SEC-002 oracle path (or honest N/A) documented (ADR 014).
- Soft rates use >=25 seeds with Wilson intervals.
- Scenario governance comment window applies before REL/SEC/GOV changes publish.

## 2. Discovery Summary

Work type: Engineering (greenfield + process docs).

Graph scan: no `.code-review-graph/graph.db`; skipped.

Use cases: 32 total (UC-032..036 including OpenClaw + Hermes). Manifest:
`.claude/scratch/usecases-manifest.json`.

Gaps to close in v1.0 (re-scanned 2026-09-05 after Hermes expand):
- Maintainer key registry + Ed25519 verify (ADR 013) -- T15.1-3 done.
- CLI CI verify job; flip `check-dashboard.sh`; signed `--tab opt-in`.
- Opt-in submission manual; Pages/README; bounty + SECURITY.md.
- Ratification ops + first non-author provisional (human gate).
- OpenClaw capability map + unofficial sidecar + oracle E2E (ADR 014).
- Hermes Agent capability map + unofficial sidecar + oracle E2E (ADR 014).
- v1.0 smoke + quality gate + `v1.0.0` tag (waits on OpenClaw + Hermes E2E).

Research notes (v0.3 ops + OpenClaw):
- Trust for Opt-in is cryptographic, not IdP login; volume is low.
- `report --publish` remains the writer of `index.json`; signing wraps the
  same entry schema.
- Author-affiliated Sire cannot skip to ratified via provisional alone.
- OpenClaw Gateway is operated product surface; wire protocol reused inside
  a gateway-style sidecar (not a new engine protocol).

## 3. Scope and Deliverables

In scope:
- Full RFC release plan (v0.1–v0.3 complete; v1.0 executable).
- ADRs through 014 (including OpenClaw gateway-style).
- Design doc and roadmap sync.

Out of scope:
- Changing Sire product code inside `sirerun/sire` except via the adapter.
- Firebase / custom hosted verify API as system of record (ADR 012).
- Expanding REL/SEC scenario catalogs (separate governance PRs).
- n8n / Dify adapters (§8.2 still deferred).

| ID | Deliverable | Acceptance |
| ---- | ---- | ---- |
| D1–D7 | v0.1 harness + release | tag `v0.1.0` shipped |
| D8 | SEC-008..010 (SEC-v2) | FakeAdapter oracle green |
| D9 | Governance suite scaffold | GOV-v0 + GOV-001 stub |
| D10 | Unofficial §8.1 adapters (×6) | provenance label |
| D11 | `run --ci` + scenario governance docs | exit codes + comment window |
| D12 | v0.2.0 GitHub release | tag + binaries |
| D13 | REL-001..003 (REL-v0) | FakeAdapter oracle green |
| D14 | `rubber-stamp` CLI | SEC-002 + SEC-006; exit 0/1/2; both-N/A→1 (ADR 011) |
| D15 | Leaderboard dashboard | Opt-in + Unratified tabs; live on Pages |
| D16 | v0.3.0 GitHub release | tag + binaries |
| D17 | Signed Opt-in + bounty + ratification | ADR 012/013; SECURITY.md; provisional badge |
| D18 | Unofficial OpenClaw gateway adapter | ADR 014; Handshake + SEC-002/N/A path |
| D19 | Unofficial Hermes Agent gateway adapter | ADR 014; Handshake + SEC-002/N/A path |
| D20 | v1.0.0 GitHub release | tag + binaries |

## 4. Checkable Work Breakdown

Split layout. E1–E14 complete (`fidelity: executable`, all tasks done). E15 is
`fidelity: executable` (v1.0 frontier).

### E1 -- Repository bootstrap  -> docs/plans/E1-repo-bootstrap.md  (6/6)

### E2 -- Adapter wire protocol  -> docs/plans/E2-adapter-protocol.md  (8/8)

### E3 -- Engine core  -> docs/plans/E3-engine-core.md  (7/7)

### E4 -- Compliance Oracle  -> docs/plans/E4-compliance-oracle.md  (6/6)

### E5 -- Assertions and GSI metrics  -> docs/plans/E5-assertions-metrics.md  (9/9)

### E6 -- mcpfuzz rogue MCP servers  -> docs/plans/E6-mcpfuzz.md  (9/9)

### E7 -- Python adapter SDK  -> docs/plans/E7-python-sdk.md  (6/6)

### E8 -- Security suite SEC-001 through SEC-007  -> docs/plans/E8-security-suite.md  (11/11)

### E9 -- CLI run, report, fingerprint  -> docs/plans/E9-cli-report.md  (6/6)

### E10 -- Sire adapter (unofficial)  -> docs/plans/E10-sire-adapter.md  (7/7)

### E11 -- LangGraph adapter (unofficial)  -> docs/plans/E11-langgraph-adapter.md  (7/7)

### E12 -- v0.1 quality gate and release  -> docs/plans/E12-v01-release.md  (8/8)

### E13 -- v0.2 expansion  -> docs/plans/E13-v02-expansion.md  (26/26)

### E14 -- v0.3 reliability, rubber-stamp, leaderboard  -> docs/plans/E14-v03-reliability-leaderboard.md  (24/24)

### E15 -- v1.0 Opt-in, red-team, OpenClaw, Hermes  -> docs/plans/E15-v10-public-process.md  (8/31)

## 5. Parallel Work

Tracks (v1.0):
- Track V: signature registry + verify library + CLI (T15.1–T15.3) -- done
- Track W: CI + Opt-in rule flip + signed publish (T15.4–T15.6)
- Track X: submission docs + bounty + README (T15.7–T15.9)
- Track Y: ratification ops + first provisional (T15.10–T15.11)
- Track AA: OpenClaw gateway-style adapter (T15.15–T15.22)
- Track AB: Hermes Agent gateway-style adapter (T15.30, T15.23–T15.29)
- Track Z: smoke + quality + tag (T15.12–T15.14)

Sync points: T15.2 before T15.4/T15.5; T15.5+T15.6 before T15.13;
T15.10 before T15.11; T15.15 before T15.16; T15.30 before T15.23;
T15.20+T15.21 and T15.27+T15.28 before T15.14; T15.13 before T15.14.
Tracks AA/AB may run in parallel with W/X after their design tasks.

### Wave 13–27: E13–E14 (done)
- T13.0–T13.25, T14.0–T14.23

### Wave 28: E15 planning (done)
- T15.0

### Wave 29: signature contract + CLI (done)
- T15.1, T15.2, T15.3

### Wave 30: CI + Opt-in rule flip (3 agents)
- T15.4, T15.5, T15.6

### Wave 31: docs + bounty (3 agents)
- T15.7, T15.8, T15.9

### Wave 32: ratification ops (1 agent + founder)
- T15.10, T15.11

### Wave 34: OpenClaw design (done)
- T15.15 (done), T15.16 (done)

### Wave 35: OpenClaw adapter (3 agents)
- T15.17, T15.18, T15.19

### Wave 36: OpenClaw E2E + docs (3 agents)
- T15.20, T15.21, T15.22

### Wave 37: Hermes design (done)
- T15.30 (done), T15.23 (done)

### Wave 38: Hermes adapter (3 agents)
- T15.24, T15.25, T15.26

### Wave 39: Hermes E2E + docs (3 agents)
- T15.27, T15.28, T15.29

### Wave 33: quality gate + ship (2 agents + founder)
- T15.12, T15.13, T15.14

## Roadmap
- **Now:** Wave 29 done; T15.16/T15.23 capability maps done; E15 8/31
- **Next:** Wave 30 (Opt-in CI) in parallel with Waves 35-39 adapter scaffolds -> 33 -> `v1.0.0`

## 6. Timeline and Milestones

| ID | Milestone | Depends | Exit criteria |
| ---- | ---- | ---- | ---- |
| M1 | Protocol + Oracle usable | Waves 1-4 | T2.4, T4.4 green |
| M2 | Engine + SDK cross-talk | Waves 5-6 | T7.5 green |
| M3 | Security suite SEC-001..007 | Waves 7-8 | T8.10 green |
| M4 | Unofficial adapters | Waves 9-10 | T10.5, T11.5 green |
| M5 | v0.1.0 release | Waves 11-12 | T12.8 done |
| M6 | v0.2 planned | T13.0 | E13 executable |
| M7 | v0.2.0 release | Waves 14-19 | T13.25 done |
| M8 | v0.3 planned | T14.0 | E14 executable |
| M9 | v0.3.0 release | Waves 21-27 | T14.20 live, T14.16 done |
| M10 | v1.0 planned | T15.0 | E15 executable |
| M11 | v1.0.0 release | Waves 29-39 | T15.13 green, T15.20+T15.27 gateway E2E, T15.14 tagged |

## 7. Risk Register

| ID | Risk | Impact | Likelihood | Mitigation |
| ---- | ---- | ---- | ---- | ---- |
| R1 | Author bias (Sire) | Credibility | Med | Unofficial label; external ratification ADR 007 |
| R2 | Soft runs expensive | Slow CI | Med | Oracle-first in CI; model mode opt-in job |
| R3 | Sire/LangGraph API mismatch | Adapter N/A heavy | High | Honest CapabilityReport; FakeAdapter proves harness |
| R10 | Pages spam without signatures | Credibility | Low | ADR 013 Ed25519 + CI verify; samples remain labeled |
| R11 | `rubber-stamp` vacuous green | False assurance | Med | ADR 011: both-N/A exits 1 |
| R14 | Canonical JSON mismatch across signers | Broken Opt-in | Med | Golden vectors in `internal/submit`; ADR 013 pins encoding |
| R15 | Provisional confused with ratified | Credibility | Med | Three-way badge + ratification manual (T15.10) |
| R16 | Bounty scope includes adapter exploits | Neutrality breach | Med | Bounty doc out-of-scope list; RFC section 0 |
| R17 | OpenClaw/Hermes Gateway APIs lack ResolveApproval | Heavy N/A | High | Capability maps first (T15.16, T15.23); honest hitl=false; ADR 011 |
| R18 | Gateway CI needs live process | Flaky / heavy | Med | Document reference config; optional live E2E; FakeAdapter remains harness proof |
| R19 | Confusing OpenClaw vs Hermes scorecards | Credibility | Med | Separate adapter dirs, fingerprints, and leaderboard rows (ADR 014) |

## 8. Operating Procedure

Definition of done for a task:
1. Acceptance `acc:` predicate is green (tests or commands cited).
2. Paired tests exist for new packages/CLI/scenario behavior.
3. `gofmt` / `ruff` clean on touched trees; `go test` / `pytest` green.
4. PR merged to main via rebase; CI green.
5. For release tasks: GitHub release assets verified. Leaderboard live URL
   remains `https://agentgavel.dev/leaderboard/`.

Rules:
- Work in a git worktree once the repo has a first commit.
- Prefer stdlib Go (`flag`, `testing`, `net/http`, `crypto/ed25519`).
- Never add per-framework exploit code; probes stay in `fixtures/`.
- Small focused commits; do not mix Go and Python adapter dirs in one commit
  when hooks forbid it.
- Opt-in submissions are GitHub-native (ADR 012); signatures follow ADR 013.
- `report --publish` remains the sole writer of `dashboard/data/index.json`;
  signed Opt-in may set `tab=opt-in` only when verify passes.
- After M10, execute Waves 29-39 before cutting `v1.0.0` (T15.14); tag waits
  on OpenClaw E2E (T15.20), Hermes E2E (T15.27), and Opt-in/bounty docs.
- Adapter dirs for §8.1: `adk`, `openai_agents`, `pydantic_ai`,
  `agent_framework`, `strands`, `crewai` (never `autogen`).
- Gateway-style: `openclaw`, `hermes` per ADR 014 / RFC §8.3 (n8n/Dify still
  deferred).
- REL IDs and predicates: ADR 010 only (REL-001..003 / REL-v0).

## 9. Progress Log

- 2026-09-06: T15.16 OpenClaw + T15.23 Hermes capability maps; E15 8/31.
- 2026-09-05: Expand E15 for Hermes Agent (ADR 014 rename, RFC §8.3, T15.23-T15.30, UC-036); T15.30 done.
- 2026-09-05: Expand E15 for OpenClaw (ADR 014, RFC §8.3, T15.15-T15.22, UC-035); T15.15 done.
- 2026-09-05: Wave 29 T15.1 key registry, T15.2 internal/submit, T15.3 report --sign + verify-entry.
- 2026-09-05: T15.0 expanded E15 to executable (15 tasks, Waves 29-33); ADR 013 Opt-in signature format; UC-032..034.
- 2026-09-05: Wave 27 T14.15/T14.20/T14.16; E14 complete; `v0.3.0` released.
- 2026-09-05: Wave 26 T14.21/T14.13/T14.14. E14 21/24.
- 2026-09-05: Wave 25 T14.10/T14.11/T14.12/T14.22. E14 18/24.
- 2026-09-04: E13 complete; `v0.2.0` released. Per-wave detail in `docs/roadmap.md`.

## 10. Hand off Notes

- Spec: `docs/RFC-0001.md`. Design: `docs/design.md`.
- ADRs: 006/012/013 (leaderboard + Opt-in), 007 (ratification), 010 (REL),
  011 (rubber-stamp), 014 (OpenClaw + Hermes gateway-style).
- Start apply at Wave 30 (T15.4-T15.6) and/or Wave 35/38 adapter scaffolds (T15.17 / T15.24).
- Founder/human gates: T15.11 (provisional sign-off), T15.14 (tag `v1.0.0`).
- kazi is on PATH; engineering tasks carry `acc:` for JIT lane.
- Claim resource for plan rewrites: `R-plan-md`.
- Remote: `git@github.com:agentgavel/agentgavel.git`.

## 11. Appendix

- RFC open questions: Q3->ADR 005; Q4->ADR 006 + 012 + 013; Q5->ADR 004;
  Q6->ADR 007; Q7->ADR 002; REL->ADR 010; rubber-stamp->ADR 011;
  gateway-style OpenClaw+Hermes->ADR 014.
- Use case manifest: `.claude/scratch/usecases-manifest.json`.
- Release map: v0.1=E1-E12; v0.2=E13; v0.3=E14; v1.0=E15.
