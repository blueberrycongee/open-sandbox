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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	allocCtx, allocCancel := chromedp.NewRemoteAllocator(ctx, cdpAddress)
	defer allocCancel()
	tabCtx, tabCancel := chromedp.NewContext(allocCtx)
	defer tabCancel()
	if err := chromedp.Run(tabCtx, chromedp.Navigate("about:blank")); err != nil {
		t.Fatalf("cdp connect failed: %v", err)
	}
}
