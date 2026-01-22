package unit

import (
	"os"
	"path/filepath"
	"testing"

	"open-sandbox/internal/config"
)

func TestEnsureWorkspace(t *testing.T) {
	workspace := config.WorkspacePath()
	if !filepath.IsAbs(workspace) {
		t.Fatalf("workspace path must be absolute: %s", workspace)
	}

	if err := config.EnsureWorkspace(); err != nil {
		t.Fatalf("ensure workspace failed: %v", err)
	}

	info, err := os.Stat(workspace)
	if err != nil {
		t.Fatalf("workspace should exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("workspace path is not a directory: %s", workspace)
	}
}
