# Context Window Management

Long-horizon session continuity guidance for both users and the MoAI orchestrator. Response streams stall (`stream_idle_partial`) near the context-window ceiling — intermittently, but predictably above the model-specific threshold.

> This digest carries the binding rules. The full text — Claude Code's graduated-compaction layers, the GLM context-window caveat, the detection heuristics, and the two-stage statusline marker — is `context-window-management-full.md` (auto-loads when this file is read).

## Context Window Targets

[ZONE:Evolvable] [HARD] Operational threshold is **model-specific**:

| Model class | Window | Handoff threshold |
|-------------|--------|-------------------|
| Opus 5 / Opus 4.8 (1M), GLM-5.3 via `moai glm` / `moai cg` (1M) | 1,000,000 tokens | **50%** (~500,000) |
| Fable (256K) | 256,000 tokens | **90%** (~230,000) |
| Sonnet / Opus standard / Haiku (200K) | 200,000 tokens | **90%** (~180,000) |

Beyond the threshold, plan a `/clear` before the next non-trivial action; `session-handoff.md` Trigger #1 reads this same table. A GLM session is a 1M session even when raw telemetry reports ~180K — trust the MoAI statusline.

## Reduction Ladder — cheaper moves before `/clear`

`/clear` discards the warm cache and re-pays the whole always-loaded prefix, so pick the cheaper move by cause first: `/btw <question>` for a side question, `/compact <instructions>` when the current task must continue, `/rewind` → summarize to keep recent turns verbatim, `/rewind` → restore a checkpoint to undo a polluted line of attempts. Continue a session that still exists with `claude --continue` / `--resume`. `/clear` remains mandatory at the thresholds below.

## User Responsibilities

[ZONE:Evolvable] [HARD] When usage crosses the model-specific threshold: save in-flight state to `.moai/specs/<SPEC-ID>/progress.md` (the orchestrator does this automatically), run `/clear`, then paste the resume message.

[ZONE:Evolvable] [HARD] When usage crosses 95% on any model, the next action MUST be `/clear` — no further large work in the current session.

## Orchestrator Responsibilities

[ZONE:Evolvable] [HARD] Pre-clear announcement: approaching the model-specific threshold, the orchestrator MUST stop initiating large tool calls and `Agent()` delegations, persist progress to `progress.md`, emit a paste-ready resume message in the `session-handoff.md` format so the next session is self-sufficient, and recommend `/clear` as a status announcement (no `AskUserQuestion` required).

Estimate usage state-file-first: `.moai/state/context-usage.json` (`raw_pct`, `stage`) when its session id matches this session; otherwise fall back to heuristics (output volume, large tool results, returned `Agent()` calls), under-estimating when uncertain.

---

Status: HARD operational rule, applies to all sessions. Digest; full text: `context-window-management-full.md`.
