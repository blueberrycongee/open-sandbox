package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"open-sandbox/internal/api"
	"open-sandbox/internal/codeexec"
	"open-sandbox/internal/config"
	"open-sandbox/pkg/types"
)

type codeExecRequest struct {
	Runtime string   `json:"runtime"`
	Args    []string `json:"args"`
}

func RegisterCodeExecRoutes(router *api.Router) {
	router.Handle(http.MethodPost, "/v1/code/exec", CodeExecHandler)
}

func CodeExecHandler(w http.ResponseWriter, r *http.Request) *api.AppError {
	var req codeExecRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
	}

	result, err := codeexec.Exec(req.Runtime, req.Args, config.WorkspacePath())
	if err != nil {
		message := err.Error()
		switch {
		case strings.Contains(message, "unsupported"):
			return api.NewAppError("bad_request", message, http.StatusBadRequest)
		case strings.Contains(message, "cannot find"):
			return api.NewAppError("runtime_not_found", message, http.StatusServiceUnavailable)
		default:
			return api.NewAppError("exec_failed", message, http.StatusInternalServerError)
		}
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
