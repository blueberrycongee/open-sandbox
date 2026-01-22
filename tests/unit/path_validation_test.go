package unit

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"open-sandbox/internal/file"
)

func TestValidateWorkspacePathBlocksPrefixEscape(t *testing.T) {
	root := t.TempDir()
	workspace := filepath.Join(root, "workspace")
	evil := filepath.Join(root, "workspace-evil", "secret.txt")

	if err := os.MkdirAll(workspace, 0755); err != nil {
		t.Fatalf("create workspace: %v", err)
	}

	if err := file.ValidateWorkspacePath(evil, workspace); err == nil {
		t.Fatalf("expected workspace escape to be rejected")
	}
}

func TestValidateWorkspacePathAllowsCleanInside(t *testing.T) {
	root := t.TempDir()
	workspace := filepath.Join(root, "workspace")
	if err := os.MkdirAll(workspace, 0755); err != nil {
		t.Fatalf("create workspace: %v", err)
	}

	path := filepath.Join(workspace, "..", filepath.Base(workspace), "file.txt")
	if err := file.ValidateWorkspacePath(path, workspace); err != nil {
		t.Fatalf("expected path to be allowed, got %v", err)
	}
}

func TestValidateWorkspacePathWindowsCaseFold(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("case-fold test only applies to Windows")
	}

	root := t.TempDir()
	workspace := filepath.Join(root, "workspace")
	if err := os.MkdirAll(workspace, 0755); err != nil {
		t.Fatalf("create workspace: %v", err)
	}

	upper := strings.ToUpper(workspace)
	path := filepath.Join(upper, "file.txt")
	if err := file.ValidateWorkspacePath(path, workspace); err != nil {
		t.Fatalf("expected case-insensitive path to be allowed, got %v", err)
	}
}
