package handlers

import (
	"net/http"

	"open-sandbox/internal/api"
	"open-sandbox/internal/browser"
	"open-sandbox/internal/mcp"
	"open-sandbox/internal/mcp/tools"
)

func RegisterMCPRoutes(router *api.Router, browserService *browser.Service) {
	registry := mcp.NewRegistry()
	registry.Register(mcp.Tool{
		Name:    "browser.navigate",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "network",
		},
		Handler: tools.BrowserNavigate(browserService),
	})
	registry.Register(mcp.Tool{
		Name:    "browser.screenshot",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Handler: tools.BrowserScreenshot(browserService),
	})
	registry.Register(mcp.Tool{
		Name:    "file.read",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Handler: tools.FileRead(),
	})
	registry.Register(mcp.Tool{
		Name:    "file.write",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Handler: tools.FileWrite(),
	})
	registry.Register(mcp.Tool{
		Name:    "file.list",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Handler: tools.FileList(),
	})
	registry.Register(mcp.Tool{
		Name:    "file.search",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Handler: tools.FileSearch(),
	})
	registry.Register(mcp.Tool{
		Name:    "file.replace",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "workspace",
		},
		Handler: tools.FileReplace(),
	})
	registry.Register(mcp.Tool{
		Name:    "shell.exec",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "exec",
		},
		Handler: tools.ShellExec(),
	})
	registry.Register(mcp.Tool{
		Name:    "code.exec",
		Version: "v1",
		Permissions: mcp.PermissionMeta{
			Allow: true,
			Scope: "exec",
		},
		Handler: tools.CodeExec(),
	})

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
