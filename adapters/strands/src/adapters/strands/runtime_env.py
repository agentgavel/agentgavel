"""Env bootstrap for AWS Strands stub vs live runtime (ADR 015 / E18)."""

from __future__ import annotations

import os

from adapters.strands.adapter import StrandsAdapter


def require_strands() -> str:
    try:
        import strands  # noqa: F401
    except ImportError as exc:
        raise RuntimeError(
            "strands-agents package required for runtime=live; "
            "install with: pip install 'agentgavel-adapter-strands[live]'"
        ) from exc
    from importlib.metadata import PackageNotFoundError, version

    for dist in ("strands-agents", "strands"):
        try:
            return version(dist)
        except PackageNotFoundError:
            continue
    return "probed"


def runtime_from_env() -> str:
    raw = (
        (
            os.environ.get("AGENTGAVEL_STRANDS_RUNTIME")
            or os.environ.get("STRANDS_RUNTIME")
            or "stub"
        )
        .strip()
        .lower()
    )
    if raw in ("stub", "live"):
        return raw
    raise RuntimeError(f"AGENTGAVEL_STRANDS_RUNTIME must be stub|live, got {raw!r}")


def adapter_from_env() -> StrandsAdapter:
    runtime = runtime_from_env()
    if runtime == "live":
        return StrandsAdapter(runtime="live", framework_version=require_strands())
    return StrandsAdapter(runtime="stub")
