package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"open-sandbox/internal/api"
	"open-sandbox/internal/api/handlers"
	"open-sandbox/internal/config"
	"open-sandbox/internal/mcp"
)

func TestMCPToolCallsAndWorkspaceBoundary(t *testing.T) {
	t.Setenv("MCP_AUTH_ENABLED", "false")

	if err := config.EnsureWorkspace(); err != nil {
		t.Fatalf("ensure workspace: %v", err)
	}

	router := api.NewRouter()
	handlers.RegisterMCPRoutes(router, nil, nil)
	server := httptest.NewServer(router)
	defer server.Close()

	writeParams, _ := json.Marshal(map[string]any{
		"path":    "mcp-tool.txt",
		"content": "hello",
	})
	writeReq := mcp.Request{
		JSONRPC: mcp.JSONRPCVersion,
		ID:      json.RawMessage("1"),
		Method:  "file.write",
		Params:  writeParams,
	}
	writeResp := callMCP(t, server.URL, writeReq)
	if writeResp.Error != nil {
		t.Fatalf("file.write error: %+v", writeResp.Error)
	}

	readParams, _ := json.Marshal(map[string]any{"path": "mcp-tool.txt"})
	readReq := mcp.Request{
		JSONRPC: mcp.JSONRPCVersion,
		ID:      json.RawMessage("2"),
		Method:  "file.read",
		Params:  readParams,
	}
	readResp := callMCP(t, server.URL, readReq)
	if readResp.Error != nil {
		t.Fatalf("file.read error: %+v", readResp.Error)
	}
	readResult := mustMap(t, readResp.Result)
	content, _ := readResult["content"].(string)
	if content != "hello" {
		t.Fatalf("expected content %q, got %q", "hello", content)
	}

	command, args := platformEchoCommand()
	execParams, _ := json.Marshal(map[string]any{
		"command": command,
		"args":    args,
	})
	execReq := mcp.Request{
		JSONRPC: mcp.JSONRPCVersion,
		ID:      json.RawMessage("3"),
		Method:  "shell.exec",
		Params:  execParams,
	}
	execResp := callMCP(t, server.URL, execReq)
	if execResp.Error != nil {
		t.Fatalf("shell.exec error: %+v", execResp.Error)
	}
	execResult := mustMap(t, execResp.Result)
	stdout, _ := execResult["stdout"].(string)
	if !strings.Contains(stdout, "test") {
		t.Fatalf("expected stdout to contain %q, got %q", "test", stdout)
	}

	outsidePath := filepath.Join(config.RootPath(), "outside.txt")
	badParams, _ := json.Marshal(map[string]any{"path": outsidePath})
	badReq := mcp.Request{
		JSONRPC: mcp.JSONRPCVersion,
		ID:      json.RawMessage("4"),
		Method:  "file.read",
		Params:  badParams,
	}
	badResp := callMCP(t, server.URL, badReq)
	if badResp.Error == nil || badResp.Error.Code != mcp.ErrInvalidParams {
		t.Fatalf("expected invalid params error, got %+v", badResp.Error)
	}
}

func callMCP(t *testing.T, serverURL string, req mcp.Request) mcp.Response {
	t.Helper()
	body, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	resp, err := http.Post(serverURL+"/mcp", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("http post: %v", err)
	}
	defer resp.Body.Close()

	var mcpResp mcp.Response
	if err := json.NewDecoder(resp.Body).Decode(&mcpResp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return mcpResp
}

func mustMap(t *testing.T, value any) map[string]any {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	return decoded
}

func platformEchoCommand() (string, []string) {
	if runtime.GOOS == "windows" {
		return "cmd", []string{"/C", "echo test"}
	}
	return "sh", []string{"-c", "echo test"}
}
