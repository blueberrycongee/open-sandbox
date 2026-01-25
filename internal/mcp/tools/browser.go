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

type browserFormInputFillParams struct {
	Selector string `json:"selector"`
	Value    string `json:"value"`
}

type browserSelectParams struct {
	Selector string `json:"selector"`
	Value    string `json:"value"`
}

type browserScrollParams struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type browserEvaluateParams struct {
	Expression string `json:"expression"`
}

type browserNewTabParams struct {
	URL string `json:"url"`
}

type browserSwitchTabParams struct {
	Index int `json:"index"`
}

type browserCloseTabParams struct {
	Index int `json:"index"`
}

type browserPressKeyParams struct {
	Keys string `json:"keys"`
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

func BrowserFormInputFill(service *browser.Service) mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		if service == nil {
			return nil, toolFailure("browser service unavailable")
		}
		var payload browserFormInputFillParams
		if err := json.Unmarshal(params, &payload); err != nil {
			return nil, invalidParams("invalid params")
		}
		if payload.Selector == "" {
			return nil, invalidParams("selector is required")
		}
		if err := service.FormInputFill(payload.Selector, payload.Value); err != nil {
			return nil, toolFailure(err.Error())
		}
		return map[string]any{"filled": true}, nil
	}
}

func BrowserSelect(service *browser.Service) mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		if service == nil {
			return nil, toolFailure("browser service unavailable")
		}
		var payload browserSelectParams
		if err := json.Unmarshal(params, &payload); err != nil {
			return nil, invalidParams("invalid params")
		}
		if payload.Selector == "" {
			return nil, invalidParams("selector is required")
		}
		if err := service.ElementSelect(payload.Selector, payload.Value); err != nil {
			return nil, toolFailure(err.Error())
		}
		return map[string]any{"selected": true}, nil
	}
}

func BrowserScroll(service *browser.Service) mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		if service == nil {
			return nil, toolFailure("browser service unavailable")
		}
		var payload browserScrollParams
		if err := json.Unmarshal(params, &payload); err != nil {
			return nil, invalidParams("invalid params")
		}
		if err := service.Scroll(payload.X, payload.Y); err != nil {
			return nil, toolFailure(err.Error())
		}
		return map[string]any{"scrolled": true}, nil
	}
}

func BrowserEvaluate(service *browser.Service) mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		if service == nil {
			return nil, toolFailure("browser service unavailable")
		}
		var payload browserEvaluateParams
		if err := json.Unmarshal(params, &payload); err != nil {
			return nil, invalidParams("invalid params")
		}
		if payload.Expression == "" {
			return nil, invalidParams("expression is required")
		}
		result, err := service.Evaluate(payload.Expression)
		if err != nil {
			return nil, toolFailure(err.Error())
		}
		return map[string]any{"result": result}, nil
	}
}

func BrowserNewTab(service *browser.Service) mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		if service == nil {
			return nil, toolFailure("browser service unavailable")
		}
		var payload browserNewTabParams
		if len(params) > 0 {
			if err := json.Unmarshal(params, &payload); err != nil {
				return nil, invalidParams("invalid params")
			}
		}
		index, err := service.NewTab(payload.URL)
		if err != nil {
			return nil, toolFailure(err.Error())
		}
		return map[string]any{"index": index}, nil
	}
}

func BrowserSwitchTab(service *browser.Service) mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		if service == nil {
			return nil, toolFailure("browser service unavailable")
		}
		var payload browserSwitchTabParams
		if err := json.Unmarshal(params, &payload); err != nil {
			return nil, invalidParams("invalid params")
		}
		if err := service.SwitchTab(payload.Index); err != nil {
			return nil, toolFailure(err.Error())
		}
		return map[string]any{"index": payload.Index}, nil
	}
}

func BrowserTabList(service *browser.Service) mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		if service == nil {
			return nil, toolFailure("browser service unavailable")
		}
		tabs, err := service.TabList()
		if err != nil {
			return nil, toolFailure(err.Error())
		}
		return map[string]any{"tabs": tabs}, nil
	}
}

func BrowserCloseTab(service *browser.Service) mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		if service == nil {
			return nil, toolFailure("browser service unavailable")
		}
		var payload browserCloseTabParams
		if err := json.Unmarshal(params, &payload); err != nil {
			return nil, invalidParams("invalid params")
		}
		if err := service.CloseTab(payload.Index); err != nil {
			return nil, toolFailure(err.Error())
		}
		return map[string]any{"closed": true}, nil
	}
}

func BrowserGetDownloadList(service *browser.Service) mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		if service == nil {
			return nil, toolFailure("browser service unavailable")
		}
		downloads := service.DownloadList()
		return map[string]any{"downloads": downloads}, nil
	}
}

func BrowserPressKey(service *browser.Service) mcp.ToolHandler {
	return func(ctx context.Context, params json.RawMessage) (any, *mcp.ErrorDetail) {
		if service == nil {
			return nil, toolFailure("browser service unavailable")
		}
		var payload browserPressKeyParams
		if err := json.Unmarshal(params, &payload); err != nil {
			return nil, invalidParams("invalid params")
		}
		if payload.Keys == "" {
			return nil, invalidParams("keys is required")
		}
		if err := service.PressKey(payload.Keys); err != nil {
			return nil, toolFailure(err.Error())
		}
		return map[string]any{"pressed": true}, nil
	}
}
