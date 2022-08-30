package goweb

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
)

// Engine contains routing and logging information for your
// app.
type Engine struct {
	server        *http.Server
	getRoutes     []*route
	putRoutes     []*route
	postRoutes    []*route
	patchRoutes   []*route
	deleteRoutes  []*route
	headRoutes    []*route
	optionsRoutes []*route

	notFoundHandler Handler

	loggers []Logger
}

var paramNameRegExp = regexp.MustCompile(`{([a-zA-Z0-9-]+):?(.*?)}`)

func (e *Engine) registerRoute(method string, path string, handler Handler) {
	rt := getRouteFromPath(path)
	rt.handler = handler
	rt.method = method
	switch method {
	case http.MethodGet:
		e.getRoutes = append(e.getRoutes, rt)
	case http.MethodPut:
		e.putRoutes = append(e.putRoutes, rt)
	case http.MethodPatch:
		e.patchRoutes = append(e.patchRoutes, rt)
	case http.MethodPost:
		e.postRoutes = append(e.postRoutes, rt)
	case http.MethodDelete:
		e.deleteRoutes = append(e.deleteRoutes, rt)
	case http.MethodHead:
		e.headRoutes = append(e.headRoutes, rt)
	case http.MethodOptions:
		e.optionsRoutes = append(e.optionsRoutes, rt)
	}
}

func getRouteFromPath(path string) *route {
	rt := new(route)
	if path[0] != '/' {
		panic("path '" + path + "' does not start with '/'")
	}
	pathRegExpStr := "^" + path + "$"
	matches := paramNameRegExp.FindAllStringSubmatch(path, -1)
	for _, match := range matches {
		if match[2] == "" {
			match[2] = `[^\/]*`
		}
		rt.paramNames = append(rt.paramNames, match[1])
		paramInfoRegExp := regexp.MustCompile(fmt.Sprintf("{%s:?(.*?)}", match[1]))
		pathRegExpStr = paramInfoRegExp.ReplaceAllString(pathRegExpStr, "("+match[2]+")")
	}
	rt.regexp = regexp.MustCompile(pathRegExpStr)
	return rt
}

// GET registers a route for method GET.
func (e *Engine) GET(path string, handler Handler) {
	e.registerRoute(http.MethodGet, path, handler)
}

// PUT registers a route for method PUT.
func (e *Engine) PUT(path string, handler Handler) {
	e.registerRoute(http.MethodPut, path, handler)
}

// POST registers a route for method POST.
func (e *Engine) POST(path string, handler Handler) {
	e.registerRoute(http.MethodPost, path, handler)
}

// PATCH registers a route for method PATCH.
func (e *Engine) PATCH(path string, handler Handler) {
	e.registerRoute(http.MethodPatch, path, handler)
}

// DELETE registers a route for method DELETE.
func (e *Engine) DELETE(path string, handler Handler) {
	e.registerRoute(http.MethodDelete, path, handler)
}

// HEAD registers a route for method HEAD.
func (e *Engine) HEAD(path string, handler Handler) {
	e.registerRoute(http.MethodHead, path, handler)
}

// OPTIONS registers a route for method OPTIONS.
func (e *Engine) OPTIONS(path string, handler Handler) {
	e.registerRoute(http.MethodOptions, path, handler)
}

// Middleware returns a new middleware chain.
func (e *Engine) Middleware(middleware ...Handler) *Middleware {
	return &Middleware{
		chain:  middleware,
		engine: e,
	}
}

// NotFound registers a handler to be called if no route is
// matched.
func (e *Engine) NotFound(handler Handler) {
	e.notFoundHandler = handler
}

// ServeHTTP implements the http.Handler interface.
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		e.serve(w, r, e.getRoutes)
	case http.MethodPost:
		e.serve(w, r, e.postRoutes)
	case http.MethodPut:
		e.serve(w, r, e.putRoutes)
	case http.MethodPatch:
		e.serve(w, r, e.patchRoutes)
	case http.MethodDelete:
		e.serve(w, r, e.deleteRoutes)
	case http.MethodHead:
		e.serve(w, r, e.headRoutes)
	case http.MethodOptions:
		e.serve(w, r, e.optionsRoutes)
	}
}

func (e *Engine) serve(w http.ResponseWriter, r *http.Request, routes []*route) {
	c := &Context{
		ResponseWriter: w,
		Request:        r,
		loggers:        e.loggers,
	}
	for _, route := range routes {
		if route.regexp.MatchString(r.URL.Path) {
			c.store = make(Map)
			matches := route.regexp.FindAllStringSubmatch(r.URL.Path, -1)
			for i := 1; i < len(matches[0]); i++ {
				c.params = append(c.params, param{
					key:   route.paramNames[i-1],
					value: matches[0][i],
				})
			}
			if res := route.handler(c); res != nil {
				res.Respond()
			}
			return
		}
	}
	if res := e.notFoundHandler(c); res != nil {
		res.Respond()
	}
}

// Run starts a server on the given port.
func (e *Engine) Run(port string) error {
	e.server = &http.Server{
		Addr:    port,
		Handler: e,
	}
	return e.server.ListenAndServe()
}

// RegisterLogger registers a logger.
func (e *Engine) RegisterLogger(logger Logger) {
	e.loggers = append(e.loggers, logger)
}

// Shutdown shuts down the server.
func (e *Engine) Shutdown() error {
	return e.server.Shutdown(context.TODO())
}
