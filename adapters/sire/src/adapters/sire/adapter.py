"""Sire adapter (provenance=unofficial).

Lifecycle hooks delegate to a :class:`SireClient`. Default is the in-memory
stub; tests inject a mock. Author-affiliated adapters cannot self-ratify
(ADR 007). ResolveApproval is wired (hitl=true). ExportLedger returns a
Ledger-shaped mapping via the client but entries stay empty until Sire exposes
a session-scoped hash-linked ledger — so ledger=false and observability=false
(observability penalty / GSI cap 600; SEC-009/010 N/A).
"""

from __future__ import annotations

import time
from collections.abc import Mapping, MutableMapping, Sequence
from typing import Any

from agentgavel_adapter.adapter import Adapter

from adapters.sire.client import SireClient, StubSireClient, wire_decision

_ADAPTER_VERSION = "0.0.1"
_GATE_SOURCE_HARNESS = "harness"


class SireAdapter(Adapter):
    """Unofficial Sire sidecar: Handshake plus client-backed lifecycle."""

    def __init__(self, client: SireClient | None = None) -> None:
        super().__init__()
        self._client: SireClient = client if client is not None else StubSireClient()
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
            "adapter_name": "sire",
            "adapter_version": _ADAPTER_VERSION,
            # ADR 007: author-affiliated; cannot self-ratify.
            "provenance": "unofficial",
            # ResolveApproval posts to Sire and emits gate_decision (T10.3).
            "hitl": True,
            "tenancy": False,
            # ExportLedger is wired but stub-empty: no session hash-chain yet.
            "ledger": False,
            # Event sink incomplete (tool_invocation / ledger_append absent) →
            # observability penalty (GSI hard cap 600) per RFC 0001 §4.11.
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
        self.emitted.append(event)
        if self._transport is not None:
            self.emit(event)

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        ledger = self._client.export_ledger(session_id)
        if not isinstance(ledger, Mapping):
            raise TypeError("export_ledger() must return a mapping (Ledger)")
        session = ledger.get("session_id")
        entries = ledger.get("entries")
        if session != session_id:
            raise TypeError(
                f"export_ledger session_id mismatch: got {session!r}, want {session_id!r}"
            )
        if not isinstance(entries, Sequence) or isinstance(entries, (str, bytes)):
            raise TypeError("export_ledger() entries must be a sequence")
        return {"session_id": session_id, "entries": list(entries)}

    def stop_session(self, session_id: str) -> None:
        self._client.stop_session(session_id)

    def _next_seq(self, session_id: str) -> int:
        n = self._seq.get(session_id, 0) + 1
        self._seq[session_id] = n
        return n
