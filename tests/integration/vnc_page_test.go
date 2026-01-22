package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"open-sandbox/internal/api"
	"open-sandbox/internal/api/handlers"
)

func TestVNCPageAccessible(t *testing.T) {
	router := api.NewRouter()
	handlers.RegisterVNCRoutes(router)

	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/vnc/index.html")
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
