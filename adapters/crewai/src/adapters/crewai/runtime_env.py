"""Env bootstrap for CrewAI stub vs live runtime (T18.6 / ADR 015)."""

from __future__ import annotations

import os

from adapters.crewai.adapter import CrewAIAdapter


def require_crewai() -> str:
    """Import crewai and return its package version; fail closed if missing."""
    try:
        import crewai  # noqa: F401
    except ImportError as exc:
        raise RuntimeError(
            "crewai package required for runtime=live; "
            "install with: pip install 'agentgavel-adapter-crewai[live]'"
        ) from exc
    from importlib.metadata import PackageNotFoundError, version

    try:
        return version("crewai")
    except PackageNotFoundError as exc:
        raise RuntimeError("crewai installed but metadata missing") from exc


def runtime_from_env() -> str:
    raw = (
        (os.environ.get("AGENTGAVEL_CREWAI_RUNTIME") or os.environ.get("CREWAI_RUNTIME") or "stub")
        .strip()
        .lower()
    )
    if raw in ("stub", "live"):
        return raw
    raise RuntimeError(f"AGENTGAVEL_CREWAI_RUNTIME must be stub|live, got {raw!r}")


def adapter_from_env() -> CrewAIAdapter:
    """Stub by default; live requires the optional crewai package."""
    runtime = runtime_from_env()
    if runtime == "live":
        ver = require_crewai()
        return CrewAIAdapter(runtime="live", framework_version=ver)
    return CrewAIAdapter(runtime="stub")
