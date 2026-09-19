"""Env bootstrap for OpenAI Agents stub vs live runtime (ADR 015 / E18)."""

from __future__ import annotations

import os

from adapters.openai_agents.adapter import OpenAIAgentsAdapter


def require_openai_agents() -> str:
    try:
        import agents  # noqa: F401  # openai-agents package import name
    except ImportError:
        try:
            import openai_agents  # noqa: F401
        except ImportError as exc:
            raise RuntimeError(
                "openai-agents package required for runtime=live; "
                "install with: pip install 'agentgavel-adapter-openai-agents[live]'"
            ) from exc
    from importlib.metadata import PackageNotFoundError, version

    for dist in ("openai-agents", "openai_agents"):
        try:
            return version(dist)
        except PackageNotFoundError:
            continue
    return "probed"


def runtime_from_env() -> str:
    raw = (
        (
            os.environ.get("AGENTGAVEL_OPENAI_AGENTS_RUNTIME")
            or os.environ.get("OPENAI_AGENTS_RUNTIME")
            or "stub"
        )
        .strip()
        .lower()
    )
    if raw in ("stub", "live"):
        return raw
    raise RuntimeError(f"AGENTGAVEL_OPENAI_AGENTS_RUNTIME must be stub|live, got {raw!r}")


def adapter_from_env() -> OpenAIAgentsAdapter:
    runtime = runtime_from_env()
    if runtime == "live":
        return OpenAIAgentsAdapter(runtime="live", framework_version=require_openai_agents())
    return OpenAIAgentsAdapter(runtime="stub")
