package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"

	"open-sandbox/internal/api"
	"open-sandbox/internal/api/handlers"
	"open-sandbox/internal/mcp"
)

func TestMCPHTTPRoundTrip(t *testing.T) {
	t.Setenv("MCP_AUTH_ENABLED", "false")

	router := api.NewRouter()
	handlers.RegisterMCPRoutes(router, nil)

	server := httptest.NewServer(router)
	defer server.Close()

	payload := map[string]any{
		"jsonrpc": mcp.JSONRPCVersion,
		"id":      1,
		"method":  "mcp.capabilities",
		"params": map[string]any{
			"protocol_version": mcp.SupportedProtocolVersion,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	resp, err := http.Post(server.URL+"/mcp", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("http post: %v", err)
	}
	defer resp.Body.Close()

	var mcpResp mcp.Response
	if err := json.NewDecoder(resp.Body).Decode(&mcpResp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if mcpResp.Error != nil {
		t.Fatalf("unexpected error: %+v", mcpResp.Error)
	}

	resultBytes, err := json.Marshal(mcpResp.Result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}
	var caps mcp.CapabilitiesResponse
	if err := json.Unmarshal(resultBytes, &caps); err != nil {
		t.Fatalf("unmarshal capabilities: %v", err)
	}
	if caps.ProtocolVersion != mcp.SupportedProtocolVersion {
		t.Fatalf("expected protocol version %q, got %q", mcp.SupportedProtocolVersion, caps.ProtocolVersion)
	}
}

func TestMCPHTTPAuthEnforced(t *testing.T) {
	t.Setenv("MCP_AUTH_ENABLED", "true")
	t.Setenv("MCP_AUTH_JWT_SECRET", "secret")

	router := api.NewRouter()
	handlers.RegisterMCPRoutes(router, nil)

	server := httptest.NewServer(router)
	defer server.Close()

	payload := map[string]any{
		"jsonrpc": mcp.JSONRPCVersion,
		"id":      1,
		"method":  "mcp.capabilities",
		"params": map[string]any{
			"protocol_version": mcp.SupportedProtocolVersion,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	resp, err := http.Post(server.URL+"/mcp", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("http post: %v", err)
	}
	defer resp.Body.Close()

	var mcpResp mcp.Response
	if err := json.NewDecoder(resp.Body).Decode(&mcpResp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if mcpResp.Error == nil || mcpResp.Error.Code != mcp.ErrUnauthorized {
		t.Fatalf("expected unauthorized error, got %+v", mcpResp.Error)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": "test"})
	signed, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, server.URL+"/mcp", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+signed)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("http post with auth: %v", err)
	}
	defer resp.Body.Close()

	var authResp mcp.Response
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if authResp.Error != nil {
		t.Fatalf("unexpected error: %+v", authResp.Error)
	}
}

func TestMCPHTTPNotificationNoResponse(t *testing.T) {
	t.Setenv("MCP_AUTH_ENABLED", "false")

	router := api.NewRouter()
	handlers.RegisterMCPRoutes(router, nil)

	server := httptest.NewServer(router)
	defer server.Close()

	payload := map[string]any{
		"jsonrpc": mcp.JSONRPCVersion,
		"id":      nil,
		"method":  "mcp.capabilities",
		"params": map[string]any{
			"protocol_version": mcp.SupportedProtocolVersion,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	resp, err := http.Post(server.URL+"/mcp", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("http post: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if len(data) != 0 {
		t.Fatalf("expected empty body, got %q", string(data))
	}
}
