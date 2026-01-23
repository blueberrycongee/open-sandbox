package integration

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"

	"open-sandbox/internal/api"
	"open-sandbox/internal/api/handlers"
	"open-sandbox/internal/mcp"
)

func TestMCPSSEEndpoint(t *testing.T) {
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

	resp, err := http.Get(server.URL + "/mcp/sse?request=" + url.QueryEscape(string(body)))
	if err != nil {
		t.Fatalf("http get: %v", err)
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

func TestMCPSSEAuthEnforced(t *testing.T) {
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

	resp, err := http.Get(server.URL + "/mcp/sse?request=" + url.QueryEscape(string(body)))
	if err != nil {
		t.Fatalf("http get: %v", err)
	}
	defer resp.Body.Close()

	mcpResp := parseSSEPayload(t, resp)
	if mcpResp.Error == nil || mcpResp.Error.Code != mcp.ErrUnauthorized {
		t.Fatalf("expected unauthorized error, got %+v", mcpResp.Error)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": "test"})
	signed, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	req, err := http.NewRequest(http.MethodGet, server.URL+"/mcp/sse?request="+url.QueryEscape(string(body)), nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+signed)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("http get with auth: %v", err)
	}
	defer resp.Body.Close()

	mcpResp = parseSSEPayload(t, resp)
	if mcpResp.Error != nil {
		t.Fatalf("unexpected error: %+v", mcpResp.Error)
	}
}

func TestMCPSSENotificationNoResponse(t *testing.T) {
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

	resp, err := http.Get(server.URL + "/mcp/sse?request=" + url.QueryEscape(string(body)))
	if err != nil {
		t.Fatalf("http get: %v", err)
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

func parseSSEPayload(t *testing.T, resp *http.Response) mcp.Response {
	t.Helper()
	scanner := bufio.NewScanner(resp.Body)
	var dataLine string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			dataLine = strings.TrimPrefix(line, "data: ")
			break
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("read body: %v", err)
	}
	if dataLine == "" {
		t.Fatalf("missing data line in sse response")
	}
	var mcpResp mcp.Response
	if err := json.Unmarshal([]byte(dataLine), &mcpResp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return mcpResp
}
