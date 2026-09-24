# Kanban Dispatch Protocol

How the **lead** session of Kanban Mode moves a card across the board: what admits work, who is told to do it, how completion is judged, and when the operator is asked to `/clear`.

> **Loading scope**: this digest is always-loaded because a session learns it is in Kanban Mode from its SessionStart context, not from a file path. The protocol itself — every [HARD] rule, the dispatch format, the worktree, verification, and integration rules — is `kanban-dispatch-full.md`; tables, walkthroughs, and rationale are in `kanban-dispatch-detail.md`. Both auto-load when this file is read, and cost nothing outside Kanban Mode.

## Scope — when this rule is live

Inert unless the SessionStart context declares **Kanban Mode** (`moai cc -k` / `moai glm -k`: one lead plus `plan` / `run` / `sync` companion sessions) or **Factory Mode** (`moai cc -f <N>`: one lead plus `lane-1..lane-N`).

[HARD] A session in either mode — lead, companion, or lane — MUST Read `.claude/rules/moai/workflow/kanban-dispatch-full.md` before its first card action (admitting, promoting, dispatching, working, verifying, integrating, or closing a card), and again after every `/clear`.

## Entry into the board is an operator act

The operator is the only source of work: the lead turns operator requests into cards (`moai todo add`) and never invents one, and promotion out of `backlog` is always the operator's pick through `AskUserQuestion`.

### The delegation channel is the queue

The queue on disk carries the delegation. A cross-session message is only a nudge, and an idle notice is not a completion signal.

## Completion is read, never trusted

A card advances only on evidence the lead read — the card's `progress.md` and its declared evidence path — never on a reply. The final verdict is the lead's.

## Isolation is entered, never provisioned

One card per worktree, entered through the launcher (`moai cc -w <name>`, `EnterWorktree`) and never with a bare `git worktree add`; no worktree is disposed before its branch has merged on the remote. Lane-local verification covers only what the card's change can affect — CI runs the full suite — and never leaves background load running.

## Integration into the release branch is self-served

A lane whose card passed verification merges its own branch into the batch's release branch, entering the release worktree to do it; the batch pull request stays with the lead.

## Boundaries — what this protocol does not do

No session spawning, no question delegation, no gate bypass: the lead addresses sessions the operator launched and asks the operator itself.

---

Classification: Evolvable operational rule — applies to Kanban and Factory Mode sessions. Digest; full text: `kanban-dispatch-full.md`.
