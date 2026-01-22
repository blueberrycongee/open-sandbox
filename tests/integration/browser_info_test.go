package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"open-sandbox/internal/api"
	"open-sandbox/internal/api/handlers"
	"open-sandbox/internal/browser"
	"open-sandbox/pkg/types"
)

func TestBrowserInfo(t *testing.T) {
	t.Setenv("SANDBOX_BROWSER_CDP", "ws://127.0.0.1:9222/devtools/browser/mock")
	service := browser.NewService(os.Getenv("SANDBOX_BROWSER_CDP"))

	router := api.NewRouter()
	handlers.RegisterBrowserRoutes(router, service)

	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/v1/browser/info")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var payload types.Response
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	data, ok := payload.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected data to be an object")
	}
	if data["cdp_address"] != "ws://127.0.0.1:9222/devtools/browser/mock" {
		t.Fatalf("unexpected cdp address: %v", data["cdp_address"])
	}
}
