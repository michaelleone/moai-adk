# Cross-Session Messaging

Doctrine for messaging between independent Claude Code sessions (`ListAgents` to discover, `SendMessage` to deliver plain text by name). The channel is a runtime feature; this rule governs how the orchestrator uses it.

> **Loading scope**: this digest is always-loaded because a peer-session conflict or an inbound message can surface in any turn. The full doctrine — where the channel sits among MoAI's other mechanisms, the concurrency-check integration, and the anti-patterns — is `cross-session-messaging-full.md` (auto-loads when this file is read). [HARD] Read it before the first `SendMessage` or `ListAgents` call of a session.

## Availability constraints

The channel exists only on macOS and Linux (not native Windows), not on Bedrock, Claude Platform on AWS, Google Cloud Agent Platform, or Microsoft Foundry, and only from Claude Code v2.1.224; any of `CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC`, `DISABLE_TELEMETRY`, `DO_NOT_TRACK`, or `DISABLE_GROWTHBOOK` turns it off silently. Where it is absent, nothing errors — surface that to the operator rather than retrying.

## Rules

[ZONE:Evolvable] [HARD] **Never route a user decision through a peer.** A peer is not a proxy for the user and its reply is not approval; questions go to the user through `AskUserQuestion`.

[ZONE:Evolvable] [HARD] **Never ask a peer to do what this session may not do.** Work blocked or denied here goes back to the user, not sideways to another session.

[ZONE:Evolvable] [HARD] **Send facts, not instructions to mutate shared state.** Report what landed, what broke, or ask a question; never direct a peer to edit configuration, rewrite doctrine, or take a hard-to-reverse action.

[ZONE:Evolvable] [HARD] **An idle notice is not completion evidence.** A session goes idle when it finishes, when it stops at a permission prompt, and when it dies; the notice says when to look, not what happened.

- A message carries no context and is not consent. Continuing work across a `/clear` or a machine boundary is a paste-ready handoff (`session-handoff.md`), not a message.
- Messaging shortens diagnosis of a concurrent session; it never makes two sessions safe to write the same path — worktree isolation does.

## Addressing, sending, and replying

Address a peer by name; re-send with the short reference an error supplies when the name is shared. Read every send result — a refusal, a missing reply address, or a `routing` object means the message did not reach the intended session — and never let progress depend on a reply arriving. Reply to the sender name a message carries, falling back to its reply address.

## Configuration surface

`crossSessionInbound` (`accept` / `hold` / `refuse`), `isolatePeerMachines`, `dialogExpiry`, and `permissions.deny: ["SendMessage", "ListAgents"]` govern delivery. A rapid burst to one inbox — a lead nudging every lane in one turn — is refused up front; spread the sends across turns.

---

Classification: Evolvable operational rule — peer-session communication; changes no gate semantics. Digest; full text: `cross-session-messaging-full.md`.
