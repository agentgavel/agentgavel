"""OpenClaw ResolveApproval decision map (T15.18).

Documents how AgentGavel wire Decisions would map onto Gateway
``exec.approval.resolve`` / ``approval.resolve`` decision strings
(``allow-once`` | ``allow-always`` | ``deny``) when a live Gateway is
available with ``operator.approvals`` scope.

This module is **not** a live Gateway client. Handshake keeps
``hitl=false`` until a probe can drive approve / deny / withhold as the
harness human without auto-approve or silent timeout approve (capability
map; ADR 014). ``withhold`` has no OpenClaw enum — that gap alone forces
honest N/A until a documented withhold surface exists.
"""

from __future__ import annotations

# AgentGavel Decision wire names (proto / adapter contract).
_DECISION_BY_INT = {
    1: "approve",
    2: "deny",
    3: "withhold",
}
_VALID_WIRE = frozenset(_DECISION_BY_INT.values())

# OpenClaw Gateway exec.approval.resolve decision strings (public docs).
OPENCLAW_ALLOW_ONCE = "allow-once"
OPENCLAW_ALLOW_ALWAYS = "allow-always"
OPENCLAW_DENY = "deny"

# Intended harness map when hitl flips true (approve → one-shot allow only;
# never allow-always from the harness — that widens host policy permanently).
_WIRE_TO_OPENCLAW: dict[str, str] = {
    "approve": OPENCLAW_ALLOW_ONCE,
    "deny": OPENCLAW_DENY,
}

# RPC method names from docs/manual/openclaw-capability-map.md.
GATEWAY_RESOLVE_METHODS = (
    "exec.approval.resolve",
    "approval.resolve",
)


class ApprovalMapError(ValueError):
    """Wire decision cannot be mapped onto an OpenClaw resolve enum."""


def wire_decision(decision: str | int) -> str:
    """Normalize a ResolveApproval decision to AgentGavel wire names."""
    if isinstance(decision, bool):
        raise ApprovalMapError(f"unknown decision {decision!r}")
    if isinstance(decision, int):
        try:
            return _DECISION_BY_INT[decision]
        except KeyError as exc:
            raise ApprovalMapError(f"unknown decision wire value {decision}") from exc
    name = str(decision).strip().lower()
    prefix = "decision_"
    if name.startswith(prefix):
        name = name[len(prefix) :]
    if name not in _VALID_WIRE:
        raise ApprovalMapError(f"unknown decision {decision!r}")
    return name


def openclaw_resolve_decision(decision: str | int) -> str:
    """Map an AgentGavel decision onto an OpenClaw Gateway resolve enum.

    Raises :class:`ApprovalMapError` for unknown wire values and for
    ``withhold`` (no Gateway equivalent — keeps hitl=false honest).
    """
    wire = wire_decision(decision)
    mapped = _WIRE_TO_OPENCLAW.get(wire)
    if mapped is None:
        raise ApprovalMapError(
            f"OpenClaw Gateway has no resolve enum for wire decision {wire!r} "
            f"(only {OPENCLAW_ALLOW_ONCE!r} / {OPENCLAW_ALLOW_ALWAYS!r} / "
            f"{OPENCLAW_DENY!r}); CapabilityReport.hitl stays false"
        )
    return mapped
