package unit

import (
	"net/http"
	"regexp"
	"testing"

	"open-sandbox/internal/api"
	"open-sandbox/pkg/types"
)

func TestNewTraceID(t *testing.T) {
	traceID := api.NewTraceID()
	if traceID == "" {
		t.Fatalf("expected trace id to be non-empty")
	}
	if len(traceID) != 32 {
		t.Fatalf("expected trace id length 32, got %d", len(traceID))
	}
	matched, err := regexp.MatchString("^[a-f0-9]{32}$", traceID)
	if err != nil {
		t.Fatalf("trace id regex error: %v", err)
	}
	if !matched {
		t.Fatalf("trace id should be hex, got %q", traceID)
	}
}

func TestErrorResponseMapping(t *testing.T) {
	appErr := api.NewAppError("bad_request", "bad request", http.StatusBadRequest)
	appErr = api.WithTraceID(appErr, "trace-123")
	resp := api.ErrorResponse(appErr)

	if resp.Status != types.StatusError {
		t.Fatalf("expected status %q, got %q", types.StatusError, resp.Status)
	}
	if resp.Error == nil {
		t.Fatalf("expected error details to be present")
	}
	if resp.Error.Code != "bad_request" {
		t.Fatalf("expected error code %q, got %q", "bad_request", resp.Error.Code)
	}
	if resp.Error.TraceID != "trace-123" {
		t.Fatalf("expected trace id %q, got %q", "trace-123", resp.Error.TraceID)
	}
}
