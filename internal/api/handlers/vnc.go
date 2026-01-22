package handlers

import (
	"net/http"

	"open-sandbox/internal/api"
)

const vncHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>open-sandbox VNC</title>
</head>
<body>
  <h1>VNC Placeholder</h1>
  <p>VNC takeover is not configured yet. Provide a VNC endpoint to enable it.</p>
</body>
</html>`

func RegisterVNCRoutes(router *api.Router) {
	router.Handle(http.MethodGet, "/vnc/index.html", VNCIndexHandler)
}

func VNCIndexHandler(w http.ResponseWriter, r *http.Request) *api.AppError {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(vncHTML))
	return nil
}
