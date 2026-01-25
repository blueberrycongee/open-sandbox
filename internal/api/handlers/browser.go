package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"open-sandbox/internal/api"
	"open-sandbox/internal/browser"
	"open-sandbox/internal/config"
	"open-sandbox/internal/file"
	"open-sandbox/pkg/types"
)

type navigateRequest struct {
	URL string `json:"url"`
}

type screenshotRequest struct {
	Path string `json:"path"`
}

type clickRequest struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type formInputFillRequest struct {
	Selector string `json:"selector"`
	Value    string `json:"value"`
}

type elementSelectRequest struct {
	Selector string `json:"selector"`
	Value    string `json:"value"`
}

type scrollRequest struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type evaluateRequest struct {
	Expression string `json:"expression"`
}

type tabNewRequest struct {
	URL string `json:"url"`
}

type tabSwitchRequest struct {
	Index int `json:"index"`
}

type tabCloseRequest struct {
	Index int `json:"index"`
}

type pressKeyRequest struct {
	Keys string `json:"keys"`
}

func RegisterBrowserRoutes(router *api.Router, service *browser.Service) {
	router.Handle(http.MethodGet, "/v1/browser/info", BrowserInfoHandler(service))
	router.Handle(http.MethodPost, "/v1/browser/navigate", BrowserNavigateHandler(service))
	router.Handle(http.MethodPost, "/v1/browser/screenshot", BrowserScreenshotHandler(service))
	router.Handle(http.MethodPost, "/v1/browser/click", BrowserClickHandler(service))
	router.Handle(http.MethodPost, "/v1/browser/form_input_fill", BrowserFormInputFillHandler(service))
	router.Handle(http.MethodPost, "/v1/browser/select", BrowserSelectHandler(service))
	router.Handle(http.MethodPost, "/v1/browser/scroll", BrowserScrollHandler(service))
	router.Handle(http.MethodPost, "/v1/browser/evaluate", BrowserEvaluateHandler(service))
	router.Handle(http.MethodPost, "/v1/browser/new_tab", BrowserNewTabHandler(service))
	router.Handle(http.MethodPost, "/v1/browser/switch_tab", BrowserSwitchTabHandler(service))
	router.Handle(http.MethodPost, "/v1/browser/close_tab", BrowserCloseTabHandler(service))
	router.Handle(http.MethodGet, "/v1/browser/tab_list", BrowserTabListHandler(service))
	router.Handle(http.MethodGet, "/v1/browser/get_download_list", BrowserDownloadListHandler(service))
	router.Handle(http.MethodPost, "/v1/browser/press_key", BrowserPressKeyHandler(service))
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
		if file.ValidateWorkspacePath(req.Path, config.WorkspacePath()) != nil &&
			file.ValidateWorkspacePath(req.Path, config.ContainerWorkspacePath) != nil {
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

func BrowserClickHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req clickRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if err := service.Click(req.X, req.Y); err != nil {
			if err == browser.ErrBrowserUnavailable {
				return api.NewAppError("browser_unavailable", "browser binary not found", http.StatusServiceUnavailable)
			}
			return api.NewAppError("click_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"clicked": true})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func BrowserFormInputFillHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req formInputFillRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if strings.TrimSpace(req.Selector) == "" {
			return api.NewAppError("bad_request", "selector is required", http.StatusBadRequest)
		}
		if err := service.FormInputFill(req.Selector, req.Value); err != nil {
			if err == browser.ErrBrowserUnavailable {
				return api.NewAppError("browser_unavailable", "browser binary not found", http.StatusServiceUnavailable)
			}
			return api.NewAppError("form_input_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"filled": true})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func BrowserSelectHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req elementSelectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if strings.TrimSpace(req.Selector) == "" {
			return api.NewAppError("bad_request", "selector is required", http.StatusBadRequest)
		}
		if err := service.ElementSelect(req.Selector, req.Value); err != nil {
			if err == browser.ErrBrowserUnavailable {
				return api.NewAppError("browser_unavailable", "browser binary not found", http.StatusServiceUnavailable)
			}
			return api.NewAppError("select_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"selected": true})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func BrowserScrollHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req scrollRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if err := service.Scroll(req.X, req.Y); err != nil {
			if err == browser.ErrBrowserUnavailable {
				return api.NewAppError("browser_unavailable", "browser binary not found", http.StatusServiceUnavailable)
			}
			return api.NewAppError("scroll_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"scrolled": true})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func BrowserEvaluateHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req evaluateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if strings.TrimSpace(req.Expression) == "" {
			return api.NewAppError("bad_request", "expression is required", http.StatusBadRequest)
		}
		result, err := service.Evaluate(req.Expression)
		if err != nil {
			if err == browser.ErrBrowserUnavailable {
				return api.NewAppError("browser_unavailable", "browser binary not found", http.StatusServiceUnavailable)
			}
			return api.NewAppError("evaluate_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"result": result})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func BrowserNewTabHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req tabNewRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		index, err := service.NewTab(req.URL)
		if err != nil {
			if err == browser.ErrBrowserUnavailable {
				return api.NewAppError("browser_unavailable", "browser binary not found", http.StatusServiceUnavailable)
			}
			return api.NewAppError("tab_new_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"index": index})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func BrowserSwitchTabHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req tabSwitchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if err := service.SwitchTab(req.Index); err != nil {
			if err == browser.ErrBrowserUnavailable {
				return api.NewAppError("browser_unavailable", "browser binary not found", http.StatusServiceUnavailable)
			}
			return api.NewAppError("tab_switch_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"index": req.Index})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func BrowserCloseTabHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req tabCloseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if err := service.CloseTab(req.Index); err != nil {
			if err == browser.ErrBrowserUnavailable {
				return api.NewAppError("browser_unavailable", "browser binary not found", http.StatusServiceUnavailable)
			}
			return api.NewAppError("tab_close_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"closed": true})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func BrowserTabListHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		tabs, err := service.TabList()
		if err != nil {
			if err == browser.ErrBrowserUnavailable {
				return api.NewAppError("browser_unavailable", "browser binary not found", http.StatusServiceUnavailable)
			}
			return api.NewAppError("tab_list_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"tabs": tabs})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func BrowserDownloadListHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		downloads := service.DownloadList()
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"downloads": downloads})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func BrowserPressKeyHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req pressKeyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if strings.TrimSpace(req.Keys) == "" {
			return api.NewAppError("bad_request", "keys is required", http.StatusBadRequest)
		}
		if err := service.PressKey(req.Keys); err != nil {
			if err == browser.ErrBrowserUnavailable {
				return api.NewAppError("browser_unavailable", "browser binary not found", http.StatusServiceUnavailable)
			}
			return api.NewAppError("press_key_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"pressed": true})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}
