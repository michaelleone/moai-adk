package cli

import (
	"os"
	"strings"

	"github.com/modu-ai/moai-adk/internal/config"
)

// isDevBuild reports a dev build (version contains "dirty", "dev" or "none"),
// where self-update is meaningless.
func isDevBuild(v string) bool {
	return strings.Contains(v, "dirty") || v == "dev" || strings.Contains(v, "none")
}

// isCustomBuild reports whether v carries semver build metadata
// ("3.1.2+adax.1"), which marks a local or fork build. compareSemver drops the
// metadata and ranks such a build equal to the release it is based on, so any
// newer upstream release would otherwise replace it without being asked.
func isCustomBuild(v string) bool {
	return strings.Contains(v, "+")
}

// autoUpdateBlocked reports whether the SessionStart auto-update must leave
// the binary at version v alone: MOAI_SKIP_BINARY_UPDATE=1 is set, or v is a
// dev or custom build. A custom build moves only on an explicit
// `moai update --binary` or `moai update --version <tag>`.
func autoUpdateBlocked(v string) bool {
	return os.Getenv(config.EnvSkipBinaryUpdate) == "1" || isDevBuild(v) || isCustomBuild(v)
}
