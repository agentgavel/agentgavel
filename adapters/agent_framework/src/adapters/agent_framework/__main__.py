"""CLI entry: ``python -m adapters.agent_framework`` serves stdio; ``--help`` prints usage."""

from __future__ import annotations

import argparse

from adapters.agent_framework.adapter import AgentFrameworkAdapter


def _build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="python -m adapters.agent_framework",
        description=(
            "Unofficial AgentGavel Microsoft Agent Framework adapter "
            "(AutoGen successor; stdio JSON-RPC). "
            "Provenance is always unofficial (ADR 007)."
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
    # Matches FakeAdapter / LangGraphAdapter default: module entry starts stdio serve.
    AgentFrameworkAdapter().serve()


if __name__ == "__main__":
    main()
