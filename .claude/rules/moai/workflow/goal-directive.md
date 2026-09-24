# Goal Directive (`/moai goal`) — Autonomous Continuation

`/moai goal "<condition>"` registers a session-scoped completion condition — mechanical conditions decided by a shell command's exit code, model conditions demonstrated in the transcript — and the `moai hook stop-goal` Stop hook blocks turn-end until it holds or a bound fires. Bounds: a turn ceiling (default 30, ending in a 5-section verdict), the runtime consecutive-block cap (`CLAUDE_CODE_STOP_HOOK_BLOCK_CAP`, default 8 — the effective bound is the smaller, and a missing verdict is not convergence), and a stagnation guard. Verbs: `/moai goal "<condition>"`, `/moai goal status [--all]`, `/moai goal clear`. It needs hooks enabled.

> This digest carries the binding rules. Full text: `goal-directive-full.md`; condition templates, the comparing-approaches table, and the native `/goal` prohibition rationale: `goal-directive-detail.md`; verb surface and progression modes: `.claude/skills/moai/workflows/goal.md`. Read them before arming a goal.

## Goal-Presentation Timing

**`/moai goal` is arm-only.** Arming starts no work: a goal armed while nothing runs spins idle turns to the ceiling, so arming is always paired with a work-starting action and never substitutes for one. A paste-ready resume therefore keeps the work-starting command (`/moai run SPEC-X`) as Block 5's primary action (`session-handoff.md` § Canonical Format).

**The goal is presented at the Implementation Kickoff Approval gate** as the autonomous vs semi-autonomous progression-mode axis, and armed only after the gate passes. Arming never authorizes run-phase entry, a PR, or a destructive operation, and never relaxes the confirm-before-hard-to-reverse boundary.

## Proactive Recommendation Triggers

Arm a goal rather than driving the work turn by turn when a condition-declared loop is the right continuation: T1 long run-phase or multi-milestone work (`run.md` § Run-phase Autonomy `ac_converge` owns it), T2 migrations across many enumerated call sites, T3 TDD or acceptance-criteria convergence, T4 work better expressed as a verifiable end-state than as `/moai loop`'s "fix what the tooling flags".

---

Classification: Evolvable orchestration guidance — applies to autonomous multi-turn continuation. Digest; full text: `goal-directive-full.md`.
