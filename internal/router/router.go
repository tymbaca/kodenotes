package router

import (
	"log"
	"net/http"
	"time"
)

type Middleware func(next http.HandlerFunc) http.HandlerFunc

type Router struct {
	routes map[string]*Route

	// Middlewares for all requests. Executes before routesMiddlewares
	//
	// Executes in order: req -> third(second(first(myhandler)))
	globalMiddlewares []Middleware

	// Map with middlewares for each pattern. Pattern is key, middleware is
	// value. Executes before localMiddlewaresMap
	//
	// Executes in order: req -> third(second(first(myhandler)))
	localMiddlewaresMap map[string][]Middleware
}

func Default() *Router {
	return &Router{
		routes:            map[string]*Route{},
		globalMiddlewares: []Middleware{Recover, Logger},
	}
}

func New(middlewares ...Middleware) *Router {
	return &Router{
		routes:            map[string]*Route{},
		globalMiddlewares: middlewares,
	}
}

func (r *Router) Run(addr string) error {
	mux := http.NewServeMux()
	for pattern, route := range r.routes {
		handler := route.GetHandlerFunc()
		for _, middleware := range r.localMiddlewaresMap[pattern] {
			handler = middleware(handler)
		}
		for _, middleware := range r.globalMiddlewares {
			handler = middleware(handler)
		}
		mux.HandleFunc(pattern, handler)
	}

	return http.ListenAndServe(addr, mux)
}

func (r *Router) Handle(method, pattern string, handler http.HandlerFunc) {
	route, ok := r.routes[pattern]
	if ok {
		// pattern already registered
		route.handlers[method] = handler
	} else {
		// new pattern
		r.routes[pattern] = NewRoute()
		route := r.routes[pattern]

		route.handlers[method] = handler
	}
}

// Get is a shortcut to r.Handle(http.MethodGet, ...)
func (r *Router) Get(pattern string, handler http.HandlerFunc) {
	r.Handle(http.MethodGet, pattern, handler)
}

// Post is a shortcut to r.Handle(http.MethodPost, ...)
func (r *Router) Post(pattern string, handler http.HandlerFunc) {
	r.Handle(http.MethodPost, pattern, handler)
}

// Post is a shortcut to r.Handle(http.MethodPut, ...)
func (r *Router) Put(pattern string, handler http.HandlerFunc) {
	r.Handle(http.MethodPut, pattern, handler)
}

// Post is a shortcut to r.Handle(http.MethodPatch, ...)
func (r *Router) Patch(pattern string, handler http.HandlerFunc) {
	r.Handle(http.MethodPatch, pattern, handler)
}

// Post is a shortcut to r.Handle(http.MethodDelete, ...)
func (r *Router) Delete(pattern string, handler http.HandlerFunc) {
	r.Handle(http.MethodDelete, pattern, handler)
}

func (r *Router) Use(middleware ...Middleware) {
	r.globalMiddlewares = append(r.globalMiddlewares, middleware...)
}

func (r *Router) UseLocal(pattern string, middleware ...Middleware) {
	r.localMiddlewaresMap[pattern] = append(r.localMiddlewaresMap[pattern], middleware...)
}

type Route struct {
	// Map with HTTP method as a key and handler as a value. Created by using Handle wethod of Router struct
	handlers map[string]http.HandlerFunc
	// Middlewares for current route.
	middlewares []Middleware
}

func NewRoute() *Route {
	return &Route{
		handlers:    map[string]http.HandlerFunc{},
		middlewares: []Middleware{},
	}
}

func (rt *Route) getHandlerWithMethod() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get corresponding handler if exists
		handler, ok := rt.handlers[r.Method]
		if ok {
			handler(w, r)
		} else {
			http.Error(w, `"detail": "method not allowed"`, http.StatusMethodNotAllowed)
		}
	}
}

func (rt *Route) GetHandlerFunc() http.HandlerFunc {
	handler := rt.getHandlerWithMethod()
	return handler
}

func example() {
	getHandler := func(w http.ResponseWriter, r *http.Request) {}
	postHandler := func(w http.ResponseWriter, r *http.Request) {}

	var timeRecorder Middleware = func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			next(w, r)
			endTime := time.Now()
			duration := startTime.Sub(endTime)
			w.Header().Set("X-Duration", duration.String())
		}
	}

	var authorize Middleware = func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			_, _, ok := r.BasicAuth()
			if !ok {
				http.Error(w, `"detail": "set basic auth"`, http.StatusBadRequest)
			}
			// check if credentials are correct...
			next(w, r)
		}
	}

	router := New(Logger)

	// Set new routes
	router.Get("/posts", getHandler)
	router.Post("/posts", postHandler)

	// Set global middleware
	router.Use(timeRecorder)

	// Set middleware for "/posts" pattern
	router.UseLocal("/posts", authorize)
	// WARN: what about '/posts/my' request? is will ignore it
	// It can be solved with http.ServeMux.

	log.Fatal(router.Run(":8080"))
}
