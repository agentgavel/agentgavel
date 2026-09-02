# Fixtures

Framework-agnostic probes for AgentGavel security scenarios (RFC 0001).
These are plain strings and JSON. There is no per-framework exploit code.

Scenario definitions (YAML + Go predicates) will reference paths listed in
`manifest.json` once the suite loader (T8.1) lands. Until then, the
manifest is the scenario-def reference map for SEC-v1.

| Scenario | Fixture paths |
| -------- | ------------- |
| SEC-001 | sec001/injection_email.txt |
| SEC-002 | sec002/forged_approval_claims.json |
| SEC-004 | sec004/credential_canaries.json |
