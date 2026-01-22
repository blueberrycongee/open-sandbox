package unit

import (
	"encoding/json"
	"testing"

	"open-sandbox/internal/mcp"
)

func TestParseRequestRejectsInvalidJSONRPC(t *testing.T) {
	payload := []byte(`{"jsonrpc":"1.0","id":1,"method":"mcp.capabilities"}`)
	_, err := mcp.ParseRequest(payload)
	if err == nil {
		t.Fatalf("expected invalid jsonrpc version error")
	}
}

func TestValidateProtocolVersionRejectsMismatch(t *testing.T) {
	params := []byte(`{"protocol_version":"999.0"}`)
	if err := mcp.ValidateProtocolVersion(params); err == nil {
		t.Fatalf("expected protocol version mismatch error")
	}
}

func TestUnifiedErrorSchemaIncludesTraceID(t *testing.T) {
	id := json.RawMessage("1")
	detail := mcp.NewErrorDetail("bad_request", "bad request", "bad_request")
	resp := mcp.NewErrorResponse(id, mcp.ErrInvalidRequest, "invalid request", detail)
	if resp.Error == nil || resp.Error.Data == nil {
		t.Fatalf("expected error data to be present")
	}
	if resp.Error.Data.TraceID == "" {
		t.Fatalf("expected trace id to be populated")
	}
	if resp.Error.Data.Code != "bad_request" {
		t.Fatalf("expected error code %q, got %q", "bad_request", resp.Error.Data.Code)
	}
	if resp.Error.Data.Message != "bad request" {
		t.Fatalf("expected error message %q, got %q", "bad request", resp.Error.Data.Message)
	}
}
