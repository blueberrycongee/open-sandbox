package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"open-sandbox/internal/api"
	"open-sandbox/internal/api/handlers"
	"open-sandbox/internal/config"
)

func TestBrowserScreenshot(t *testing.T) {
	if err := config.EnsureWorkspace(); err != nil {
		t.Fatalf("ensure workspace: %v", err)
	}

	service := startBrowserService(t)
	router := api.NewRouter()
	handlers.RegisterBrowserRoutes(router, service)

	page := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html><body>ok</body></html>"))
	}))
	defer page.Close()

	server := httptest.NewServer(router)
	defer server.Close()

	navigateBody, err := json.Marshal(map[string]string{"url": page.URL})
	if err != nil {
		t.Fatalf("marshal navigate: %v", err)
	}
	navigateResp, err := http.Post(server.URL+"/v1/browser/navigate", "application/json", bytes.NewReader(navigateBody))
	if err != nil {
		t.Fatalf("navigate request failed: %v", err)
	}
	if navigateResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(navigateResp.Body)
		navigateResp.Body.Close()
		t.Fatalf("navigate status %d: %s", navigateResp.StatusCode, string(bodyBytes))
	}
	navigateResp.Body.Close()

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
