"""LangGraph adapter scaffold (provenance=unofficial; no graph tools yet).

Capability flags stay honest for stubs: HITL/ledger/observability are false
until T11.2+ maps real LangGraph APIs. Unofficial until a maintainer signs
off (ADR 007).
"""

from __future__ import annotations

from typing import Any, Mapping
from uuid import uuid4

from agentgavel_adapter.adapter import Adapter

_ADAPTER_VERSION = "0.0.1"


class LangGraphAdapter(Adapter):
    """Minimal LangGraph sidecar: Handshake + no-op lifecycle stubs."""

    def handshake(
        self,
        engine_protocol_version: str,
        *,
        engine_version: str | None = None,
    ) -> Mapping[str, Any]:
        del engine_version  # reserved for future negotiation
        return {
            "adapter_protocol_version": engine_protocol_version or "1.0",
            "adapter_name": "langgraph",
            "adapter_version": _ADAPTER_VERSION,
            # ADR 007: unofficial until maintainer ratification.
            "provenance": "unofficial",
            # Honest stubs — real capabilities land with T11.2+.
            "hitl": False,
            "tenancy": False,
            "ledger": False,
            "observability": False,
            "context_mode": "none",
            "framework_name": "langgraph",
            "framework_version": "unknown",
        }

    def start_session(self, config: Mapping[str, Any]) -> Mapping[str, Any]:
        del config  # T11.2 maps SessionConfig onto LangGraph
        return {"id": f"langgraph-sess-{uuid4().hex[:12]}"}

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        del session_id, task
        return None

    def resolve_approval(
        self,
        session_id: str,
        approval_id: str,
        decision: str | int,
        *,
        principal: str | None = None,
    ) -> None:
        del session_id, approval_id, decision, principal
        return None

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        return {"session_id": session_id, "entries": []}

    def stop_session(self, session_id: str) -> None:
        del session_id
        return None
