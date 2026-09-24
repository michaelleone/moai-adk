# Cache-Aware Execution

Prompt-caching-aware ordering rules for orchestrator execution. Anthropic prompt caching is a **prefix match** over the rendered request (`tools` → `system` → `messages`): cache reads cost ~0.1× the base input price, writes 1.25×, and the 5-minute TTL is **idle-based** — a gap over 5 minutes (typically a blocking `AskUserQuestion` wait) expires the cache and the next turn re-writes the whole prefix. These rules govern WHEN and IN WHAT ORDER the orchestrator acts; they change no gate semantics.

> **Loading scope**: this digest is always-loaded because the directives bind ordering decisions on any non-trivial turn. Full directive text: `cache-aware-execution-full.md`; cited numbers and worked examples: `cache-aware-execution-reference.md`.

## Directives

1. **Front-load user gates** [ZONE:Evolvable] — ask intent-drain gates early, in the small-context Clarify stage; batch unavoidable late gates into consecutive rounds so the expiry is paid once.
2. **Stagger-spawn parallel same-type agents** [ZONE:Evolvable] — spawn one, then the remaining N−1 once it is producing output, so they read its cache instead of all paying the cold write.
3. **Defer session-loaded file edits to task end** [ZONE:Evolvable] — edits to `.claude/rules/`, `CLAUDE.md`, output styles, or always-loaded skills invalidate the prefix; batch them at the end of a task or just before a `/clear`.
4. **Consider `/clear` before large batches** [ZONE:Evolvable] — when a large multi-spawn batch is next and the context is bloated with finished unrelated work.
5. **Inherit the session model on spawns** [ZONE:Evolvable] — caches are model-scoped; override only when the task needs another tier.
6. **Pass files by `@`-mention, not by name** [ZONE:Evolvable] [HARD] — one deterministic load beats a fetch-retry cycle.
7. **Keep command output bounded** [ZONE:Evolvable] [HARD] — quiet flags, targeted queries, or redirect-to-file with the exit code and a bounded tail.
8. **Prefer the quiet form of routine commands** [ZONE:Evolvable] [HARD] — no spinners, banners, or color noise.
9. **Weigh session length as a cost axis** [ZONE:Evolvable] [HARD] — every fresh session re-pays the always-loaded prefix at write price; justify a split like a `/clear`.
10. **A mid-session model or effort switch busts the cache** [ZONE:Evolvable] [HARD] — switch at a natural boundary.

## Non-goals

These directives NEVER justify skipping, weakening, or reordering an approval gate's semantics — only the placement and batching of questions. The orchestrator does not place `cache_control` markers; it controls ordering, spawn timing, and edit timing.

---

Classification: Evolvable operational rule — execution ordering only; gate semantics unchanged. Digest; full text: `cache-aware-execution-full.md`.
