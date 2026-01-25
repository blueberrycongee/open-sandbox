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
	Name         string         `json:"name"`
	Version      string         `json:"version"`
	Permissions  PermissionMeta `json:"permissions"`
	InputSchema  JSONSchema     `json:"inputSchema"`
	OutputSchema JSONSchema     `json:"outputSchema,omitempty"`
	Schema       ToolSchema     `json:"schema,omitempty"`
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

func (registry *Registry) Unregister(name string) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	delete(registry.tools, name)
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
		inputSchema := normalizeSchema(tool.Schema.Input)
		outputSchema := normalizeOptionalSchema(tool.Schema.Output)
		infos = append(infos, ToolInfo{
			Name:         tool.Name,
			Version:      tool.Version,
			Permissions:  tool.Permissions,
			InputSchema:  inputSchema,
			OutputSchema: outputSchema,
			Schema:       tool.Schema,
		})
	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name < infos[j].Name
	})
	return infos
}

func normalizeSchema(schema JSONSchema) JSONSchema {
	if schema == nil || len(schema) == 0 {
		return JSONSchema{"type": "object"}
	}
	if _, ok := schema["type"]; ok {
		return schema
	}
	clone := make(JSONSchema, len(schema)+1)
	for key, value := range schema {
		clone[key] = value
	}
	clone["type"] = "object"
	return clone
}

func normalizeOptionalSchema(schema JSONSchema) JSONSchema {
	if schema == nil || len(schema) == 0 {
		return nil
	}
	return normalizeSchema(schema)
}
