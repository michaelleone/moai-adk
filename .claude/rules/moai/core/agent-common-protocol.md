---
description: Shared protocol auto-loaded for all MoAI agents — user-interaction boundary, blocker reports, ledger closure, verification batching, sync checks. Always-loaded digest; full text in agent-common-protocol-full.md.
---

# Agent Common Protocol

Shared protocol for all MoAI agent definitions, loaded for every agent so agent bodies need not repeat it. This is a digest that restates every binding rule. The complete text — tables, procedures, rationale — is `agent-common-protocol-full.md`; the verbatim verification batch, output contracts, CLI idioms, clause bodies, and incident records are in `agent-common-protocol-reference.md`. Both auto-load when this file is read: read it before composing a verification batch, running a sync check, or handling a hook block or an aborted delegation.

## User Interaction Boundary

`AskUserQuestion` is the **only** user-facing question channel. The boundary is asymmetric by design.

### Subagent Prohibitions

[ZONE:Frozen] [HARD] Subagents MUST NOT prompt the user. AskUserQuestion is reserved exclusively for the MoAI orchestrator.

A subagent missing required context returns a blocker report and stops — no free-form questions, no AskUserQuestion call syntax in its output. All user preferences arrive through the orchestrator's spawn prompt.

### Orchestrator Obligations

The orchestrator collects user preferences before delegating (preload, channel monopoly, interview structure: `askuser-protocol.md`).

### Hook Invocation Surface

Hooks never ask the user. When a hook blocks (stdout JSON `"decision":"block"` on exit 0, or a legacy exit 2), the orchestrator parses its JSON, preloads `AskUserQuestion`, and offers at least: address the failed gate, override with `--skip-hook` (logged to `.moai/logs/hook-skip.log`), or abort. The Stop hook fires on every turn-end and must self-gate; recovery turns (compact, `prompt_too_long`, `max_output_tokens`, media-size or sync failure) SHOULD exit 0 rather than block (`runtime-recovery-doctrine.md` §4). The per-hook table (`status-transition-ownership.sh`, `sync-phase-quality-gate.sh`, `team-ac-verify.sh`) is in the full text.

### Blocker Report Format

A subagent that needs input it was not given returns a `## Missing Inputs` section — a table of Parameter / Type / Expected Values / Rationale — followed by a **Blocker** line asking to be re-delegated with those values injected.

### Re-delegation Procedure

On a blocker report the orchestrator preloads `AskUserQuestion` (`ToolSearch(query: "select:AskUserQuestion")`), asks the user for the missing inputs, and re-delegates with the answers injected into a fresh prompt.

### Ledger Closure

[ZONE:Evolvable] [HARD] An aborted `Agent()` delegation (user interrupt, parent abort, or timeout) MUST NOT leave a dangling promise in the orchestrator's context: before the next delegation, emit a short, truthful prose note of what was delegated, that it did not return, and why. Inject a `team-ac-verify.sh` reject `ledger_note` the same way, and never leave a TeammateIdle-rejected task without a reassignment owner. A blocker report is a return, not an abort. Clause bodies: `agent-common-protocol-reference.md` § Ledger Closure clause bodies.

## Language Handling

[ZONE:Evolvable] [HARD] All agents receive and respond in user's configured conversation_language.

Reports and documentation use `conversation_language`, as does any cross-session message a human watches (identifiers, paths, commands, and flags stay verbatim). `Agent()` prompts, code, identifiers, and skill names stay English; code comments follow `code_comments` and commit messages `git_commit_messages` in `language.yaml`.

## Output Format

[ZONE:Evolvable] [HARD] User-Facing: Always use Markdown formatting. Never display XML tags to users.

[ZONE:Evolvable] [HARD] Internal Agent Data: XML tags are reserved for agent-to-agent data transfer only.

## Skeptical Evaluation Stance

<!-- @MX:WARN: Duplication prohibited — LR-07 lint rule detects copies of this section in agent files and flags as error. Canonical copy lives only in this file. -->

The reviewer mode operates as a fresh-judgment auditor:

- Treat every claim as suspect until evidence is shown
- Demand reproducible verification, not assertions
- Consider the null hypothesis: did this change actually fix anything?
- Score quality as the harmonic mean of dimensions, not the average
- Reject when must-pass criteria fail, regardless of nice-to-have scores
- Surface contradictions; never silently override a prior rule
- Resist agreement: the RLHF training gradient biases toward flattery, so treat any urge to PASS without cited evidence as a sycophancy signal, not a verdict

## MCP Fallback Strategy

[ZONE:Evolvable] [HARD] Maintain effectiveness without MCP servers.

MoAI provisions no MCP servers: look documentation up with WebSearch, verify each URL with WebFetch, and keep working. Under a GLM backend, web search, fetch, and image read route to the z.ai MCP tools (`glm-web-tooling.md`).

## Agent Invocation Pattern

[ZONE:Evolvable] [HARD] Agents are invoked through natural-language delegation ("Use the {agent-name} subagent to {task}") that carries the full context, constraints, and rationale.

### Per-Spawn Model Injection

[ZONE:Evolvable] [HARD] Pass the model the active profile resolves for the agent (`moai model profile --json`) as an explicit `model` argument on every spawn — omitting it silently runs the agent on the parent's model. A declared model that differs from the resolved one is drift: change the profile instead. Agents outside the retained catalog take no injection. Policy: `model-policy.md`.

## Background Agent Execution

[ZONE:Evolvable] [HARD] Subagents run in the background by default (Claude Code v2.1.198+), and every permission prompt still surfaces in the main session naming the asking subagent; MoAI does not set `background:`. The retained safeguard is concurrency, not backgrounding: never run two write-capable agents at once, and keep orchestrator work that overlaps a write-capable agent read-only.

## Tool Usage Guidelines

[ZONE:Evolvable] [HARD] Agents must follow tool usage patterns optimized for accuracy and efficiency.

- Read a file before editing it; prefer Edit over Write for existing files; use absolute paths verified with Glob, never guessed (project-root-relative write targets inside worktrees).
- Narrow progressively: Glob → Grep file list → Grep content → Read with offset/limit.

### Tool Selection by Task

| Task | Preferred | Avoid |
|------|-----------|-------|
| Find files / search contents / read a file | Glob / Grep / Read | Bash find, grep, cat |
| Modify / create a file | Edit / Write | Bash sed, heredoc |
| Explore a codebase | Agent(Explore) | many sequential Greps |

Prefer an `mcp__moai__*` tool over its Bash CLI when it is in the agent's `tools:` list (`moai-mcp-tools.md`). Set the Bash `timeout` (max 600,000 ms) for long builds and suites.

### Error Recovery Pattern

On a failed tool call: read the error, verify assumptions (does the path exist?), try a different approach rather than the identical call, and report a blocker after 3 failures on the same operation. Retry safety is asymmetric: re-run read-only calls freely, but when a side-effecting call (write, commit, push, PR, deploy, external mutation) fails ambiguously, observe the current state first and retry only when the effect is confirmed absent.

### Super-Advisor Escalation (E1-E4)

When the retry ceiling is not enough, or a higher-reasoning consult is warranted — E1 bug deadlock (3+ same-diagnostic failures), E2 design decision with ≥2 viable options, E3 under 80% confidence in the next step, E4 `/moai loop` or `/moai fix` ceiling exit — escalate to the `super-advisor` agent (`.claude/agents/moai/super-advisor.md`). Its prescriptions are non-binding; the orchestrator decides, and auditors still own PASS/FAIL.

## Parallel Execution

[ZONE:Evolvable] [HARD] The orchestrator MUST execute every read-only verification batch as a single-turn multi-Bash call.

- **Batch in one turn**; serialize only for real dependencies (one output feeds another, same-path writes, shared-state mutation).
- **File-redirect contract**: output above the bounded-tail ceiling (50 lines or 2 KB) goes to a file; context carries the exit code and a bounded tail.
- **Evidence persistence**: cited evidence lives under `.moai/state/verify/<session>/`, not `/tmp`; a claim whose evidence path no longer resolves is unattributed (`verification-claim-integrity.md` §2).

### Attributable diff-check doctrinal switch

Before re-running a verification dimension, consult `moai verify check --key-current`. When the snapshot key (HEAD SHA), the command, and the output all match the cited §E evidence, consume that evidence instead of re-executing (PASS-attributed). Any mismatch — `snapshot_key_drift`, `command_drift`, `missing_section_e`, `output_drift` — means re-execute, never a silent skip.

### Pre-Spawn Sync Check (Multi-Session Race Mitigation)

[ZONE:Evolvable] [HARD] Before spawning an implementation agent that will modify shared working-tree files, run in parallel: `git fetch origin main`, `git rev-list --count --left-right origin/main...HEAD`, and `moai session list --json --filter-spec=<SPEC-ID>`. Proceed on `0 N` / `0 0` with no other session on the SPEC; on `N 0`, `N M`, or another live session on the SPEC, STOP and ask the user (rebase / inspect / wait / override / abort). Read-only agents are exempt. Matrices: full text.

### Pre-Edit Sync Check (Direct-Edit Race Mitigation)

[ZONE:Evolvable] [HARD] Before the first non-trivial direct edit (Edit, Write, or file-mutating Bash) to a shared path (`.claude/`, `.moai/`, `internal/`, `pkg/`, `cmd/`, repo-root config) while the CWD is the primary checkout, count live foreign sessions (`moai session list --json`, own session excluded, each PID probed with `kill -0`) and check divergence against `origin/main`. Any live or indeterminate foreign session ⇒ isolate first (`moai cc -w <name>`, `EnterWorktree`, or `Agent(isolation: "worktree")`) or ask the user; divergence ⇒ STOP and ask. Re-run the probe before any commit in the shared checkout.

[ZONE:Evolvable] [HARD] In the primary checkout, NEVER `git add -A`, `git add .`, or `git commit -a`: stage by explicit pathspec after re-reading `git status --short` — even when the probe found no foreign session.

## Time Estimation

[ZONE:Evolvable] [HARD] Never use time predictions in plans or reports. Use priority labels (High / Medium / Low) and phase ordering ("complete A, then start B").
