# Session Handoff Protocol

Long-running session continuity: clean transitions across context boundaries via paste-ready resume messages.

> **Loading scope**: this digest is always-loaded because a handoff trigger — notably an explicit session-end request — can fire from any context. The complete format specification (6-block skeleton, per-block field rules, localization table, auto-resume flow, V0 abort gate) is `session-handoff-full.md`; examples, the full 4-locale table, and the anti-pattern catalogues are in `session-handoff-examples.md`. Both auto-load when this file is read. [HARD] Read the full text before composing a resume message.

## When To Generate (5 Triggers)

[ZONE:Evolvable] [HARD] The orchestrator MUST emit a paste-ready resume message when ANY of these activate: (1) context usage crosses the model-specific threshold (`context-window-management.md` § Context Window Targets); (2) a SPEC phase completes within a multi-SPEC workflow; (3) the user asks to end the session; (4) a PR is created while more SPECs remain in the current Epic; (5) a multi-milestone task reaches a stable checkpoint. When none apply, emit a brief completion confirmation.

[ZONE:Evolvable] [HARD] Every emission also pipes the cut-line-bounded block to `moai handoff save --stdin --spec <ID> --phase <phase> [...]` (flags: full text § Emission-Time Save Obligation). **Fail-open**: when the CLI is absent or the save fails, emit the paste-ready surface unchanged — the save never blocks, delays, or alters the handoff, and is never retried in a loop.

## Canonical Format (Verbatim Spec)

[ZONE:Evolvable] [HARD] The resume message is one fenced text block, **bounded by cut-line markers inside the fence** — top `✂──── 여기부터 복사 ────✂`, bottom `✂──── 여기까지 복사 ────✂`. The `✂` (U+2702) and `─` (U+2500) characters stay verbatim in every locale; only the marker text translates per `conversation_language`. Between the markers, six blocks:

1. `ultrathink. <SPEC-ID> <phase> <entering verb>.` — plus a `mode:` seed line only when the seeded mode is not `serial`
2. `applied lessons: <memory files>` and `source_session_id: <UUID>`
3. `Preconditions:` header
4. Numbered, independently verifiable preconditions `<N>) <action> → <expected outcome>` (at most 4)
5. `Run:` — a single work-starting primary action; never a bare `/moai goal`
6. `After merge:` (PR flow) or `Follow-up:` (trunk) — exactly one next action; omitted on a single-SPEC close with nothing queued

Every seed is a seed, not a permission grant: Implementation Kickoff Approval remains mandatory.

### Localization Table

Marker text and block headers translate per `conversation_language` (read from `.moai/config/sections/language.yaml`); ja / zh are in `session-handoff-examples.md`, and any other locale uses the English skeleton with naturally translated labels.

| Element | English | Korean |
|---------|---------|--------|
| Cut-line top / bottom text | `Copy from here` / `Copy to here` | `여기부터 복사` / `여기까지 복사` |
| Block 1 entering verb | `entering` | `진입` |
| Block 3 / 5 headers | `Preconditions:` / `Run:` | `전제 검증:` / `실행:` |
| Block 6 headers | `After merge:` / `Follow-up:` | `머지 후:` / `후속:` |
| Memory heading | `## Next Session Entry Point` | `## 다음 세션 시작점` |

### Field-by-Field Specification

The per-block rules — the `mode:` enum and its couplings (`fanout` appends `fan out subagents (<read-only investigation scope>)`, `agent-team` appends `--team`, `sweep` appends `ultracode`), the `source_session_id` fallback line, and the fixed Block 1 line order — are in the full text and `session-handoff-examples.md`.

## Worktree-Anchored Resume Pattern

[ZONE:Evolvable] [HARD] When the work happened inside a worktree, the resume message MUST prepend Block 0 (cwd anchoring) — `moai cc -w <name>`, `moai cc -w <abs-path>`, or `EnterWorktree(<path>)`, never a bare `cd` — and Block 4 gains `0) git rev-parse --show-toplevel → <worktree-path>`.

## Auto-Memory Integration (Mandatory)

[ZONE:Evolvable] [HARD] When generating a resume message, the orchestrator MUST also save it verbatim to a memory project entry (`project_<epic>_<spec>_<status>.md`) under a `## Next Session Entry Point (paste-ready resume message)` heading (locale variant per the table), add a one-line `MEMORY.md` index entry, and mark superseded entries `[SUPERSEDED by <new-file>]`.

## Output Surface (User-Facing)

[ZONE:Evolvable] [HARD] Emitting means **rendering** the cut-line-bounded block in the response body of the turn that generates it, together with the memory file path and a one-sentence summary of what the next session continues. Persisting it (memory file, `moai handoff save`) and merely citing the path is the named anti-pattern *reference-instead-of-render*.

### Pre-emit self-check (emission surface) — 3 items

- [ ] Is the cut-line-bounded block rendered in THIS response body — not only written to memory or persisted via the CLI?
- [ ] Are all three surface items present: the block, the memory file path, and the one-sentence continuation summary?
- [ ] Does the completion report avoid claiming the handoff was delivered when only the persistence steps ran?

## Diet Constraints

[ZONE:Evolvable] [HARD] A resume message is the next session's minimum executable context, not an audit trail: preconditions are one-line verifiable commands with a strict criterion (≤ 200 chars each), and Block 5 is one primary action, never a nested multi-phase plan. The V0 precondition uses lsof + cwd cross-validation; when it shows competing live sessions, spawning implementation agents is prohibited and the session ends (full text § V0 Abort Gate Doctrine).

## Auto-Injected Resume Flow (mode=auto)

Where `.moai/config/sections/handoff.yaml` sets `handoff.mode: auto`, the saved record is claimed and injected at the next `/clear` only, so the user sends one message; the injected preconditions are verified first. Under the default `manual` mode the record is inert. The manual paste path stays complete on its own, save and injection failures degrade to it silently, and Implementation Kickoff Approval is unchanged in both modes.

---

Status: HARD operational rule, applies to all multi-phase MoAI workflows. Digest; full text: `session-handoff-full.md` — the SSOT for the render surface in `.claude/output-styles/moai/moai.md` §8.
