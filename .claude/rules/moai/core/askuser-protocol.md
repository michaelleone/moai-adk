---
description: Canonical reference for AskUserQuestion-only interaction protocol, ToolSearch deferred-tool preload procedure, and Socratic interview standards
---

# AskUserQuestion Protocol — Canonical Reference

> This file is the **single source of truth** for AskUserQuestion interaction rules.
> Cross-referenced by: CLAUDE.md §8, moai-constitution.md §MoAI Orchestrator, agent-common-protocol.md §User Interaction Boundary, output-styles/moai/moai.md §3/§10.
>
> **Loading scope**: this digest is always-loaded because the orchestrator may compose an `AskUserQuestion` on any non-trivial turn. It restates every binding rule. The full text (worked patterns, self-check lists) is `askuser-protocol-full.md`; the recommendation evidence base, preview catalogue, and encoding root cause are in `askuser-protocol-reference.md`. Both auto-load when this file is read.

## Channel Monopoly

**AskUserQuestion is the only user-facing question channel.** Every question to the user — clarification, preference or decision, Socratic interview rounds, branch or workflow selection, conflict resolution — goes through an `AskUserQuestion` call. Free-form interrogative prose, including a prose question followed by a `- A: / - B:` option list, is **prohibited**. Exceptions: the tool is technically unavailable, or a status statement merely ends with a question mark. Users who want a free-form answer pick the automatic **"Other"** option.

## ToolSearch Preload Procedure

`AskUserQuestion` is a **deferred tool**: its schema is not loaded at start, and calling it unloaded fails with `InputValidationError`. Immediately before **every** call — again in each new turn — invoke `ToolSearch(query: "select:AskUserQuestion")`. Every deferred tool needs the same `select:` preload (`select:AskUserQuestion,TaskCreate` for several).

## Socratic Interview Structure

On a Stage 1 Clarify trigger, interview through sequential rounds (preload → ask; each round narrows on earlier answers; the last round confirms). Constraints: at most 4 questions per call and 4 options per question; the first option carries the `(권장)` / `(Recommended)` suffix; all text in `conversation_language`; never repeat a question; continue until intent is 100% clear; obtain explicit final confirmation before irreversible actions.

## Option Description Standards

Every option has a `description` that lets the user judge it without outside context: the immediate result, side effects and risks, and quantities where they apply (tokens, files, latency). Descriptions stay neutral — the recommendation is signalled **only** by the label suffix on the first option.

## Recommendation Placement Principles

The recommended option is the rational default the user has actually been observed to choose, never a policy default the system wants to push. Ask only when the outcome is genuinely uncertain; order questions by descending information gain; on a cold start fall back to the static default and disclose it in the description; state the precondition under which the recommendation holds ("Recommended when …"); weaken the recommendation for proficient users, and place no inferred-preference label while proficiency is unknown.

## Preview Field Standards

Use `preview` only when options differ structurally or quantitatively; it complements `description` and never replaces it. [HARD] **Single-select only** — `preview` is silently dropped under `multiSelect: true`. Keep it within about 12 lines and neutral in tone.

## Report-Before-Ask Gate

[ZONE:Evolvable] [HARD] A decision-type `AskUserQuestion` whose options derive from investigation results (agent fan-out returns, audits, verification batches, multi-source evidence) MUST be preceded, in the same turn's response body, by a substantive findings report: each source named with its quantified findings, every identifier used in the options explained, the findings in the body rather than only in option previews.

### Requested-Deliverable Primacy (user requirement analysis first)

[ZONE:Evolvable] [HARD] When the user's latest message asks for a report, analysis, or explanation, that deliverable IS the turn's terminal output: complete it and end the turn WITHOUT appending a decision-type `AskUserQuestion`. Pending pipeline questions surface in a later turn.

- [HARD] **Preview-as-report substitution** is prohibited: option previews and descriptions never carry the findings alone.
- [HARD] **Report-promise fulfillment**: a consolidated report promised earlier in the task is rendered before any later decision question.
- The gate does not apply to pure clarify rounds before any investigation, confirmation gates on already-reported context, blocker re-delegation rounds, or preference questions with no investigative basis.

## Orchestrator–Subagent Boundary

The `AskUserQuestion` interaction channel is **asymmetric** by design.

### Orchestrator Obligations

The orchestrator uses `AskUserQuestion` exclusively, preloads it before each call, collects preferences before delegating, and answers a subagent's blocker report with a question round and a fresh re-delegation (`agent-common-protocol.md` § Blocker Report Format / § Re-delegation Procedure).

### Subagent Prohibitions

- [ZONE:Frozen] [HARD] Subagents MUST NOT invoke AskUserQuestion
- [ZONE:Frozen] [HARD] Subagents MUST NOT output free-form prose questions directed at the user
- [ZONE:Frozen] [HARD] Subagents MUST NOT embed AskUserQuestion call syntax in their response body

## Ambiguity Triggers and Exceptions

This section is the **single source of truth** for Stage 1 Clarify trigger conditions (CLAUDE.md §7 Rule 5 and §8 refer here).

- **Triggers** (any one): a pronoun or demonstrative without a clear referent; a multi-interpretable action verb without scope ("clean up", "improve", "fix"); unclear boundaries; potential conflict with existing state.
- **Exceptions**: a single-line typo or format fix; a bug fix with a reproduction; a file read at an explicit path; a command with all arguments given; continuation of already-confirmed work.
- **Unknowns lens**: known-unknowns → interview round; unknown-knowns → `Agent(Explore)` reconnaissance, then confirm; suspected unknown-unknowns → Blind Spot Pass.

## Blind Spot Pass

An optional pre-plan technique for a user working in unfamiliar territory: a read-only `Agent(Explore)` scan of the domain, with the likely unknown-unknowns surfaced through one `AskUserQuestion` round before the SPEC is authored. It is a judgement call, not an automatic gate, and the subagent never prompts the user itself.

## Free-form Circumvention Prohibition

No free-form question, markdown option list, or trailing "Should I …?" line substitutes for `AskUserQuestion`.

### Completion-Report Next-Step Discipline

[ZONE:Evolvable] [HARD] A completion report MUST NOT end with a free-form prose next-step question, in any `conversation_language`. It has exactly two valid closes: route a genuine next-step decision through `AskUserQuestion`, or close with no question at all.

## Non-ASCII Tool-Call Encoding

Text in `conversation_language` inside any tool-call payload — `AskUserQuestion` fields, Bash commands, Write/Edit content — MUST be written as native UTF-8; hand-authored `\uXXXX` escapes are **PROHIBITED** (a malformed escape corrupts the JSON into `InputValidationError`). After an `Invalid tool parameters` rejection of a non-ASCII payload, re-author it from the intended text as native UTF-8 rather than repairing the escapes; persistent recurrence warrants `/clear` with a paste-ready resume.

---

Version: 1.3.0
Classification: Canonical Reference — do not duplicate content; cross-reference this file instead. Digest; full text: `askuser-protocol-full.md`.
