package integration

import (
	"encoding/json"
	"runtime"
	"strings"
	"testing"

	"open-sandbox/internal/config"
	"open-sandbox/internal/mcp"
)

func TestMCPToolsCallOverHTTP(t *testing.T) {
	t.Setenv("MCP_AUTH_ENABLED", "false")
	if err := config.EnsureWorkspace(); err != nil {
		t.Fatalf("ensure workspace: %v", err)
	}

	server := newMCPTestServer(t)

	writeBody := buildMCPRequest(t, mcp.MethodToolsCall, map[string]any{
		"name": "file.write",
		"arguments": map[string]any{
			"path":    "mcp-tools-call.txt",
			"content": "hello",
		},
	})
	writeResp := postMCPRequest(t, server.URL, writeBody)
	if writeResp.Error != nil {
		t.Fatalf("tools/call file.write error: %+v", writeResp.Error)
	}

	readBody := buildMCPRequest(t, mcp.MethodToolsCall, map[string]any{
		"name": "file.read",
		"arguments": map[string]any{
			"path": "mcp-tools-call.txt",
		},
	})
	readResp := postMCPRequest(t, server.URL, readBody)
	if readResp.Error != nil {
		t.Fatalf("tools/call file.read error: %+v", readResp.Error)
	}
	readResult := decodeToolCallResult(t, readResp)
	content, _ := readResult["content"].(string)
	if content != "hello" {
		t.Fatalf("expected content %q, got %q", "hello", content)
	}

	command, args := platformEchoCommandToolsCall()
	execBody := buildMCPRequest(t, mcp.MethodToolsCall, map[string]any{
		"name": "shell.exec",
		"arguments": map[string]any{
			"command": command,
			"args":    args,
		},
	})
	execResp := postMCPRequest(t, server.URL, execBody)
	if execResp.Error != nil {
		t.Fatalf("tools/call shell.exec error: %+v", execResp.Error)
	}
	execResult := decodeToolCallResult(t, execResp)
	stdout, _ := execResult["Stdout"].(string)
	if !strings.Contains(stdout, "test") {
		t.Fatalf("expected stdout to contain %q, got %q", "test", stdout)
	}
}

func decodeToolCallResult(t *testing.T, resp mcp.Response) map[string]any {
	t.Helper()
	var callResult mcp.ToolCallResult
	resultBytes, err := json.Marshal(resp.Result)
	if err != nil {
		t.Fatalf("marshal tools/call result: %v", err)
	}
	if err := json.Unmarshal(resultBytes, &callResult); err != nil {
		t.Fatalf("unmarshal tools/call result: %v", err)
	}
	payloadBytes, err := json.Marshal(callResult.Result)
	if err != nil {
		t.Fatalf("marshal tool result: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		t.Fatalf("unmarshal tool result: %v", err)
	}
	return payload
}

func platformEchoCommandToolsCall() (string, []string) {
	if runtime.GOOS == "windows" {
		return "cmd", []string{"/C", "echo test"}
	}
	return "sh", []string{"-c", "echo test"}
}
