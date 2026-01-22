package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
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

	command := "cmd"
	args := []string{"/c", "echo", "test"}
	if runtime.GOOS != "windows" {
		command = "/usr/bin/sh"
		args = []string{"-c", "echo test"}
	}

	body, err := json.Marshal(map[string]any{
		"command": command,
		"args":    args,
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
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status 200, got %d: %s", resp.StatusCode, string(bodyBytes))
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
