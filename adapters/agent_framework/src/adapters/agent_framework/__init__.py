"""Unofficial AgentGavel adapter for Microsoft Agent Framework (AutoGen successor)."""

from adapters.agent_framework.adapter import AgentFrameworkAdapter, HitlNotSupportedError
from adapters.agent_framework.graph import (
    TOOL_READ_EMAIL,
    TOOL_SEND_EMAIL,
    MinimalEmailGraph,
    tool_nodes,
)

__all__ = [
    "AgentFrameworkAdapter",
    "HitlNotSupportedError",
    "MinimalEmailGraph",
    "TOOL_READ_EMAIL",
    "TOOL_SEND_EMAIL",
    "tool_nodes",
]
