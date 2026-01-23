package unit

import (
	"context"
	"encoding/json"
	"testing"

	"open-sandbox/internal/mcp"
)

func TestToolsCallRoutesToHandler(t *testing.T) {
	registry := mcp.NewRegistry()
	registry.Register(mcp.Tool{
		Name:    "echo",
		Version: "v1",
		Handler: func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
			var payload map[string]any
			if err := json.Unmarshal(params, &payload); err != nil {
				return nil, &mcp.ErrorDetail{Code: "invalid_params", Message: "invalid params", Kind: mcp.KindInvalidParams}
			}
			return payload, nil
		},
	})

	server := mcp.NewServer(registry, nil, nil)
	params, err := json.Marshal(map[string]any{
		"name": "echo",
		"arguments": map[string]any{
			"message": "hello",
		},
	})
	if err != nil {
		t.Fatalf("marshal params: %v", err)
	}
	req := mcp.Request{
		JSONRPC: mcp.JSONRPCVersion,
		ID:      json.RawMessage("1"),
		Method:  mcp.MethodToolsCall,
		Params:  params,
	}
	resp := server.HandleRequest(context.Background(), req)
	if resp.Error != nil {
		t.Fatalf("tools/call error: %+v", resp.Error)
	}
	var callResult mcp.ToolCallResult
	resultBytes, err := json.Marshal(resp.Result)
	if err != nil {
		t.Fatalf("marshal tools/call result: %v", err)
	}
	if err := json.Unmarshal(resultBytes, &callResult); err != nil {
		t.Fatalf("unmarshal tools/call result: %v", err)
	}
	if callResult.StructuredContent == nil {
		t.Fatalf("expected structuredContent to be present")
	}
	payloadBytes, err := json.Marshal(callResult.StructuredContent)
	if err != nil {
		t.Fatalf("marshal tool result: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		t.Fatalf("unmarshal tool result: %v", err)
	}
	if payload["message"] != "hello" {
		t.Fatalf("expected message %q, got %v", "hello", payload["message"])
	}

	missingParams, err := json.Marshal(map[string]any{
		"name": "missing",
	})
	if err != nil {
		t.Fatalf("marshal missing params: %v", err)
	}
	missingReq := mcp.Request{
		JSONRPC: mcp.JSONRPCVersion,
		ID:      json.RawMessage("2"),
		Method:  mcp.MethodToolsCall,
		Params:  missingParams,
	}
	missingResp := server.HandleRequest(context.Background(), missingReq)
	if missingResp.Error == nil || missingResp.Error.Code != mcp.ErrMethodNotFound {
		t.Fatalf("expected method not found error, got %+v", missingResp.Error)
	}
}
