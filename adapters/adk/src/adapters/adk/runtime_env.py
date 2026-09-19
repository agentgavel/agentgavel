"""Env bootstrap for Google ADK stub vs live runtime (ADR 015 / E18)."""

from __future__ import annotations

import os

from adapters.adk.adapter import AdkAdapter

_LIVE_PACKAGES = ("google-adk", "google_adk", "adk")


def require_adk() -> str:
    """Import google-adk (or alias) and return version; fail closed if missing."""
    last: Exception | None = None
    for name in ("google.adk", "google_adk", "adk"):
        try:
            __import__(name)
            break
        except ImportError as exc:
            last = exc
    else:
        raise RuntimeError(
            "google-adk package required for runtime=live; "
            "install with: pip install 'agentgavel-adapter-adk[live]'"
        ) from last

    from importlib.metadata import PackageNotFoundError, version

    for dist in _LIVE_PACKAGES:
        try:
            return version(dist)
        except PackageNotFoundError:
            continue
    return "probed"


def runtime_from_env() -> str:
    raw = (
        (os.environ.get("AGENTGAVEL_ADK_RUNTIME") or os.environ.get("ADK_RUNTIME") or "stub")
        .strip()
        .lower()
    )
    if raw in ("stub", "live"):
        return raw
    raise RuntimeError(f"AGENTGAVEL_ADK_RUNTIME must be stub|live, got {raw!r}")


def adapter_from_env() -> AdkAdapter:
    runtime = runtime_from_env()
    if runtime == "live":
        return AdkAdapter(runtime="live", framework_version=require_adk())
    return AdkAdapter(runtime="stub")
