package unit

import (
	"bytes"
	"encoding/json"
	"testing"

	"open-sandbox/internal/mcp"
)

func TestServeStdioReturnsInvalidRequestError(t *testing.T) {
	server := mcp.NewServer(mcp.NewRegistry(), nil, nil)
	input := bytes.NewBufferString("{\"jsonrpc\":\"2.0\",\"id\":1}\n")
	output := &bytes.Buffer{}

	if err := server.ServeStdio(input, output); err != nil {
		t.Fatalf("serve stdio: %v", err)
	}

	var resp mcp.Response
	if err := json.NewDecoder(output).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Error == nil || resp.Error.Code != mcp.ErrInvalidRequest {
		t.Fatalf("expected invalid request error, got %+v", resp.Error)
	}
}
