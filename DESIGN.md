# AgentGavel marketing site — DESIGN.md

## Direction (adopted 2026-09-05)

Canonical home page is the long-form editorial site from
`~/Downloads/index.html`: cinematic gravel-style hero + scrollable content
(About → Benchmark → Methodology → Architecture → Docs/FAQ).

Shipped as a single self-contained `dashboard/index.html` (inline CSS + JS).

## Brand overlays on the adopted file

- Seal mark rings `#F5F5F5`, stop-bar `#C23B22` (brand README — not the
  white stop-bar in the Downloads draft)
- Favicon + `theme-color`
- Leaderboard links in primary nav, mobile menu, and footer → `leaderboard/`

## Surfaces

| URL | File |
| --- | --- |
| `/` (`agentgavel.dev`) | `dashboard/index.html` |
| `/leaderboard/` | `dashboard/leaderboard/` |
| `/data/*` | `report --publish` target (unchanged) |

## Temporary hero video

Still using the CloudFront portal loop. Replacement prompt:

```
Cinematic dark void, near-black environment, subtle ground fog.
A solitary vertical aperture of cold white light stands in the distance —
a hard gate / chokepoint, not a doorway for comfort: a clean luminous slit
with a faint vermillion (#C23B22) horizontal stop-bar interrupting the light
midway, like a seal blocking passage. Thin concentric ring reflections on
wet ground. Slow push-in toward the gate; smoke drifts low; no people, no
faces, no gavel, no scales, no AI sparkles, no purple neon, no HUD text.
Photoreal, restrained, premium infrastructure mood. Seamless loop.
Camera locked, shallow depth, anamorphic softness at edges.
```
