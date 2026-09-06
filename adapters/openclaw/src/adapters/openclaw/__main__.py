"""CLI entry: ``python -m adapters.openclaw`` serves stdio; ``--help`` prints usage."""

from __future__ import annotations

import argparse

from adapters.openclaw.adapter import OpenClawAdapter


def _build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="python -m adapters.openclaw",
        description=(
            "Unofficial AgentGavel OpenClaw adapter (stdio JSON-RPC, gateway-style). "
            "Provenance is always unofficial (ADR 007 / ADR 014)."
        ),
    )
    parser.add_argument(
        "--serve",
        action="store_true",
        help="run the stdio JSON-RPC serve loop (default with no flags)",
    )
    return parser


def main(argv: list[str] | None = None) -> None:
    """Parse argv; ``--help`` prints usage, otherwise start stdio serve()."""
    _build_parser().parse_args(argv)
    OpenClawAdapter().serve()


if __name__ == "__main__":
    main()
