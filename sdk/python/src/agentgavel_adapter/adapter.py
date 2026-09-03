"""AgentGavel adapter base class (proto/adapter.proto contract).

Subclasses implement session lifecycle hooks. JSON-RPC stdio transport is
out of scope for this scaffold (see follow-on SDK tasks).
"""

from __future__ import annotations

from typing import Any, Mapping, MutableMapping, Sequence


class Adapter:
    """Base class for Python AgentGavel sidecar adapters.

    Method shapes mirror the AgentGavelAdapter RPCs in ``proto/adapter.proto``.
    Concrete adapters override these hooks; the default implementations raise
    ``NotImplementedError`` until behavior is filled in.
    """

    def handshake(
        self,
        engine_protocol_version: str,
        *,
        engine_version: str | None = None,
    ) -> Mapping[str, Any]:
        """Negotiate protocol version and return a capability report.

        Corresponds to ``Handshake`` / ``CapabilityReport``.
        """
        raise NotImplementedError("handshake")

    def start_session(
        self,
        config: Mapping[str, Any],
    ) -> Mapping[str, Any]:
        """Create a run with model binding and MCP fixture endpoints.

        Corresponds to ``StartSession`` / ``SessionConfig`` → ``SessionId``.
        ``config`` keys match SessionConfig field names (snake_case).
        Returns a mapping with at least ``id`` (session id).
        """
        raise NotImplementedError("start_session")

    def submit_task(
        self,
        session_id: str,
        task: Mapping[str, Any],
    ) -> None:
        """Deliver the scenario task prompt to the target.

        Corresponds to ``SubmitTask`` / ``SubmitTaskRequest``.
        ``task`` keys match Task field names (``id``, ``prompt``, ``metadata``).
        """
        raise NotImplementedError("submit_task")

    def resolve_approval(
        self,
        session_id: str,
        approval_id: str,
        decision: str | int,
        *,
        principal: str | None = None,
    ) -> None:
        """Act as the human approver for a pending HITL gate.

        Corresponds to ``ResolveApproval`` / ``ResolveApprovalRequest``.
        ``decision`` is a Decision enum name or wire value
        (approve / deny / withhold).
        """
        raise NotImplementedError("resolve_approval")

    def export_ledger(
        self,
        session_id: str,
    ) -> Mapping[str, Any]:
        """Return the audit ledger for the session.

        Corresponds to ``ExportLedger`` / ``Ledger``.
        Returns a mapping with ``session_id`` and ``entries``.
        """
        raise NotImplementedError("export_ledger")

    def stop_session(
        self,
        session_id: str,
    ) -> None:
        """Tear down the session.

        Corresponds to ``StopSession``.
        """
        raise NotImplementedError("stop_session")

    def emit_events(
        self,
        session_id: str,
    ) -> Sequence[Mapping[str, Any]] | MutableMapping[str, Any]:
        """Optional hook for adapter-pushed Events stream items.

        Corresponds to ``Events``. Transport framing is deferred; subclasses
        may override when they need to surface queued events.
        """
        raise NotImplementedError("emit_events")


# Alias matching the proto service name.
AgentGavelAdapter = Adapter
