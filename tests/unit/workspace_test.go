package unit

import (
	"os"
	"path/filepath"
	"testing"

	"open-sandbox/internal/config"
)

func TestEnsureWorkspace(t *testing.T) {
	if !filepath.IsAbs(config.HostWorkspacePath) {
		t.Fatalf("workspace path must be absolute: %s", config.HostWorkspacePath)
	}

	if err := config.EnsureWorkspace(); err != nil {
		t.Fatalf("ensure workspace failed: %v", err)
	}

	info, err := os.Stat(config.HostWorkspacePath)
	if err != nil {
		t.Fatalf("workspace should exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("workspace path is not a directory: %s", config.HostWorkspacePath)
	}
}
