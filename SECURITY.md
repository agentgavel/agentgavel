# Security Policy

AgentGavel is an adversarial benchmarking harness for AI agent governance
and security. This policy covers vulnerabilities in **AgentGavel itself**
(the harness), not flaws in third-party frameworks under test.

For bounty scope, safe harbor, and what is in / out of scope, see
[docs/manual/harness-bounty.md](docs/manual/harness-bounty.md).

## Supported versions

Security fixes target the latest release on `main` and the most recent
tagged release. Older tags may not receive backports.

## Reporting a vulnerability

Please **do not** open a public GitHub issue, discussion, or pull request
for security-sensitive findings.

**Preferred path — coordinated disclosure via GitHub Security Advisories:**

1. Open a private report at
   [github.com/agentgavel/agentgavel/security/advisories/new](https://github.com/agentgavel/agentgavel/security/advisories/new)
   (Security → Report a vulnerability).
2. Include: affected version or commit SHA, description, reproduction
   steps or proof of concept, and expected impact.
3. We aim to acknowledge within **3 business days** and share a triage
   decision within **10 business days**.

We follow **coordinated disclosure**: fixes are developed privately,
validated, released, and then disclosed publicly (typically via a GitHub
Security Advisory) with credit to the reporter when they want it.

## Harness red-team bounty

AgentGavel invites adversarial review of the harness. Scope, out-of-scope
rules, and safe-harbor terms are in
[docs/manual/harness-bounty.md](docs/manual/harness-bounty.md).
