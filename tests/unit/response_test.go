package unit

import (
	"encoding/json"
	"testing"

	"open-sandbox/pkg/types"
)

func TestResponseJSONSerialization(t *testing.T) {
	resp := types.Response{
		Status: types.StatusOK,
		Data: map[string]string{
			"result": "ok",
		},
	}

	raw, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if decoded["status"] != types.StatusOK {
		t.Fatalf("expected status %q, got %v", types.StatusOK, decoded["status"])
	}

	if _, ok := decoded["data"]; !ok {
		t.Fatalf("expected data field to be present")
	}
}

func TestErrorResponseJSONSerialization(t *testing.T) {
	resp := types.Response{
		Status: types.StatusError,
		Error: &types.ErrorDetail{
			Code:    "bad_request",
			Message: "bad request",
			TraceID: "trace-1",
		},
	}

	raw, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if decoded["status"] != types.StatusError {
		t.Fatalf("expected status %q, got %v", types.StatusError, decoded["status"])
	}

	errorValue, ok := decoded["error"].(map[string]any)
	if !ok {
		t.Fatalf("expected error field to be an object")
	}

	if errorValue["code"] != "bad_request" {
		t.Fatalf("expected error.code to be %q, got %v", "bad_request", errorValue["code"])
	}
	if errorValue["trace_id"] != "trace-1" {
		t.Fatalf("expected error.trace_id to be %q, got %v", "trace-1", errorValue["trace_id"])
	}
}
