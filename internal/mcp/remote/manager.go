package remote

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"open-sandbox/internal/mcp"
)

type Manager struct {
	mu             sync.Mutex
	path           string
	servers        map[string]ServerConfig
	registeredTool map[string]map[string]struct{}
	client         ToolClient
}

func NewManager(path string) (*Manager, error) {
	return NewManagerWithClient(path, NewClient())
}

type ToolClient interface {
	ToolsList(ctx context.Context, cfg ServerConfig) (mcp.ToolsListResult, error)
	ToolsCall(ctx context.Context, cfg ServerConfig, name string, args json.RawMessage) (mcp.ToolCallResult, error)
}

func NewManagerWithClient(path string, client ToolClient) (*Manager, error) {
	manager := &Manager{
		path:           path,
		servers:        make(map[string]ServerConfig),
		registeredTool: make(map[string]map[string]struct{}),
		client:         client,
	}
	if err := manager.load(); err != nil {
		return nil, err
	}
	return manager, nil
}

func (manager *Manager) List() []ServerConfig {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	servers := make([]ServerConfig, 0, len(manager.servers))
	for _, server := range manager.servers {
		servers = append(servers, server)
	}
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Name < servers[j].Name
	})
	return servers
}

func (manager *Manager) Get(name string) (ServerConfig, bool) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	server, ok := manager.servers[name]
	return server, ok
}

func (manager *Manager) Upsert(server ServerConfig) error {
	if strings.TrimSpace(server.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(server.URL) == "" {
		return errors.New("url is required")
	}
	if server.Transport == "" {
		server.Transport = "http"
	}
	manager.mu.Lock()
	manager.servers[server.Name] = server
	manager.mu.Unlock()
	return manager.save()
}

func (manager *Manager) Delete(name string) error {
	manager.mu.Lock()
	delete(manager.servers, name)
	manager.mu.Unlock()
	return manager.save()
}

func (manager *Manager) SyncRegistry(ctx context.Context, registry *mcp.Registry) error {
	manager.mu.Lock()
	servers := make([]ServerConfig, 0, len(manager.servers))
	for _, server := range manager.servers {
		servers = append(servers, server)
	}
	knownServers := make(map[string]struct{}, len(manager.servers))
	for name := range manager.servers {
		knownServers[name] = struct{}{}
	}
	manager.mu.Unlock()

	manager.mu.Lock()
	for name, tools := range manager.registeredTool {
		if _, ok := knownServers[name]; ok {
			continue
		}
		for toolName := range tools {
			registry.Unregister(toolName)
		}
		delete(manager.registeredTool, name)
	}
	manager.mu.Unlock()

	for _, server := range servers {
		manager.mu.Lock()
		prevTools := manager.registeredTool[server.Name]
		manager.mu.Unlock()

		tools, err := manager.client.ToolsList(ctx, server)
		if err != nil {
			continue
		}
		nextTools := make(map[string]struct{})
		for _, tool := range tools.Tools {
			if !toolAllowed(server, tool.Name) {
				continue
			}
			registeredName := "ext." + server.Name + "." + tool.Name
			handler := buildRemoteToolHandler(manager.client, server, tool.Name)
			schema := mcp.ToolSchema{Input: tool.InputSchema, Output: tool.OutputSchema}
			if tool.Schema.Input != nil || tool.Schema.Output != nil {
				schema = tool.Schema
			}
			registry.Register(mcp.Tool{
				Name:    registeredName,
				Version: tool.Version,
				Permissions: mcp.PermissionMeta{
					Allow: true,
					Scope: "external",
				},
				Schema:  schema,
				Handler: handler,
			})
			nextTools[registeredName] = struct{}{}
		}
		for toolName := range prevTools {
			if _, ok := nextTools[toolName]; !ok {
				registry.Unregister(toolName)
			}
		}
		manager.mu.Lock()
		manager.registeredTool[server.Name] = nextTools
		manager.mu.Unlock()
	}
	return nil
}

func buildRemoteToolHandler(client ToolClient, server ServerConfig, name string) mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		result, err := client.ToolsCall(ctx, server, name, params)
		if err != nil {
			detail := mcp.NewErrorDetail(mcp.KindToolError, err.Error(), mcp.KindToolError)
			return nil, &detail
		}
		return result, nil
	}
}

func toolAllowed(server ServerConfig, name string) bool {
	if len(server.ToolAllow) > 0 {
		for _, allowed := range server.ToolAllow {
			if allowed == name {
				return true
			}
		}
		return false
	}
	for _, denied := range server.ToolDeny {
		if denied == name {
			return false
		}
	}
	return true
}

func (manager *Manager) load() error {
	if _, err := os.Stat(manager.path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	raw, err := os.ReadFile(manager.path)
	if err != nil {
		return err
	}
	var config ConfigFile
	if err := json.Unmarshal(raw, &config); err != nil {
		return err
	}
	for _, server := range config.Servers {
		if strings.TrimSpace(server.Name) == "" {
			continue
		}
		manager.servers[server.Name] = server
	}
	return nil
}

func (manager *Manager) save() error {
	dir := filepath.Dir(manager.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	manager.mu.Lock()
	servers := make([]ServerConfig, 0, len(manager.servers))
	for _, server := range manager.servers {
		servers = append(servers, server)
	}
	manager.mu.Unlock()
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Name < servers[j].Name
	})
	payload, err := json.MarshalIndent(ConfigFile{Servers: servers}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(manager.path, payload, 0644)
}
