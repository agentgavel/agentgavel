"""LangGraph-style interrupt / HITL bridge for ResolveApproval (T11.3).

When interrupt support is enabled, gated tool nodes pause before side effects
and wait for a harness ``ResolveApproval``. That maps to LangGraph's
``interrupt()`` + resume pattern without pulling the langgraph PyPI stack.

When support is disabled, CapabilityReport must keep ``hitl=false`` — never
claim HITL while ResolveApproval is a no-op.
"""

from __future__ import annotations

from collections.abc import Mapping, MutableMapping
from dataclasses import dataclass, field
from typing import Any
from uuid import uuid4

# AgentGavel Decision wire names (proto/adapter.proto).
_DECISION_BY_INT = {
    1: "approve",
    2: "deny",
    3: "withhold",
}
_VALID_DECISIONS = frozenset(_DECISION_BY_INT.values())

# Tools that must not run side effects until a genuine HITL decision.
# send_email is the SEC outbound side-effect; read_email stays ungated.
DEFAULT_GATED_TOOLS = frozenset({"send_email"})


class HitlNotSupportedError(RuntimeError):
    """ResolveApproval called while CapabilityReport.hitl is false."""


class InterruptError(RuntimeError):
    """Interrupt / approval state is inconsistent."""


def wire_decision(decision: str | int) -> str:
    """Normalize a ResolveApproval decision to AgentGavel wire names."""
    if isinstance(decision, bool):
        raise InterruptError(f"unknown decision {decision!r}")
    if isinstance(decision, int):
        try:
            return _DECISION_BY_INT[decision]
        except KeyError as exc:
            raise InterruptError(f"unknown decision wire value {decision}") from exc
    name = str(decision).strip().lower()
    prefix = "decision_"
    if name.startswith(prefix):
        name = name[len(prefix) :]
    if name not in _VALID_DECISIONS:
        raise InterruptError(f"unknown decision {decision!r}")
    return name


@dataclass
class PendingInterrupt:
    """One LangGraph-style interrupt waiting on ResolveApproval."""

    session_id: str
    approval_id: str
    tool_name: str
    arguments: dict[str, Any]
    call_id: str
    resolved: bool = False
    decision: str | None = None
    principal: str | None = None


@dataclass
class InterruptSupport:
    """In-process interrupt registry. Presence ⇒ hitl capability is real.

    Construct with ``enabled=True`` (default) to advertise and honor HITL.
    Pass ``enabled=False`` (or use :func:`disabled_interrupt_support`) so
    CapabilityReport.hitl stays false and ResolveApproval refuses loudly.
    """

    enabled: bool = True
    gated_tools: frozenset[str] = field(default_factory=lambda: DEFAULT_GATED_TOOLS)
    _pending: dict[str, PendingInterrupt] = field(default_factory=dict, init=False, repr=False)

    def is_gated(self, tool_name: str) -> bool:
        return self.enabled and tool_name in self.gated_tools

    def request(
        self,
        session_id: str,
        tool_name: str,
        arguments: Mapping[str, Any],
        call_id: str,
        *,
        approval_id: str | None = None,
    ) -> PendingInterrupt:
        """Record a pending interrupt; caller must not execute the tool yet."""
        if not self.enabled:
            raise HitlNotSupportedError(
                "interrupt support disabled; CapabilityReport.hitl is false"
            )
        if tool_name not in self.gated_tools:
            raise InterruptError(f"tool {tool_name!r} is not HITL-gated")
        aid = approval_id or f"lg-appr-{uuid4().hex[:12]}"
        if aid in self._pending and not self._pending[aid].resolved:
            raise InterruptError(f"approval_id already pending: {aid}")
        pending = PendingInterrupt(
            session_id=session_id,
            approval_id=aid,
            tool_name=tool_name,
            arguments=dict(arguments),
            call_id=call_id,
        )
        self._pending[aid] = pending
        return pending

    def get(self, approval_id: str) -> PendingInterrupt | None:
        return self._pending.get(approval_id)

    def pending_for_session(self, session_id: str) -> list[PendingInterrupt]:
        return [p for p in self._pending.values() if p.session_id == session_id and not p.resolved]

    def resolve(
        self,
        session_id: str,
        approval_id: str,
        decision: str | int,
        *,
        principal: str | None = None,
    ) -> PendingInterrupt:
        """Apply harness decision to a pending interrupt (LangGraph resume)."""
        if not self.enabled:
            raise HitlNotSupportedError(
                "interrupt support disabled; CapabilityReport.hitl is false"
            )
        pending = self._pending.get(approval_id)
        if pending is None:
            raise InterruptError(f"unknown approval_id: {approval_id}")
        if pending.session_id != session_id:
            raise InterruptError(
                f"session mismatch for {approval_id}: "
                f"got {session_id!r}, want {pending.session_id!r}"
            )
        if pending.resolved:
            raise InterruptError(f"approval already resolved: {approval_id}")
        pending.decision = wire_decision(decision)
        pending.principal = principal
        pending.resolved = True
        return pending

    def clear_session(self, session_id: str) -> None:
        dead = [aid for aid, p in self._pending.items() if p.session_id == session_id]
        for aid in dead:
            del self._pending[aid]


def disabled_interrupt_support() -> InterruptSupport:
    """Factory for the honest hitl=false path."""
    return InterruptSupport(enabled=False)


def gate_decision_event(
    *,
    session_id: str,
    seq: int,
    unix_ms: int,
    approval_id: str,
    decision: str,
    principal: str | None = None,
) -> MutableMapping[str, Any]:
    """Build a protocol ``gate_decision`` Event (source=harness, genuine_hitl)."""
    gate: dict[str, Any] = {
        "approval_id": approval_id,
        "source": "harness",
        "decision": decision,
        "genuine_hitl": True,
    }
    if principal:
        gate["principal"] = principal
    return {
        "session_id": session_id,
        "seq": seq,
        "unix_ms": unix_ms,
        "gate_decision": gate,
    }
