"""Sire adapter scaffold (provenance=unofficial; no real Sire client yet).

Capability flags stay honest for stubs: HITL/ledger/observability are false
until T10.2+ maps real Sire APIs. Author-affiliated adapters cannot self-ratify
(ADR 007).
"""

from __future__ import annotations

from typing import Any, Mapping
from uuid import uuid4

from agentgavel_adapter.adapter import Adapter

_ADAPTER_VERSION = "0.0.1"


class SireAdapter(Adapter):
    """Minimal Sire sidecar: Handshake + no-op lifecycle stubs."""

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
            # Honest stubs — real capabilities land with T10.2+.
            "hitl": False,
            "tenancy": False,
            "ledger": False,
            "observability": False,
            "context_mode": "none",
            "framework_name": "sire",
            "framework_version": "unknown",
        }

    def start_session(self, config: Mapping[str, Any]) -> Mapping[str, Any]:
        del config  # T10.2 maps SessionConfig onto Sire
        return {"id": f"sire-sess-{uuid4().hex[:12]}"}

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
