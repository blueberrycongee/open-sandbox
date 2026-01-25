package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"open-sandbox/internal/api"
	"open-sandbox/internal/browser"
	"open-sandbox/pkg/types"
)

const vncHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>open-sandbox VNC</title>
  <style>
    body { margin: 0; font-family: sans-serif; background: #111; color: #eee; }
    #toolbar { padding: 8px 12px; background: #1e1e1e; display: flex; gap: 12px; align-items: center; }
    #screen { display: block; margin: 0 auto; max-width: 100vw; max-height: calc(100vh - 48px); cursor: crosshair; }
    #status { font-size: 12px; opacity: 0.8; }
  </style>
</head>
<body>
  <div id="toolbar">
    <strong>VNC Takeover</strong>
    <span id="status">connecting...</span>
  </div>
  <img id="screen" alt="live screen" />
  <script>
    const statusEl = document.getElementById('status');
    const screen = document.getElementById('screen');
    let lastUpdate = 0;

    async function refresh() {
      try {
        const ts = Date.now();
        screen.src = '/vnc/screen.png?ts=' + ts;
        lastUpdate = ts;
        statusEl.textContent = 'live';
      } catch (err) {
        statusEl.textContent = 'error';
      }
    }

    screen.addEventListener('click', async (event) => {
      const rect = screen.getBoundingClientRect();
      const scaleX = screen.naturalWidth / rect.width;
      const scaleY = screen.naturalHeight / rect.height;
      const x = (event.clientX - rect.left) * scaleX;
      const y = (event.clientY - rect.top) * scaleY;
      await fetch('/vnc/click', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ x, y })
      });
    });

    setInterval(refresh, 800);
    refresh();
  </script>
</body>
</html>`

type clickRequest struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type keyboardRequest struct {
	Keys string `json:"keys"`
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

type tabSwitchRequest struct {
	Index int `json:"index"`
}

type tabCloseRequest struct {
	Index int `json:"index"`
}

type tabNewRequest struct {
	URL string `json:"url"`
}

func RegisterVNCRoutes(router *api.Router, service *browser.Service) {
	router.Handle(http.MethodGet, "/vnc/index.html", VNCIndexHandler())
	router.Handle(http.MethodGet, "/vnc/screen.png", VNCScreenHandler(service))
	router.Handle(http.MethodPost, "/vnc/click", VNCClickHandler(service))
	router.Handle(http.MethodPost, "/vnc/keyboard", VNCKeyboardHandler(service))
	router.Handle(http.MethodPost, "/vnc/form_input_fill", VNCFormInputFillHandler(service))
	router.Handle(http.MethodPost, "/vnc/element_select", VNCElementSelectHandler(service))
	router.Handle(http.MethodPost, "/vnc/scroll", VNCScrollHandler(service))
	router.Handle(http.MethodPost, "/vnc/evaluate", VNCEvaluateHandler(service))
	router.Handle(http.MethodPost, "/vnc/tab/new", VNCTabNewHandler(service))
	router.Handle(http.MethodPost, "/vnc/tab/switch", VNCTabSwitchHandler(service))
	router.Handle(http.MethodGet, "/vnc/tab/list", VNCTabListHandler(service))
	router.Handle(http.MethodPost, "/vnc/tab/close", VNCTabCloseHandler(service))
	router.Handle(http.MethodGet, "/vnc/downloads", VNCDownloadListHandler(service))
}

func VNCIndexHandler() api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(vncHTML))
		return nil
	}
}

func VNCScreenHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		data, err := service.ScreenshotPNG()
		if err != nil {
			return api.NewAppError("vnc_unavailable", err.Error(), http.StatusServiceUnavailable)
		}
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
		return nil
	}
}

func VNCClickHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req clickRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if err := service.Click(req.X, req.Y); err != nil {
			return api.NewAppError("click_failed", err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}

func VNCKeyboardHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req keyboardRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if req.Keys == "" {
			return api.NewAppError("bad_request", "keys is required", http.StatusBadRequest)
		}
		if err := service.PressKey(req.Keys); err != nil {
			return api.NewAppError("keyboard_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"pressed": true})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func VNCFormInputFillHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req formInputFillRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if req.Selector == "" {
			return api.NewAppError("bad_request", "selector is required", http.StatusBadRequest)
		}
		if err := service.FormInputFill(req.Selector, req.Value); err != nil {
			return api.NewAppError("form_input_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"filled": true})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func VNCElementSelectHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req elementSelectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if req.Selector == "" {
			return api.NewAppError("bad_request", "selector is required", http.StatusBadRequest)
		}
		if err := service.ElementSelect(req.Selector, req.Value); err != nil {
			return api.NewAppError("select_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"selected": true})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func VNCScrollHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req scrollRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if err := service.Scroll(req.X, req.Y); err != nil {
			return api.NewAppError("scroll_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"scrolled": true})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func VNCEvaluateHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req evaluateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if req.Expression == "" {
			return api.NewAppError("bad_request", "expression is required", http.StatusBadRequest)
		}
		result, err := service.Evaluate(req.Expression)
		if err != nil {
			return api.NewAppError("evaluate_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"result": result})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func VNCTabNewHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req tabNewRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		index, err := service.NewTab(req.URL)
		if err != nil {
			return api.NewAppError("tab_new_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"index": index})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func VNCTabSwitchHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req tabSwitchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if err := service.SwitchTab(req.Index); err != nil {
			return api.NewAppError("tab_switch_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"index": req.Index})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func VNCTabListHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		tabs, err := service.TabList()
		if err != nil {
			return api.NewAppError("tab_list_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"tabs": tabs})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func VNCTabCloseHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		var req tabCloseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return api.NewAppError("bad_request", "invalid request body", http.StatusBadRequest)
		}
		if err := service.CloseTab(req.Index); err != nil {
			return api.NewAppError("tab_close_failed", err.Error(), http.StatusInternalServerError)
		}
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"closed": true})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}

func VNCDownloadListHandler(service *browser.Service) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		downloads := service.DownloadList()
		if err := api.WriteJSON(w, http.StatusOK, types.Ok(map[string]any{"downloads": downloads})); err != nil {
			return api.NewAppError(api.CodeInternalError, "internal error", http.StatusInternalServerError)
		}
		return nil
	}
}
