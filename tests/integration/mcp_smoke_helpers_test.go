package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"open-sandbox/internal/api"
	"open-sandbox/internal/api/handlers"
	"open-sandbox/internal/mcp"
)

func newMCPTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	t.Setenv("MCP_AUTH_ENABLED", "false")

	router := api.NewRouter()
	handlers.RegisterMCPRoutes(router, nil)

	server := httptest.NewServer(router)
	t.Cleanup(server.Close)
	return server
}

func buildMCPRequest(t *testing.T, method string, params any) []byte {
	t.Helper()
	payload := map[string]any{
		"jsonrpc": mcp.JSONRPCVersion,
		"id":      1,
		"method":  method,
		"params":  params,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	return body
}

func postMCPRequest(t *testing.T, serverURL string, body []byte) mcp.Response {
	t.Helper()
	resp, err := http.Post(serverURL+"/mcp", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("http post: %v", err)
	}
	defer resp.Body.Close()

	return decodeMCPResponse(t, resp.Body)
}

func decodeMCPResponse(t *testing.T, reader io.Reader) mcp.Response {
	t.Helper()
	var mcpResp mcp.Response
	if err := json.NewDecoder(reader).Decode(&mcpResp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return mcpResp
}
