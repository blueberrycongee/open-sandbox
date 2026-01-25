package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"open-sandbox/internal/api"
	"open-sandbox/internal/mcp"
	"open-sandbox/internal/mcp/remote"
	"open-sandbox/pkg/types"
)

func RegisterExternalMCPRoutes(router *api.Router, manager *remote.Manager, registry *mcp.Registry) {
	router.Handle(http.MethodGet, "/v1/mcp/servers", ListExternalMCPHandler(manager))
	router.Handle(http.MethodPost, "/v1/mcp/servers", UpsertExternalMCPHandler(manager, registry))
	router.HandlePrefix(http.MethodGet, "/v1/mcp/servers/", ExternalMCPServerHandler(manager))
	router.HandlePrefix(http.MethodPut, "/v1/mcp/servers/", UpdateExternalMCPHandler(manager, registry))
	router.HandlePrefix(http.MethodDelete, "/v1/mcp/servers/", DeleteExternalMCPHandler(manager, registry))
	router.HandlePrefix(http.MethodPost, "/v1/mcp/servers/", RefreshExternalMCPHandler(manager, registry))
}

func ListExternalMCPHandler(manager *remote.Manager) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		payload := map[string]any{"servers": manager.List()}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(payload)); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func UpsertExternalMCPHandler(manager *remote.Manager, registry *mcp.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req remote.ServerConfig
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if err := manager.Upsert(req); err != nil {
			return api.NewAppError("bad_request", err.Error(), http.StatusBadRequest)
		}
		_ = manager.SyncRegistry(r.Context(), registry)
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"name": req.Name})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func ExternalMCPServerHandler(manager *remote.Manager) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		name := strings.TrimPrefix(r.URL.Path, "/v1/mcp/servers/")
		if name == "" {
			return api.NewAppError("bad_request", "name is required", http.StatusBadRequest)
		}
		if strings.Contains(name, "/") {
			return api.NewAppError("bad_request", "invalid path", http.StatusBadRequest)
		}
		server, ok := manager.Get(name)
		if !ok {
			return api.NewAppError("not_found", "server not found", http.StatusNotFound)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(server)); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func UpdateExternalMCPHandler(manager *remote.Manager, registry *mcp.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		name := strings.TrimPrefix(r.URL.Path, "/v1/mcp/servers/")
		if name == "" || strings.Contains(name, "/") {
			return api.NewAppError("bad_request", "name is required", http.StatusBadRequest)
		}
		var req remote.ServerConfig
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		req.Name = name
		if err := manager.Upsert(req); err != nil {
			return api.NewAppError("bad_request", err.Error(), http.StatusBadRequest)
		}
		_ = manager.SyncRegistry(r.Context(), registry)
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"name": req.Name})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func DeleteExternalMCPHandler(manager *remote.Manager, registry *mcp.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		name := strings.TrimPrefix(r.URL.Path, "/v1/mcp/servers/")
		if name == "" || strings.Contains(name, "/") {
			return api.NewAppError("bad_request", "name is required", http.StatusBadRequest)
		}
		if err := manager.Delete(name); err != nil {
			return api.NewAppError("bad_request", err.Error(), http.StatusBadRequest)
		}
		_ = manager.SyncRegistry(r.Context(), registry)
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"name": name})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func RefreshExternalMCPHandler(manager *remote.Manager, registry *mcp.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		name := strings.TrimPrefix(r.URL.Path, "/v1/mcp/servers/")
		if strings.HasSuffix(name, "/refresh") {
			name = strings.TrimSuffix(name, "/refresh")
		}
		if strings.Contains(name, "/") {
			return api.NewAppError("bad_request", "invalid path", http.StatusBadRequest)
		}
		if err := manager.SyncRegistry(r.Context(), registry); err != nil {
			return api.NewAppError("refresh_failed", err.Error(), http.StatusInternalServerError)
		}
		payload := map[string]any{"refreshed": true}
		if name != "" {
			payload["name"] = name
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(payload)); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}
