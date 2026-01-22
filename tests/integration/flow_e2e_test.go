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
	"open-sandbox/internal/browser"
	"open-sandbox/internal/config"
)

func TestEndToEndFlow(t *testing.T) {
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

	service := browser.NewService("ws://127.0.0.1:9222/devtools/browser/mock")
	router := api.NewRouter()
	handlers.RegisterBrowserRoutes(router, service)
	handlers.RegisterFileRoutes(router)
	handlers.RegisterCodeExecRoutes(router)

	server := httptest.NewServer(router)
	defer server.Close()

	screenshotPath := filepath.Join(config.HostWorkspacePath, "screenshots", "example.png")
	_ = os.Remove(screenshotPath)

	screenshotBody, _ := json.Marshal(map[string]string{"path": screenshotPath})
	resp, err := http.Post(server.URL+"/v1/browser/screenshot", "application/json", bytes.NewReader(screenshotBody))
	if err != nil {
		t.Fatalf("screenshot request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("screenshot status %d", resp.StatusCode)
	}
	resp.Body.Close()

	readBody, _ := json.Marshal(map[string]string{"path": screenshotPath})
	resp, err = http.Post(server.URL+"/v1/file/read", "application/json", bytes.NewReader(readBody))
	if err != nil {
		t.Fatalf("file read failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("file read status %d", resp.StatusCode)
	}
	resp.Body.Close()

	outputPath := filepath.Join(config.HostWorkspacePath, "output.txt")
	_ = os.Remove(outputPath)

	var args []string
	switch runtime {
	case "python":
		args = []string{"-c", "open(r'" + outputPath + "', 'w').write('ok')"}
	case "node":
		args = []string{"-e", "require('fs').writeFileSync('" + outputPath + "', 'ok')"}
	}

	codeBody, _ := json.Marshal(map[string]any{
		"runtime": runtime,
		"args":    args,
	})
	resp, err = http.Post(server.URL+"/v1/code/exec", "application/json", bytes.NewReader(codeBody))
	if err != nil {
		t.Fatalf("code exec failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("code exec status %d", resp.StatusCode)
	}
	resp.Body.Close()

	readOutputBody, _ := json.Marshal(map[string]string{"path": outputPath})
	resp, err = http.Post(server.URL+"/v1/file/read", "application/json", bytes.NewReader(readOutputBody))
	if err != nil {
		t.Fatalf("output read failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("output read status %d", resp.StatusCode)
	}
	resp.Body.Close()
}
