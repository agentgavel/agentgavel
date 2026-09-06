"""Hermes Agent adapter (provenance=unofficial).

Gateway-style sidecar (ADR 014) for Hermes Agent
(https://github.com/nousresearch/hermes-agent). Speaks AgentGavel wire
protocol; drives Hermes via a documented control plane — not an in-process
graph SDK.

Intended control plane (pinned for T15.24+; see
docs/manual/hermes-capability-map.md): OpenAI-compatible API server
(``gateway/platforms/api_server.py``) — HTTP + SSE, language-agnostic,
``GET /v1/capabilities``, ``POST /v1/runs``, ``POST /v1/runs/{id}/approval``.
TUI gateway JSON-RPC and ACP are alternatives; this sidecar targets the API
server path.

T15.25: ResolveApproval posts to Hermes ``POST /v1/runs/{run_id}/approval``
(``choice`` once/deny; withhold leaves the gate pending). CapabilityReport
``hitl=true``. Unofficial until ADR 007 ratification. Stub client needs no
live Hermes.

T15.26: Events mapping (tool_invocation / gate_decision) + ExportLedger
honesty (``ledger=false`` while entries stay empty). ``observability=true``
when frames map via ``ingest_hermes_event``.
"""

from __future__ import annotations

import time
from collections.abc import Mapping, MutableMapping, Sequence
from typing import Any

from agentgavel_adapter.adapter import Adapter

from adapters.hermes.client import HermesClient, StubHermesClient, wire_decision
from adapters.hermes.events import (
    assert_tool_invocation_order,
    empty_ledger,
    make_event,
    map_hermes_frame,
)

_ADAPTER_VERSION = "0.0.1"
_GATE_SOURCE_HARNESS = "harness"


class HitlNotSupportedError(RuntimeError):
    """ResolveApproval called while CapabilityReport.hitl is false.

    Kept for API stability; T15.25 wires hitl=true via the API-server path.
    """


class HermesAdapter(Adapter):
    """Unofficial Hermes Agent sidecar: Handshake + ResolveApproval + Events."""

    def __init__(self, client: HermesClient | None = None) -> None:
        super().__init__()
        self._client: HermesClient = client if client is not None else StubHermesClient()
        self._sessions: set[str] = set()
        self._seq: dict[str, int] = {}
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
            "adapter_name": "hermes",
            "adapter_version": _ADAPTER_VERSION,
            # ADR 007 / ADR 014: unofficial until ratification.
            "provenance": "unofficial",
            # ResolveApproval → POST /v1/runs/{run_id}/approval (T15.25).
            "hitl": True,
            "tenancy": False,
            # Trajectories ≠ hash-linked wire Ledger (capability map / T15.26).
            "ledger": False,
            # API/SSE frames map to tool_invocation + gate_decision (T15.26).
            "observability": True,
            "context_mode": "none",
            "framework_name": "hermes",
            # No live Hermes required for unit tests; probe sets real version later.
            "framework_version": "unknown",
        }

    def start_session(self, config: Mapping[str, Any]) -> Mapping[str, Any]:
        session_id = self._client.start_session(config)
        self._sessions.add(session_id)
        self._seq.setdefault(session_id, 0)
        return {"id": session_id}

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        self._client.submit_task(session_id, task)

    def resolve_approval(
        self,
        session_id: str,
        approval_id: str,
        decision: str | int,
        *,
        principal: str | None = None,
    ) -> None:
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        wire = wire_decision(decision)
        self._client.resolve_approval(
            session_id,
            approval_id,
            wire,
            principal=principal,
        )
        gate: dict[str, Any] = {
            "approval_id": approval_id,
            "source": _GATE_SOURCE_HARNESS,
            "decision": wire,
            "genuine_hitl": True,
        }
        if principal:
            gate["principal"] = principal
        event: dict[str, Any] = {
            "session_id": session_id,
            "seq": self._next_seq(session_id),
            "unix_ms": int(time.time() * 1000),
            "gate_decision": gate,
        }
        self._record_event(event)

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        # Honest empty: trajectories / session export are not wire Ledger.
        ledger = empty_ledger(session_id)
        entries = ledger["entries"]
        if not isinstance(entries, Sequence) or isinstance(entries, (str, bytes)):
            raise TypeError("export_ledger() entries must be a sequence")
        return {"session_id": session_id, "entries": list(entries)}

    def stop_session(self, session_id: str) -> None:
        self._client.stop_session(session_id)
        self._sessions.discard(session_id)
        self._seq.pop(session_id, None)

    def ingest_hermes_event(
        self,
        session_id: str,
        frame: Mapping[str, Any],
    ) -> list[MutableMapping[str, Any]]:
        """Map a Hermes API/SSE frame and emit AgentGavel Events (T15.26).

        Used by unit tests and by a future live ``GET /v1/runs/{id}/events``
        subscriber. Unknown frames are ignored (empty list).
        """
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        recorded: list[MutableMapping[str, Any]] = []
        for kind_payload in map_hermes_frame(frame):
            tool = kind_payload.get("tool_invocation")
            gate = kind_payload.get("gate_decision")
            event = make_event(
                session_id,
                self._next_seq(session_id),
                tool_invocation_payload=tool if isinstance(tool, Mapping) else None,
                gate_decision_payload=gate if isinstance(gate, Mapping) else None,
            )
            self._record_event(event)
            recorded.append(event)
        return recorded

    def _next_seq(self, session_id: str) -> int:
        n = self._seq.get(session_id, 0) + 1
        self._seq[session_id] = n
        return n

    def _record_event(self, event: MutableMapping[str, Any]) -> None:
        self.emitted.append(event)
        if self._transport is not None:
            self.emit(event)

    def assert_emitted_tool_order(self) -> None:
        """Test helper: before/after ordering for recorded tool_invocation."""
        assert_tool_invocation_order(self.emitted)
