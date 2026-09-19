"""Env bootstrap for OpenClaw stub vs live runtime (T17.9 / ADR 015)."""

from __future__ import annotations

from adapters.openclaw.adapter import OpenClawAdapter
from adapters.openclaw.gateway import (
    GatewayProbeError,
    gateway_url_from_env,
    probe_gateway,
)


def adapter_from_env() -> OpenClawAdapter:
    """Stub by default; live only after a successful Gateway health probe.

    ``hitl`` remains false (withhold unmapped). Missing URL → stub. Set URL
    with a failing probe → fail closed (never silent stub-as-live).
    """
    url = gateway_url_from_env()
    if not url:
        return OpenClawAdapter(runtime="stub")
    try:
        report = probe_gateway(url)
    except GatewayProbeError:
        raise
    return OpenClawAdapter(
        runtime="live",
        framework_version=str(report.get("framework_version") or "probed"),
        gateway_url=url,
    )
