package handlers

import (
	"net/http"

	"open-sandbox/internal/api"
)

const codeServerHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>open-sandbox Code Server</title>
</head>
<body>
  <h1>Code Server Placeholder</h1>
  <p>code-server is not configured yet. Provide a code-server endpoint to enable it.</p>
</body>
</html>`

func RegisterCodeServerRoutes(router *api.Router) {
	router.Handle(http.MethodGet, "/code-server/", CodeServerHandler)
}

func CodeServerHandler(w http.ResponseWriter, r *http.Request) *api.AppError {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(codeServerHTML))
	return nil
}
