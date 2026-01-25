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

type fakeRemoteClient struct {
	toolsByServer map[string][]mcp.ToolInfo
	errByServer   map[string]error
}

func (client *fakeRemoteClient) ToolsList(ctx context.Context, cfg remote.ServerConfig) (mcp.ToolsListResult, error) {
	if err := client.errByServer[cfg.Name]; err != nil {
		return mcp.ToolsListResult{}, err
	}
	return mcp.ToolsListResult{Tools: client.toolsByServer[cfg.Name]}, nil
}

func (client *fakeRemoteClient) ToolsCall(ctx context.Context, cfg remote.ServerConfig, name string, args json.RawMessage) (mcp.ToolCallResult, error) {
	return mcp.ToolCallResult{}, errors.New("not implemented")
}

func TestRemoteMCPIncrementalSync(t *testing.T) {
	temp := t.TempDir()
	path := filepath.Join(temp, "mcp.json")
	client := &fakeRemoteClient{
		toolsByServer: map[string][]mcp.ToolInfo{
			"s1": {
				{Name: "toolA"},
				{Name: "toolB"},
			},
		},
		errByServer: map[string]error{},
	}
	manager, err := remote.NewManagerWithClient(path, client)
	if err != nil {
		t.Fatalf("manager init: %v", err)
	}
	if err := manager.Upsert(remote.ServerConfig{Name: "s1", URL: "http://example.com"}); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	registry := mcp.NewRegistry()
	registry.Register(mcp.Tool{Name: "file.read", Version: "v1"})

	if err := manager.SyncRegistry(context.Background(), registry); err != nil {
		t.Fatalf("sync: %v", err)
	}
	if _, ok := registry.Get("ext.s1.toolA"); !ok {
		t.Fatalf("expected toolA registered")
	}
	if _, ok := registry.Get("ext.s1.toolB"); !ok {
		t.Fatalf("expected toolB registered")
	}
	if _, ok := registry.Get("file.read"); !ok {
		t.Fatalf("expected local tool preserved")
	}

	client.toolsByServer["s1"] = []mcp.ToolInfo{
		{Name: "toolB"},
		{Name: "toolC"},
	}
	if err := manager.SyncRegistry(context.Background(), registry); err != nil {
		t.Fatalf("sync: %v", err)
	}
	if _, ok := registry.Get("ext.s1.toolA"); ok {
		t.Fatalf("expected toolA removed")
	}
	if _, ok := registry.Get("ext.s1.toolB"); !ok {
		t.Fatalf("expected toolB retained")
	}
	if _, ok := registry.Get("ext.s1.toolC"); !ok {
		t.Fatalf("expected toolC added")
	}
	if _, ok := registry.Get("file.read"); !ok {
		t.Fatalf("expected local tool preserved")
	}
}

func TestRemoteMCPSyncKeepsToolsOnError(t *testing.T) {
	temp := t.TempDir()
	path := filepath.Join(temp, "mcp.json")
	client := &fakeRemoteClient{
		toolsByServer: map[string][]mcp.ToolInfo{
			"s1": {{Name: "toolA"}},
		},
		errByServer: map[string]error{},
	}
	manager, err := remote.NewManagerWithClient(path, client)
	if err != nil {
		t.Fatalf("manager init: %v", err)
	}
	if err := manager.Upsert(remote.ServerConfig{Name: "s1", URL: "http://example.com"}); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	registry := mcp.NewRegistry()
	if err := manager.SyncRegistry(context.Background(), registry); err != nil {
		t.Fatalf("sync: %v", err)
	}
	if _, ok := registry.Get("ext.s1.toolA"); !ok {
		t.Fatalf("expected toolA registered")
	}

	client.errByServer["s1"] = errors.New("boom")
	if err := manager.SyncRegistry(context.Background(), registry); err != nil {
		t.Fatalf("sync: %v", err)
	}
	if _, ok := registry.Get("ext.s1.toolA"); !ok {
		t.Fatalf("expected toolA retained on error")
	}
}
