package unit

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"

	"open-sandbox/internal/mcp"
	"open-sandbox/internal/mcp/remote"
)

type fakeRemoteFilterClient struct {
	toolsByServer map[string][]mcp.ToolInfo
	errByServer   map[string]error
}

func (client *fakeRemoteFilterClient) ToolsList(ctx context.Context, cfg remote.ServerConfig) (mcp.ToolsListResult, error) {
	if err := client.errByServer[cfg.Name]; err != nil {
		return mcp.ToolsListResult{}, err
	}
	return mcp.ToolsListResult{Tools: client.toolsByServer[cfg.Name]}, nil
}

func (client *fakeRemoteFilterClient) ToolsCall(ctx context.Context, cfg remote.ServerConfig, name string, args json.RawMessage) (mcp.ToolCallResult, error) {
	return mcp.ToolCallResult{}, errors.New("not implemented")
}

func TestRemoteMCPToolAllowGlob(t *testing.T) {
	temp := t.TempDir()
	path := filepath.Join(temp, "mcp.json")
	client := &fakeRemoteFilterClient{
		toolsByServer: map[string][]mcp.ToolInfo{
			"s1": {
				{Name: "toolA"},
				{Name: "toolB"},
				{Name: "other"},
			},
		},
		errByServer: map[string]error{},
	}
	manager, err := remote.NewManagerWithClient(path, client)
	if err != nil {
		t.Fatalf("manager init: %v", err)
	}
	if err := manager.Upsert(remote.ServerConfig{
		Name:          "s1",
		URL:           "http://example.com",
		ToolAllowGlob: []string{"tool*"},
	}); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	registry := mcp.NewRegistry()
	if err := manager.SyncRegistry(context.Background(), registry); err != nil {
		t.Fatalf("sync: %v", err)
	}
	if _, ok := registry.Get("ext.s1.toolA"); !ok {
		t.Fatalf("expected toolA registered")
	}
	if _, ok := registry.Get("ext.s1.toolB"); !ok {
		t.Fatalf("expected toolB registered")
	}
	if _, ok := registry.Get("ext.s1.other"); ok {
		t.Fatalf("expected other filtered")
	}
}

func TestRemoteMCPToolDenyGlob(t *testing.T) {
	temp := t.TempDir()
	path := filepath.Join(temp, "mcp.json")
	client := &fakeRemoteFilterClient{
		toolsByServer: map[string][]mcp.ToolInfo{
			"s1": {
				{Name: "secret.read"},
				{Name: "public.read"},
			},
		},
		errByServer: map[string]error{},
	}
	manager, err := remote.NewManagerWithClient(path, client)
	if err != nil {
		t.Fatalf("manager init: %v", err)
	}
	if err := manager.Upsert(remote.ServerConfig{
		Name:         "s1",
		URL:          "http://example.com",
		ToolDenyGlob: []string{"secret*"},
	}); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	registry := mcp.NewRegistry()
	if err := manager.SyncRegistry(context.Background(), registry); err != nil {
		t.Fatalf("sync: %v", err)
	}
	if _, ok := registry.Get("ext.s1.secret.read"); ok {
		t.Fatalf("expected secret filtered")
	}
	if _, ok := registry.Get("ext.s1.public.read"); !ok {
		t.Fatalf("expected public registered")
	}
}

func TestRemoteMCPToolAllowOverridesDeny(t *testing.T) {
	temp := t.TempDir()
	path := filepath.Join(temp, "mcp.json")
	client := &fakeRemoteFilterClient{
		toolsByServer: map[string][]mcp.ToolInfo{
			"s1": {
				{Name: "toolA"},
				{Name: "toolB"},
				{Name: "toolC"},
			},
		},
		errByServer: map[string]error{},
	}
	manager, err := remote.NewManagerWithClient(path, client)
	if err != nil {
		t.Fatalf("manager init: %v", err)
	}
	if err := manager.Upsert(remote.ServerConfig{
		Name:         "s1",
		URL:          "http://example.com",
		ToolAllow:    []string{"toolA"},
		ToolDenyGlob: []string{"tool*"},
	}); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	registry := mcp.NewRegistry()
	if err := manager.SyncRegistry(context.Background(), registry); err != nil {
		t.Fatalf("sync: %v", err)
	}
	if _, ok := registry.Get("ext.s1.toolA"); !ok {
		t.Fatalf("expected toolA allowed")
	}
	if _, ok := registry.Get("ext.s1.toolB"); ok {
		t.Fatalf("expected toolB filtered")
	}
	if _, ok := registry.Get("ext.s1.toolC"); ok {
		t.Fatalf("expected toolC filtered")
	}
}
