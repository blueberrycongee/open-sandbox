package tools

import (
	"context"
	"encoding/json"

	"open-sandbox/internal/mcp"
	"open-sandbox/internal/shell"
)

type shellExecParams struct {
	Command    string   `json:"command"`
	Args       []string `json:"args"`
	WorkingDir string   `json:"working_dir"`
}

func ShellExec() mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		var payload shellExecParams
		if err := json.Unmarshal(params, &payload); err != nil {
			return nil, invalidParams("invalid params")
		}
		if payload.Command == "" {
			return nil, invalidParams("command is required")
		}
		workingDir, errDetail := resolveWorkspaceDir(payload.WorkingDir)
		if errDetail != nil {
			return nil, errDetail
		}
		result, err := shell.Exec(payload.Command, payload.Args, workingDir)
		if err != nil {
			return nil, toolFailure(err.Error())
		}
		return result, nil
	}
}
