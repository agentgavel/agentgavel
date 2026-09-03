"""Python adapter callbacks for AgentGavel sidecars."""

from agentgavel_adapter.adapter import Adapter, AgentGavelAdapter
from agentgavel_adapter.transport import (
    METHOD_EVENT_NOTIFY,
    METHOD_HANDSHAKE,
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
    "METHOD_HANDSHAKE",
    "StdioConn",
    "TransportLoop",
    "serve_adapter",
]
__version__ = "0.1.0"
