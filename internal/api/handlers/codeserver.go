package handlers

import (
	"net/http"

	"open-sandbox/internal/api"
)

func RegisterCodeServerRoutes(router *api.Router, target string) {
	router.HandlePrefix(http.MethodGet, "/code-server", CodeServerProxyHandler(target))
	router.HandlePrefix(http.MethodPost, "/code-server", CodeServerProxyHandler(target))
	router.HandlePrefix(http.MethodPut, "/code-server", CodeServerProxyHandler(target))
	router.HandlePrefix(http.MethodPatch, "/code-server", CodeServerProxyHandler(target))
	router.HandlePrefix(http.MethodDelete, "/code-server", CodeServerProxyHandler(target))
}

func CodeServerProxyHandler(target string) api.HandlerFunc {
	proxy, err := buildReverseProxy(target, "/code-server")
	if err != nil {
		return func(w http.ResponseWriter, r *http.Request) *api.AppError {
			return api.NewAppError("service_unavailable", err.Error(), http.StatusServiceUnavailable)
		}
	}

	return proxyHandler(proxy)
}
