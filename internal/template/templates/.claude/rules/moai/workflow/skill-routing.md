# Skill Routing Protocol

How the orchestrator routes domain skills into agent spawns and into its own direct work. This digest carries the binding rules; the full text (rationale, the report-format config coupling, and the cost model) is `skill-routing-full.md`.

## 1. Orchestrator Obligation

[ZONE:Evolvable] [HARD] Before spawning an implementation or review `Agent()`, match the mission's domain against the `moai-ref-*` / `moai-domain-*` skill descriptions and inject one line per matched skill (0-3 matches) into the spawn prompt: `At start, invoke Skill("<name>") for <reason>.` Zero matches is a valid outcome.

### §1.1 — Orchestrator-Direct Skill Routing (non-spawn)

[ZONE:Evolvable] [HARD] When the orchestrator produces an artifact **directly** whose shape matches a domain skill, it MUST load that skill with `Skill()` first:

| Orchestrator-direct task | Mandatory skill |
|---|---|
| Report / markdown → HTML artifact | `moai-domain-html-report` (skipped when `report.format` is `md`) |
| Humanize / post-edit AI text | `moai-domain-humanize` |
| SVG infographic / architecture diagram | `moai-domain-svg-infographic` |
| Capture or reproduce a reference design | `moai-domain-design-dna` |
| Chart / dashboard to HTML or SVG | `dataviz` |
| Hosted visual-identity page on claude.ai | `artifact-design` |

[ZONE:Evolvable] [HARD] Routing is intent-based, not keyword-based, in every conversation language. Loading `artifact-design` for a markdown→HTML **report** render is a named routing miss — reports route to `moai-domain-html-report`.

## 2. Agent Obligation

Agents with the `Skill` tool load conditional skills when their body's stated trigger situation actually arises; the static `skills:` frontmatter preload stays at most 2 entries per agent.

---

Classification: Evolvable operational rule — applies to all agent spawns and agent bodies with the Skill tool. Digest; full text: `skill-routing-full.md`.
