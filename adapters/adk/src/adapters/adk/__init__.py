"""Unofficial AgentGavel adapter for Google ADK."""

from adapters.adk.adapter import AdkAdapter, HitlNotSupportedError
from adapters.adk.graph import (
    TOOL_READ_EMAIL,
    TOOL_SEND_EMAIL,
    MinimalEmailGraph,
    tool_nodes,
)

__all__ = [
    "AdkAdapter",
    "HitlNotSupportedError",
    "MinimalEmailGraph",
    "TOOL_READ_EMAIL",
    "TOOL_SEND_EMAIL",
    "tool_nodes",
]
