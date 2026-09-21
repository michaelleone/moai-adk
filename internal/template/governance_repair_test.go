package template

import (
	"io/fs"
	"strings"
	"testing"
)

// TestTemplateConstitutionKeepsLessonsStoreLive pins the lessons-loader repair:
// run/context-loading.md § Lessons Loading reads a project's lessons.md, so the
// constitution must not declare that file superseded. The two surfaces
// contradicted each other, and `moai update` restored the contradiction.
func TestTemplateConstitutionKeepsLessonsStoreLive(t *testing.T) {
	fsys, err := EmbeddedTemplates()
	if err != nil {
		t.Fatalf("EmbeddedTemplates: %v", err)
	}
	const rule = ".claude/rules/moai/core/moai-constitution.md"
	data, err := fs.ReadFile(fsys, rule)
	if err != nil {
		t.Fatalf("read %s: %v", rule, err)
	}
	for i, line := range strings.Split(string(data), "\n") {
		if strings.Contains(line, "lessons.md") && strings.Contains(strings.ToLower(line), "superseded") {
			t.Errorf("%s:%d declares lessons.md superseded, but Lessons Loading still reads it:\n%s", rule, i+1, line)
		}
	}
	if !strings.Contains(string(data), "`lessons.md` in that directory is also a live lesson store") {
		t.Errorf("%s does not state that lessons.md is a live lesson store", rule)
	}
	loader, err := fs.ReadFile(fsys, ".claude/skills/moai/workflows/run/context-loading.md")
	if err != nil {
		t.Fatalf("read context-loading.md: %v", err)
	}
	if !strings.Contains(string(loader), "memory/lessons.md") {
		t.Errorf("context-loading.md no longer reads lessons.md; revisit the constitution line")
	}
}

// TestTemplateSyncDeliveryMergesWithConfiguredMethod pins the merge-method
// repair: the auto-merge instruction resolves git_strategy.{mode}.merge_method
// instead of prescribing squash, which contradicted any project configured to
// merge.
func TestTemplateSyncDeliveryMergesWithConfiguredMethod(t *testing.T) {
	fsys, err := EmbeddedTemplates()
	if err != nil {
		t.Fatalf("EmbeddedTemplates: %v", err)
	}
	const wf = ".claude/skills/moai/workflows/sync/delivery.md"
	data, err := fs.ReadFile(fsys, wf)
	if err != nil {
		t.Fatalf("read %s: %v", wf, err)
	}
	text := string(data)
	if strings.Contains(text, "gh pr merge --squash") {
		t.Errorf("%s prescribes `gh pr merge --squash`; it must use the configured merge_method", wf)
	}
	if got := strings.Count(text, "Execute `gh pr merge --<merge_method> --delete-branch`"); got != 2 {
		t.Errorf("%s: %d merge instructions resolve merge_method, want 2", wf, got)
	}
	if !strings.Contains(text, "`git_strategy.{mode}.merge_method` read from `.moai/config/sections/git-strategy.yaml`") {
		t.Errorf("%s does not say where merge_method is read from", wf)
	}
}
