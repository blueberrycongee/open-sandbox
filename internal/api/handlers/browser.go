package handlers

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"

	"open-sandbox/internal/api"
	"open-sandbox/internal/browser"
	"open-sandbox/internal/config"
	"open-sandbox/pkg/types"
)

type navigateRequest struct {
	URL string `json:"url"`
}

type screenshotRequest struct {
	Path string `json:"path"`
}

func RegisterBrowserRoutes(router *api.Router, service *browser.Service) {
	router.Handle(http.MethodGet, "/v1/browser/info", BrowserInfoHandler(service))
	router.Handle(http.MethodPost, "/v1/browser/navigate", BrowserNavigateHandler(service))
	router.Handle(http.MethodPost, "/v1/browser/screenshot", BrowserScreenshotHandler(service))
}

func BrowserInfoHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		cdpAddress, err := service.Info()
		if err != nil {
			if err == browser.ErrBrowserUnavailable {
				return api.NewAppError("browser_unavailable", "browser binary not found", http.StatusServiceUnavailable)
			}
			return api.NewAppError("browser_unavailable", err.Error(), http.StatusServiceUnavailable)
		}

		payload := map[string]string{
			"cdp_address": cdpAddress,
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(payload)); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func BrowserNavigateHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req navigateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if strings.TrimSpace(req.URL) == "" {
			return api.NewAppError("bad_request", "url is required", http.StatusBadRequest)
		}

		if err := service.Navigate(req.URL); err != nil {
			if err == browser.ErrBrowserUnavailable {
				return api.NewAppError("browser_unavailable", "browser binary not found", http.StatusServiceUnavailable)
			}
			return api.NewAppError("navigate_failed", err.Error(), http.StatusInternalServerError)
		}

		payload := map[string]string{
			"navigated_to": req.URL,
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(payload)); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func BrowserScreenshotHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req screenshotRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if strings.TrimSpace(req.Path) == "" {
			return api.NewAppError("bad_request", "path is required", http.StatusBadRequest)
		}
		if !filepath.IsAbs(req.Path) {
			return api.NewAppError("bad_request", "path must be absolute", http.StatusBadRequest)
		}
		if !strings.HasPrefix(req.Path, config.HostWorkspacePath) && !strings.HasPrefix(req.Path, config.ContainerWorkspacePath) {
			return api.NewAppError("bad_request", "path must be within workspace", http.StatusBadRequest)
		}

		if err := service.Screenshot(req.Path); err != nil {
			if err == browser.ErrBrowserUnavailable {
				return api.NewAppError("browser_unavailable", "browser binary not found", http.StatusServiceUnavailable)
			}
			return api.NewAppError(api.CodeInternalError, "screenshot failed", http.StatusInternalServerError)
		}

		payload := map[string]string{
			"path": req.Path,
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(payload)); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}
