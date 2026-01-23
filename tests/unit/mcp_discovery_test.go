package unit

import (
	"context"
	"encoding/json"
	"testing"

	"open-sandbox/internal/mcp"
)

func TestStandardDiscoveryMethods(t *testing.T) {
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
	})

	server := mcp.NewServer(registry, nil, nil)

	initReq := mcp.Request{
		JSONRPC: mcp.JSONRPCVersion,
		ID:      json.RawMessage("1"),
		Method:  mcp.MethodInitialize,
		Params:  json.RawMessage(`{"protocolVersion":"1.0"}`),
	}
	initResp := server.HandleRequest(context.Background(), initReq)
	if initResp.Error != nil {
		t.Fatalf("initialize error: %+v", initResp.Error)
	}
	var initResult mcp.InitializeResult
	initBytes, err := json.Marshal(initResp.Result)
	if err != nil {
		t.Fatalf("marshal init result: %v", err)
	}
	if err := json.Unmarshal(initBytes, &initResult); err != nil {
		t.Fatalf("unmarshal init result: %v", err)
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

	listReq := mcp.Request{
		JSONRPC: mcp.JSONRPCVersion,
		ID:      json.RawMessage("2"),
		Method:  mcp.MethodToolsList,
		Params:  json.RawMessage(`{}`),
	}
	listResp := server.HandleRequest(context.Background(), listReq)
	if listResp.Error != nil {
		t.Fatalf("tools/list error: %+v", listResp.Error)
	}
	var listResult mcp.ToolsListResult
	listBytes, err := json.Marshal(listResp.Result)
	if err != nil {
		t.Fatalf("marshal list result: %v", err)
	}
	if err := json.Unmarshal(listBytes, &listResult); err != nil {
		t.Fatalf("unmarshal list result: %v", err)
	}
	if len(listResult.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(listResult.Tools))
	}
	if len(listResult.Tools[0].InputSchema) == 0 {
		t.Fatalf("expected input schema to be present")
	}
}
