package gorouter

import (
	"fmt"
	"net/http"
	"regexp"
)

// Handler .
type Handler func(http.ResponseWriter, *http.Request, Params)

// Route .
type Route struct {
	re         *regexp.Regexp
	handler    Handler
	paramNames []string
	method     string
}

// Router .
type Router struct {
	getRoutes     []*Route
	putRoutes     []*Route
	postRoutes    []*Route
	patchRoutes   []*Route
	deleteRoutes  []*Route
	optionsRoutes []*Route
	headRoutes    []*Route
	staticDir     string
}

// New .
func New() *Router {
	return new(Router)
}

var paramNameRegExp = regexp.MustCompile(`{([a-zA-Z0-9-_]+):?(.*?)}`)

func (r *Router) registerRoute(method string, path string, handler func(http.ResponseWriter, *http.Request, Params)) {
	if len(path) == 0 {
		panic("Path can not be empty")
	}
	if path[0] != '/' {
		panic("Path '" + path + "' does not start with '/'")
	}
	rt := new(Route)
	pathRegExpStr := "^" + path + "$"
	matches := paramNameRegExp.FindAllStringSubmatch(path, -1)
	for _, match := range matches {
		if match[2] == "" {
			panic("Regular expressions are required for all route parameters")
		}
		rt.paramNames = append(rt.paramNames, match[1])
		paramInfoRegExp := regexp.MustCompile(fmt.Sprintf(`{%s:?(.*?)}`, match[1]))
		pathRegExpStr = paramInfoRegExp.ReplaceAllString(pathRegExpStr, "("+match[2]+")")
	}
	rt.re = regexp.MustCompile(pathRegExpStr)
	rt.handler = handler
	rt.method = method
	switch method {
	case "GET":
		r.getRoutes = append(r.getRoutes, rt)
	case "PUT":
		r.putRoutes = append(r.putRoutes, rt)
	case "PATCH":
		r.patchRoutes = append(r.patchRoutes, rt)
	case "POST":
		r.postRoutes = append(r.postRoutes, rt)
	case "DELETE":
		r.deleteRoutes = append(r.deleteRoutes, rt)
	case "OPTIONS":
		r.optionsRoutes = append(r.optionsRoutes, rt)
	case "HEAD":
		r.headRoutes = append(r.headRoutes, rt)
	default:
		panic("can not register route for method '" + method + "'")
	}
}

// Get .
func (r *Router) Get(path string, handler Handler) {
	r.registerRoute("GET", path, handler)
}

// Put .
func (r *Router) Put(path string, handler Handler) {
	r.registerRoute("PUT", path, handler)
}

// Post .
func (r *Router) Post(path string, handler Handler) {
	r.registerRoute("POST", path, handler)
}

// Patch .
func (r *Router) Patch(path string, handler Handler) {
	r.registerRoute("PATCH", path, handler)
}

// Delete .
func (r *Router) Delete(path string, handler Handler) {
	r.registerRoute("DELETE", path, handler)
}

// Options .
func (r *Router) Options(path string, handler Handler) {
	r.registerRoute("OPTIONS", path, handler)
}

// Head .
func (r *Router) Head(path string, handler Handler) {
	r.registerRoute("HEAD", path, handler)
}

// StaticFiles .
func (r *Router) StaticFiles(dir string) {
	r.staticDir = dir
}

// ServeHTTP .
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		r.serve(w, req, r.getRoutes)
	case "POST":
		r.serve(w, req, r.postRoutes)
	case "PUT":
		r.serve(w, req, r.putRoutes)
	case "PATCH":
		r.serve(w, req, r.patchRoutes)
	case "DELETE":
		r.serve(w, req, r.deleteRoutes)
	case "OPTIONS":
		r.serve(w, req, r.optionsRoutes)
	case "HEAD":
		r.serve(w, req, r.headRoutes)
	}
}

func (r *Router) serve(w http.ResponseWriter, req *http.Request, routes []*Route) {
	for _, route := range routes {
		if route.re.MatchString(req.URL.Path) {
			var params Params
			matches := route.re.FindAllStringSubmatch(req.URL.Path, -1)
			for i := 1; i < len(matches[0]); i++ {
				params = append(params, &Param{
					Key:   route.paramNames[i-1],
					Value: matches[0][i],
				})
			}
			route.handler(w, req, params)
			return
		}
	}
	if r.staticDir != "" {
		http.ServeFile(w, req, r.staticDir+"/"+req.URL.Path)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Page not found"))
}
