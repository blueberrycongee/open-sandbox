package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"open-sandbox/internal/mcp"
)

func TestMCPDiscoveryOverStdio(t *testing.T) {
	registry := mcp.NewRegistry()
	registry.Register(mcp.Tool{
		Name:    "file.read",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema: mcp.ToolSchema{
			Input: mcp.JSONSchema{
				"type": "object",
			},
			Output: mcp.JSONSchema{
				"type": "object",
			},
		},
		Handler: func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
			return map[string]any{"content": "ok"}, nil
		},
	})

	server := mcp.NewServer(registry, nil, nil)

	initReq := buildRawRequest(t, mcp.MethodInitialize, map[string]any{
		"protocolVersion": mcp.SupportedProtocolVersion,
	}, json.RawMessage("1"))
	listReq := buildRawRequest(t, mcp.MethodToolsList, map[string]any{}, json.RawMessage("2"))

	var input bytes.Buffer
	input.Write(initReq)
	input.WriteByte('\n')
	input.Write(listReq)
	input.WriteByte('\n')

	var output bytes.Buffer
	if err := server.ServeStdio(&input, &output); err != nil {
		t.Fatalf("serve stdio: %v", err)
	}

	decoder := json.NewDecoder(&output)
	var initResp mcp.Response
	if err := decoder.Decode(&initResp); err != nil {
		t.Fatalf("decode initialize response: %v", err)
	}
	if initResp.Error != nil {
		t.Fatalf("initialize error: %+v", initResp.Error)
	}
	var initResult mcp.InitializeResult
	initBytes, err := json.Marshal(initResp.Result)
	if err != nil {
		t.Fatalf("marshal initialize result: %v", err)
	}
	if err := json.Unmarshal(initBytes, &initResult); err != nil {
		t.Fatalf("unmarshal initialize result: %v", err)
	}
	if initResult.ProtocolVersion != mcp.SupportedProtocolVersion {
		t.Fatalf("expected protocol version %q, got %q", mcp.SupportedProtocolVersion, initResult.ProtocolVersion)
	}
	if initResult.ServerInfo.Name == "" || initResult.ServerInfo.Version == "" {
		t.Fatalf("expected serverInfo to be populated")
	}
	if initResult.Capabilities.Tools == nil {
		t.Fatalf("expected capabilities.tools to be present")
	}

	var listResp mcp.Response
	if err := decoder.Decode(&listResp); err != nil {
		t.Fatalf("decode tools/list response: %v", err)
	}
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
	if len(listResult.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(listResult.Tools))
	}
	if listResult.Tools[0].Name != "file.read" {
		t.Fatalf("expected file.read tool, got %q", listResult.Tools[0].Name)
	}
	if len(listResult.Tools[0].Schema.Input) == 0 {
		t.Fatalf("expected input schema to be present")
	}
}

func buildRawRequest(t *testing.T, method string, params any, id json.RawMessage) []byte {
	t.Helper()
	payload := map[string]any{
		"jsonrpc": mcp.JSONRPCVersion,
		"id":      id,
		"method":  method,
		"params":  params,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	return body
}
