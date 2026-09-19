"""Env bootstrap for LangGraph stub vs live runtime (T17.3 / ADR 015)."""

from __future__ import annotations

import os

from adapters.langgraph.adapter import LangGraphAdapter
from adapters.langgraph.live_graph import langgraph_package_version, require_langgraph


def runtime_from_env() -> str:
    """Return stub|live from AGENTGAVEL_LANGGRAPH_RUNTIME (default stub)."""
    raw = (
        (
            os.environ.get("AGENTGAVEL_LANGGRAPH_RUNTIME")
            or os.environ.get("LANGGRAPH_RUNTIME")
            or "stub"
        )
        .strip()
        .lower()
    )
    if raw in ("stub", "live"):
        return raw
    raise RuntimeError(f"AGENTGAVEL_LANGGRAPH_RUNTIME must be stub|live, got {raw!r}")


def adapter_from_env(*, hitl: bool = True) -> LangGraphAdapter:
    """Build adapter; live mode fails closed if langgraph is not installed."""
    runtime = runtime_from_env()
    if runtime == "live":
        require_langgraph()
        version = langgraph_package_version()
        return LangGraphAdapter(hitl=hitl, runtime="live", framework_version=version)
    return LangGraphAdapter(hitl=hitl, runtime="stub")
