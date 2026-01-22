package handlers

import (
	"encoding/json"
	"net/http"

	"open-sandbox/internal/api"
	"open-sandbox/internal/config"
	"open-sandbox/internal/file"
	"open-sandbox/pkg/types"
)

type fileReadRequest struct {
	Path string `json:"path"`
}

type fileWriteRequest struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type fileSearchRequest struct {
	Path  string `json:"path"`
	Query string `json:"query"`
}

type fileReplaceRequest struct {
	Path    string `json:"path"`
	Search  string `json:"search"`
	Replace string `json:"replace"`
}

func RegisterFileRoutes(router *api.Router) {
	router.Handle(http.MethodPost, "/v1/file/read", FileReadHandler)
	router.Handle(http.MethodPost, "/v1/file/write", FileWriteHandler)
	router.Handle(http.MethodGet, "/v1/file/list", FileListHandler)
	router.Handle(http.MethodPost, "/v1/file/search", FileSearchHandler)
	router.Handle(http.MethodPost, "/v1/file/replace", FileReplaceHandler)
}

func FileReadHandler(w http.ResponseWriter, r *http.Request) *api.AppError {
	var req fileReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
	}
	if err := file.ValidateWorkspacePath(req.Path, config.WorkspacePath()); err != nil {
		return api.NewAppError("bad_request", err.Error(), http.StatusBadRequest)
	}

	content, err := file.Read(req.Path)
	if err != nil {
		return api.NewAppError("read_failed", err.Error(), http.StatusInternalServerError)
	}

	payload := map[string]string{
		"path":    req.Path,
		"content": content,
	}
	if err := api.WriteJSON(w, http.StatusOK, types.Ok(payload)); err != nil {
		return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
	}
	return nil
}

func FileWriteHandler(w http.ResponseWriter, r *http.Request) *api.AppError {
	var req fileWriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
	}
	if err := file.ValidateWorkspacePath(req.Path, config.WorkspacePath()); err != nil {
		return api.NewAppError("bad_request", err.Error(), http.StatusBadRequest)
	}
	if err := file.Write(req.Path, req.Content); err != nil {
		return api.NewAppError("write_failed", err.Error(), http.StatusInternalServerError)
	}

	payload := map[string]string{
		"path": req.Path,
	}
	if err := api.WriteJSON(w, http.StatusOK, types.Ok(payload)); err != nil {
		return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
	}
	return nil
}

func FileListHandler(w http.ResponseWriter, r *http.Request) *api.AppError {
	path := r.URL.Query().Get("path")
	if err := file.ValidateWorkspacePath(path, config.WorkspacePath()); err != nil {
		return api.NewAppError("bad_request", err.Error(), http.StatusBadRequest)
	}

	entries, err := file.List(path)
	if err != nil {
		return api.NewAppError("list_failed", err.Error(), http.StatusInternalServerError)
	}

	payload := map[string]any{
		"path":    path,
		"entries": entries,
	}
	if err := api.WriteJSON(w, http.StatusOK, types.Ok(payload)); err != nil {
		return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
	}
	return nil
}

func FileSearchHandler(w http.ResponseWriter, r *http.Request) *api.AppError {
	var req fileSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
	}
	if err := file.ValidateWorkspacePath(req.Path, config.WorkspacePath()); err != nil {
		return api.NewAppError("bad_request", err.Error(), http.StatusBadRequest)
	}

	matches, err := file.Search(req.Path, req.Query)
	if err != nil {
		return api.NewAppError("search_failed", err.Error(), http.StatusInternalServerError)
	}

	payload := map[string]any{
		"path":    req.Path,
		"matches": matches,
	}
	if err := api.WriteJSON(w, http.StatusOK, types.Ok(payload)); err != nil {
		return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
	}
	return nil
}

func FileReplaceHandler(w http.ResponseWriter, r *http.Request) *api.AppError {
	var req fileReplaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
	}
	if err := file.ValidateWorkspacePath(req.Path, config.WorkspacePath()); err != nil {
		return api.NewAppError("bad_request", err.Error(), http.StatusBadRequest)
	}

	count, err := file.Replace(req.Path, req.Search, req.Replace)
	if err != nil {
		return api.NewAppError("replace_failed", err.Error(), http.StatusInternalServerError)
	}

	payload := map[string]any{
		"path":         req.Path,
		"replacements": count,
	}
	if err := api.WriteJSON(w, http.StatusOK, types.Ok(payload)); err != nil {
		return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
	}
	return nil
}
