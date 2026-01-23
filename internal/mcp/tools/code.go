package tools

import (
	"context"
	"encoding/json"

	"open-sandbox/internal/codeexec"
	"open-sandbox/internal/mcp"
)

type codeExecParams struct {
	Runtime    string   `json:"runtime"`
	Args       []string `json:"args"`
	WorkingDir string   `json:"working_dir"`
}

func CodeExec() mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		var payload codeExecParams
		if err := json.Unmarshal(params, &payload); err != nil {
			return nil, invalidParams("invalid params")
		}
		if payload.Runtime == "" {
			return nil, invalidParams("runtime is required")
		}
		workingDir, errDetail := resolveWorkspaceDir(payload.WorkingDir)
		if errDetail != nil {
			return nil, errDetail
		}
		result, err := codeexec.Exec(payload.Runtime, payload.Args, workingDir)
		if err != nil {
			return nil, toolFailure(err.Error())
		}
		return result, nil
	}
}
