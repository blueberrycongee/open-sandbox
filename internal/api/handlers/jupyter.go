package handlers

import (
	"net/http"

	"open-sandbox/internal/api"
)

func RegisterJupyterRoutes(router *api.Router, target string) {
	router.HandlePrefix(http.MethodGet, "/jupyter", JupyterProxyHandler(target))
	router.HandlePrefix(http.MethodPost, "/jupyter", JupyterProxyHandler(target))
	router.HandlePrefix(http.MethodPut, "/jupyter", JupyterProxyHandler(target))
	router.HandlePrefix(http.MethodPatch, "/jupyter", JupyterProxyHandler(target))
	router.HandlePrefix(http.MethodDelete, "/jupyter", JupyterProxyHandler(target))
}

func JupyterProxyHandler(target string) api.HandlerFunc {
	proxy, err := buildReverseProxy(target, "/jupyter")
	if err != nil {
		return func(w http.ResponseWriter, r *http.Request) *api.AppError {
			return api.NewAppError("service_unavailable", err.Error(), http.StatusServiceUnavailable)
		}
	}

	return proxyHandler(proxy)
}
