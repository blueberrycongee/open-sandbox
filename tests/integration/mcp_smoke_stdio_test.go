package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"open-sandbox/internal/mcp"
)

func TestMCPSmokeStdio(t *testing.T) {
	repoRoot := repoRootPath(t)
	tempRoot := t.TempDir()
	workspace := filepath.Join(tempRoot, "workspace")
	if err := os.MkdirAll(workspace, 0755); err != nil {
		t.Fatalf("create workspace: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", "./cmd/mcp")
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(),
		"SANDBOX_ROOT="+tempRoot,
		"SANDBOX_WORKSPACE="+workspace,
		"MCP_AUTH_ENABLED=false",
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("start cmd/mcp: %v", err)
	}

	initReq := buildRawRequest(t, mcp.MethodInitialize, map[string]any{
		"protocolVersion": mcp.SupportedProtocolVersion,
	}, json.RawMessage("1"))
	listReq := buildRawRequest(t, mcp.MethodToolsList, map[string]any{}, json.RawMessage("2"))
	writeReq := buildRawRequest(t, mcp.MethodToolsCall, map[string]any{
		"name": "file.write",
		"arguments": map[string]any{
			"path":    "mcp-smoke-stdio.txt",
			"content": "smoke",
		},
	}, json.RawMessage("3"))
	readReq := buildRawRequest(t, mcp.MethodToolsCall, map[string]any{
		"name": "file.read",
		"arguments": map[string]any{
			"path": "mcp-smoke-stdio.txt",
		},
	}, json.RawMessage("4"))

	if _, err := stdin.Write(append(initReq, '\n')); err != nil {
		t.Fatalf("write initialize: %v", err)
	}
	if _, err := stdin.Write(append(listReq, '\n')); err != nil {
		t.Fatalf("write tools/list: %v", err)
	}
	if _, err := stdin.Write(append(writeReq, '\n')); err != nil {
		t.Fatalf("write tools/call file.write: %v", err)
	}
	if _, err := stdin.Write(append(readReq, '\n')); err != nil {
		t.Fatalf("write tools/call file.read: %v", err)
	}
	if err := stdin.Close(); err != nil {
		t.Fatalf("close stdin: %v", err)
	}

	decoder := json.NewDecoder(stdout)
	initResp := decodeResponse(t, decoder)
	if initResp.Error != nil {
		t.Fatalf("initialize error: %+v", initResp.Error)
	}
	listResp := decodeResponse(t, decoder)
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
	if content != "smoke" {
		t.Fatalf("expected content %q, got %q", "smoke", content)
	}

	if err := cmd.Wait(); err != nil {
		t.Fatalf("cmd/mcp failed: %v: %s", err, stderr.String())
	}
	if ctx.Err() != nil {
		t.Fatalf("cmd/mcp timed out: %v", ctx.Err())
	}
}

func repoRootPath(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(cwd, "..", ".."))
	return root
}
