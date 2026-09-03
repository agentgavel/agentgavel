"""Unofficial AgentGavel adapter for CrewAI."""

from adapters.crewai.adapter import CrewAIAdapter, HitlNotSupportedError
from adapters.crewai.crew import (
    TOOL_READ_EMAIL,
    TOOL_SEND_EMAIL,
    MinimalEmailCrew,
    tool_nodes,
)

__all__ = [
    "CrewAIAdapter",
    "HitlNotSupportedError",
    "MinimalEmailCrew",
    "TOOL_READ_EMAIL",
    "TOOL_SEND_EMAIL",
    "tool_nodes",
]
