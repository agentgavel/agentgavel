"""AgentGavel adapter base class (proto/adapter.proto contract).

Subclasses implement session lifecycle hooks. The SDK owns JSON-RPC stdio
transport (ADR 002) and Event buffering via ``emit()``.
"""

from __future__ import annotations

import sys
from typing import TYPE_CHECKING, Any, BinaryIO, Mapping, MutableMapping, Sequence, TextIO

if TYPE_CHECKING:
    from agentgavel_adapter.transport import TransportLoop


class Adapter:
    """Base class for Python AgentGavel sidecar adapters.

    Method shapes mirror the AgentGavelAdapter RPCs in ``proto/adapter.proto``.
    Concrete adapters override these hooks; the default implementations raise
    ``NotImplementedError`` until behavior is filled in.

    Call ``serve()`` to run the stdio JSON-RPC loop. Inside hooks (or framework
    callbacks), call ``emit(event)`` to push Event notifications to the engine.
    """

    def __init__(self) -> None:
        self._transport: TransportLoop | None = None

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

        Corresponds to ``Events``. Prefer ``emit()`` from framework hooks;
        subclasses may override when they need to surface queued events.
        """
        raise NotImplementedError("emit_events")

    def emit(self, event: Mapping[str, Any]) -> None:
        """Buffer and send an Event notification to the engine over stdio.

        Must be called while ``serve()`` is running (transport attached).
        ``event`` keys match the Event message in ``proto/adapter.proto``
        (snake_case), including one of the kind payloads such as
        ``tool_invocation``.
        """
        transport = self._transport
        if transport is None:
            raise RuntimeError("emit() requires an active serve() transport")
        transport.emit(event)

    def serve(
        self,
        reader: BinaryIO | TextIO | None = None,
        writer: BinaryIO | TextIO | None = None,
    ) -> None:
        """Run the JSON-RPC stdio loop until the peer closes the input stream.

        Defaults to process stdin/stdout buffers. Dispatches ``Handshake``;
        other RPC methods are wired in a follow-on release.
        """
        from agentgavel_adapter.transport import serve_adapter

        serve_adapter(
            self,
            reader if reader is not None else sys.stdin.buffer,
            writer if writer is not None else sys.stdout.buffer,
        )

    def _attach_transport(self, loop: TransportLoop) -> None:
        self._transport = loop

    def _detach_transport(self, loop: TransportLoop) -> None:
        if self._transport is loop:
            self._transport = None


# Alias matching the proto service name.
AgentGavelAdapter = Adapter
