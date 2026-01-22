package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"open-sandbox/internal/api"
	"open-sandbox/internal/api/handlers"
)

func TestShellExecEcho(t *testing.T) {
	router := api.NewRouter()
	handlers.RegisterShellRoutes(router)

	server := httptest.NewServer(router)
	defer server.Close()

	body, err := json.Marshal(map[string]any{
		"command": "cmd",
		"args":    []string{"/c", "echo", "test"},
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	resp, err := http.Post(server.URL+"/v1/shell/exec", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	data, ok := payload["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data payload")
	}
	stdout, _ := data["stdout"].(string)
	if !strings.Contains(stdout, "test") {
		t.Fatalf("expected stdout to contain test, got %q", stdout)
	}
}
