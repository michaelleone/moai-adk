# Native-Idiom & Register Policy (Non-English Locales)

> The language-quality invariant for non-English output. This digest carries the binding rules; the full text — why calques survive, the register table, the illustrative hazard list, and the pre-emit self-check — is `native-idiom-and-register-full.md`. Cross-references `moai-constitution.md` § Response Language and the `moai-domain-humanize` skill.

## The Invariant

[ZONE:Evolvable] [HARD] When `conversation_language ≠ en`, every user-facing surface — chat replies, reports, README, docs-site, generated sites, `AskUserQuestion` text — MUST read as natural native prose, NOT as English mapped one-to-one onto the target language. Translation-style calques (word-for-word carry-over of English syntax, metaphor, and figurative stock such as "pillars", "axes", or "budget defense") are prohibited; native idiom is required. The policy is conditional: English sessions are unaffected.

**Two registers.** Chat uses the colloquial native register; artifacts (reports, README, docs-site) use clean native written register — professional, never colloquial, never calqued. Coined brand terms and established loanwords are not calques.

## Mechanism — when to invoke humanize

[ZONE:Evolvable] [HARD] Heavy non-English artifacts (multi-paragraph reports, README rewrites, docs-site pages, generated sites) MUST pass through the `moai-domain-humanize` skill as a final phase before delivery, scoped to the active locale's module. Single-turn chat replies apply this rule inline.

---

Classification: Always-loaded language-quality invariant — non-English conditional. Digest; full text: `native-idiom-and-register-full.md`.
