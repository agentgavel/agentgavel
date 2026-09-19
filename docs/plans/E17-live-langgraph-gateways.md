# E17 -- Live LangGraph and gateway product runs

Acceptance: At least one non-Sire product adapter runs with runtime=live
against real LangGraph (or a live OpenClaw/Hermes gateway) and can publish
Unratified with honest CapabilityReport.
fidelity: executable

Intent: After E16 dogfoods Sire live + ADR 015 disclosure, promote the
LangGraph sidecar from stub to real `langgraph` interrupt/resume (preferred
for ranking credibility after #8992), and optionally live OpenClaw/Hermes
when a reference gateway is operated. Soft model-mode campaigns stay out
of this epic.

Exit criteria: one live non-Sire Unratified row with runtime=live; README
runtime section updated; stub path remains available for CI.

Learnings from E16 (inform T17.*):
- Env-gated live bootstrap with stub default (Sire client_from_env pattern).
- `runtime` must flip only on the proven live path; framework_version must
  come from the real package or probe, never a stub string while live.
- Local Oracle works for adapters that call tools in-process; cloud control
  planes need a publicly reachable Oracle URL (Sire Wave 43).
- Publish remains Unratified + unofficial until ADR 007 / T15.11.

## Wave 45 -- Plan + live LangGraph bootstrap

- [x] T17.0 PLAN: expand E17 to executable fidelity (informed by E16 live Sire learnings)  Owner: pool  Est: 1h  kind: plan  delivers: [plans/E17-live-langgraph-gateways.md at fidelity: executable]  deps: [T16.17]  acc: [parse_plan sees E17 with >= 5 tasks, every task has acceptance criteria, deps resolve, fidelity flipped to executable]  completed: 2026-09-19

- [x] T17.1 Optional `langgraph` extra + import probe; stub default unchanged  Owner: pool  Est: 60m  kind: agent  verifies: [UC-040]  deps: [T17.0]  acc: [pyproject optional-dependencies live includes langgraph; without extra or env, Handshake runtime=stub and CI pytest green without langgraph installed]  completed: 2026-09-19

- [x] T17.2 Live email graph using real langgraph StateGraph + interrupt/resume  Owner: pool  Est: 120m  kind: agent  verifies: [UC-040]  deps: [T17.1]  acc: [when live path selected and langgraph importable, gated send_email pauses via langgraph interrupt and resumes on ResolveApproval; unit test with MemorySaver covers approve and deny]  completed: 2026-09-19

- [x] T17.3 Env/bootstrap runtime=live + honest framework_version from package metadata  Owner: pool  Est: 60m  kind: agent  verifies: [UC-037, UC-040]  deps: [T17.2]  acc: [AGENTGAVEL_LANGGRAPH_RUNTIME=live (or adapter_from_env) sets Handshake runtime=live and framework_version from importlib.metadata.version('langgraph'); missing package fails closed with clear error, never silent stub-as-live]  completed: 2026-09-19

- [x] T17.4 Lint/format LangGraph live bootstrap  Owner: pool  Est: 20m  kind: agent  verifies: [infrastructure]  deps: [T17.1, T17.2, T17.3]  acc: [ruff check adapters/langgraph exits 0]  completed: 2026-09-19

## Wave 46 -- Live LangGraph oracle evidence

- [x] T17.5 Document LangGraph live ops in benchmark-ops + adapter README  Owner: pool  Est: 45m  kind: agent  verifies: [UC-040]  lane: agent  deps: [T17.3]  acc: [benchmark-ops.md has copy-paste install live extra, Oracle, AgentGavel run --adapter langgraph; README states stub vs live]  completed: 2026-09-19

- [x] T17.6 Run live LangGraph oracle SEC-001 + SEC-007 (local Oracle) and record summary  Owner: pool  Est: 90m  kind: agent  verifies: [UC-040, UC-001]  deps: [T17.3, T17.5]  acc: [summary.json has runtime=live provenance=unofficial; SEC-001 and SEC-007 rows present; path cited in benchmark-ops or scratch pointer]  completed: 2026-09-19

- [x] T17.7 Publish Unratified LangGraph live scorecard sample (or document PR path)  Owner: pool  Est: 60m  kind: agent  verifies: [UC-038, UC-040]  deps: [T17.6, T16.2]  acc: [dashboard data entry or documented report --publish PR has runtime=live provenance=unofficial framework_name=langgraph; check-dashboard.sh exits 0]  completed: 2026-09-19

- [x] T17.8 Append Tier-3 live LangGraph notes to docs/devlog.md  Owner: pool  Est: 20m  kind: agent  verifies: [UC-040]  lane: agent  deps: [T17.6]  delivers: [devlog entry for first live LangGraph oracle runs]  completed: 2026-09-19

## Wave 47 -- Gateway live probes (optional if no operated gateway)

- [x] T17.9 OpenClaw live probe scaffolding: env-gated Gateway client; hitl stays false until withhold map exists  Owner: pool  Est: 90m  kind: agent  verifies: [UC-041]  deps: [T17.0]  acc: [without AGENTGAVEL_OPENCLAW_GATEWAY_URL, stub/runtime=stub; with URL + health fail, fail closed; CapabilityReport never claims hitl=true until withhold mapped]  completed: 2026-09-19

- [x] T17.10 Hermes live probe scaffolding mirroring T17.9 honesty rules  Owner: pool  Est: 75m  kind: agent  verifies: [UC-041]  deps: [T17.9]  acc: [Hermes Handshake runtime=live only after successful probe; otherwise stub; pytest covers both]  completed: 2026-09-19

- [ ] T17.11 Human: operate or point at a reference OpenClaw/Hermes gateway for one live run  Owner: founder  Est: 60m  kind: human  verifies: [UC-041]  deps: [T17.9]  delivers: [scratch pointer to gateway URL; never commit secrets]  blocked: needs operated gateway

## Wave 48 -- E17 quality gate

- [x] T17.12 make test + make lint + check-dashboard green on E17 branch  Owner: pool  Est: 45m  kind: agent  verifies: [infrastructure]  deps: [T17.4, T17.7]  acc: [make test && make lint && bash scripts/check-dashboard.sh all exit 0]  completed: 2026-09-19

- [x] T17.13 Update README Quick Start / runtime section for live LangGraph  Owner: pool  Est: 30m  kind: agent  verifies: [UC-040]  lane: agent  deps: [T17.5, T17.7]  acc: [README links LangGraph live ops and ADR 015; stub default still documented]  completed: 2026-09-19

- [x] T17.14 Mark E17 complete in plan.md when exit criteria met  Owner: pool  Est: 15m  kind: agent  verifies: [infrastructure]  deps: [T17.12, T17.13]  acc: [plan.md frontier advances; E17 checkbox count reflects done; Progress Log notes first non-Sire live Unratified]  completed: 2026-09-19
