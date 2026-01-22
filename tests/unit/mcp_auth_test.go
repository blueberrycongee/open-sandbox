package unit

import (
	"testing"

	"open-sandbox/internal/mcp"
)

func TestAuthConfigDisabledAllowsMissingKeys(t *testing.T) {
	t.Setenv("MCP_AUTH_ENABLED", "false")
	config := mcp.LoadAuthConfig()
	if config.Enabled {
		t.Fatalf("expected auth to be disabled")
	}
	if err := config.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAuthConfigEnabledRequiresKey(t *testing.T) {
	t.Setenv("MCP_AUTH_ENABLED", "true")
	config := mcp.LoadAuthConfig()
	if !config.Enabled {
		t.Fatalf("expected auth to be enabled")
	}
	if err := config.Validate(); err == nil {
		t.Fatalf("expected validation error when no key provided")
	}
}

func TestAuthConfigEnabledWithSecretPasses(t *testing.T) {
	t.Setenv("MCP_AUTH_ENABLED", "true")
	t.Setenv("MCP_AUTH_JWT_SECRET", "secret")
	config := mcp.LoadAuthConfig()
	if err := config.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}
