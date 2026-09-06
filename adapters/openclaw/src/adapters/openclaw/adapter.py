"""OpenClaw adapter (provenance=unofficial).

T15.17 scaffold: Handshake + session lifecycle no-ops (no live Gateway).
T15.18 ResolveApproval: honest ``hitl=false`` N/A — no live Gateway in the
reference path, and AgentGavel ``withhold`` has no ``exec.approval.resolve``
enum (only ``allow-once`` / ``allow-always`` / ``deny``). Intended map lives
in :mod:`adapters.openclaw.approvals`; ResolveApproval refuses loudly so
SEC-002/005/006 score N/A (never silent Fail / stub green). T15.19 may add
Events/ExportLedger additively. See
``docs/manual/openclaw-capability-map.md`` and ADR 014.
"""

from __future__ import annotations

from collections.abc import Mapping
from typing import Any
from uuid import uuid4

from agentgavel_adapter.adapter import Adapter

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
            # Conservative until T15.19 wires Gateway event frames.
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
        # Honest empty: no hash-linked ledger projection yet (ledger=false).
        return {"session_id": session_id, "entries": []}

    def stop_session(self, session_id: str) -> None:
        self._sessions.pop(session_id, None)
