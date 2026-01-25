package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"open-sandbox/internal/api"
	"open-sandbox/internal/api/handlers"
	"open-sandbox/internal/mcp"
)

func TestMCPStreamableHTTPResponse(t *testing.T) {
	t.Setenv("MCP_AUTH_ENABLED", "false")

	router := api.NewRouter()
	handlers.RegisterMCPRoutes(router, nil, nil)

	server := httptest.NewServer(router)
	defer server.Close()

	payload := map[string]any{
		"jsonrpc": mcp.JSONRPCVersion,
		"id":      1,
		"method":  mcp.MethodInitialize,
		"params": map[string]any{
			"protocolVersion": mcp.SupportedProtocolVersion,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, server.URL+"/mcp", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http post: %v", err)
	}
	defer resp.Body.Close()

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/event-stream") {
		t.Fatalf("expected event-stream content type, got %q", resp.Header.Get("Content-Type"))
	}

	mcpResp := parseSSEPayload(t, resp)
	if mcpResp.Error != nil {
		t.Fatalf("unexpected error: %+v", mcpResp.Error)
	}
}
