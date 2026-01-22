package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"open-sandbox/internal/api"
	"open-sandbox/internal/api/handlers"
)

func TestVNCPageAccessible(t *testing.T) {
	service := startBrowserService(t)
	router := api.NewRouter()
	handlers.RegisterVNCRoutes(router, service)

	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/vnc/index.html")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Fatalf("expected text/html content type, got %q", contentType)
	}

	imageResp, err := http.Get(server.URL + "/vnc/screen.png")
	if err != nil {
		t.Fatalf("screen request failed: %v", err)
	}
	defer imageResp.Body.Close()
	if imageResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(imageResp.Body)
		t.Fatalf("expected screen status 200, got %d: %s", imageResp.StatusCode, string(bodyBytes))
	}
	if !strings.Contains(imageResp.Header.Get("Content-Type"), "image/png") {
		t.Fatalf("expected image/png content type")
	}

	clickBody, err := json.Marshal(map[string]float64{"x": 1, "y": 1})
	if err != nil {
		t.Fatalf("marshal click: %v", err)
	}
	clickResp, err := http.Post(server.URL+"/vnc/click", "application/json", bytes.NewReader(clickBody))
	if err != nil {
		t.Fatalf("click request failed: %v", err)
	}
	clickResp.Body.Close()
	if clickResp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected click status 204, got %d", clickResp.StatusCode)
	}
}
