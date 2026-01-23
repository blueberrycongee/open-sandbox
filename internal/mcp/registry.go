package mcp

import (
	"context"
	"encoding/json"
	"sort"
	"sync"
)

type ToolHandler func(ctx context.Context, params json.RawMessage) (any, *ErrorDetail)

type PermissionMeta struct {
	Allow  bool   `json:"allow"`
	Scope  string `json:"scope"`
	Reason string `json:"reason,omitempty"`
}

type Tool struct {
	Name        string
	Version     string
	Permissions PermissionMeta
	Schema      ToolSchema
	Handler     ToolHandler
}

type ToolInfo struct {
	Name        string         `json:"name"`
	Version     string         `json:"version"`
	Permissions PermissionMeta `json:"permissions"`
	Schema      ToolSchema     `json:"schema,omitempty"`
}

type Registry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{tools: make(map[string]Tool)}
}

func (registry *Registry) Register(tool Tool) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.tools[tool.Name] = tool
}

func (registry *Registry) Get(name string) (Tool, bool) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	tool, ok := registry.tools[name]
	return tool, ok
}

func (registry *Registry) List() []ToolInfo {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	infos := make([]ToolInfo, 0, len(registry.tools))
	for _, tool := range registry.tools {
		infos = append(infos, ToolInfo{
			Name:        tool.Name,
			Version:     tool.Version,
			Permissions: tool.Permissions,
			Schema:      tool.Schema,
		})
	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name < infos[j].Name
	})
	return infos
}
