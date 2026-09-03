# Fixtures

Framework-agnostic probes for AgentGavel security scenarios (RFC 0001).
These are plain strings and JSON. There is no per-framework exploit code.

Scenario definitions reference paths listed in `manifest.json` (SEC-v2).

| Scenario | Fixture paths |
| -------- | ------------- |
| SEC-001 | sec001/injection_email.txt |
| SEC-002 | sec002/forged_approval_claims.json |
| SEC-004 | sec004/credential_canaries.json |
| SEC-008 | sec008/tenant_canaries.json |

Reliability scenario definitions reference paths listed in `manifest.json`
under `rel_scenarios` (REL-v0, ADR 010).

| Scenario | Fixture paths |
| -------- | ------------- |
| REL-001 | rel001/demotion_probe.json |
| REL-002 | rel002/drifted_replay_plan.json |
| REL-003 | rel003/ledger_gap_corpus.json |
