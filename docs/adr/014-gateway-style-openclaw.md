# ADR 014: Gateway-Style Adapters (OpenClaw First)

## Status
Accepted

## Date
2026-09-05

## Context
RFC §8.2 deferred workflow/app-builder platforms (n8n, Dify) because they
are not code-first libraries in the shape of the §5.3 sidecar contract.
OpenClaw is a self-hosted **Gateway product**: channels, sessions, policy,
credentials, and exec approvals live in a trusted control plane, with
movable execution in sandboxes/nodes. It has strong governance surface
(exec approvals, tool policy, audit) that AgentGavel is meant to measure,
but it is not a LangGraph-style in-process SDK.

v1.0 previously scoped only signed Opt-in, bounty, and ratification.
Founder direction: expand v1.0 to include OpenClaw as a benchmark target.

## Decision

1. **v1.0 includes an unofficial OpenClaw adapter** under `adapters/openclaw/`,
   provenance `unofficial` until ADR 007 ratification.
2. **Architecture class: gateway-style.** The sidecar still speaks the
   existing AgentGavel wire protocol (Handshake, StartSession, SubmitTask,
   ResolveApproval, ExportLedger, Events). Inside the sidecar, the adapter
   drives OpenClaw via its Gateway/control API (or documented CLI/RPC),
   not by importing an agent graph library.
3. **Under test** is OpenClaw's **bundled agent runtime + Gateway policy
   plane** in a documented reference config (default personal posture vs
   enterprise/sandbox posture are separate fingerprints if both are run).
   Vendor harness plugins (Codex / Claude Agent SDK / Copilot) are out of
   scope for the first OpenClaw adapter unless CapabilityReport can name
   them honestly as the execution backend.
4. **Honest N/A.** If a Gateway API cannot expose programmatic
   `ResolveApproval` or ledger export, CapabilityReport sets `hitl` /
   `ledger` false and SEC/REL scenarios N/A with the observability penalty
   -- never silent Fail or stubbed green.
5. **§8.2 remains deferred for n8n/Dify.** OpenClaw is the first accepted
   gateway-style target; other platforms still need their own design notes
   before tasks are cut.
6. **No per-framework exploit code.** Probes stay in `fixtures/`; the
   adapter only maps protocol methods to Gateway surfaces.

## Consequences
Positive: v1.0 measures a high-relevance operated agent product; proves the
wire contract can wrap a Gateway without inventing a second engine.
Negative: adapter complexity and CI cost (Gateway process under test);
more N/A risk until Gateway APIs are mapped; release `v1.0.0` waits on the
OpenClaw E2E gate (T15.20) in addition to Opt-in/bounty work.

## Alternatives considered
- **Defer OpenClaw to post-v1.0** -- rejected by founder scope expand.
- **New engine protocol for gateways** -- premature; reuse §5.3 first.
- **Benchmark only via Claude/Codex plugins** -- wrong system under test;
  those are nested harnesses, not OpenClaw's policy plane.
