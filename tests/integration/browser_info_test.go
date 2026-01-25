package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"

	"open-sandbox/internal/api"
	"open-sandbox/internal/api/handlers"
	"open-sandbox/pkg/types"
)

func TestBrowserInfo(t *testing.T) {
	service := startBrowserService(t)

	router := api.NewRouter()
	handlers.RegisterBrowserRoutes(router, service)

	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/v1/browser/info")
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

	data, ok := payload.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected data to be an object")
	}
	cdpAddress, _ := data["cdp_address"].(string)
	if !strings.HasPrefix(cdpAddress, "ws://") && !strings.HasPrefix(cdpAddress, "wss://") {
		t.Fatalf("unexpected cdp address: %v", data["cdp_address"])
	}
	cdpURL, _ := data["cdp_url"].(string)
	if !strings.HasPrefix(cdpURL, "ws://") && !strings.HasPrefix(cdpURL, "wss://") {
		t.Fatalf("unexpected cdp url: %v", data["cdp_url"])
	}
	if cdpAddress != "" && cdpURL != "" && cdpAddress != cdpURL {
		t.Fatalf("cdp_address and cdp_url mismatch: %v vs %v", cdpAddress, cdpURL)
	}
	userAgent, _ := data["user_agent"].(string)
	if strings.TrimSpace(userAgent) == "" {
		t.Fatalf("unexpected user agent: %v", data["user_agent"])
	}
	vncURL, _ := data["vnc_url"].(string)
	if !strings.Contains(vncURL, "/vnc/index.html") {
		t.Fatalf("unexpected vnc url: %v", data["vnc_url"])
	}
	viewport, _ := data["viewport"].(map[string]any)
	if viewport == nil {
		t.Fatalf("missing viewport: %v", data["viewport"])
	}
	width, _ := viewport["width"].(float64)
	height, _ := viewport["height"].(float64)
	if width <= 0 || height <= 0 {
		t.Fatalf("unexpected viewport: %v", viewport)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	allocCtx, allocCancel := chromedp.NewRemoteAllocator(ctx, cdpURL)
	defer allocCancel()
	tabCtx, tabCancel := chromedp.NewContext(allocCtx)
	defer tabCancel()
	if err := chromedp.Run(tabCtx, chromedp.Navigate("about:blank")); err != nil {
		t.Fatalf("cdp connect failed: %v", err)
	}
}
