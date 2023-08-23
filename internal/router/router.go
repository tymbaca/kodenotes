package router

import "net/http"

type Middleware func(http.HandlerFunc) http.HandlerFunc

type Router struct {
	routes map[string]Route
	// Middlewares for all requests. Executes before Route middlewares
	middlewares []Middleware
}

func Default() *Router {
	return &Router{
		routes:      nil,
		middlewares: {Logger},
	}
}

func New(middlewares ...Middleware) *Router {
	return &Router{
		routes:      nil,
		middlewares: middlewares, // TODO: in which order thay will be used?
	}
}

type Route struct {
	// Map with HTTP method as a key and handler as a value. Created by using Handle wethod of Router struct
	handlers map[string]http.HandlerFunc
	// Middlewares for current route.
	middlewares []Middleware
}

func (route *Route) getHandlerWithMethod() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler, ok := route.handlers[r.Method]
		if !ok {
			http.Error(w, `"detail": "method not allowed"`, http.StatusMethodNotAllowed)
		}
		handler(w, r)
	}
}

func (rt *Route) GetHandlerFunc() http.HandlerFunc {
	handler := rt.getHandlerWithMethod()
	for _, middleware := range rt.middlewares {
		handler = middleware(handler)
	}
	return handler
}

func (r *Router) Handle(method, pattern string, handler http.HandlerFunc, middlewares ...Middleware) {
	route, ok := r.routes[pattern]
	if ok {
		// pattern already registered
		route.handlers[method] = handler
		routes
	}
}

// Get is a shortcut to r.Handle(http.MethodGet, ...)
func (r *Router) Get(pattern string, handler http.HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodGet, pattern, handler, middlewares...)
}

// Post is a shortcut to r.Handle(http.MethodPost, ...)
func (r *Router) Post(pattern string, handler http.HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodPost, pattern, handler, middlewares...)
}

// Post is a shortcut to r.Handle(http.MethodPut, ...)
func (r *Router) Put(pattern string, handler http.HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodPut, pattern, handler, middlewares...)
}

// Post is a shortcut to r.Handle(http.MethodPatch, ...)
func (r *Router) Patch(pattern string, handler http.HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodPatch, pattern, handler, middlewares...)
}

// Post is a shortcut to r.Handle(http.MethodDelete, ...)
func (r *Router) Delete(pattern string, handler http.HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodDelete, pattern, handler, middlewares...)
}

func GetHandler(handler http.HandlerFunc) http.HandlerFunc {
	return
}

func main() {
	router := New(Logger)
}
