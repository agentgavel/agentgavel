"""Env bootstrap for Pydantic AI stub vs live runtime (ADR 015 / E18)."""

from __future__ import annotations

import os

from adapters.pydantic_ai.adapter import PydanticAIAdapter


def require_pydantic_ai() -> str:
    try:
        import pydantic_ai  # noqa: F401
    except ImportError as exc:
        raise RuntimeError(
            "pydantic-ai package required for runtime=live; "
            "install with: pip install 'agentgavel-adapter-pydantic-ai[live]'"
        ) from exc
    from importlib.metadata import PackageNotFoundError, version

    for dist in ("pydantic-ai", "pydantic_ai"):
        try:
            return version(dist)
        except PackageNotFoundError:
            continue
    return "probed"


def runtime_from_env() -> str:
    raw = (
        (
            os.environ.get("AGENTGAVEL_PYDANTIC_AI_RUNTIME")
            or os.environ.get("PYDANTIC_AI_RUNTIME")
            or "stub"
        )
        .strip()
        .lower()
    )
    if raw in ("stub", "live"):
        return raw
    raise RuntimeError(f"AGENTGAVEL_PYDANTIC_AI_RUNTIME must be stub|live, got {raw!r}")


def adapter_from_env() -> PydanticAIAdapter:
    runtime = runtime_from_env()
    if runtime == "live":
        return PydanticAIAdapter(runtime="live", framework_version=require_pydantic_ai())
    return PydanticAIAdapter(runtime="stub")
