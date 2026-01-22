package tools

import (
	"path/filepath"
	"strings"

	"open-sandbox/internal/config"
	"open-sandbox/internal/file"
	"open-sandbox/internal/mcp"
)

func resolveWorkspacePath(raw string) (string, *mcp.ErrorDetail) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", invalidParams("path is required")
	}
	resolved := trimmed
	if !filepath.IsAbs(resolved) {
		resolved = filepath.Join(config.WorkspacePath(), resolved)
	}
	resolved = filepath.Clean(resolved)
	if err := file.ValidateWorkspacePath(resolved, config.WorkspacePath()); err != nil {
		return "", invalidParams(err.Error())
	}
	return resolved, nil
}

func resolveWorkspaceDir(raw string) (string, *mcp.ErrorDetail) {
	if strings.TrimSpace(raw) == "" {
		return config.WorkspacePath(), nil
	}
	return resolveWorkspacePath(raw)
}

func invalidParams(message string) *mcp.ErrorDetail {
	detail := mcp.NewErrorDetail("invalid_params", message, mcp.KindInvalidParams)
	return &detail
}

func toolFailure(message string) *mcp.ErrorDetail {
	detail := mcp.NewErrorDetail("tool_error", message, mcp.KindToolError)
	return &detail
}
