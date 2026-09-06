"""OpenClaw adapter (provenance=unofficial).

T15.17 scaffold: Handshake + session lifecycle no-ops (no live Gateway).
T15.18 ResolveApproval: honest ``hitl=false`` N/A — no live Gateway in the
reference path, and AgentGavel ``withhold`` has no ``exec.approval.resolve``
enum (only ``allow-once`` / ``allow-always`` / ``deny``). Intended map lives
in :mod:`adapters.openclaw.approvals`; ResolveApproval refuses loudly so
SEC-002/005/006 score N/A (never silent Fail / stub green). T15.19 may add
Events/ExportLedger additively. See
``docs/manual/openclaw-capability-map.md`` and ADR 014.

T15.19: Event builders + ExportLedger honesty. Without a live Gateway
subscription the Events stream is incomplete → ``observability=false``.
Audit activity is not AgentGavel hash-linked → ``ledger=false``. ResolveApproval
remains owned by T15.18 (hitl=false / HitlNotSupportedError until wired).
"""

from __future__ import annotations

from collections.abc import Mapping, MutableMapping, Sequence
from typing import Any
from uuid import uuid4

from agentgavel_adapter.adapter import Adapter

from adapters.openclaw.events import (
    GATE_SOURCE_HARNESS,
    build_gate_decision,
    build_tool_invocation,
    empty_ledger,
    make_event,
    map_gateway_frame,
)

_ADAPTER_VERSION = "0.0.1"
# No live Gateway probe yet — do not invent a package version.
_FRAMEWORK_VERSION = "unprobed"


class HitlNotSupportedError(RuntimeError):
    """Raised when ResolveApproval is called while hitl is false.

    Loud failure is required so SEC-002 cannot silently Pass/Fail when the
    adapter has no real OpenClaw approval surface (ADR 011 / ADR 014).
    """


class OpenClawAdapter(Adapter):
    """Unofficial OpenClaw gateway-style sidecar: Handshake + lifecycle stubs."""

    def __init__(self) -> None:
        super().__init__()
        self._sessions: dict[str, dict[str, Any]] = {}
        self._seq: dict[str, int] = {}
        # Buffered Events for tests / future Gateway subscription (T15.19).
        self.emitted: list[MutableMapping[str, Any]] = []

    def handshake(
        self,
        engine_protocol_version: str,
        *,
        engine_version: str | None = None,
    ) -> Mapping[str, Any]:
        del engine_version  # reserved for future negotiation
        return {
            "adapter_protocol_version": engine_protocol_version or "1.0",
            "adapter_name": "openclaw",
            "adapter_version": _ADAPTER_VERSION,
            # ADR 007 / ADR 014: unofficial until ratification.
            "provenance": "unofficial",
            # T15.18: documented N/A — no live Gateway + withhold gap (see
            # adapters.openclaw.approvals and README). Never stub green HITL.
            "hitl": False,
            "tenancy": False,
            # Audit RPC is not AgentGavel hash-linked ledger (capability map).
            "ledger": False,
            # No live Gateway event subscription yet: helpers can emit shapes,
            # but before/after completeness is unproven → observability penalty.
            "observability": False,
            "context_mode": "none",
            "framework_name": "openclaw",
            "framework_version": _FRAMEWORK_VERSION,
        }

    def start_session(self, config: Mapping[str, Any]) -> Mapping[str, Any]:
        # Scaffold: track session locally; no Gateway sessions.create yet.
        session_id = f"openclaw-sess-{uuid4().hex[:12]}"
        self._sessions[session_id] = dict(config)
        return {"id": session_id}

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        # Scaffold no-op: live Gateway sessions.send lands with later probes.
        # No tool_invocation Events without a Gateway subscription (N/A path).
        del task

    def resolve_approval(
        self,
        session_id: str,
        approval_id: str,
        decision: str | int,
        *,
        principal: str | None = None,
    ) -> None:
        del approval_id, decision, principal
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        # T15.18 honest N/A: do not call Gateway or emit gate_decision.
        # CapabilityReport.hitl is false → SEC-002/005/006 N/A (engine).
        raise HitlNotSupportedError(
            "OpenClaw ResolveApproval N/A (T15.18): no live Gateway and "
            "withhold has no exec.approval.resolve enum; "
            "CapabilityReport.hitl is false (SEC-002/005/006 score N/A)"
        )

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        # Honest empty: OpenClaw audit.activity.list is metadata-only and not a
        # hash-linked AgentGavel Ledger (ledger=false; SEC-009/010 N/A).
        ledger = empty_ledger(session_id)
        entries = ledger.get("entries")
        if not isinstance(entries, Sequence) or isinstance(entries, (str, bytes)):
            raise TypeError("export_ledger() entries must be a sequence")
        return {"session_id": session_id, "entries": list(entries)}

    def stop_session(self, session_id: str) -> None:
        self._sessions.pop(session_id, None)

    # --- T15.19 Events helpers (additive; T15.18 owns ResolveApproval body) ---

    def _next_seq(self, session_id: str) -> int:
        n = self._seq.get(session_id, 0) + 1
        self._seq[session_id] = n
        return n

    def _buffer_and_emit(self, event: MutableMapping[str, Any]) -> None:
        self.emitted.append(event)
        if self._transport is not None:
            self.emit(event)

    def emit_tool_invocation(
        self,
        session_id: str,
        tool_name: str,
        tool_id: str,
        phase: str,
        *,
        arguments: Mapping[str, Any] | None = None,
        outcome: str | None = None,
        error: str | None = None,
        refused: bool = False,
    ) -> Mapping[str, Any]:
        """Buffer (and transport-emit) a ``tool_invocation`` Event.

        Used when a live Gateway ``session.tool`` frame is available. Scaffold
        lifecycle does not call this (observability stays false).
        """
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        inv = build_tool_invocation(
            tool_name,
            tool_id,
            phase,
            arguments=arguments,
            outcome=outcome,
            error=error,
            refused=refused,
        )
        event = make_event(
            session_id,
            self._next_seq(session_id),
            tool_invocation_payload=inv,
        )
        self._buffer_and_emit(event)
        return event

    def emit_gate_decision(
        self,
        session_id: str,
        approval_id: str,
        decision: str | int,
        *,
        source: str = GATE_SOURCE_HARNESS,
        principal: str | None = None,
        genuine_hitl: bool = False,
    ) -> Mapping[str, Any]:
        """Buffer (and transport-emit) a ``gate_decision`` Event.

        Intended for T15.18 ResolveApproval and Gateway approval frames.
        Does not change CapabilityReport.hitl by itself.
        """
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        gate = build_gate_decision(
            approval_id,
            decision,
            source=source,
            principal=principal,
            genuine_hitl=genuine_hitl,
        )
        event = make_event(
            session_id,
            self._next_seq(session_id),
            gate_decision_payload=gate,
        )
        self._buffer_and_emit(event)
        return event

    def ingest_gateway_frame(
        self,
        session_id: str,
        frame: Mapping[str, Any],
    ) -> Mapping[str, Any] | None:
        """Map a Gateway event frame into the Events stream, or None if N/A."""
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        # Peek seq only when the frame maps; avoid burning seq on N/A frames.
        provisional = self._seq.get(session_id, 0) + 1
        event = map_gateway_frame(session_id, provisional, frame)
        if event is None:
            return None
        self._seq[session_id] = provisional
        self._buffer_and_emit(event)
        return event
