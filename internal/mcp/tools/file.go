package tools

import (
	"context"
	"encoding/json"

	"open-sandbox/internal/file"
	"open-sandbox/internal/mcp"
)

type filePathParams struct {
	Path string `json:"path"`
}

type fileWriteParams struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type fileSearchParams struct {
	Path  string `json:"path"`
	Query string `json:"query"`
}

type fileReplaceParams struct {
	Path    string `json:"path"`
	Search  string `json:"search"`
	Replace string `json:"replace"`
}

func FileRead() mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		var payload filePathParams
		if err := json.Unmarshal(params, &payload); err != nil {
			return nil, invalidParams("invalid params")
		}
		path, errDetail := resolveWorkspacePath(payload.Path)
		if errDetail != nil {
			return nil, errDetail
		}
		content, err := file.Read(path)
		if err != nil {
			return nil, toolFailure(err.Error())
		}
		return map[string]any{"content": content}, nil
	}
}

func FileWrite() mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		var payload fileWriteParams
		if err := json.Unmarshal(params, &payload); err != nil {
			return nil, invalidParams("invalid params")
		}
		path, errDetail := resolveWorkspacePath(payload.Path)
		if errDetail != nil {
			return nil, errDetail
		}
		if err := file.Write(path, payload.Content); err != nil {
			return nil, toolFailure(err.Error())
		}
		return map[string]any{"path": path}, nil
	}
}

func FileList() mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		var payload filePathParams
		if len(params) > 0 {
			if err := json.Unmarshal(params, &payload); err != nil {
				return nil, invalidParams("invalid params")
			}
		}
		path, errDetail := resolveWorkspaceDir(payload.Path)
		if errDetail != nil {
			return nil, errDetail
		}
		entries, err := file.List(path)
		if err != nil {
			return nil, toolFailure(err.Error())
		}
		return map[string]any{"entries": entries}, nil
	}
}

func FileSearch() mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		var payload fileSearchParams
		if err := json.Unmarshal(params, &payload); err != nil {
			return nil, invalidParams("invalid params")
		}
		path, errDetail := resolveWorkspacePath(payload.Path)
		if errDetail != nil {
			return nil, errDetail
		}
		matches, err := file.Search(path, payload.Query)
		if err != nil {
			return nil, toolFailure(err.Error())
		}
		return map[string]any{"matches": matches}, nil
	}
}

func FileReplace() mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		var payload fileReplaceParams
		if err := json.Unmarshal(params, &payload); err != nil {
			return nil, invalidParams("invalid params")
		}
		path, errDetail := resolveWorkspacePath(payload.Path)
		if errDetail != nil {
			return nil, errDetail
		}
		count, err := file.Replace(path, payload.Search, payload.Replace)
		if err != nil {
			return nil, toolFailure(err.Error())
		}
		return map[string]any{"count": count}, nil
	}
}
