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
	registeredTool map[string]struct{}
	client         *Client
}

func NewManager(path string) (*Manager, error) {
	manager := &Manager{
		path:           path,
		servers:        make(map[string]ServerConfig),
		registeredTool: make(map[string]struct{}),
		client:         NewClient(),
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
	manager.mu.Unlock()

	for toolName := range manager.registeredTool {
		registry.Unregister(toolName)
	}
	manager.registeredTool = make(map[string]struct{})

	for _, server := range servers {
		tools, err := manager.client.ToolsList(ctx, server)
		if err != nil {
			continue
		}
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
			manager.registeredTool[registeredName] = struct{}{}
		}
	}
	return nil
}

func buildRemoteToolHandler(client *Client, server ServerConfig, name string) mcp.ToolHandler {
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
