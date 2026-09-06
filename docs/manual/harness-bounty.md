# Harness Red-Team Bounty

AgentGavel invites researchers to find security and integrity defects in
the **harness** — the engine that runs scenarios, oracles, scores
results, and publishes Opt-in leaderboard entries. The goal is to keep
published scores trustworthy under adversarial pressure.

This program is **not** a bounty for breaking third-party agent
frameworks. Per-framework exploit kits remain out of scope (see
[CONTRIBUTING.md](../../CONTRIBUTING.md) and RFC §0).

Disclosure channel: [SECURITY.md](../../SECURITY.md) → GitHub private
vulnerability reporting /
[Security Advisories](https://github.com/agentgavel/agentgavel/security/advisories/new).

## In scope

Findings that undermine harness integrity or published trust, including:

| Area | Examples |
| --- | --- |
| **Engine** | Sidecar/session isolation breaks; result-store tampering that changes verdicts without detection; runaway or unsafe execution paths owned by the harness |
| **Oracle** | Hard/Soft misclassification that can be forced without changing adapter behavior; oracle bypass that marks Soft as Hard (or the reverse) for scoring advantage |
| **Scoring** | GSI / Hard-Soft aggregation bugs that inflate or deflate scores incorrectly; fingerprint or provenance fields that can be forged through the harness |
| **Dashboard verify** | Accepting `tab=opt-in` `sample=false` entries without a valid signature against an `active` registry key (ADR 012 / ADR 013); schema or check-dashboard gaps that let unsigned Opt-in rows publish |
| **CI signature checks** | Opt-in verify jobs or scripts that pass on tampered payloads, revoked keys, or mismatched `framework` / `key_id` bindings |

Reports should include a minimal reproduction against this repository
(commit SHA or tag) and explain how a published scorecard or Opt-in entry
could be wrong or untrusted as a result.

## Out of scope

| Category | Why |
| --- | --- |
| **Per-framework exploit kits** | AgentGavel scenarios use framework-agnostic fixtures and deterministic predicates. Building or submitting exploit code aimed at a specific product (LangGraph, OpenClaw, Hermes, Sire, etc.) is not in scope for this bounty. |
| **Social engineering** | Phishing maintainers, spoofing GitHub identities, coercion, or physical attacks are out of scope. |
| **Third-party framework bugs** | Vulnerabilities in frameworks under test belong to those projects' disclosure processes, not this bounty. |
| **Denial of service against GitHub / Pages** | Abuse of GitHub infrastructure, rate limits, or public Pages hosting is not rewarded. |
| **Issues that require maintainer secret keys** | Stealing or guessing a registered maintainer's private Ed25519 key outside the harness is out of scope; forging signatures *without* such keys (canonicalization / verify bugs) is in scope. |

## Safe harbor

If you make a good-faith effort to comply with this policy:

1. We will not pursue civil or criminal claims against you for research
   conducted under this bounty.
2. We will treat your report as made in good faith if you avoid privacy
   violations, destruction of data, and disruption of production services
   you do not own.
3. Test only against **local checkouts**, **your own forks**, and
   **pull-request / CI paths you initiate**. Do not attack other
   researchers' or maintainers' accounts, or unrelated infrastructure.
4. Stop testing and report immediately if you encounter data that is not
   yours (secrets, private keys, personal data).

Safe harbor does **not** cover social engineering, illegal access to
third-party systems, or publishing exploit details before coordinated
disclosure closes.

## Coordinated disclosure

1. **Report privately** via
   [GitHub Security Advisories](https://github.com/agentgavel/agentgavel/security/advisories/new)
   (see [SECURITY.md](../../SECURITY.md)). Do not file a public issue for
   in-scope security findings.
2. **Acknowledge** — we aim to reply within 3 business days.
3. **Triage** — we aim to confirm or decline within 10 business days.
4. **Fix** — maintainers develop and validate a fix privately when the
   report is accepted.
5. **Disclose** — after a fix (or agreed mitigation) ships, we publish a
   GitHub Security Advisory (and credit you if you want it). We ask
   reporters to wait until that advisory is public before sharing PoCs.

Non-security integrity bugs that are clearly public (documentation typos,
obvious Soft-rate UI copy) may use ordinary issues or PRs. When in doubt,
use private reporting.

## Related docs

- [SECURITY.md](../../SECURITY.md) — disclosure path
- [ADR 012](../adr/012-github-native-signed-submissions.md) — GitHub-native Opt-in
- [ADR 013](../adr/013-opt-in-signature-format.md) — Ed25519 signature contract
- [ADR 003](../adr/003-hard-soft-oracle.md) — Compliance Oracle
- [ADR 004](../adr/004-gsi-scoring.md) — scoring
