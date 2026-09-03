"""Sire adapter (provenance=unofficial).

Lifecycle hooks delegate to a :class:`SireClient`. Default is the in-memory
stub; tests inject a mock. Author-affiliated adapters cannot self-ratify
(ADR 007). HITL/ledger/observability stay false until T10.3/T10.4 wire them.
"""

from __future__ import annotations

from collections.abc import Mapping
from typing import Any

from agentgavel_adapter.adapter import Adapter

from adapters.sire.client import SireClient, StubSireClient

_ADAPTER_VERSION = "0.0.1"


class SireAdapter(Adapter):
    """Unofficial Sire sidecar: Handshake plus client-backed start/submit/stop."""

    def __init__(self, client: SireClient | None = None) -> None:
        super().__init__()
        self._client: SireClient = client if client is not None else StubSireClient()

    def handshake(
        self,
        engine_protocol_version: str,
        *,
        engine_version: str | None = None,
    ) -> Mapping[str, Any]:
        del engine_version  # reserved for future negotiation
        return {
            "adapter_protocol_version": engine_protocol_version or "1.0",
            "adapter_name": "sire",
            "adapter_version": _ADAPTER_VERSION,
            # ADR 007: author-affiliated; cannot self-ratify.
            "provenance": "unofficial",
            # Honest: ResolveApproval/ledger/events are still stubs (T10.3/T10.4).
            "hitl": False,
            "tenancy": False,
            "ledger": False,
            "observability": False,
            "context_mode": "none",
            "framework_name": "sire",
            "framework_version": "unknown",
        }

    def start_session(self, config: Mapping[str, Any]) -> Mapping[str, Any]:
        session_id = self._client.start_session(config)
        return {"id": session_id}

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        self._client.submit_task(session_id, task)

    def resolve_approval(
        self,
        session_id: str,
        approval_id: str,
        decision: str | int,
        *,
        principal: str | None = None,
    ) -> None:
        # T10.3 wires POST /approvals/{approvalId}/decide.
        del session_id, approval_id, decision, principal
        return None

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        return {"session_id": session_id, "entries": []}

    def stop_session(self, session_id: str) -> None:
        self._client.stop_session(session_id)
