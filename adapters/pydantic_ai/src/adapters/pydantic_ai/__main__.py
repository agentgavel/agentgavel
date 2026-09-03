"""CLI entry: ``python -m adapters.pydantic_ai`` serves stdio; ``--help`` prints usage."""

from __future__ import annotations

import argparse

from adapters.pydantic_ai.adapter import PydanticAIAdapter


def _build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="python -m adapters.pydantic_ai",
        description=(
            "Unofficial AgentGavel Pydantic AI adapter (stdio JSON-RPC). "
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
    PydanticAIAdapter().serve()


if __name__ == "__main__":
    main()
