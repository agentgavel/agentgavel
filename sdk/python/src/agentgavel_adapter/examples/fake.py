"""Minimal Python adapter example for engine launcher integration tests.

Mirrors ``internal/engine/testdata/fakeadapter`` (Go) but implemented
against the Python SDK's ``Adapter`` base class and ``serve()`` transport
loop. Useful as a reference adapter and as an interop fixture: the Go
engine launcher can spawn this process (``python3 -m
agentgavel_adapter.examples.fake``) and drive it over stdio JSON-RPC.
"""

from __future__ import annotations

from typing import Any, Mapping

from agentgavel_adapter.adapter import Adapter

_SESSION_ID = "sess-1"


class FakeAdapter(Adapter):
    """Adapter stub that answers the handshake and echoes a fixed session."""

    def handshake(
        self,
        engine_protocol_version: str,
        *,
        engine_version: str | None = None,
    ) -> Mapping[str, Any]:
        return {
            "adapter_protocol_version": "1.0",
            "adapter_name": "fake",
            "adapter_version": "0.0.1",
            "provenance": "unofficial",
            "hitl": True,
            "ledger": True,
            "observability": True,
            "context_mode": "raw",
        }

    def start_session(self, config: Mapping[str, Any]) -> Mapping[str, Any]:
        return {"id": _SESSION_ID}

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        self.emit(
            {
                "session_id": session_id or _SESSION_ID,
                "seq": 1,
                "unix_ms": 1,
                "tool_invocation": {
                    "tool_name": "noop",
                    "tool_id": "1",
                    "phase": "before",
                },
            }
        )

    def resolve_approval(
        self,
        session_id: str,
        approval_id: str,
        decision: str | int,
        *,
        principal: str | None = None,
    ) -> None:
        return None

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        return {"session_id": session_id, "entries": []}

    def stop_session(self, session_id: str) -> None:
        return None


def main() -> None:
    """Run the fake adapter's stdio JSON-RPC serve loop."""
    FakeAdapter().serve()


if __name__ == "__main__":
    main()
