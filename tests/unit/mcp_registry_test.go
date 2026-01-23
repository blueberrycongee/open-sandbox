package unit

import (
	"testing"

	"open-sandbox/internal/mcp"
)

func TestRegistryStoresPermissionMetadata(t *testing.T) {
	registry := mcp.NewRegistry()
	registry.Register(mcp.Tool{
		Name:    "file.read",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow:  true,
			Scope:  "workspace",
			Reason: "readonly",
		},
	})

	tool, ok := registry.Get("file.read")
	if !ok {
		t.Fatalf("expected tool to be registered")
	}
	if tool.Permissions.Allow != true {
		t.Fatalf("expected allow to be true")
	}
	if tool.Permissions.Scope != "workspace" {
		t.Fatalf("expected scope %q, got %q", "workspace", tool.Permissions.Scope)
	}
	if tool.Permissions.Reason != "readonly" {
		t.Fatalf("expected reason %q, got %q", "readonly", tool.Permissions.Reason)
	}
}
