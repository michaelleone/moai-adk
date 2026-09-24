package template

// Always-loaded instruction budget for the distributed template.
//
// Claude Code (2.1.281+) warns at startup when the instruction files that load
// into every session — CLAUDE.md files, their @-imports, and rule files without
// a top-level `paths:` key — exceed a combined character budget of
// max(120,000, per-file limit): 150,000 characters on a 1M-context model and
// 120,000 on a 200K-context model. It counts UTF-16 code units, not bytes.
//
// A project initialised from these templates adds its own CLAUDE.local.md and
// the user's ~/.claude instructions on top of this surface, so the template's
// own share is held well under the 120,000 floor. Long doctrine lives in
// paths-scoped companions: a `<name>-full.md` beside each always-loaded digest.

import (
	"io/fs"
	"path"
	"strings"
	"testing"
	"unicode/utf16"
)

// templateAlwaysLoadedCharBudget caps the template's always-loaded surface:
// the root CLAUDE.md plus every .claude/rules/**/*.md without a top-level
// `paths:` key. 100,000 leaves 20,000 characters under Claude Code's
// 120,000-character floor for project-local and user-level instructions.
const templateAlwaysLoadedCharBudget = 100_000

// instructionChars counts content the way Claude Code's instruction-size check
// does: UTF-16 code units.
func instructionChars(b []byte) int {
	return len(utf16.Encode([]rune(string(b))))
}

// ruleHasPathsScope reports whether a rule file's frontmatter carries a
// top-level `paths:` key, which makes Claude Code load it only when a matching
// file is read. No frontmatter, or no key, means always-loaded.
func ruleHasPathsScope(content string) bool {
	if !strings.HasPrefix(content, "---\n") {
		return false
	}
	for _, line := range strings.Split(content, "\n")[1:] {
		if strings.HasPrefix(line, "---") {
			return false
		}
		if strings.HasPrefix(line, "paths:") {
			return true
		}
	}
	return false
}

func TestTemplateAlwaysLoadedInstructionBudget(t *testing.T) {
	fsys, err := EmbeddedTemplates()
	if err != nil {
		t.Fatalf("EmbeddedTemplates: %v", err)
	}

	claude, err := fs.ReadFile(fsys, "CLAUDE.md")
	if err != nil {
		t.Fatalf("read CLAUDE.md: %v", err)
	}
	total := instructionChars(claude)
	rules := 0

	err = fs.WalkDir(fsys, ".claude/rules", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || path.Ext(p) != ".md" {
			return nil
		}
		data, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		if ruleHasPathsScope(string(data)) {
			return nil
		}
		rules++
		total += instructionChars(data)
		return nil
	})
	if err != nil {
		t.Fatalf("walk .claude/rules: %v", err)
	}

	if total > templateAlwaysLoadedCharBudget {
		t.Errorf("template always-loaded surface is %d chars (CLAUDE.md + %d unscoped rules), over the %d budget; "+
			"move detail into a paths-scoped companion instead of growing an always-loaded file",
			total, rules, templateAlwaysLoadedCharBudget)
	}
	t.Logf("template always-loaded surface: %d chars (CLAUDE.md + %d unscoped rules), budget %d",
		total, rules, templateAlwaysLoadedCharBudget)
}

// TestTemplateFullTextCompanionsAreScoped pins the digest / full-text split:
// every `<name>-full.md` is paths-scoped to its digest `<name>.md`, and that
// digest exists, stays always-loaded, and names its companion.
func TestTemplateFullTextCompanionsAreScoped(t *testing.T) {
	fsys, err := EmbeddedTemplates()
	if err != nil {
		t.Fatalf("EmbeddedTemplates: %v", err)
	}

	companions := 0
	err = fs.WalkDir(fsys, ".claude/rules", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(p, "-full.md") {
			return nil
		}
		companions++
		full, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		stubPath := strings.TrimSuffix(p, "-full.md") + ".md"
		stubName := path.Base(stubPath)
		if !ruleHasPathsScope(string(full)) {
			t.Errorf("%s must carry a paths: scope so it never joins the always-loaded surface", p)
		}
		if !strings.Contains(string(full), `paths: "**/`+stubName+`"`) {
			t.Errorf("%s must be scoped to its digest (paths: \"**/%s\")", p, stubName)
		}
		stub, err := fs.ReadFile(fsys, stubPath)
		if err != nil {
			t.Errorf("%s has no digest at %s: %v", p, stubPath, err)
			return nil
		}
		if ruleHasPathsScope(string(stub)) {
			t.Errorf("%s is the always-loaded digest of %s and must not be paths-scoped", stubPath, p)
		}
		if !strings.Contains(string(stub), path.Base(p)) {
			t.Errorf("%s must point readers to its full text %s", stubPath, path.Base(p))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk .claude/rules: %v", err)
	}
	if companions == 0 {
		t.Fatal("no *-full.md companions found; the digest / full-text split is missing")
	}
}
