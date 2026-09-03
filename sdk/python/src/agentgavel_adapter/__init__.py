"""Python adapter callbacks for AgentGavel sidecars."""

from agentgavel_adapter.adapter import Adapter, AgentGavelAdapter
from agentgavel_adapter.transport import (
    METHOD_EVENT_NOTIFY,
    METHOD_EXPORT_LEDGER,
    METHOD_HANDSHAKE,
    METHOD_RESOLVE_APPROVAL,
    METHOD_START_SESSION,
    METHOD_STOP_SESSION,
    METHOD_SUBMIT_TASK,
    EventBuffer,
    StdioConn,
    TransportLoop,
    serve_adapter,
)

__all__ = [
    "Adapter",
    "AgentGavelAdapter",
    "EventBuffer",
    "METHOD_EVENT_NOTIFY",
    "METHOD_EXPORT_LEDGER",
    "METHOD_HANDSHAKE",
    "METHOD_RESOLVE_APPROVAL",
    "METHOD_START_SESSION",
    "METHOD_STOP_SESSION",
    "METHOD_SUBMIT_TASK",
    "StdioConn",
    "TransportLoop",
    "serve_adapter",
]
__version__ = "0.1.0"
