package cli

import (
	"context"
	"testing"

	"github.com/modu-ai/moai-adk/internal/update"
	"github.com/modu-ai/moai-adk/pkg/version"
)

// recordingChecker offers a newer release and counts how often it is asked.
type recordingChecker struct{ calls int }

func (c *recordingChecker) CheckLatest(context.Context) (*update.VersionInfo, error) {
	c.calls++
	return &update.VersionInfo{Version: "3.1.3"}, nil
}

func (c *recordingChecker) IsUpdateAvailable(string) (bool, *update.VersionInfo, error) {
	c.calls++
	return true, &update.VersionInfo{Version: "3.1.3"}, nil
}

// recordingOrch pretends to install the offered release and counts installs.
type recordingOrch struct{ installs int }

func (o *recordingOrch) Update(context.Context) (*update.UpdateResult, error) {
	o.installs++
	return &update.UpdateResult{PreviousVersion: version.GetVersion(), NewVersion: "3.1.3"}, nil
}

// withFakeUpdater points the package deps at a checker that always offers
// 3.1.3, with HOME in a temp dir so the update cache starts empty.
func withFakeUpdater(t *testing.T, ver string) (*recordingChecker, *recordingOrch) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	t.Setenv("MOAI_SKIP_BINARY_UPDATE", "")

	origVersion, origDeps := version.Version, deps
	t.Cleanup(func() { version.Version, deps = origVersion, origDeps })

	checker, orch := &recordingChecker{}, &recordingOrch{}
	version.Version = ver
	deps = &Dependencies{UpdateChecker: checker, UpdateOrch: orch}
	return checker, orch
}

func TestIsCustomBuild(t *testing.T) {
	cases := map[string]bool{
		"3.1.2+adax.1":  true,
		"v3.1.2+adax.2": true,
		"3.1.2":         false,
		"v3.1.2":        false,
		"3.1.3-rc1":     false,
	}
	for v, want := range cases {
		if got := isCustomBuild(v); got != want {
			t.Errorf("isCustomBuild(%q) = %v, want %v", v, got, want)
		}
	}
}

// The SessionStart auto-update must never replace a custom (fork) build,
// even when a newer upstream release is on offer.
func TestBuildAutoUpdateFunc_SkipsCustomBuild(t *testing.T) {
	checker, orch := withFakeUpdater(t, "3.1.2+adax.2")

	res, err := buildAutoUpdateFunc()(context.Background())
	if err != nil {
		t.Fatalf("auto-update returned error: %v", err)
	}
	if res == nil || res.Updated {
		t.Fatalf("auto-update result = %+v, want Updated=false", res)
	}
	if checker.calls != 0 || orch.installs != 0 {
		t.Errorf("custom build reached the updater: checks=%d installs=%d, want 0/0", checker.calls, orch.installs)
	}
}

// MOAI_SKIP_BINARY_UPDATE=1 must stop the SessionStart auto-update too, not
// only the `moai update` binary step (upstream modu-ai/moai-adk#1714).
func TestBuildAutoUpdateFunc_HonoursSkipEnv(t *testing.T) {
	checker, orch := withFakeUpdater(t, "3.1.2")
	t.Setenv("MOAI_SKIP_BINARY_UPDATE", "1")

	res, err := buildAutoUpdateFunc()(context.Background())
	if err != nil {
		t.Fatalf("auto-update returned error: %v", err)
	}
	if res == nil || res.Updated {
		t.Fatalf("auto-update result = %+v, want Updated=false", res)
	}
	if checker.calls != 0 || orch.installs != 0 {
		t.Errorf("skip env ignored: checks=%d installs=%d, want 0/0", checker.calls, orch.installs)
	}
}

// Control: a stock release build still auto-updates, so the two tests above
// prove the guard and not a broken fake.
func TestBuildAutoUpdateFunc_StockBuildStillUpdates(t *testing.T) {
	_, orch := withFakeUpdater(t, "3.1.2")

	res, err := buildAutoUpdateFunc()(context.Background())
	if err != nil {
		t.Fatalf("auto-update returned error: %v", err)
	}
	if res == nil || !res.Updated || orch.installs != 1 {
		t.Fatalf("stock build: result = %+v installs=%d, want Updated=true and 1 install", res, orch.installs)
	}
}

// A bare `moai update` skips the binary step for a custom build; an explicit
// `--binary` still reaches it.
func TestShouldSkipBinaryUpdate_CustomBuild(t *testing.T) {
	withFakeUpdater(t, "3.1.2+adax.2")

	cmd := updateCmd
	t.Cleanup(func() { _ = cmd.Flags().Set("binary", "false") })

	if !shouldSkipBinaryUpdate(cmd) {
		t.Error("bare update: shouldSkipBinaryUpdate = false for a custom build, want true")
	}
	if err := cmd.Flags().Set("binary", "true"); err != nil {
		t.Fatal(err)
	}
	if shouldSkipBinaryUpdate(cmd) {
		t.Error("--binary: shouldSkipBinaryUpdate = true for a custom build, want false")
	}
}
