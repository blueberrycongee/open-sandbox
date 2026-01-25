package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

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

type actionEnvelope struct {
	ActionType string `json:"action_type"`
}

type moveToAction struct {
	ActionType string  `json:"action_type"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
}

type moveRelAction struct {
	ActionType string  `json:"action_type"`
	XOffset    float64 `json:"x_offset"`
	YOffset    float64 `json:"y_offset"`
}

type clickAction struct {
	ActionType string   `json:"action_type"`
	X          *float64 `json:"x"`
	Y          *float64 `json:"y"`
	Button     string   `json:"button"`
	NumClicks  int      `json:"num_clicks"`
}

type mouseButtonAction struct {
	ActionType string `json:"action_type"`
	Button     string `json:"button"`
}

type rightClickAction struct {
	ActionType string   `json:"action_type"`
	X          *float64 `json:"x"`
	Y          *float64 `json:"y"`
}

type doubleClickAction struct {
	ActionType string   `json:"action_type"`
	X          *float64 `json:"x"`
	Y          *float64 `json:"y"`
}

type dragToAction struct {
	ActionType string  `json:"action_type"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
}

type dragRelAction struct {
	ActionType string  `json:"action_type"`
	XOffset    float64 `json:"x_offset"`
	YOffset    float64 `json:"y_offset"`
}

type scrollAction struct {
	ActionType string `json:"action_type"`
	DX         int    `json:"dx"`
	DY         int    `json:"dy"`
}

type typingAction struct {
	ActionType   string `json:"action_type"`
	Text         string `json:"text"`
	UseClipboard *bool  `json:"use_clipboard"`
}

type pressAction struct {
	ActionType string `json:"action_type"`
	Key        string `json:"key"`
}

type keyAction struct {
	ActionType string `json:"action_type"`
	Key        string `json:"key"`
}

type hotkeyAction struct {
	ActionType string   `json:"action_type"`
	Keys       []string `json:"keys"`
}

type waitAction struct {
	ActionType string  `json:"action_type"`
	Duration   float64 `json:"duration"`
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
	router.Handle(http.MethodPost, "/v1/browser/actions", BrowserActionsHandler(service))
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

func BrowserActionsHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		body = bytes.TrimSpace(body)
		if len(body) == 0 {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}

		if body[0] == '[' {
			var items []json.RawMessage
			if err := json.Unmarshal(body, &items); err != nil {
				return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
			}
			performed := make([]string, 0, len(items))
			for _, item := range items {
				actionName, err := executeAction(service, item)
				if err != nil {
					return api.NewAppError("action_failed", err.Error(), http.StatusInternalServerError)
				}
				performed = append(performed, actionName)
			}
			payload := map[string]any{
				"status":            "success",
				"action_performed":  "BATCH",
				"actions_performed": performed,
			}
			if err := api.WriteJSON(w, http.StatusOK, types.Ok(payload)); err != nil {
				return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
			}
			return nil
		}

		actionName, err := executeAction(service, body)
		if err != nil {
			return api.NewAppError("action_failed", err.Error(), http.StatusInternalServerError)
		}
		payload := map[string]any{
			"status":           "success",
			"action_performed": actionName,
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(payload)); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func executeAction(service *browser.Service, raw json.RawMessage) (string, error) {
	var envelope actionEnvelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return "", err
	}
	switch envelope.ActionType {
	case "MOVE_TO":
		var payload moveToAction
		if err := json.Unmarshal(raw, &payload); err != nil {
			return "", err
		}
		return envelope.ActionType, service.MoveTo(payload.X, payload.Y)
	case "MOVE_REL":
		var payload moveRelAction
		if err := json.Unmarshal(raw, &payload); err != nil {
			return "", err
		}
		return envelope.ActionType, service.MoveRel(payload.XOffset, payload.YOffset)
	case "CLICK":
		var payload clickAction
		if err := json.Unmarshal(raw, &payload); err != nil {
			return "", err
		}
		button := normalizeButton(payload.Button)
		count := payload.NumClicks
		x, y, ok := actionCoords(service, payload.X, payload.Y)
		if !ok {
			return "", errors.New("mouse position unknown")
		}
		return envelope.ActionType, service.ClickAt(x, y, button, count)
	case "MOUSE_DOWN":
		var payload mouseButtonAction
		if err := json.Unmarshal(raw, &payload); err != nil {
			return "", err
		}
		return envelope.ActionType, service.MouseDown(normalizeButton(payload.Button))
	case "MOUSE_UP":
		var payload mouseButtonAction
		if err := json.Unmarshal(raw, &payload); err != nil {
			return "", err
		}
		return envelope.ActionType, service.MouseUp(normalizeButton(payload.Button))
	case "RIGHT_CLICK":
		var payload rightClickAction
		if err := json.Unmarshal(raw, &payload); err != nil {
			return "", err
		}
		x, y, ok := actionCoords(service, payload.X, payload.Y)
		if !ok {
			return "", errors.New("mouse position unknown")
		}
		return envelope.ActionType, service.ClickAt(x, y, "right", 1)
	case "DOUBLE_CLICK":
		var payload doubleClickAction
		if err := json.Unmarshal(raw, &payload); err != nil {
			return "", err
		}
		x, y, ok := actionCoords(service, payload.X, payload.Y)
		if !ok {
			return "", errors.New("mouse position unknown")
		}
		return envelope.ActionType, service.ClickAt(x, y, "left", 2)
	case "DRAG_TO":
		var payload dragToAction
		if err := json.Unmarshal(raw, &payload); err != nil {
			return "", err
		}
		return envelope.ActionType, service.DragTo(payload.X, payload.Y)
	case "DRAG_REL":
		var payload dragRelAction
		if err := json.Unmarshal(raw, &payload); err != nil {
			return "", err
		}
		return envelope.ActionType, service.DragRel(payload.XOffset, payload.YOffset)
	case "SCROLL":
		var payload scrollAction
		if err := json.Unmarshal(raw, &payload); err != nil {
			return "", err
		}
		return envelope.ActionType, service.ScrollWheel(float64(payload.DX), float64(payload.DY))
	case "TYPING":
		var payload typingAction
		if err := json.Unmarshal(raw, &payload); err != nil {
			return "", err
		}
		if strings.TrimSpace(payload.Text) == "" {
			return "", errors.New("text is required")
		}
		return envelope.ActionType, service.PressKey(payload.Text)
	case "PRESS":
		var payload pressAction
		if err := json.Unmarshal(raw, &payload); err != nil {
			return "", err
		}
		if strings.TrimSpace(payload.Key) == "" {
			return "", errors.New("key is required")
		}
		return envelope.ActionType, service.PressSingleKey(payload.Key)
	case "KEY_DOWN":
		var payload keyAction
		if err := json.Unmarshal(raw, &payload); err != nil {
			return "", err
		}
		if strings.TrimSpace(payload.Key) == "" {
			return "", errors.New("key is required")
		}
		return envelope.ActionType, service.KeyDown(payload.Key)
	case "KEY_UP":
		var payload keyAction
		if err := json.Unmarshal(raw, &payload); err != nil {
			return "", err
		}
		if strings.TrimSpace(payload.Key) == "" {
			return "", errors.New("key is required")
		}
		return envelope.ActionType, service.KeyUp(payload.Key)
	case "HOTKEY":
		var payload hotkeyAction
		if err := json.Unmarshal(raw, &payload); err != nil {
			return "", err
		}
		if len(payload.Keys) == 0 {
			return "", errors.New("keys are required")
		}
		return envelope.ActionType, service.Hotkey(payload.Keys)
	case "WAIT":
		var payload waitAction
		if err := json.Unmarshal(raw, &payload); err != nil {
			return "", err
		}
		if payload.Duration <= 0 {
			return "", errors.New("duration must be positive")
		}
		time.Sleep(time.Duration(payload.Duration * float64(time.Second)))
		return envelope.ActionType, nil
	default:
		return "", errors.New("unsupported action type")
	}
}

func actionCoords(service *browser.Service, x *float64, y *float64) (float64, float64, bool) {
	if x != nil && y != nil {
		return *x, *y, true
	}
	return service.MousePosition()
}

func normalizeButton(button string) string {
	switch strings.ToLower(strings.TrimSpace(button)) {
	case "right":
		return "right"
	case "middle":
		return "middle"
	default:
		return "left"
	}
}
