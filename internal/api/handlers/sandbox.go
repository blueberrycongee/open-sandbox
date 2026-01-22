package handlers

import (
	"net/http"

	"open-sandbox/internal/api"
	"open-sandbox/pkg/types"
)

func RegisterSandboxRoutes(router *api.Router) {
	router.Handle(http.MethodGet, "/v1/sandbox", SandboxIndexHandler)
}

func SandboxIndexHandler(w http.ResponseWriter, r *http.Request) *api.AppError {
	payload := map[string]any{
		"status": "ok",
		"capabilities": map[string]string{
			"browser":     "/v1/browser",
			"vnc":         "/vnc/index.html",
			"shell":       "/v1/shell",
			"file":        "/v1/file",
			"code_exec":   "/v1/code",
			"jupyter":     "/jupyter",
			"code_server": "/code-server/",
		},
	}

	if err := api.WriteJSON(w, http.StatusOK, types.Ok(payload)); err != nil {
		return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
	}

	return nil
}
