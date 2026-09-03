# E13 -- v0.2 expansion

Acceptance: SEC-008..010 ship against FakeAdapter; governance suite scaffold
loads; unofficial adapters for all six RFC §8.1 targets (Google ADK, OpenAI
Agents SDK, Pydantic AI, Microsoft Agent Framework, AWS Strands Agents,
CrewAI) exist with READMEs; `AgentGavel run --ci` produces machine-readable
output and documented exit codes; scenario governance comment-window docs are
published.
fidelity: executable

## Learnings from v0.1 (bind into tasks)

- CapabilityReport honesty over false fails: missing hitl / ledger / tenancy →
  N/A, never silent Fail (Sire/LangGraph).
- Event ordering: `tool_invocation` before/after must stay ordered; adapters
  that cannot emit both take the observability penalty path.
- `ExportLedger` may be empty when `ledger=false`; do not invent receipts.
- Provenance stays `unofficial` until ADR 007 ratification; author-affiliated
  and third-party scaffolds alike ship unofficial.
- SEC-008 semantic canary: ADR 009 — deterministic string-variant in CI;
  optional LLM judge local-only; attestation mode → semantic N/A.
- Prefer LangGraph-style lightweight stubs over heavy framework deps in CI
  when the real package pulls a large tree; document the dependency choice in
  the adapter README (same pattern as adapters/langgraph).
- §8.1 ordering is by ResolveApproval mapping fidelity, not adoption size.
  Adapter dirs: `adk`, `openai_agents`, `pydantic_ai`, `agent_framework`,
  `strands`, `crewai`.

## Wave 1 -- suite version + SEC-008

- [x] T13.0 PLAN: expand E13 to executable fidelity (informed by v0.1 adapter and suite learnings)  Owner: pool  Est: 1h  kind: plan  delivers: [plans/E13-v02-expansion.md at fidelity: executable]  deps: [T12.8]  acc: [parse_plan sees E13 with >= 5 tasks, every task has acceptance criteria, deps resolve, fidelity flipped to executable]  completed: 2026-09-04

- [x] T13.1 Bump security suite catalog to SEC-v2 listing SEC-001..010  Owner: pool  Est: 45m  kind: agent  verifies: [UC-001, UC-018]  acc: [go test ./suites/security -run Load lists SEC-001 through SEC-010 and reports version SEC-v2]  completed: 2026-09-04
  - deps: [T13.0, T8.1]

- [x] T13.2 Fixtures for SEC-008: dual-tenant stores, exact nonce, semantic canary fact  Owner: pool  Est: 60m  kind: agent  verifies: [UC-018]  acc: [fixtures/ files exist for T1/T2 canaries; no framework-specific exploit code]  completed: 2026-09-04
  - deps: [T13.0, T8.2]

- [x] T13.3 Implement SEC-008 predicate (exact + string-variant; N/A without tenancy; Catastrophic on leak)  Owner: pool  Est: 90m  kind: agent  verifies: [UC-018]  acc: [go test ./suites/security -run SEC008 covers leak=Catastrophic, clean=100, tenancy=false→N/A; ADR 009 matcher used]  completed: 2026-09-04
  - deps: [T13.1, T13.2, T5.3]

## Wave 2 -- SEC-009 / SEC-010 + E2E

- [x] T13.4 Implement SEC-009 audit tamper resistance (append-only hash-chain verifier)  Owner: pool  Est: 90m  kind: agent  verifies: [UC-019]  acc: [go test ./suites/security -run SEC009 detects truncation/reorder; clean chain scores 100; ledger=false→N/A]  completed: 2026-09-04
  - deps: [T13.1, T5.2]

- [x] T13.5 Implement SEC-010 plan/receipt binding (substitution, arg swap, replay)  Owner: pool  Est: 90m  kind: agent  verifies: [UC-019]  acc: [go test ./suites/security -run SEC010 detects three divergence classes; missing receipts→N/A]  completed: 2026-09-04
  - deps: [T13.1, T5.2]

- [x] T13.6 End-to-end FakeAdapter oracle run for SEC-008..010  Owner: pool  Est: 60m  kind: agent  verifies: [UC-001, UC-018, UC-019]  acc: [go test ./suites/security -run SuiteOracleFakeV2 -count=1 exits 0 with SEC-008..010 rows]  completed: 2026-09-04
  - deps: [T13.3, T13.4, T13.5, T8.10]

## Wave 3 -- governance scaffold + CI mode + process docs

- [x] T13.7 Scaffold suites/governance with GOV-v0 loader and GOV-001 policy-ceiling stub  Owner: pool  Est: 75m  kind: agent  verifies: [UC-023]  acc: [go test ./suites/governance -run Load lists GOV-001; stub is N/A unless capability policy_ceiling=true]  completed: 2026-09-04
  - deps: [T13.0, T8.1]

- [x] T13.8 Add AgentGavel run --ci: non-interactive, machine-readable summary, exit-code map  Owner: pool  Est: 90m  kind: agent  verifies: [UC-020]  acc: [go test ./cmd/AgentGavel -run CIMode asserts JSON summary path and Fail→exit 1 Catastrophic→exit 2]  completed: 2026-09-04
  - deps: [T9.2, T9.3, T13.1]

- [x] T13.9 Document scenario governance comment-window process for post-v0.1 catalog changes  Owner: pool  Est: 45m  kind: agent  delivers: [docs/manual/scenario-governance.md]  verifies: [UC-024]  acc: [doc states comment window, how to propose SEC/GOV changes, and that unpublished drafts do not affect published scores]  completed: 2026-09-04
  - deps: [T13.0]

## Wave 4 -- §8.1 adapter scaffolds (6 parallel)

Retargeted 2026-09-04 after RFC §8.1 merge (PR #86 / dec-0675). AutoGen
replaced by Microsoft Agent Framework; four frameworks added; CrewAI kept.
Directory names follow the list below. Each scaffold mirrors LangGraph:
pyproject, Adapter subclass, `__main__` stdio serve, unofficial README,
provenance=unofficial.

- [x] T13.10 Scaffold adapters/adk (Google ADK; pyproject, Adapter, unofficial README)  Owner: pool  Est: 60m  kind: agent  verifies: [UC-025]  lane: agent  acc: [python -m adapters.adk starts stdio Handshake with provenance=unofficial]  completed: 2026-09-04
  - deps: [T7.3, T13.0]

- [x] T13.11 Scaffold adapters/openai_agents (OpenAI Agents SDK)  Owner: pool  Est: 60m  kind: agent  verifies: [UC-026]  lane: agent  acc: [python -m adapters.openai_agents starts stdio Handshake with provenance=unofficial]  completed: 2026-09-04
  - deps: [T7.3, T13.0]

- [x] T13.12 Scaffold adapters/pydantic_ai (Pydantic AI)  Owner: pool  Est: 60m  kind: agent  verifies: [UC-027]  lane: agent  acc: [python -m adapters.pydantic_ai starts stdio Handshake with provenance=unofficial]  completed: 2026-09-04
  - deps: [T7.3, T13.0]

- [x] T13.13 Scaffold adapters/agent_framework (Microsoft Agent Framework; AutoGen successor)  Owner: pool  Est: 60m  kind: agent  verifies: [UC-028]  lane: agent  acc: [python -m adapters.agent_framework starts stdio Handshake with provenance=unofficial]  completed: 2026-09-04
  - deps: [T7.3, T13.0]

- [x] T13.14 Scaffold adapters/strands (AWS Strands Agents)  Owner: pool  Est: 60m  kind: agent  verifies: [UC-029]  lane: agent  acc: [python -m adapters.strands starts stdio Handshake with provenance=unofficial]  completed: 2026-09-04
  - deps: [T7.3, T13.0]

- [x] T13.15 Scaffold adapters/crewai (CrewAI)  Owner: pool  Est: 60m  kind: agent  verifies: [UC-030]  lane: agent  acc: [python -m adapters.crewai starts stdio Handshake with provenance=unofficial]  completed: 2026-09-04
  - deps: [T7.3, T13.0]

## Wave 5 -- §8.1 adapter tool paths (6 parallel)

Minimal Oracle tool path + honest CapabilityReport (hitl/tenancy/ledger).
Follow T11.2 LangGraph pattern: record tool_invocation against Oracle;
capabilities must match real support (N/A path over false Fail).

- [x] T13.16 Google ADK minimal tool path + honest CapabilityReport  Owner: pool  Est: 90m  kind: agent  verifies: [UC-025, UC-003]  lane: agent  acc: [pytest adapters/adk records a tool_invocation against Oracle; capabilities match real support]  completed: 2026-09-04
  - deps: [T13.10, T11.2]

- [x] T13.17 OpenAI Agents SDK minimal tool path + honest CapabilityReport  Owner: pool  Est: 90m  kind: agent  verifies: [UC-026, UC-003]  lane: agent  acc: [pytest adapters/openai_agents records a tool_invocation against Oracle; capabilities match real support]  completed: 2026-09-04
  - deps: [T13.11, T11.2]

- [x] T13.18 Pydantic AI minimal tool path + honest CapabilityReport  Owner: pool  Est: 90m  kind: agent  verifies: [UC-027, UC-003]  lane: agent  acc: [pytest adapters/pydantic_ai records a tool_invocation against Oracle; capabilities match real support]  completed: 2026-09-04
  - deps: [T13.12, T11.2]

- [x] T13.19 Microsoft Agent Framework minimal tool path + honest CapabilityReport  Owner: pool  Est: 90m  kind: agent  verifies: [UC-028, UC-003]  lane: agent  acc: [pytest adapters/agent_framework records a tool_invocation against Oracle; capabilities match real support]  completed: 2026-09-04
  - deps: [T13.13, T11.2]

- [x] T13.20 AWS Strands Agents minimal tool path + honest CapabilityReport  Owner: pool  Est: 90m  kind: agent  verifies: [UC-029, UC-003]  lane: agent  acc: [pytest adapters/strands records a tool_invocation against Oracle; capabilities match real support]  completed: 2026-09-04
  - deps: [T13.14, T11.2]

- [x] T13.21 CrewAI minimal tool path + honest CapabilityReport  Owner: pool  Est: 90m  kind: agent  verifies: [UC-030, UC-003]  lane: agent  acc: [pytest adapters/crewai records a tool_invocation against Oracle; capabilities match real support]  completed: 2026-09-04
  - deps: [T13.15, T11.2]

## Wave 6 -- smoke + quality gate

- [x] T13.22 Add docs/manual/v0.2-smoke.md for all six §8.1 adapters + --ci  Owner: pool  Est: 45m  kind: agent  delivers: [docs/manual/v0.2-smoke.md]  verifies: [UC-020, UC-025, UC-026, UC-027, UC-028, UC-029, UC-030]  lane: agent  acc: [smoke doc lists commands for each adapter and expected unofficial provenance]  completed: 2026-09-04
  - deps: [T13.8, T13.16, T13.17, T13.18, T13.19, T13.20, T13.21]

- [x] T13.23 gofmt / ruff clean on suites/governance, new SEC files, all six adapters  Owner: pool  Est: 30m  kind: agent  verifies: [infrastructure]  lane: agent  acc: [gofmt -l suites/governance suites/security is empty; ruff check adapters/adk adapters/openai_agents adapters/pydantic_ai adapters/agent_framework adapters/strands adapters/crewai exits 0]  completed: 2026-09-04
  - deps: [T13.7, T13.6, T13.16, T13.17, T13.18, T13.19, T13.20, T13.21]

- [x] T13.24 Full make test && make lint on clean tree after E13 landings  Owner: pool  Est: 45m  kind: agent  verifies: [infrastructure]  lane: agent  acc: [make test and make lint exit 0]  completed: 2026-09-04
  - deps: [T13.23, T13.8, T13.9, T13.22]

- [ ] T13.25 Cut v0.2.0 tag after main green (founder)  Owner: pool  Est: 30m  kind: human  verifies: [infrastructure]  acc: [git tag v0.2.0 exists on origin and release assets uploaded]
  - deps: [T13.24]
  - Note: human confirms public v0.2 messaging and tag push.
