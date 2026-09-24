package quality

// astgrep_gate_walk_test.go covers which files the suppression policy check
// reads: in a git work tree, the tracked and untracked-but-not-ignored files,
// minus the excluded roots the scan's findings are filtered from. It must never
// read gitignored content, such as a nested clone parked under
// .claude/worktrees/ or a virtualenv, whose fixtures would otherwise block
// every commit of the enclosing project.

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

const unpairedSuppressionSrc = "package x\n\n// ast-grep-ignore\nx := dangerousOp()\n"

// writeTree writes each relative path under root with the given content.
func writeTree(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", rel, err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
	}
}

// gitOnlyPath returns a PATH value that resolves git and nothing else, so a
// test can use a git work tree while the sg scanner stays unavailable. It
// isolates git from the owner's global and system config, whose excludes
// file would otherwise change what counts as ignored.
func gitOnlyPath(t *testing.T) (gitBin, pathValue string) {
	t.Helper()
	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git binary not in PATH; the git-listing path is not mechanically verifiable")
	}
	binDir := t.TempDir()
	if err := os.Symlink(gitBin, filepath.Join(binDir, "git")); err != nil {
		t.Fatalf("symlink git: %v", err)
	}
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	return gitBin, binDir
}

// initGitProject makes root a git work tree and stages the given paths.
func initGitProject(t *testing.T, gitBin, root string, stage ...string) {
	t.Helper()
	run := func(args ...string) {
		cmd := exec.Command(gitBin, args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
	}
	run("init", "-q")
	if len(stage) > 0 {
		run(append([]string{"add", "--"}, stage...)...)
	}
}

// TestRunAstGrepGateV2_GitignoredNestedCheckoutDoesNotBlock reproduces the
// adax-v3 failure: a gitignored clone of moai-adk parked at
// .claude/worktrees/moai-adk-fork failed every commit of the enclosing project
// on its own fixture internal/astgrep/testdata/fixtures/go/suppressed.go.
// A second gitignored clone outside any excluded root proves the check
// honours .gitignore itself, not only the excluded-root list. A tracked
// violation must still block.
func TestRunAstGrepGateV2_GitignoredNestedCheckoutDoesNotBlock(t *testing.T) {
	gitBin, pathValue := gitOnlyPath(t)

	projectDir := t.TempDir()
	writeTree(t, projectDir, map[string]string{
		".gitignore": ".claude/worktrees/\nparked/\n",
		"pkg/ok.go":  "package pkg\n\nfunc ok() {}\n",
		".claude/worktrees/moai-adk-fork/internal/astgrep/testdata/fixtures/go/suppressed.go": unpairedSuppressionSrc,
		"parked/moai-adk-fork/internal/foo/bar.go":                                            unpairedSuppressionSrc,
	})
	initGitProject(t, gitBin, projectDir, ".gitignore", "pkg/ok.go")
	t.Setenv("PATH", pathValue)

	cfg := &AstGrepGateConfig{
		Enabled:      true,
		RulesDir:     ".moai/config/astgrep-rules",
		BlockOnError: true,
	}

	passed, output := RunAstGrepGateV2(context.Background(), projectDir, cfg)
	if !passed {
		t.Fatalf("gate blocked on gitignored content: %q", output)
	}

	// Negative control: the same violation in project code still blocks.
	writeTree(t, projectDir, map[string]string{"pkg/bad.go": unpairedSuppressionSrc})
	passed, output = RunAstGrepGateV2(context.Background(), projectDir, cfg)
	if passed {
		t.Fatal("gate passed although untracked, unignored pkg/bad.go has an unpaired suppression")
	}
	if !strings.Contains(output, filepath.Join("pkg", "bad.go")) {
		t.Errorf("output should name pkg/bad.go, got: %q", output)
	}
	for _, ignored := range []string{"moai-adk-fork", "parked"} {
		if strings.Contains(output, ignored) {
			t.Errorf("output names gitignored path %q: %q", ignored, output)
		}
	}
}

// relSorted returns walkSourceFiles' result relative to root, slash-separated
// and sorted, for order-independent comparison.
func relSorted(t *testing.T, root string, files []string) []string {
	t.Helper()
	out := make([]string, 0, len(files))
	for _, f := range files {
		rel, err := filepath.Rel(root, f)
		if err != nil {
			t.Fatalf("rel %s: %v", f, err)
		}
		out = append(out, filepath.ToSlash(rel))
	}
	sort.Strings(out)
	return out
}

// excludedRootsTree holds one checkable file plus one file in every root the
// suppression check must skip whether or not git ignores it.
var excludedRootsTree = map[string]string{
	"pkg/a.go":                              unpairedSuppressionSrc,
	"pkg/a_test.go":                         unpairedSuppressionSrc,
	"pkg/testdata/fixture.go":               unpairedSuppressionSrc,
	"testdata/root_fixture.go":              unpairedSuppressionSrc,
	"vendor/v/v.go":                         unpairedSuppressionSrc,
	"node_modules/m/m.js":                   unpairedSuppressionSrc,
	"lib/__pycache__/c.py":                  unpairedSuppressionSrc,
	".claude/worktrees/wt/internal/x/wt.go": unpairedSuppressionSrc,
	"docs/readme.md":                        unpairedSuppressionSrc,
}

// TestWalkSourceFiles_OutsideGitPrunesExcludedRoots covers the fallback walk:
// without git it must still skip every excluded root.
func TestWalkSourceFiles_OutsideGitPrunesExcludedRoots(t *testing.T) {
	t.Setenv("PATH", "") // no git: force the directory walk

	root := t.TempDir()
	writeTree(t, root, excludedRootsTree)

	got := relSorted(t, root, walkSourceFiles(context.Background(), root))
	want := []string{"pkg/a.go"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("walkSourceFiles = %v, want %v", got, want)
	}
}

// TestWalkSourceFiles_GitListsTrackedAndUnignoredOnly covers the git listing:
// tracked and untracked-unignored files are read; gitignored files and nested
// repositories are not, and the excluded roots stay excluded even when git
// does not ignore them.
func TestWalkSourceFiles_GitListsTrackedAndUnignoredOnly(t *testing.T) {
	gitBin, pathValue := gitOnlyPath(t)

	root := t.TempDir()
	writeTree(t, root, excludedRootsTree)
	writeTree(t, root, map[string]string{
		".gitignore":           "build/\n.venv/\n",
		"pkg/untracked.go":     unpairedSuppressionSrc,
		"build/gen.go":         unpairedSuppressionSrc,
		".venv/lib/site.py":    unpairedSuppressionSrc,
		"nested/internal/n.go": unpairedSuppressionSrc,
		"nested/.git/HEAD":     "ref: refs/heads/main\n",
	})
	initGitProject(t, gitBin, root, ".gitignore", "pkg/a.go")
	initGitProject(t, gitBin, filepath.Join(root, "nested"))
	t.Setenv("PATH", pathValue)

	got := relSorted(t, root, walkSourceFiles(context.Background(), root))
	want := []string{"pkg/a.go", "pkg/untracked.go"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("walkSourceFiles = %v, want %v", got, want)
	}
}

// TestWalkSourceFiles_ProjectUnderExcludedRootIsChecked pins that excluded
// roots match relative to the project, not the absolute path: a project that
// is itself a worktree under .claude/worktrees/ must still be checked.
func TestWalkSourceFiles_ProjectUnderExcludedRootIsChecked(t *testing.T) {
	t.Setenv("PATH", "")

	root := filepath.Join(t.TempDir(), ".claude", "worktrees", "feature")
	writeTree(t, root, map[string]string{"pkg/a.go": unpairedSuppressionSrc})

	got := relSorted(t, root, walkSourceFiles(context.Background(), root))
	if strings.Join(got, ",") != "pkg/a.go" {
		t.Errorf("walkSourceFiles = %v, want [pkg/a.go]", got)
	}
}
