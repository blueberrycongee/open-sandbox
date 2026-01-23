package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"open-sandbox/internal/api"
	"open-sandbox/internal/api/handlers"
	"open-sandbox/internal/browser"
	"open-sandbox/internal/config"
)

func main() {
	if err := config.EnsureWorkspace(); err != nil {
		log.Fatalf("workspace init failed: %v", err)
	}
	if err := os.MkdirAll(config.CachePath(), 0755); err != nil {
		log.Fatalf("cache init failed: %v", err)
	}

	addr := os.Getenv("SANDBOX_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	router := api.NewRouter()
	handlers.RegisterSandboxRoutes(router)

	browserService := browser.NewService(browser.Config{
		BinaryPath:             os.Getenv("SANDBOX_BROWSER_BIN"),
		UserDataDir:            filepath.Join(config.CachePath(), "chrome-profile"),
		RemoteDebuggingHost:    getenv("SANDBOX_CDP_HOST", "127.0.0.1"),
		RemoteDebuggingPort:    getenvInt("SANDBOX_CDP_PORT", 9222),
		ExistingWebSocketDebug: os.Getenv("SANDBOX_BROWSER_CDP"),
		Headless:               getenvBool("SANDBOX_BROWSER_HEADLESS", false),
	})
	handlers.RegisterBrowserRoutes(router, browserService)
	handlers.RegisterVNCRoutes(router, browserService)
	handlers.RegisterShellRoutes(router)
	handlers.RegisterFileRoutes(router)
	handlers.RegisterCodeExecRoutes(router)
	handlers.RegisterJupyterRoutes(router, os.Getenv("SANDBOX_JUPYTER_URL"))
	handlers.RegisterCodeServerRoutes(router, os.Getenv("SANDBOX_CODESERVER_URL"))
	handlers.RegisterMCPRoutes(router, browserService)

	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server stopped: %v", err)
	}
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getenvInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}

func getenvBool(key string, fallback bool) bool {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	switch raw {
	case "1", "true", "TRUE", "True", "yes", "YES", "Yes":
		return true
	case "0", "false", "FALSE", "False", "no", "NO", "No":
		return false
	default:
		return fallback
	}
}
