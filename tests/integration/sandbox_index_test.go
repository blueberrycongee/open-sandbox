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

func TestSandboxIndex(t *testing.T) {
	router := api.NewRouter()
	handlers.RegisterSandboxRoutes(router)

	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/v1/sandbox")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var payload types.Response
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload.Status != types.StatusOK {
		t.Fatalf("expected status %q, got %q", types.StatusOK, payload.Status)
	}
	if payload.Data == nil {
		t.Fatalf("expected data payload")
	}
}
