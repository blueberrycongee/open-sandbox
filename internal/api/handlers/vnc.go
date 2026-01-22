package handlers

import (
	"encoding/json"
	"net/http"

	"open-sandbox/internal/api"
	"open-sandbox/internal/browser"
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

func RegisterVNCRoutes(router *api.Router, service *browser.Service) {
	router.Handle(http.MethodGet, "/vnc/index.html", VNCIndexHandler())
	router.Handle(http.MethodGet, "/vnc/screen.png", VNCScreenHandler(service))
	router.Handle(http.MethodPost, "/vnc/click", VNCClickHandler(service))
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
