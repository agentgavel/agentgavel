"""CLI entry: ``python -m adapters.hermes`` serves stdio; ``--help`` prints usage."""

from __future__ import annotations

import argparse

from adapters.hermes.adapter import HermesAdapter
from adapters.hermes.client import client_from_env


def _build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="python -m adapters.hermes",
        description=(
            "Unofficial AgentGavel Hermes Agent adapter (stdio JSON-RPC). "
            "Provenance is always unofficial (ADR 007 / ADR 014). "
            "Set AGENTGAVEL_HERMES_API_BASE for live capabilities probe."
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
    HermesAdapter(client=client_from_env()).serve()


if __name__ == "__main__":
    main()
