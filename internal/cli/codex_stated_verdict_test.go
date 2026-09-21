package cli

import "testing"

// Backport of upstream f505955a9 (card t178), part (a) only: codex's verdict
// synthesis honours a verdict the review states in its own words. On v3.1.2 the
// synthesis read only the review-mode "- [P1]" finding bullets, so a review
// whose prose states that it FAILS, with no finding bullet, was synthesized as
// "pass" (the U1 defect recorded by SPEC-ADAX-V6-READINESS-001).

// The three qualification cases, verbatim from
// SPEC-ADAX-V6-READINESS-001 evidence/fixtures/verdict-normalisation/.
const (
	u1PositiveCase = "# Review of change set `parse-guard`\n\nVerdict: PASS\n\n" +
		"Scope read: the diff to `internal/example/parse.go` and `internal/example/parse_test.go`.\n\n" +
		"The change guards the slice index in `ParsePair` with a length check before it reads the second\n" +
		"field, and returns a wrapped error that every caller already handles. The new table test covers the\n" +
		"empty input, the input without a separator and the ordinary pair. No correctness, security or scope\n" +
		"issue was found, and the change is sound as written.\n"

	u1BulletNegativeCase = "# Review of change set `parse-guard`\n\nVerdict: FAIL\n\n" +
		"Scope read: the diff to `internal/example/parse.go` and `internal/example/parse_test.go`.\n\n" +
		"- [P1] `internal/example/parse.go:42` reads `parts[1]` without checking `len(parts)`, so an input\n" +
		"  with no separator panics at runtime.\n" +
		"- [P2] `internal/example/parse_test.go` has no case for the input without a separator, so the panic\n" +
		"  above is not caught by the suite.\n\n" +
		"The change must not merge until the P1 finding is fixed.\n"

	u1NegativeCase = "# Review of change set `parse-guard`\n\nVerdict: FAIL\n\n" +
		"Scope read: the diff to `internal/example/parse.go` and `internal/example/parse_test.go`.\n\n" +
		"This review FAILS the change. `ParsePair` in `internal/example/parse.go` reads the second field of the\n" +
		"split result without checking how many fields the split produced, so an input without a separator\n" +
		"panics at runtime. The new test does not cover that input, so the suite cannot catch the panic. The\n" +
		"change must not merge until the index is guarded and the missing case is tested.\n"
)

// TestSynthesizeReviewOutput_U1StatedFailWithoutBullets is the U1 regression:
// a stated "Verdict: FAIL" with no finding bullet must synthesize as fail.
func TestSynthesizeReviewOutput_U1StatedFailWithoutBullets(t *testing.T) {
	if got := synthesizeReviewOutput(u1NegativeCase).Verdict; got != "fail" {
		t.Errorf("u1-negative synthesized Verdict = %q, want \"fail\"", got)
	}
}

// TestSynthesizeReviewOutput_U1Controls keeps the two controls of the
// qualification set where they were.
func TestSynthesizeReviewOutput_U1Controls(t *testing.T) {
	if got := synthesizeReviewOutput(u1PositiveCase).Verdict; got != "pass" {
		t.Errorf("positive synthesized Verdict = %q, want \"pass\"", got)
	}
	if got := synthesizeReviewOutput(u1BulletNegativeCase).Verdict; got != "fail" {
		t.Errorf("bullet-negative synthesized Verdict = %q, want \"fail\"", got)
	}
}

// TestSynthesizeReviewOutput_VerdictLineDirections is upstream's table
// (f505955a9): both directions plus the conflict case. A stated pass must NOT
// override finding bullets, and prose that merely mentions a verdict is not one.
func TestSynthesizeReviewOutput_VerdictLineDirections(t *testing.T) {
	cases := []struct {
		name, text, want string
	}{
		{"stated fail, no bullets", "Verdict: fail — merge blocked.", "fail"},
		{"stated FAIL uppercase", "VERDICT: FAIL\nreasons follow", "fail"},
		{"stated fail, markdown bold", "**Verdict:** fail\n\nfindings follow", "fail"},
		{"stated pass, no bullets", "Verdict: pass — no blocking findings.", "pass"},
		{"stated pass but bullets present", "Verdict: pass\n- [P1] secret at vuln.go:7", "fail"},
		{"no verdict line, bullets", "- [P2] minor style issue", "fail"},
		{"no verdict line, clean", "The change introduces no blocking issues.", "pass"},
		{"the word verdict in prose only", "I could not reach a verdict on the caching layer.", "pass"},
	}
	for _, c := range cases {
		if got := synthesizeReviewOutput(c.text).Verdict; got != c.want {
			t.Errorf("%s: Verdict = %q, want %q (text: %q)", c.name, got, c.want, c.text)
		}
	}
}
