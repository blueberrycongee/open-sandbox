package integration

import (
	"bytes"
	"encoding/json"
	"runtime"
	"strings"
	"testing"

	"open-sandbox/internal/config"
	"open-sandbox/internal/mcp"
	"open-sandbox/internal/mcp/tools"
)

func TestMCPToolsCallOverStdio(t *testing.T) {
	if err := config.EnsureWorkspace(); err != nil {
		t.Fatalf("ensure workspace: %v", err)
	}

	registry := mcp.NewRegistry()
	registry.Register(mcp.Tool{
		Name:    "file.write",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Handler: tools.FileWrite(),
	})
	registry.Register(mcp.Tool{
		Name:    "file.read",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Handler: tools.FileRead(),
	})
	registry.Register(mcp.Tool{
		Name:    "shell.exec",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "exec",
		},
		Handler: tools.ShellExec(),
	})

	server := mcp.NewServer(registry, nil, nil)

	writeReq := buildToolsCallRequest(t, json.RawMessage("1"), "file.write", map[string]any{
		"path":    "mcp-tools-call-stdio.txt",
		"content": "hello",
	})
	readReq := buildToolsCallRequest(t, json.RawMessage("2"), "file.read", map[string]any{
		"path": "mcp-tools-call-stdio.txt",
	})
	command, args := platformEchoCommandStdio()
	execReq := buildToolsCallRequest(t, json.RawMessage("3"), "shell.exec", map[string]any{
		"command": command,
		"args":    args,
	})

	var input bytes.Buffer
	input.Write(writeReq)
	input.WriteByte('\n')
	input.Write(readReq)
	input.WriteByte('\n')
	input.Write(execReq)
	input.WriteByte('\n')

	var output bytes.Buffer
	if err := server.ServeStdio(&input, &output); err != nil {
		t.Fatalf("serve stdio: %v", err)
	}

	decoder := json.NewDecoder(&output)
	writeResp := decodeResponse(t, decoder)
	if writeResp.Error != nil {
		t.Fatalf("tools/call file.write error: %+v", writeResp.Error)
	}

	readResp := decodeResponse(t, decoder)
	if readResp.Error != nil {
		t.Fatalf("tools/call file.read error: %+v", readResp.Error)
	}
	readResult := decodeToolResult(t, readResp)
	content, _ := readResult["content"].(string)
	if content != "hello" {
		t.Fatalf("expected content %q, got %q", "hello", content)
	}

	execResp := decodeResponse(t, decoder)
	if execResp.Error != nil {
		t.Fatalf("tools/call shell.exec error: %+v", execResp.Error)
	}
	execResult := decodeToolResult(t, execResp)
	stdout, _ := execResult["Stdout"].(string)
	if !strings.Contains(stdout, "test") {
		t.Fatalf("expected stdout to contain %q, got %q", "test", stdout)
	}
}

func buildToolsCallRequest(t *testing.T, id json.RawMessage, name string, arguments map[string]any) []byte {
	t.Helper()
	payload := map[string]any{
		"jsonrpc": mcp.JSONRPCVersion,
		"id":      id,
		"method":  mcp.MethodToolsCall,
		"params": map[string]any{
			"name":      name,
			"arguments": arguments,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	return body
}

func decodeResponse(t *testing.T, decoder *json.Decoder) mcp.Response {
	t.Helper()
	var resp mcp.Response
	if err := decoder.Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return resp
}

func decodeToolResult(t *testing.T, resp mcp.Response) map[string]any {
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

func platformEchoCommandStdio() (string, []string) {
	if runtime.GOOS == "windows" {
		return "cmd", []string{"/C", "echo test"}
	}
	return "sh", []string{"-c", "echo test"}
}
