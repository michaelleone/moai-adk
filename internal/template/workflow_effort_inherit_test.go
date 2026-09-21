package template

import (
	"io/fs"
	"path"
	"regexp"
	"strings"
	"testing"
)

// workflowOptionPin matches an effort or model option literal inside a
// workflow agent() options object, e.g. `effort: 'xhigh'` or `model: 'opus'`.
var workflowOptionPin = regexp.MustCompile(`\b(effort|model):\s*'([a-z0-9-]+)'`)

// TestTemplateWorkflowAgentsInheritSessionEffort pins that the bundled workflow
// scripts leave a spawned agent's effort and model to the session. A literal in
// the script wins over every config default and is re-derived by `moai update`,
// so a pin here overrides a project's model/effort policy on every update: a
// `low`/`medium` pin runs an agent below the session, an `xhigh` pin above it.
// The one tolerated literal is `effort: 'high'` — the floor every session
// already runs at, so it can neither undercut nor exceed one.
func TestTemplateWorkflowAgentsInheritSessionEffort(t *testing.T) {
	fsys, err := EmbeddedTemplates()
	if err != nil {
		t.Fatalf("EmbeddedTemplates: %v", err)
	}
	scripts, err := fs.Glob(fsys, ".claude/workflows/*.js")
	if err != nil || len(scripts) == 0 {
		t.Fatalf("no embedded workflow scripts found (err=%v)", err)
	}
	for _, name := range scripts {
		data, err := fs.ReadFile(fsys, name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		for i, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "//") {
				continue // prose about a pin is not a pin
			}
			for _, m := range workflowOptionPin.FindAllStringSubmatch(line, -1) {
				if m[1] == "effort" && m[2] == "high" {
					continue
				}
				t.Errorf("%s:%d pins %s: '%s' — omit it so the agent inherits the session",
					path.Base(name), i+1, m[1], m[2])
			}
		}
	}
}
