package main

import (
	"log"
	"net/http"
	"os"

	"open-sandbox/internal/api"
	"open-sandbox/internal/api/handlers"
	"open-sandbox/internal/browser"
	"open-sandbox/internal/config"
)

func main() {
	if err := config.EnsureWorkspace(); err != nil {
		log.Fatalf("workspace init failed: %v", err)
	}

	addr := os.Getenv("SANDBOX_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	router := api.NewRouter()
	handlers.RegisterSandboxRoutes(router)
	handlers.RegisterBrowserRoutes(router, browser.NewService(os.Getenv("SANDBOX_BROWSER_CDP")))
	handlers.RegisterVNCRoutes(router)
	handlers.RegisterShellRoutes(router)
	handlers.RegisterFileRoutes(router)
	handlers.RegisterCodeExecRoutes(router)
	handlers.RegisterJupyterRoutes(router)
	handlers.RegisterCodeServerRoutes(router)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Printf("listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server stopped: %v", err)
	}
}
