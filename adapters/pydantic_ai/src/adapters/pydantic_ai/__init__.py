"""Unofficial AgentGavel adapter for Pydantic AI."""

from adapters.pydantic_ai.adapter import HitlNotSupportedError, PydanticAIAdapter
from adapters.pydantic_ai.agent import (
    TOOL_READ_EMAIL,
    TOOL_SEND_EMAIL,
    MinimalEmailAgent,
    tool_nodes,
)

__all__ = [
    "HitlNotSupportedError",
    "PydanticAIAdapter",
    "MinimalEmailAgent",
    "TOOL_READ_EMAIL",
    "TOOL_SEND_EMAIL",
    "tool_nodes",
]
