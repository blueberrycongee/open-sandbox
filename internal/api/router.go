package api

import (
	"net/http"
	"strings"
)

type Router struct {
	routes   map[string]map[string]HandlerFunc
	prefixes []prefixRoute
}

type prefixRoute struct {
	prefix  string
	methods map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]HandlerFunc),
	}
}

func (router *Router) Handle(method, path string, handler HandlerFunc) {
	if router.routes[path] == nil {
		router.routes[path] = make(map[string]HandlerFunc)
	}
	router.routes[path][method] = handler
}

func (router *Router) HandlePrefix(method, prefix string, handler HandlerFunc) {
	for i := range router.prefixes {
		if router.prefixes[i].prefix == prefix {
			if router.prefixes[i].methods == nil {
				router.prefixes[i].methods = make(map[string]HandlerFunc)
			}
			router.prefixes[i].methods[method] = handler
			return
		}
	}

	entry := prefixRoute{
		prefix:  prefix,
		methods: map[string]HandlerFunc{method: handler},
	}
	router.prefixes = append(router.prefixes, entry)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	methods := router.routes[r.URL.Path]
	if methods == nil {
		handler := router.matchPrefix(r.URL.Path, r.Method)
		if handler == nil {
			WriteAppError(w, NewAppError(CodeNotFound, "not found", http.StatusNotFound))
			return
		}
		WrapHandler(handler).ServeHTTP(w, r)
		return
	}
	handler := methods[r.Method]
	if handler == nil {
		WriteAppError(w, NewAppError(CodeMethodNotAllowed, "method not allowed", http.StatusMethodNotAllowed))
		return
	}
	WrapHandler(handler).ServeHTTP(w, r)
}

func (router *Router) matchPrefix(path string, method string) HandlerFunc {
	var best HandlerFunc
	bestLen := 0
	for _, entry := range router.prefixes {
		if !strings.HasPrefix(path, entry.prefix) {
			continue
		}
		if len(entry.prefix) <= bestLen {
			continue
		}
		handler := entry.methods[method]
		if handler == nil {
			continue
		}
		best = handler
		bestLen = len(entry.prefix)
	}
	return best
}
