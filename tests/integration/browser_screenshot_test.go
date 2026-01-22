package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"open-sandbox/internal/api"
	"open-sandbox/internal/api/handlers"
	"open-sandbox/internal/browser"
	"open-sandbox/internal/config"
)

func TestBrowserScreenshot(t *testing.T) {
	if err := config.EnsureWorkspace(); err != nil {
		t.Fatalf("ensure workspace: %v", err)
	}

	service := browser.NewService("ws://127.0.0.1:9222/devtools/browser/mock")
	router := api.NewRouter()
	handlers.RegisterBrowserRoutes(router, service)

	server := httptest.NewServer(router)
	defer server.Close()

	targetPath := filepath.Join(config.HostWorkspacePath, "screenshots", "example.png")
	_ = os.Remove(targetPath)

	body, err := json.Marshal(map[string]string{"path": targetPath})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	resp, err := http.Post(server.URL+"/v1/browser/screenshot", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if _, err := os.Stat(targetPath); err != nil {
		t.Fatalf("expected screenshot file to exist: %v", err)
	}
}
