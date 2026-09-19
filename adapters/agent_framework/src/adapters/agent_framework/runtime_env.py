"""Env bootstrap for Microsoft Agent Framework stub vs live (ADR 015 / E18)."""

from __future__ import annotations

import os

from adapters.agent_framework.adapter import AgentFrameworkAdapter


def require_agent_framework() -> str:
    try:
        import agent_framework  # noqa: F401
    except ImportError:
        try:
            import microsoft_agents  # noqa: F401
        except ImportError as exc:
            raise RuntimeError(
                "agent-framework package required for runtime=live; "
                "install with: pip install 'agentgavel-adapter-agent-framework[live]'"
            ) from exc
    from importlib.metadata import PackageNotFoundError, version

    for dist in ("agent-framework", "agent_framework", "microsoft-agents"):
        try:
            return version(dist)
        except PackageNotFoundError:
            continue
    return "probed"


def runtime_from_env() -> str:
    raw = (
        (
            os.environ.get("AGENTGAVEL_AGENT_FRAMEWORK_RUNTIME")
            or os.environ.get("AGENT_FRAMEWORK_RUNTIME")
            or "stub"
        )
        .strip()
        .lower()
    )
    if raw in ("stub", "live"):
        return raw
    raise RuntimeError(f"AGENTGAVEL_AGENT_FRAMEWORK_RUNTIME must be stub|live, got {raw!r}")


def adapter_from_env() -> AgentFrameworkAdapter:
    runtime = runtime_from_env()
    if runtime == "live":
        return AgentFrameworkAdapter(runtime="live", framework_version=require_agent_framework())
    return AgentFrameworkAdapter(runtime="stub")
