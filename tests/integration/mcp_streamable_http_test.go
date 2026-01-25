package integration

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"

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

	body := buildStreamableBody(t, []map[string]any{
		{
			"jsonrpc": mcp.JSONRPCVersion,
			"id":      1,
			"method":  mcp.MethodInitialize,
			"params": map[string]any{
				"protocolVersion": mcp.SupportedProtocolVersion,
			},
		},
		{
			"jsonrpc": mcp.JSONRPCVersion,
			"id":      2,
			"method":  mcp.MethodCapabilities,
			"params": map[string]any{
				"protocol_version": mcp.SupportedProtocolVersion,
			},
		},
	})

	req, err := http.NewRequest(http.MethodPost, server.URL+"/mcp/stream", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-ndjson")
	req.Header.Set("Accept", "application/x-ndjson")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http post: %v", err)
	}
	defer resp.Body.Close()

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/x-ndjson") {
		t.Fatalf("expected ndjson content type, got %q", resp.Header.Get("Content-Type"))
	}

	responses := parseNDJSONResponses(t, resp)
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
	for _, mcpResp := range responses {
		if mcpResp.Error != nil {
			t.Fatalf("unexpected error: %+v", mcpResp.Error)
		}
	}
}

func TestMCPStreamableHTTPInvalidRequest(t *testing.T) {
	t.Setenv("MCP_AUTH_ENABLED", "false")

	router := api.NewRouter()
	handlers.RegisterMCPRoutes(router, nil, nil)

	server := httptest.NewServer(router)
	defer server.Close()

	body := []byte("{invalid-json}\n")
	req, err := http.NewRequest(http.MethodPost, server.URL+"/mcp/stream", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-ndjson")
	req.Header.Set("Accept", "application/x-ndjson")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http post: %v", err)
	}
	defer resp.Body.Close()

	responses := parseNDJSONResponses(t, resp)
	if len(responses) != 1 {
		t.Fatalf("expected 1 response, got %d", len(responses))
	}
	if responses[0].Error == nil || responses[0].Error.Code != mcp.ErrInvalidRequest {
		t.Fatalf("expected invalid request error, got %+v", responses[0].Error)
	}
}

func TestMCPStreamableHTTPAuthEnforced(t *testing.T) {
	t.Setenv("MCP_AUTH_ENABLED", "true")
	t.Setenv("MCP_AUTH_JWT_SECRET", "secret")

	router := api.NewRouter()
	handlers.RegisterMCPRoutes(router, nil, nil)

	server := httptest.NewServer(router)
	defer server.Close()

	body := buildStreamableBody(t, []map[string]any{
		{
			"jsonrpc": mcp.JSONRPCVersion,
			"id":      1,
			"method":  mcp.MethodCapabilities,
			"params": map[string]any{
				"protocol_version": mcp.SupportedProtocolVersion,
			},
		},
	})

	req, err := http.NewRequest(http.MethodPost, server.URL+"/mcp/stream", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-ndjson")
	req.Header.Set("Accept", "application/x-ndjson")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http post: %v", err)
	}
	defer resp.Body.Close()

	unauth := parseNDJSONResponses(t, resp)
	if len(unauth) != 1 {
		t.Fatalf("expected 1 response, got %d", len(unauth))
	}
	if unauth[0].Error == nil || unauth[0].Error.Code != mcp.ErrUnauthorized {
		t.Fatalf("expected unauthorized error, got %+v", unauth[0].Error)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": "test"})
	signed, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	req, err = http.NewRequest(http.MethodPost, server.URL+"/mcp/stream", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-ndjson")
	req.Header.Set("Accept", "application/x-ndjson")
	req.Header.Set("Authorization", "Bearer "+signed)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http post with auth: %v", err)
	}
	defer resp.Body.Close()

	authResp := parseNDJSONResponses(t, resp)
	if len(authResp) != 1 {
		t.Fatalf("expected 1 response, got %d", len(authResp))
	}
	if authResp[0].Error != nil {
		t.Fatalf("unexpected error: %+v", authResp[0].Error)
	}
}

func buildStreamableBody(t *testing.T, payloads []map[string]any) []byte {
	t.Helper()
	var buf bytes.Buffer
	for _, payload := range payloads {
		line, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal request: %v", err)
		}
		buf.Write(line)
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

func parseNDJSONResponses(t *testing.T, resp *http.Response) []mcp.Response {
	t.Helper()
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 1024), 1024*1024)
	var responses []mcp.Response
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var mcpResp mcp.Response
		if err := json.Unmarshal([]byte(line), &mcpResp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		responses = append(responses, mcpResp)
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan response: %v", err)
	}
	return responses
}
