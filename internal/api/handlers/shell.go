package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"open-sandbox/internal/api"
	"open-sandbox/internal/config"
	"open-sandbox/internal/shell"
	"open-sandbox/pkg/types"
)

type shellExecRequest struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

func RegisterShellRoutes(router *api.Router) {
	router.Handle(http.MethodPost, "/v1/shell/exec", ShellExecHandler)
}

func ShellExecHandler(w http.ResponseWriter, r *http.Request) *api.AppError {
	var req shellExecRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
	}
	if strings.TrimSpace(req.Command) == "" {
		return api.NewAppError("bad_request", "command is required", http.StatusBadRequest)
	}

	result, err := shell.Exec(req.Command, req.Args, config.WorkspacePath())
	if err != nil {
		return api.NewAppError("exec_failed", err.Error(), http.StatusInternalServerError)
	}

	payload := map[string]any{
		"stdout":    result.Stdout,
		"stderr":    result.Stderr,
		"exit_code": result.ExitCode,
	}
	if err := api.WriteJSON(w, http.StatusOK, types.Ok(payload)); err != nil {
		return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
	}
	return nil
}
