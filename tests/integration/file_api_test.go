package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"open-sandbox/internal/api"
	"open-sandbox/internal/api/handlers"
	"open-sandbox/internal/config"
)

func TestFileCRUDSearchReplace(t *testing.T) {
	if err := config.EnsureWorkspace(); err != nil {
		t.Fatalf("ensure workspace: %v", err)
	}

	router := api.NewRouter()
	handlers.RegisterFileRoutes(router)

	server := httptest.NewServer(router)
	defer server.Close()

	targetPath := filepath.Join(config.WorkspacePath(), "test-file.txt")

	writeBody, err := json.Marshal(map[string]string{
		"path":    targetPath,
		"content": "hello world",
	})
	if err != nil {
		t.Fatalf("marshal write: %v", err)
	}

	resp, err := http.Post(server.URL+"/v1/file/write", "application/json", bytes.NewReader(writeBody))
	if err != nil {
		t.Fatalf("write request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("write status %d", resp.StatusCode)
	}
	resp.Body.Close()

	readBody, err := json.Marshal(map[string]string{"path": targetPath})
	if err != nil {
		t.Fatalf("marshal read: %v", err)
	}
	resp, err = http.Post(server.URL+"/v1/file/read", "application/json", bytes.NewReader(readBody))
	if err != nil {
		t.Fatalf("read request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("read status %d", resp.StatusCode)
	}
	resp.Body.Close()

	listResp, err := http.Get(server.URL + "/v1/file/list?path=" + filepath.ToSlash(config.WorkspacePath()))
	if err != nil {
		t.Fatalf("list request failed: %v", err)
	}
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("list status %d", listResp.StatusCode)
	}
	listResp.Body.Close()

	searchBody, err := json.Marshal(map[string]string{
		"path":  targetPath,
		"query": "world",
	})
	if err != nil {
		t.Fatalf("marshal search: %v", err)
	}
	resp, err = http.Post(server.URL+"/v1/file/search", "application/json", bytes.NewReader(searchBody))
	if err != nil {
		t.Fatalf("search request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("search status %d", resp.StatusCode)
	}
	resp.Body.Close()

	replaceBody, err := json.Marshal(map[string]string{
		"path":    targetPath,
		"search":  "world",
		"replace": "there",
	})
	if err != nil {
		t.Fatalf("marshal replace: %v", err)
	}
	resp, err = http.Post(server.URL+"/v1/file/replace", "application/json", bytes.NewReader(replaceBody))
	if err != nil {
		t.Fatalf("replace request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("replace status %d", resp.StatusCode)
	}
	resp.Body.Close()
}
