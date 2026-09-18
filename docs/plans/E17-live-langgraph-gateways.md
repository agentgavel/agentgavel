# E17 -- Live LangGraph and gateway product runs

Acceptance: At least one non-Sire product adapter runs with runtime=live
against real LangGraph (or a live OpenClaw/Hermes gateway) and can publish
Unratified with honest CapabilityReport.
fidelity: outline

Intent: After E16 dogfoods Sire live + ADR 015 disclosure, promote the
LangGraph sidecar from stub to real `langgraph` interrupt/resume (preferred
for ranking credibility after #8992), and optionally live OpenClaw/Hermes
when a reference gateway is operated. Soft model-mode campaigns stay out
of this epic.

Exit criteria: one live non-Sire Unratified row with runtime=live; README
runtime section updated; stub path remains available for CI.

- [ ] T17.0 PLAN: expand E17 to executable fidelity (informed by E16 live Sire learnings)  Owner: pool  Est: 1h  kind: plan  delivers: [plans/E17-live-langgraph-gateways.md at fidelity: executable]  deps: [T16.17]  acc: [parse_plan sees E17 with >= 5 tasks, every task has acceptance criteria, deps resolve, fidelity flipped to executable]
