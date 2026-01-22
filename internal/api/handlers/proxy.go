package handlers

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"

	"open-sandbox/internal/api"
)

var errMissingTarget = errors.New("proxy target is not configured")

func buildReverseProxy(target string, basePath string) (*httputil.ReverseProxy, error) {
	parsed, err := validateTargetURL(target)
	if err != nil {
		return nil, err
	}

	targetQuery := parsed.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = parsed.Scheme
		req.URL.Host = parsed.Host
		req.Host = parsed.Host

		reqPath := strings.TrimPrefix(req.URL.Path, basePath)
		req.URL.Path = joinURLPath(parsed.Path, reqPath)

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}

	proxy := &httputil.ReverseProxy{
		Director: director,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			api.WriteAppError(w, api.NewAppError("proxy_error", err.Error(), http.StatusBadGateway))
		},
	}
	return proxy, nil
}

func joinURLPath(basePath string, reqPath string) string {
	cleanReq := strings.TrimPrefix(reqPath, "/")
	if basePath == "" {
		if cleanReq == "" {
			return "/"
		}
		return "/" + cleanReq
	}
	if cleanReq == "" {
		return basePath
	}
	return path.Join(basePath, cleanReq)
}

func proxyHandler(proxy *httputil.ReverseProxy) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *api.AppError {
		proxy.ServeHTTP(w, r)
		return nil
	}
}

func validateTargetURL(target string) (*url.URL, error) {
	if strings.TrimSpace(target) == "" {
		return nil, errMissingTarget
	}
	parsed, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, errors.New("proxy target must include scheme and host")
	}
	return parsed, nil
}
