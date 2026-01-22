package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"open-sandbox/internal/api"
	"open-sandbox/internal/api/handlers"
)

func TestCodeServerEntry(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<html>code-server</html>"))
	}))
	defer upstream.Close()

	router := api.NewRouter()
	handlers.RegisterCodeServerRoutes(router, upstream.URL)

	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/code-server/?folder=/workspace")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Fatalf("expected text/html content type, got %q", contentType)
	}
}
