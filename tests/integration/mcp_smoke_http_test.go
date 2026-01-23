package integration

import (
	"encoding/json"
	"testing"

	"open-sandbox/internal/config"
	"open-sandbox/internal/mcp"
)

func TestMCPSmokeHTTP(t *testing.T) {
	t.Setenv("MCP_AUTH_ENABLED", "false")
	if err := config.EnsureWorkspace(); err != nil {
		t.Fatalf("ensure workspace: %v", err)
	}

	server := newMCPTestServer(t)

	initBody := buildMCPRequest(t, mcp.MethodInitialize, map[string]any{
		"protocolVersion": mcp.SupportedProtocolVersion,
	})
	initResp := postMCPRequest(t, server.URL, initBody)
	if initResp.Error != nil {
		t.Fatalf("initialize error: %+v", initResp.Error)
	}

	listBody := buildMCPRequest(t, mcp.MethodToolsList, map[string]any{})
	listResp := postMCPRequest(t, server.URL, listBody)
	if listResp.Error != nil {
		t.Fatalf("tools/list error: %+v", listResp.Error)
	}
	var listResult mcp.ToolsListResult
	listBytes, err := json.Marshal(listResp.Result)
	if err != nil {
		t.Fatalf("marshal tools/list result: %v", err)
	}
	if err := json.Unmarshal(listBytes, &listResult); err != nil {
		t.Fatalf("unmarshal tools/list result: %v", err)
	}
	if len(listResult.Tools) == 0 {
		t.Fatalf("expected tools list to be non-empty")
	}

	writeBody := buildMCPRequest(t, mcp.MethodToolsCall, map[string]any{
		"name": "file.write",
		"arguments": map[string]any{
			"path":    "mcp-smoke-http.txt",
			"content": "smoke",
		},
	})
	writeResp := postMCPRequest(t, server.URL, writeBody)
	if writeResp.Error != nil {
		t.Fatalf("tools/call file.write error: %+v", writeResp.Error)
	}

	readBody := buildMCPRequest(t, mcp.MethodToolsCall, map[string]any{
		"name": "file.read",
		"arguments": map[string]any{
			"path": "mcp-smoke-http.txt",
		},
	})
	readResp := postMCPRequest(t, server.URL, readBody)
	if readResp.Error != nil {
		t.Fatalf("tools/call file.read error: %+v", readResp.Error)
	}
	readResult := decodeToolCallResult(t, readResp)
	content, _ := readResult["content"].(string)
	if content != "smoke" {
		t.Fatalf("expected content %q, got %q", "smoke", content)
	}
}
