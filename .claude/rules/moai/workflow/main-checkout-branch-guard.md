# Main-Checkout Branch Guard

Branch-state isolation rules for the primary project checkout. The checkout is **shared** — several sessions, teammates, hooks, and background tools can operate on the same working tree at once — so a `git switch` in one session changes what every other session sees, mid-operation, with no signal to either side.

> **Loading scope**: this digest is always-loaded because the guard binds any turn that performs git work. The full text — why the race is silent, verification commands, and the opt-in PreToolUse enforcement (deny semantics, exemptions, fail-open norm) — is `main-checkout-branch-guard-full.md`.

## Rules

[ZONE:Evolvable] [HARD] The orchestrator MUST NOT change branch state in the primary project checkout. Forbidden there: `git checkout <branch>` / `git switch`; `git checkout -b` / `git switch -c` / `git branch`; `git reset --hard` / `git checkout -- <path>`; `git stash` (repository-global — it absorbs other sessions' uncommitted work); `git rebase` / `git merge` onto the checked-out branch.

Permitted there: read-only inspection (`git status`, `log`, `diff`, `rev-parse`, `show`, `branch -vv`), `git fetch`, commits to the branch already checked out staged by explicit pathspec (never `git add -A`), and `git push` of that branch.

Work that needs a different branch goes in a worktree — `git worktree add -b <branch> <worktree-path> origin/main` — driven with `git -C <worktree-path>` rather than `cd`, and removed once merged.

## Staleness Rule

[ZONE:Evolvable] [HARD] Re-read branch and commit state (`git rev-parse --short HEAD`, `git branch --show-current`) **immediately before** any commit or push — never rely on a value read earlier in the turn or on session-start context. If either differs from what the turn assumed, stop and report the divergence.

Treat concurrency as the default: an empty or stale session registry does NOT establish that no other session is active. A `BRANCH_GUARD_VIOLATION:` deny means move the work to a worktree; it does not mean the exemption is broken.

---

Classification: Evolvable operational rule — branch-state isolation; changes no gate semantics. Digest; full text: `main-checkout-branch-guard-full.md`.
