package api

import "net/http"

type Router struct {
	routes map[string]map[string]HandlerFunc
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

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	methods := router.routes[r.URL.Path]
	if methods == nil {
		WriteAppError(w, NewAppError(CodeNotFound, "not found", http.StatusNotFound))
		return
	}
	handler := methods[r.Method]
	if handler == nil {
		WriteAppError(w, NewAppError(CodeMethodNotAllowed, "method not allowed", http.StatusMethodNotAllowed))
		return
	}
	WrapHandler(handler).ServeHTTP(w, r)
}
