package goweb

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

const (
	methodGET     = "GET"
	methodPUT     = "PUT"
	methodPATCH   = "PATCH"
	methodPOST    = "POST"
	methodDELETE  = "DELETE"
	methodOPTIONS = "OPTIONS"
	methodHEAD    = "HEAD"
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
	optionsRoutes []*route
	headRoutes    []*route

	notFoundHandler Handler

	logger Logger

	redirectWWW bool
}

var paramNameRegExp = regexp.MustCompile(`{([a-zA-Z0-9-]+):?(.*?)}`)

func (e *Engine) registerRoute(method string, path string, handler Handler) {
	rt := getRouteFromPath(path)
	rt.handler = handler
	rt.method = method
	switch method {
	case methodGET:
		e.getRoutes = append(e.getRoutes, rt)
	case methodPUT:
		e.putRoutes = append(e.putRoutes, rt)
	case methodPATCH:
		e.patchRoutes = append(e.patchRoutes, rt)
	case methodPOST:
		e.postRoutes = append(e.postRoutes, rt)
	case methodDELETE:
		e.deleteRoutes = append(e.deleteRoutes, rt)
	case methodOPTIONS:
		e.optionsRoutes = append(e.optionsRoutes, rt)
	case methodHEAD:
		e.headRoutes = append(e.headRoutes, rt)
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
	e.registerRoute(methodGET, path, handler)
}

// PUT registers a route for method PUT.
func (e *Engine) PUT(path string, handler Handler) {
	e.registerRoute(methodPUT, path, handler)
}

// POST registers a route for method POST.
func (e *Engine) POST(path string, handler Handler) {
	e.registerRoute(methodPOST, path, handler)
}

// PATCH registers a route for method PATCH.
func (e *Engine) PATCH(path string, handler Handler) {
	e.registerRoute(methodPATCH, path, handler)
}

// DELETE registers a route for method DELETE.
func (e *Engine) DELETE(path string, handler Handler) {
	e.registerRoute(methodDELETE, path, handler)
}

// OPTIONS registers a route for method OPTIONS.
func (e *Engine) OPTIONS(path string, handler Handler) {
	e.registerRoute(methodOPTIONS, path, handler)
}

// HEAD registers a route for method HEAD.
func (e *Engine) HEAD(path string, handler Handler) {
	e.registerRoute(methodHEAD, path, handler)
}

// Middleware returns a new middleware chain.
func (e *Engine) Middleware(middleware ...Handler) *Middleware {
	return &Middleware{
		chain:  middleware,
		engine: e,
	}
}

// ServeFiles will serve files from the given directory
// with the given path.
func (e *Engine) ServeFiles(path string, directory string) {
	e.registerRoute(methodGET, path+"{name:.+}", func(c *Context) Responder {
		return &FileResponse{
			path:    directory + "/" + c.Param("name"),
			context: c,
		}
	})
	e.registerRoute(methodGET, path, func(c *Context) Responder {
		return &FileResponse{
			path:    directory + "/index.html",
			context: c,
		}
	})
}

// GzipAndServeFiles will serve files from the given directory
// with the given path. If the file size is greater than the given
// size, a gzipped version of the file will be created and served.
func (e *Engine) GzipAndServeFiles(path string, directory string, size int64) {
	e.registerRoute(methodGET, path+"{name:.+}", func(c *Context) Responder {
		return &FileResponse{
			path:    directory + "/" + c.Param("name"),
			context: c,
			gzip:    size,
		}
	})
	e.registerRoute(methodGET, path, func(c *Context) Responder {
		return &FileResponse{
			path:    directory + "/index.html",
			context: c,
			gzip:    size,
		}
	})
}

// NotFound registers a handler to be called if no route is
// matched.
func (e *Engine) NotFound(handler Handler) {
	e.notFoundHandler = handler
}

// ServeHTTP implements the http.Handler interface.
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e.redirectWWW && strings.HasPrefix(r.Host, "www.") {
		w.Header().Set("Connection", "close")
		var scheme string
		if r.TLS == nil {
			scheme = "http"
		} else {
			scheme = "https"
		}
		http.Redirect(
			w,
			r,
			scheme+"://"+strings.TrimPrefix(r.Host, "www.")+r.URL.String(),
			http.StatusMovedPermanently,
		)
		return
	}
	switch r.Method {
	case methodGET:
		e.serve(w, r, e.getRoutes)
	case methodPOST:
		e.serve(w, r, e.postRoutes)
	case methodPUT:
		e.serve(w, r, e.putRoutes)
	case methodPATCH:
		e.serve(w, r, e.patchRoutes)
	case methodDELETE:
		e.serve(w, r, e.deleteRoutes)
	case methodOPTIONS:
		e.serve(w, r, e.optionsRoutes)
	case methodHEAD:
		e.serve(w, r, e.headRoutes)
	}
}

func (e *Engine) serve(w http.ResponseWriter, r *http.Request, routes []*route) {
	c := &Context{
		ResponseWriter: w,
		Request:        r,
		logger:         e.logger,
		status:         http.StatusOK,
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

// Shutdown shuts down the server.
func (e *Engine) Shutdown() error {
	return e.server.Shutdown(context.TODO())
}
