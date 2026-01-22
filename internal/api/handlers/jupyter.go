package handlers

import (
	"net/http"

	"open-sandbox/internal/api"
)

const jupyterHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>open-sandbox Jupyter</title>
</head>
<body>
  <h1>Jupyter Placeholder</h1>
  <p>Jupyter Lab is not configured yet. Provide a Jupyter endpoint to enable it.</p>
</body>
</html>`

func RegisterJupyterRoutes(router *api.Router) {
	router.Handle(http.MethodGet, "/jupyter", JupyterHandler)
}

func JupyterHandler(w http.ResponseWriter, r *http.Request) *api.AppError {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(jupyterHTML))
	return nil
}
