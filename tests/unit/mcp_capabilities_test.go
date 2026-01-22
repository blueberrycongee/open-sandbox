package unit

import (
	"testing"

	"open-sandbox/internal/mcp"
)

func TestCapabilitiesPayloadIncludesVersionsAndPermissions(t *testing.T) {
	registry := mcp.NewRegistry()
	registry.Register(mcp.Tool{
		Name:    "browser.navigate",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "network",
		},
	})

	caps := mcp.BuildCapabilities(registry)
	if caps.ProtocolVersion == "" {
		t.Fatalf("expected protocol version to be set")
	}
	if len(caps.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(caps.Tools))
	}
	tool := caps.Tools[0]
	if tool.Name != "browser.navigate" {
		t.Fatalf("expected tool name %q, got %q", "browser.navigate", tool.Name)
	}
	if tool.Version == "" {
		t.Fatalf("expected tool version to be set")
	}
	if tool.Permissions.Scope == "" {
		t.Fatalf("expected permissions scope to be set")
	}
}
