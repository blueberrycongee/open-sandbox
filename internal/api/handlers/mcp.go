package handlers

import (
	"net/http"

	"open-sandbox/internal/api"
	"open-sandbox/internal/browser"
	"open-sandbox/internal/mcp"
	"open-sandbox/internal/mcp/tools"
)

func NewMCPRegistry(browserService *browser.Service) *mcp.Registry {
	registry := mcp.NewRegistry()
	browserNavigateSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"url": map[string]any{"type": "string"},
			},
			"required": []string{"url"},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"navigated": map[string]any{"type": "boolean"},
			},
			"required": []string{"navigated"},
		},
	}
	browserScreenshotSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{"type": "string"},
			},
			"required": []string{"path"},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{"type": "string"},
			},
			"required": []string{"path"},
		},
	}
	browserClickSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"x": map[string]any{"type": "number"},
				"y": map[string]any{"type": "number"},
			},
			"required": []string{"x", "y"},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"clicked": map[string]any{"type": "boolean"},
			},
			"required": []string{"clicked"},
		},
	}
	browserFormInputFillSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"selector": map[string]any{"type": "string"},
				"value":    map[string]any{"type": "string"},
			},
			"required": []string{"selector", "value"},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"filled": map[string]any{"type": "boolean"},
			},
			"required": []string{"filled"},
		},
	}
	browserSelectSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"selector": map[string]any{"type": "string"},
				"value":    map[string]any{"type": "string"},
			},
			"required": []string{"selector", "value"},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"selected": map[string]any{"type": "boolean"},
			},
			"required": []string{"selected"},
		},
	}
	browserScrollSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"x": map[string]any{"type": "number"},
				"y": map[string]any{"type": "number"},
			},
			"required": []string{"x", "y"},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"scrolled": map[string]any{"type": "boolean"},
			},
			"required": []string{"scrolled"},
		},
	}
	browserEvaluateSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"expression": map[string]any{"type": "string"},
			},
			"required": []string{"expression"},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"result": map[string]any{},
			},
			"required": []string{"result"},
		},
	}
	browserNewTabSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"url": map[string]any{"type": "string"},
			},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"index": map[string]any{"type": "integer"},
			},
			"required": []string{"index"},
		},
	}
	browserSwitchTabSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"index": map[string]any{"type": "integer"},
			},
			"required": []string{"index"},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"index": map[string]any{"type": "integer"},
			},
			"required": []string{"index"},
		},
	}
	browserTabListSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"tabs": map[string]any{"type": "array"},
			},
			"required": []string{"tabs"},
		},
	}
	browserCloseTabSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"index": map[string]any{"type": "integer"},
			},
			"required": []string{"index"},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"closed": map[string]any{"type": "boolean"},
			},
			"required": []string{"closed"},
		},
	}
	browserDownloadListSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"downloads": map[string]any{"type": "array"},
			},
			"required": []string{"downloads"},
		},
	}
	browserPressKeySchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"keys": map[string]any{"type": "string"},
			},
			"required": []string{"keys"},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"pressed": map[string]any{"type": "boolean"},
			},
			"required": []string{"pressed"},
		},
	}
	fileReadSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{"type": "string"},
			},
			"required": []string{"path"},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"content": map[string]any{"type": "string"},
			},
			"required": []string{"content"},
		},
	}
	fileWriteSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"path":    map[string]any{"type": "string"},
				"content": map[string]any{"type": "string"},
			},
			"required": []string{"path", "content"},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{"type": "string"},
			},
			"required": []string{"path"},
		},
	}
	fileListSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{"type": "string"},
			},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"entries": map[string]any{"type": "array"},
			},
			"required": []string{"entries"},
		},
	}
	fileSearchSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"path":  map[string]any{"type": "string"},
				"query": map[string]any{"type": "string"},
			},
			"required": []string{"path", "query"},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"matches": map[string]any{"type": "array"},
			},
			"required": []string{"matches"},
		},
	}
	fileReplaceSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"path":    map[string]any{"type": "string"},
				"search":  map[string]any{"type": "string"},
				"replace": map[string]any{"type": "string"},
			},
			"required": []string{"path", "search", "replace"},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"count": map[string]any{"type": "integer"},
			},
			"required": []string{"count"},
		},
	}
	shellExecSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"command":     map[string]any{"type": "string"},
				"args":        map[string]any{"type": "array"},
				"working_dir": map[string]any{"type": "string"},
			},
			"required": []string{"command"},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"stdout":    map[string]any{"type": "string"},
				"stderr":    map[string]any{"type": "string"},
				"exit_code": map[string]any{"type": "integer"},
			},
			"required": []string{"stdout", "stderr", "exit_code"},
		},
	}
	codeExecSchema := mcp.ToolSchema{
		Input: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"runtime":     map[string]any{"type": "string"},
				"args":        map[string]any{"type": "array"},
				"working_dir": map[string]any{"type": "string"},
			},
			"required": []string{"runtime"},
		},
		Output: mcp.JSONSchema{
			"type": "object",
			"properties": map[string]any{
				"stdout":    map[string]any{"type": "string"},
				"stderr":    map[string]any{"type": "string"},
				"exit_code": map[string]any{"type": "integer"},
			},
			"required": []string{"stdout", "stderr", "exit_code"},
		},
	}
	registry.Register(mcp.Tool{
		Name:    "browser.navigate",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "network",
		},
		Schema:  browserNavigateSchema,
		Handler: tools.BrowserNavigate(browserService),
	})
	registry.Register(mcp.Tool{
		Name:    "browser.screenshot",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  browserScreenshotSchema,
		Handler: tools.BrowserScreenshot(browserService),
	})
	registry.Register(mcp.Tool{
		Name:    "browser.click",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  browserClickSchema,
		Handler: tools.BrowserClick(browserService),
	})
	registry.Register(mcp.Tool{
		Name:    "browser_navigate",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "network",
		},
		Schema:  browserNavigateSchema,
		Handler: tools.BrowserNavigate(browserService),
	})
	registry.Register(mcp.Tool{
		Name:    "browser_screenshot",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  browserScreenshotSchema,
		Handler: tools.BrowserScreenshot(browserService),
	})
	registry.Register(mcp.Tool{
		Name:    "browser_click",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  browserClickSchema,
		Handler: tools.BrowserClick(browserService),
	})
	registry.Register(mcp.Tool{
		Name:    "browser_form_input_fill",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  browserFormInputFillSchema,
		Handler: tools.BrowserFormInputFill(browserService),
	})
	registry.Register(mcp.Tool{
		Name:    "browser_select",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  browserSelectSchema,
		Handler: tools.BrowserSelect(browserService),
	})
	registry.Register(mcp.Tool{
		Name:    "browser_scroll",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  browserScrollSchema,
		Handler: tools.BrowserScroll(browserService),
	})
	registry.Register(mcp.Tool{
		Name:    "browser_evaluate",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  browserEvaluateSchema,
		Handler: tools.BrowserEvaluate(browserService),
	})
	registry.Register(mcp.Tool{
		Name:    "browser_new_tab",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  browserNewTabSchema,
		Handler: tools.BrowserNewTab(browserService),
	})
	registry.Register(mcp.Tool{
		Name:    "browser_switch_tab",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  browserSwitchTabSchema,
		Handler: tools.BrowserSwitchTab(browserService),
	})
	registry.Register(mcp.Tool{
		Name:    "browser_tab_list",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  browserTabListSchema,
		Handler: tools.BrowserTabList(browserService),
	})
	registry.Register(mcp.Tool{
		Name:    "browser_close_tab",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  browserCloseTabSchema,
		Handler: tools.BrowserCloseTab(browserService),
	})
	registry.Register(mcp.Tool{
		Name:    "browser_get_download_list",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  browserDownloadListSchema,
		Handler: tools.BrowserGetDownloadList(browserService),
	})
	registry.Register(mcp.Tool{
		Name:    "browser_press_key",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  browserPressKeySchema,
		Handler: tools.BrowserPressKey(browserService),
	})
	registry.Register(mcp.Tool{
		Name:    "file.read",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  fileReadSchema,
		Handler: tools.FileRead(),
	})
	registry.Register(mcp.Tool{
		Name:    "file.write",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  fileWriteSchema,
		Handler: tools.FileWrite(),
	})
	registry.Register(mcp.Tool{
		Name:    "file.list",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  fileListSchema,
		Handler: tools.FileList(),
	})
	registry.Register(mcp.Tool{
		Name:    "file.search",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  fileSearchSchema,
		Handler: tools.FileSearch(),
	})
	registry.Register(mcp.Tool{
		Name:    "file.replace",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Schema:  fileReplaceSchema,
		Handler: tools.FileReplace(),
	})
	registry.Register(mcp.Tool{
		Name:    "shell.exec",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "exec",
		},
		Schema:  shellExecSchema,
		Handler: tools.ShellExec(),
	})
	registry.Register(mcp.Tool{
		Name:    "code.exec",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "exec",
		},
		Schema:  codeExecSchema,
		Handler: tools.CodeExec(),
	})
	return registry
}

func RegisterMCPRoutes(router *api.Router, browserService *browser.Service) {
	registry := NewMCPRegistry(browserService)

	auth, authErr := mcp.NewAuthenticator(mcp.LoadAuthConfig())
	server := mcp.NewServer(registry, auth, authErr)

	router.Handle("POST", "/mcp", func(w http.ResponseWriter, r *http.Request) *api.AppError {
		server.ServeHTTP(w, r)
		return nil
	})
	router.Handle("GET", "/mcp/sse", func(w http.ResponseWriter, r *http.Request) *api.AppError {
		server.ServeSSE(w, r)
		return nil
	})
}
