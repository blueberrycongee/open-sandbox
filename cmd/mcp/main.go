package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"open-sandbox/internal/api/handlers"
	"open-sandbox/internal/browser"
	"open-sandbox/internal/config"
	"open-sandbox/internal/mcp"
	"open-sandbox/internal/mcp/remote"
)

func main() {
	if err := config.EnsureWorkspace(); err != nil {
		fmt.Fprintf(os.Stderr, "workspace init failed: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(config.CachePath(), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "cache init failed: %v\n", err)
		os.Exit(1)
	}

	browserConfig := browser.DefaultConfig()
	browserConfig.BinaryPath = os.Getenv("SANDBOX_BROWSER_BIN")
	browserConfig.UserDataDir = filepath.Join(config.CachePath(), "chrome-profile")
	browserConfig.RemoteDebuggingHost = getenv("SANDBOX_CDP_HOST", browserConfig.RemoteDebuggingHost)
	browserConfig.RemoteDebuggingPort = getenvInt("SANDBOX_CDP_PORT", browserConfig.RemoteDebuggingPort)
	browserConfig.ExistingWebSocketDebug = os.Getenv("SANDBOX_BROWSER_CDP")
	browserConfig.Headless = getenvBool("SANDBOX_BROWSER_HEADLESS", browserConfig.Headless)
	browserConfig.DownloadDir = getenv("SANDBOX_BROWSER_DOWNLOAD_DIR", filepath.Join(config.WorkspacePath(), "Downloads"))

	browserService := browser.NewService(browserConfig)
	remoteManager, err := remote.NewManager(config.MCPServersPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "external mcp config load failed: %v\n", err)
	}
	registry := handlers.NewMCPRegistry(browserService, remoteManager)
	server := mcp.NewServer(registry, nil, nil)

	if err := server.ServeStdio(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "stdio server failed: %v\n", err)
		os.Exit(1)
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
