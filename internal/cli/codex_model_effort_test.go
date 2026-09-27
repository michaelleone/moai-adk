package cli

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/modu-ai/moai-adk/internal/config"
	"github.com/modu-ai/moai-adk/internal/template"
)

// workflow.codex.model / workflow.codex.effort — the Codex-only model + effort
// surface. A maintainer who wants codex on a specific model must be able to say
// so WITHOUT writing a Codex model id into llm.agent_overrides["sync-auditor"]:
// that per-agent override is the same key the Claude sync-auditor spawn resolves
// through, so a Codex id there would mis-route the Claude auditor.
//
// Precedence under test (codexSSOTModelEffort, inherited by every caller):
//   - model:  workflow.codex.model (verbatim, no servability filter)
//     → else the SSOT model when codex can serve it → else empty;
//   - effort: workflow.codex.effort
//     → else the SSOT effort ONLY when the SSOT model was used → else empty;
//   - both Codex-only keys empty ⇒ byte-identical to the SSOT-only behaviour;
//   - an explicit per-call `model` still wins over everything.

// codexModelEffortYAML is a workflow.yaml body carrying the Codex-only keys.
// An empty value omits that key entirely.
func codexModelEffortYAML(model, effort string) string {
	body := "workflow:\n  codex:\n    review_gate:\n      enabled: false\n"
	if model != "" {
		body += "    model: " + model + "\n"
	}
	if effort != "" {
		body += "    effort: " + effort + "\n"
	}
	return body
}

// (i) The Codex-only keys win even when llm.yaml maps the audit agent key to a
// Claude model, they reach the params actually transmitted to codex, and the
// Claude sync-auditor resolution is left untouched.
func TestCodexModelEffort_CodexOnlyKeysWinOverClaudeSSOT(t *testing.T) {
	root := writeCodexLLMFixture(t, "opus", "xhigh")
	writeCodexWorkflowYAML(t, root, codexModelEffortYAML("gpt-6-astra", "high"))

	got := codexSSOTModelEffort(root)
	want := config.ModelEffort{Model: "gpt-6-astra", Effort: "high"}
	if got != want {
		t.Fatalf("codexSSOTModelEffort = %+v, want %+v (workflow.codex.* must win over the Claude SSOT cell)", got, want)
	}

	// The same values reach the wire on the turn/start path.
	sess := withCodexSession(t, codexSessionScript("clean"))
	if _, err := runCodexReviewRPC(context.Background(), "/fake/codex", codexMethodTurnStart, map[string]any{
		"prompt": "review this",
		"cwd":    root,
	}); err != nil {
		t.Fatalf("rpc: %v", err)
	}
	if m, _ := sentParams(t, sess.sent, 1)["model"].(string); m != "gpt-6-astra" {
		t.Errorf("thread/start model = %q, want gpt-6-astra", m)
	}
	turn := sentParams(t, sess.sent, 2)
	if m, _ := turn["model"].(string); m != "gpt-6-astra" {
		t.Errorf("turn/start model = %q, want gpt-6-astra", m)
	}
	if e, _ := turn["effort"].(string); e != "high" {
		t.Errorf("turn/start effort = %q, want high", e)
	}

	// Separation: the Claude sync-auditor still resolves to its own cell.
	llm, err := loadLLMSectionOnly(filepath.Join(root, ".moai", "config", "sections"))
	if err != nil {
		t.Fatalf("load llm.yaml: %v", err)
	}
	if me, _ := template.ResolveAgentModelEffort(llm, codexAuditAgentKey); me.Model != "opus" || me.Effort != "xhigh" {
		t.Errorf("Claude %s resolution = %+v, want {opus xhigh} (the Codex-only keys must not touch it)", codexAuditAgentKey, me)
	}
}

// (ii) A Codex-only model outside the codex-servable families is still sent:
// the maintainer chose it deliberately, so the servability filter (which
// protects callers who did NOT choose) does not apply to this key.
func TestCodexModelEffort_CodexOnlyModelBypassesServabilityFilter(t *testing.T) {
	const private = "acme-private-reviewer"
	if codexServableModel(private) {
		t.Fatalf("precondition: %q must be outside the codex-servable families for this test to mean anything", private)
	}
	root := writeCodexLLMFixture(t, "opus", "high")
	writeCodexWorkflowYAML(t, root, codexModelEffortYAML(private, ""))

	got := codexSSOTModelEffort(root)
	if got.Model != private {
		t.Errorf("model = %q, want %q sent verbatim (no servability filter on workflow.codex.model)", got.Model, private)
	}
	if got.Effort != "" {
		t.Errorf("effort = %q, want empty (no workflow.codex.effort, and the SSOT cell was not used)", got.Effort)
	}
}

// A Codex-only model replaces the SSOT cell whole: the SSOT effort is not
// borrowed onto a model it was never paired with.
func TestCodexModelEffort_CodexOnlyModelDoesNotBorrowSSOTEffort(t *testing.T) {
	root := writeCodexLLMFixture(t, "gpt-5-codex", "medium")
	writeCodexWorkflowYAML(t, root, codexModelEffortYAML("gpt-6-astra", ""))

	got := codexSSOTModelEffort(root)
	want := config.ModelEffort{Model: "gpt-6-astra"}
	if got != want {
		t.Errorf("codexSSOTModelEffort = %+v, want %+v (SSOT effort applies only when the SSOT model is used)", got, want)
	}
}

// (iii) Both Codex-only keys empty ⇒ identical to the existing SSOT behaviour,
// for every way of leaving them empty. The expected cells are the ones the
// existing SSOT tests pin (TestCodexSession_ResolvedModelReachesTransmittedParams
// and TestCodexSession_NonCodexModelNotTransmitted).
func TestCodexModelEffort_BothKeysEmptyMatchesSSOT(t *testing.T) {
	ssotCases := []struct {
		name          string
		model, effort string
		want          config.ModelEffort
	}{
		{"servable-ssot", "gpt-5-codex", "high", config.ModelEffort{Model: "gpt-5-codex", Effort: "high"}},
		{"claude-ssot", "opus", "high", config.ModelEffort{}},
	}
	workflowCases := []struct {
		name string
		body string
	}{
		{"absent-file", ""},
		{"absent-keys", codexModelEffortYAML("", "")},
		{"explicit-empty", "workflow:\n  codex:\n    model: \"\"\n    effort: \"\"\n"},
		{"whitespace-only", "workflow:\n  codex:\n    model: \"   \"\n    effort: \"  \"\n"},
		{"garbled", "workflow: [this is not: a mapping\n"},
	}
	for _, sc := range ssotCases {
		for _, wc := range workflowCases {
			t.Run(sc.name+"/"+wc.name, func(t *testing.T) {
				root := writeCodexLLMFixture(t, sc.model, sc.effort)
				writeCodexWorkflowYAML(t, root, wc.body)
				if got := codexSSOTModelEffort(root); got != sc.want {
					t.Errorf("codexSSOTModelEffort = %+v, want %+v (unchanged SSOT behaviour)", got, sc.want)
				}
			})
		}
	}
}

// (iv) Effort-only: the Codex-only effort applies; the model comes from the
// SSOT when codex can serve it, otherwise it stays empty.
func TestCodexModelEffort_EffortOnly(t *testing.T) {
	cases := []struct {
		name          string
		model, effort string
		want          config.ModelEffort
	}{
		{"servable-ssot-model", "gpt-5-codex", "medium", config.ModelEffort{Model: "gpt-5-codex", Effort: "high"}},
		{"claude-ssot-model", "opus", "xhigh", config.ModelEffort{Effort: "high"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := writeCodexLLMFixture(t, tc.model, tc.effort)
			writeCodexWorkflowYAML(t, root, codexModelEffortYAML("", "high"))
			if got := codexSSOTModelEffort(root); got != tc.want {
				t.Errorf("codexSSOTModelEffort = %+v, want %+v", got, tc.want)
			}
		})
	}
}

// (v) An explicit per-call `model` still wins over the Codex-only key; the
// Codex-only effort still applies alongside it.
func TestCodexModelEffort_ExplicitParamModelStillWins(t *testing.T) {
	root := writeCodexLLMFixture(t, "opus", "high")
	writeCodexWorkflowYAML(t, root, codexModelEffortYAML("gpt-6-astra", "high"))

	got := resolveCodexModelEffort(map[string]any{"cwd": root, "model": "o4-mini"})
	want := config.ModelEffort{Model: "o4-mini", Effort: "high"}
	if got != want {
		t.Errorf("resolveCodexModelEffort = %+v, want %+v (explicit model wins)", got, want)
	}
}

// (vi) codex_setup reports the RESOLVED model and effort — the values codex
// would actually be sent — under `model` and `effort`, always present.
func TestCodexSetup_ReportsResolvedModelAndEffort(t *testing.T) {
	cases := []struct {
		name         string
		llmModel     string
		llmEffort    string
		workflowBody string
		wantModel    string
		wantEffort   string
	}{
		{"codex-only-keys", "opus", "xhigh", codexModelEffortYAML("gpt-6-astra", "high"), "gpt-6-astra", "high"},
		{"servable-ssot", "gpt-5-codex", "medium", "", "gpt-5-codex", "medium"},
		{"nothing-resolves", "opus", "high", "", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := writeCodexLLMFixture(t, tc.llmModel, tc.llmEffort)
			writeCodexWorkflowYAML(t, root, tc.workflowBody)
			withCodexLookPath(t, func(string) (string, error) { return "", errFakeLookPath })

			res, err := handleCodexSetup(context.Background(), mcp.CallToolRequest{})
			if err != nil {
				t.Fatalf("handleCodexSetup: %v", err)
			}
			got := structuredMap(t, res)
			for key, want := range map[string]string{"model": tc.wantModel, "effort": tc.wantEffort} {
				raw, ok := got[key]
				if !ok {
					t.Errorf("codex_setup result carries no %q key (it must always be present)", key)
					continue
				}
				if raw != want {
					t.Errorf("codex_setup %s = %v, want %q", key, raw, want)
				}
			}
		})
	}
}

// (vii) The reader is fail-OPEN: every unreadable or unusable state reads as
// "unset" (empty values), never an error, so a broken workflow.yaml leaves the
// SSOT path in charge rather than blanking it. Values are trimmed.
func TestReadCodexModelEffort_FailOpen(t *testing.T) {
	cases := []struct {
		name       string
		body       string
		wantModel  string
		wantEffort string
	}{
		{"absent-file", "", "", ""},
		{"garbled-yaml", "workflow: [this is not: a mapping\n", "", ""},
		{"wrong-type", "workflow:\n  codex:\n    model: [a, b]\n    effort: high\n", "", ""},
		{"absent-codex-block", "workflow:\n  multi:\n    review_gate:\n      enabled: false\n", "", ""},
		{"absent-keys", codexModelEffortYAML("", ""), "", ""},
		{"both-set-trimmed", "workflow:\n  codex:\n    model: \"  gpt-6-astra \"\n    effort: \" high\"\n", "gpt-6-astra", "high"},
		{"flat-shape-not-honoured", "codex:\n  model: gpt-6-astra\n  effort: high\n", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := writeGateWorkflowYAML(t, tc.body)
			model, effort := readCodexModelEffort(dir)
			if model != tc.wantModel || effort != tc.wantEffort {
				t.Errorf("readCodexModelEffort = (%q, %q), want (%q, %q)", model, effort, tc.wantModel, tc.wantEffort)
			}
		})
	}
	if model, effort := readCodexModelEffort(""); model != "" || effort != "" {
		t.Errorf("empty projectDir: readCodexModelEffort = (%q, %q), want empty", model, effort)
	}
	if model, effort := readCodexModelEffort(filepath.Join(t.TempDir(), "does-not-exist")); model != "" || effort != "" {
		t.Errorf("missing project dir: readCodexModelEffort = (%q, %q), want empty", model, effort)
	}
}

// The hand-rolled reader must agree with the real config loader on the key path
// (the sibling readers carry the same guard, for the same reason: a reader that
// reads the wrong nesting level silently never fires).
func TestReadCodexModelEffort_AgreesWithConfigLoader(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
	}{
		{"both-set", "workflow:\n    codex:\n        model: gpt-6-astra\n        effort: high\n"},
		{"model-only", "workflow:\n    codex:\n        model: gpt-6-astra\n"},
		{"absent-keys", "workflow:\n    codex:\n        review_gate:\n            enabled: true\n"},
		{"absent-file", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := writeGateWorkflowYAML(t, tc.body)
			cfg, err := config.NewLoader().Load(filepath.Join(dir, ".moai"))
			if err != nil {
				t.Fatalf("config load: %v", err)
			}
			model, effort := readCodexModelEffort(dir)
			if model != cfg.Workflow.Codex.Model || effort != cfg.Workflow.Codex.Effort {
				t.Errorf("reader = (%q, %q), config loader = (%q, %q) (schema drift)",
					model, effort, cfg.Workflow.Codex.Model, cfg.Workflow.Codex.Effort)
			}
		})
	}
}

// The distributed default leaves both keys unset, so a project that never opts
// in keeps the SSOT-only behaviour.
func TestCodexModelEffort_DistributedDefaultIsUnset(t *testing.T) {
	c := config.NewDefaultWorkflowConfig().Codex
	if c.Model != "" || c.Effort != "" {
		t.Errorf("default workflow.codex = {model:%q effort:%q}, want both empty (unset)", c.Model, c.Effort)
	}
}

// The typed probe (shared with the web console) carries the same resolved
// values the codex_setup JSON reports.
func TestProbeCodexSetup_CarriesResolvedModelAndEffort(t *testing.T) {
	root := writeCodexLLMFixture(t, "opus", "xhigh")
	writeCodexWorkflowYAML(t, root, codexModelEffortYAML("gpt-6-astra", "high"))
	withCodexLookPath(t, func(string) (string, error) { return "", errFakeLookPath })

	s := ProbeCodexSetup(context.Background())
	if s.Model != "gpt-6-astra" || s.Effort != "high" {
		t.Errorf("ProbeCodexSetup = {Model:%q Effort:%q}, want {gpt-6-astra high}", s.Model, s.Effort)
	}
}
