# Verification-Claim Integrity

Doctrine establishing the **"no unobserved-verification-claim" invariant** for all MoAI actors. It is a policy-layer norm: a runtime layer may detect one shape of violation, but the norm binds every actor regardless.

> This digest carries the binding rules. The full text — each binding surface in detail, the report-section definitions, and two worked hazard examples — is `verification-claim-integrity-full.md` (auto-loads when this file is read).

## 1. The Invariant — no unobserved-claim (verification, defect, OR premise)

[ZONE:Evolvable] [HARD] An actor MUST NOT assert a verification, a completion, **a defect / debt / drift, OR the premise underlying a recommendation** it did not actually verify with the domain's mechanical tooling.

> **Evidence absent ≠ evidence of success — NOR of failure.**

- A pass is valid only when the actor ran the command and observed its output; an unran command, a skipped step, or a silent assumption is a gap, never a pass.
- A defect inferred from text patterns, grep matches, or file absence is a hypothesis until the domain's dedicated tool (audit, lint, type-check, coverage) confirms it.
- **Reachability is not justification.** A reference existing does not prove the referent is live, and an originating task still reading as in-service does not prove its feature survived. Before recommending retention — above all against a user's instruction — verify that the producer still exists and that no completed retirement covers it.

### 1.1 Binding scope — ALL FOUR surfaces

1. **Orchestrator self-report** — completion reports and verification matrices: every PASS row is an actually-observed command output.
2. **Manager-agent completion report** — each reported result is the verbatim output of a command the agent actually ran.
3. **Defect / debt / drift identification claim** — valid only on the output of the domain's dedicated tool.
4. **Recommendation-premise claim** — the stated reason for or against an action carries the same evidence burden as a defect claim.

## 2. Baseline-Integrity Attribution

[ZONE:Evolvable] [HARD] Every verification claim MUST be attributed to an actually-measured baseline — the command that was run plus the output that was observed, in this run, against this tree. A figure carried over from another task, package, or point in time is not a baseline; report it as a Gap.

## 3. The 5-Section Evidence-Bearing Report Format

[ZONE:Evolvable] [HARD] Verification and completion reports — on every binding surface — SHOULD carry these five sections:

### 3.1 Claim

What is asserted, one discrete claim per row or sentence.

### 3.2 Evidence

The command that was run **plus its verbatim output** — never a summary.

### 3.3 Baseline-attribution

What the claim was measured against, per §2.

### 3.4 Gaps

What was NOT observed. An empty Gaps section asserts that nothing was left unobserved; when in doubt, name the gap.

### 3.5 Residual-risk

What could still be wrong despite the evidence (flaky tests, environment-specific behavior, deferred criteria).

---

Version: 1.2.0
Classification: Canonical Reference (policy-layer codification) — do not duplicate cross-referenced content; cross-reference this file instead. Digest; full text: `verification-claim-integrity-full.md`.
