# AgentGavel marketing site — DESIGN.md

## Direction

Adapt the gravel.html cinematic stage (dark void, Manrope, white pills,
full-bleed looping video, measured hero layout) for AgentGavel while keeping
the approved brand sacred: seal + vermillion stop-bar `#C23B22`, near-white
rings on dark, no gavel illustration, no purple, no glow orbs.

## Locked choices (2026-09-05)

1. Marketing site with scrolled sections; leaderboard linked at `/leaderboard/`
2. Brand sacred — seal/lockup + stop accent; template craft for layout/motion
3. Keep the portal video as temporary hero plate; replace with AgentGavel-
   specific footage when ready (prompt below)
4. Ship HTML directly (no mockup-only pass)

## Surfaces

| URL (Pages) | File |
| --- | --- |
| `/` | `dashboard/index.html` + `site.css` / `site.js` |
| `/leaderboard/` | `dashboard/leaderboard/` (existing Opt-in / Unratified UI) |
| `/data/*` | Unchanged publish target for `report --publish` |

## Tokens

- Stage black `#050505`, ink `#fafafa`, muted `#a7a6a6`, nav `#b6b5b5`
- Stop / Hard accent `#C23B22` (one saturated hue; Hard pillar + SEC ids)
- Pill CTA white on black; ghost text secondary
- Type: Manrope 400/500/600

## Composition

- First viewport = one composition: seal mark, nav, one headline, one sub,
  CTA pair, full-bleed video, bottom strip. No cards in the hero.
- Below the fold: Why → Hard/Soft → Suites → Neutrality → CTA — each section
  one job, typography + hairlines, no marketing card grid.

## Video replacement prompt

Use with Seedance 2.0 / Higgsfield (or equivalent cinematic text-to-video).
Generate a seamless loop, ~8–12s, 16:9 or wider, dark, no UI chrome, no logos
burned in (we overlay the brand in HTML).

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

Negative cues: purple glow, cyberpunk city, robot, judge gavel, courtroom,
busy particles, UI overlays, readable text, cartoon style.

After generation: host the MP4 (CloudFront or `dashboard/media/hero.mp4`) and
swap the `<source>` in `dashboard/index.html`. Keep the existing plate fades.
