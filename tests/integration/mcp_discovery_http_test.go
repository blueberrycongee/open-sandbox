package integration

import (
	"encoding/json"
	"testing"

	"open-sandbox/internal/mcp"
)

func TestMCPDiscoveryOverHTTP(t *testing.T) {
	server := newMCPTestServer(t)

	initBody := buildMCPRequest(t, mcp.MethodInitialize, map[string]any{
		"protocol_version": mcp.SupportedProtocolVersion,
	})
	initResp := postMCPRequest(t, server.URL, initBody)
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
	var fileTool *mcp.ToolInfo
	for i := range listResult.Tools {
		if listResult.Tools[i].Name == "file.read" {
			fileTool = &listResult.Tools[i]
			break
		}
	}
	if fileTool == nil {
		t.Fatalf("expected file.read tool in list")
	}
	if fileTool.Version == "" {
		t.Fatalf("expected file.read version to be set")
	}
	if len(fileTool.Schema.Input) == 0 {
		t.Fatalf("expected file.read input schema to be present")
	}
}
