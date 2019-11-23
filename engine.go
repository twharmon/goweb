package goweb

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/crypto/acme/autocert"
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
	getRoutes     []*route
	putRoutes     []*route
	postRoutes    []*route
	patchRoutes   []*route
	deleteRoutes  []*route
	optionsRoutes []*route
	headRoutes    []*route

	notFoundHandler Handler

	loggers []Logger

	redirectWWW bool
}

var paramNameRegExp = regexp.MustCompile(`{([a-zA-Z0-9-]+):?(.*?)}`)

func (e *Engine) registerRoute(method string, path string, handler Handler) {
	if len(path) == 0 {
		panic("path can not be empty")
	}
	rt := new(route)
	if path[0] != '/' {
		parts := strings.Split(path, "/")
		if len(parts) == 1 {
			panic("path '" + path + "' does not start with '/'")
		}
		rt.host = parts[0]
		path = "/" + strings.Join(parts[1:], "/")
	}
	pathRegExpStr := "^" + path + "$"
	matches := paramNameRegExp.FindAllStringSubmatch(path, -1)
	for _, match := range matches {
		if match[2] == "" {
			match[2] = `[^\\]*`
		}
		rt.paramNames = append(rt.paramNames, match[1])
		paramInfoRegExp := regexp.MustCompile(fmt.Sprintf("{%s:?(.*?)}", match[1]))
		pathRegExpStr = paramInfoRegExp.ReplaceAllString(pathRegExpStr, "("+match[2]+")")
	}
	rt.regexp = regexp.MustCompile(pathRegExpStr)
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

// ServeFiles will serve files from the given directory
// with the given path.
func (e *Engine) ServeFiles(path string, directory string) {
	e.registerRoute(methodGET, path+"{name:.+}", func(c *Context) Responder {
		http.ServeFile(c.writer, c.request, directory+"/"+c.Param("name"))
		return nil
	})
	e.registerRoute(methodGET, path, func(c *Context) Responder {
		http.ServeFile(c.writer, c.request, directory+"/index.html")
		return nil
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
		writer:  w,
		request: r,
		engine:  e,
	}
	for _, route := range routes {
		hostMatches := route.host == "" || route.host == r.Host
		if hostMatches && route.regexp.MatchString(r.URL.Path) {
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

// RedirectWWW redirects all requests to the non-www host.
func (e *Engine) RedirectWWW() {
	e.redirectWWW = true
}

// AddCustomLogger adds a Logger to the Engine.
func (e *Engine) AddCustomLogger(logger Logger) {
	e.loggers = append(e.loggers, logger)
}

// AddStdLogger adds a logger that will output log messages
// with log.Println.
func (e *Engine) AddStdLogger(level int) {
	e.loggers = append(e.loggers, &stdLogger{
		minLevel: level,
	})
}

// AddSlackLogger adds a logger that will output log
// messages to a slack webhook. In order to use this
// channel, you must set the environment variable
// SLACK_LOG_WEBHOOK.
func (e *Engine) AddSlackLogger(level int) {
	e.loggers = append(e.loggers, &slackLogger{
		minLevel: level,
	})
}

// Run starts a server on the given port.
func (e *Engine) Run(port string) error {
	return http.ListenAndServe(port, e)
}

// RunTLS starts a server on port :443.
func (e *Engine) RunTLS(config *TLSConfig) error {
	m := &autocert.Manager{
		Cache:  autocert.DirCache(config.CertDir),
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(_ context.Context, host string) error {
			return config.HostPolicy(host)
		},
	}
	if config.RedirectHTTP {
		go http.ListenAndServe(":http", m.HTTPHandler(nil))
	}
	s := &http.Server{
		Addr:      ":https",
		TLSConfig: m.TLSConfig(),
		Handler:   e,
	}
	return s.ListenAndServeTLS("", "")
}
