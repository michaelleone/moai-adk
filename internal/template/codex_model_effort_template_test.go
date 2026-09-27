package template

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestCodexModelEffortDocumentedInTemplateWorkflowYAML pins discoverability of
// the Codex-only model/effort keys: a user cannot set a key that appears nowhere
// in the config they were given. Both keys ship, both ship empty (unset), so a
// fresh project keeps resolving codex through the llm.yaml SSOT.
func TestCodexModelEffortDocumentedInTemplateWorkflowYAML(t *testing.T) {
	path := filepath.Join(repoRootFromTemplatePkg(t),
		"internal", "template", "templates", ".moai", "config", "sections", "workflow.yaml")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var doc struct {
		Workflow struct {
			Codex map[string]any `yaml:"codex"`
		} `yaml:"workflow"`
	}
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	for _, key := range []string{"model", "effort"} {
		v, ok := doc.Workflow.Codex[key]
		if !ok {
			t.Errorf("template workflow.yaml: workflow.codex.%s is missing", key)
			continue
		}
		if s, isString := v.(string); !isString || s != "" {
			t.Errorf("template workflow.yaml: workflow.codex.%s = %#v, want the empty string (unset)", key, v)
		}
	}
}
