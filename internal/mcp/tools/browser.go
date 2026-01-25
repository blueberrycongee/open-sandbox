package tools

import (
	"context"
	"encoding/json"

	"open-sandbox/internal/browser"
	"open-sandbox/internal/mcp"
)

type browserNavigateParams struct {
	URL string `json:"url"`
}

type browserScreenshotParams struct {
	Path string `json:"path"`
}

type browserClickParams struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func BrowserNavigate(service *browser.Service) mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		if service == nil {
			return nil, toolFailure("browser service unavailable")
		}
		var payload browserNavigateParams
		if err := json.Unmarshal(params, &payload); err != nil {
			return nil, invalidParams("invalid params")
		}
		if payload.URL == "" {
			return nil, invalidParams("url is required")
		}
		if err := service.Navigate(payload.URL); err != nil {
			return nil, toolFailure(err.Error())
		}
		return map[string]any{"navigated": true}, nil
	}
}

func BrowserScreenshot(service *browser.Service) mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		if service == nil {
			return nil, toolFailure("browser service unavailable")
		}
		var payload browserScreenshotParams
		if err := json.Unmarshal(params, &payload); err != nil {
			return nil, invalidParams("invalid params")
		}
		path, errDetail := resolveWorkspacePath(payload.Path)
		if errDetail != nil {
			return nil, errDetail
		}
		if err := service.Screenshot(path); err != nil {
			return nil, toolFailure(err.Error())
		}
		return map[string]any{"path": path}, nil
	}
}

func BrowserClick(service *browser.Service) mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		if service == nil {
			return nil, toolFailure("browser service unavailable")
		}
		var payload browserClickParams
		if err := json.Unmarshal(params, &payload); err != nil {
			return nil, invalidParams("invalid params")
		}
		if err := service.Click(payload.X, payload.Y); err != nil {
			return nil, toolFailure(err.Error())
		}
		return map[string]any{"clicked": true}, nil
	}
}
