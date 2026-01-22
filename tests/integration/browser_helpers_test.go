package integration

import (
	"os"
	"path/filepath"
	"testing"

	"open-sandbox/internal/browser"
	"open-sandbox/internal/config"
)

func startBrowserService(t *testing.T) *browser.Service {
	t.Helper()

	cfg := browser.DefaultConfig()
	cfg.BinaryPath = os.Getenv("SANDBOX_BROWSER_BIN")
	cfg.ExistingWebSocketDebug = os.Getenv("SANDBOX_BROWSER_CDP")
	cfg.UserDataDir = filepath.Join(config.CacheRoot, "chrome-profile-test")
	cfg.RemoteDebuggingPort = 0
	cfg.Headless = true

	if err := os.MkdirAll(config.CacheRoot, 0755); err != nil {
		t.Fatalf("cache dir setup failed: %v", err)
	}

	service := browser.NewService(cfg)
	if err := service.Start(); err != nil {
		if err == browser.ErrBrowserUnavailable {
			t.Skip("browser binary not available for integration test")
		}
		t.Fatalf("start browser: %v", err)
	}
	t.Cleanup(service.Close)
	return service
}
