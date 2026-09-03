## 2026-09-03 -- v0.1 plan execution (stub)

Stub for the RFC-0001 / `docs/plan.md` v0.1 execution pass across epics E1–E12.

- **Start:** _YYYY-MM-DD — fill when Wave execution formally opens (or note already underway)._
- **Finish:** _YYYY-MM-DD — fill when v0.1 acceptance criteria are met; do not claim finish yet._

Scope reminder: protocol, oracle, engine, suites, CLI, CI, and docs/infra tasks tracked under E1–E12. Update this section in place when start/finish dates are known; leave historical Wave entries below untouched.

## 2026-09-03 -- Wave 4 protocol session, metrics, engine scaffold

Session Client lifecycle + version reject, tool_invocation order checks, attestation leak scanner, GSI/grades/catastrophic caps, Wilson intervals, engine Scenario interfaces, seed scheduler (25 seeds/4 workers), run artifact writer.

# Devlog

## 2026-09-03 -- Wave 3 protocol codec, oracle CLI, assertions

Landed internal/protocol types+JSON codec+stdio Handshake, CapabilityReport N/A
helpers, AgentGavel oracle --listen with /healthz, and assertions for tool
invocation, gate genuineness, credential leak encodings, and partial effects.

## 2026-09-03 -- Wave 2 Makefile, proto, Oracle handlers, mcpfuzz modes


Landed Makefile (GOWORK=off), GitHub Actions CI, proto/adapter.proto (seven
RPCs), OpenAI and Anthropic Compliance Oracle handlers with MissingDirective
4xx, and mcpfuzz modes toxic-output through masquerade.

## 2026-09-02 -- Wave 1 bootstrap (T1.1 T1.2 T1.5 T8.2)


Landed Go module + AgentGavel stub CLI, Apache-2.0 LICENSE/README/CONTRIBUTING,
root .gitignore (scratch/worktrees/venvs/binaries), and SEC-001/002/004
fixtures under fixtures/. Parent directory go.work breaks builds unless
GOWORK=off. Sonnet subagent dispatch failed on usage limits; coordinator
completed Wave 1 directly on wave-1-integration.

## 2026-09-02 -- RFC 0001 implementation plan seeded

Created the initial engineering plan to implement AgentGavel per
`docs/RFC-0001.md`. Repo was documentation-only at this time (no Go module yet).

Open questions Q3-Q7 were resolved as planning defaults in ADRs 002 and 004-007.
Execution should start at Wave 1 in `docs/plan.md` after an initial commit
enables worktrees.
