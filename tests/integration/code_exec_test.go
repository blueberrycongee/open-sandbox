package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"open-sandbox/internal/api"
	"open-sandbox/internal/api/handlers"
	"open-sandbox/internal/config"
)

func TestCodeExec(t *testing.T) {
	runtime := ""
	if _, err := exec.LookPath("python"); err == nil {
		runtime = "python"
	} else if _, err := exec.LookPath("node"); err == nil {
		runtime = "node"
	} else {
		t.Skip("python or node runtime not available")
	}

	if err := config.EnsureWorkspace(); err != nil {
		t.Fatalf("ensure workspace: %v", err)
	}

	router := api.NewRouter()
	handlers.RegisterCodeExecRoutes(router)

	server := httptest.NewServer(router)
	defer server.Close()

	outputPath := filepath.Join(config.WorkspacePath(), "output.txt")
	_ = os.Remove(outputPath)

	var args []string
	switch runtime {
	case "python":
		args = []string{"-c", "open(r'" + outputPath + "', 'w').write('ok')"}
	case "node":
		args = []string{"-e", "require('fs').writeFileSync('" + outputPath + "', 'ok')"}
	}

	body, err := json.Marshal(map[string]any{
		"runtime": runtime,
		"args":    args,
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	resp, err := http.Post(server.URL+"/v1/code/exec", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if string(content) != "ok" {
		t.Fatalf("unexpected output content: %q", string(content))
	}
}
