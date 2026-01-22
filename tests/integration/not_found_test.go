package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"open-sandbox/internal/api"
	"open-sandbox/internal/api/handlers"
	"open-sandbox/pkg/types"
)

func TestNotFoundReturnsUnifiedError(t *testing.T) {
	router := api.NewRouter()
	handlers.RegisterSandboxRoutes(router)

	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/v1/does-not-exist")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}

	var payload types.Response
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload.Status != types.StatusError {
		t.Fatalf("expected status %q, got %q", types.StatusError, payload.Status)
	}
	if payload.Error == nil {
		t.Fatalf("expected error details")
	}
	if payload.Error.Code != api.CodeNotFound {
		t.Fatalf("expected error code %q, got %q", api.CodeNotFound, payload.Error.Code)
	}
	if payload.Error.TraceID == "" {
		t.Fatalf("expected trace id")
	}
}
