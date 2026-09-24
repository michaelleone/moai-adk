# moai-mcp Tool Catalogue

> The self-hosted `moai` MCP server (`.mcp.json` → `{command: "moai", args: ["mcp-server"]}`) exposes 21 tools, each prefixed `mcp__moai__` at the call site. This digest carries the MCP-over-CLI rule and the tool families; the full catalogue — purpose, consumer agents, and CLI equivalent per tool — is `moai-mcp-tools-full.md` (auto-loads when this file is read).

## MCP-over-CLI rule

Prefer the MCP tool when it is in the calling agent's `tools:` list. The MCP path and the Bash CLI back the same implementation; the MCP path returns structured output, avoids shell-quoting hazards, and is lower-latency inside a subagent. Use the Bash CLI only when the MCP tool is absent from the agent's `tools:` list, or from the main session when the CLI form reads more naturally inline.

## Tool families

| Family | Tools |
|---|---|
| SPEC lifecycle | `spec_progress`, `spec_audit`, `spec_drift` |
| Verification snapshots | `verify_snapshot`, `verify_trend` |
| Goal + session | `goal_arm`, `goal_status`, `session_list` |
| Cross-model audit | `audit_multi`, `codex_audit`, `glm_audit`, `audit_cache` |
| Codex delegation | `codex_task`, `codex_setup`, `codex_job_status`, `codex_job_result`, `codex_job_cancel` |
| GLM delegation | `glm_task`, `glm_job_status`, `glm_job_result`, `glm_job_cancel` |

`goal_arm` is wired to no agent: arming an autonomous loop is an orchestrator concern, and agents read `goal_status` instead. Audit and delegation backends are optional and fail open — an unavailable backend yields `inconclusive`, never an error.

---

Classification: Evolvable reference rule — the MCP tool surface map. Digest; full catalogue: `moai-mcp-tools-full.md`. Update both whenever a tool is added, removed, or renamed on `moai mcp-server` (producer: `internal/cli/mcp_server.go`).
