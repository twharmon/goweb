package goweb

import (
	"fmt"
	"net/http"
)

// Middleware contains a set of Handler functions that will
// be applied in the same order in which they were
// registered.
type Middleware struct {
	chain  []Handler
	engine *Engine
}

func (m *Middleware) apply(handler Handler) Handler {
	return func(c *Context) Responder {
		for _, mw := range m.chain {
			if res := mw(c); res != nil {
				return res
			}
		}
		return handler(c)
	}
}

// Middleware returns a new middleware chain.
func (m *Middleware) Middleware(middleware ...Handler) *Middleware {
	return &Middleware{
		chain:  append(m.chain, middleware...),
		engine: m.engine,
	}
}

// GET registers a route for method GET.
func (m *Middleware) GET(path string, handler Handler) {
	m.engine.GET(path, m.apply(handler))
}

// PUT registers a route for method PUT.
func (m *Middleware) PUT(path string, handler Handler) {
	m.engine.PUT(path, m.apply(handler))
}

// POST registers a route for method POST.
func (m *Middleware) POST(path string, handler Handler) {
	m.engine.POST(path, m.apply(handler))
}

// PATCH registers a route for method PATCH.
func (m *Middleware) PATCH(path string, handler Handler) {
	m.engine.PATCH(path, m.apply(handler))
}

// DELETE registers a route for method DELETE.
func (m *Middleware) DELETE(path string, handler Handler) {
	m.engine.DELETE(path, m.apply(handler))
}

// HEAD registers a route for method HEAD.
func (m *Middleware) HEAD(path string, handler Handler) {
	m.engine.HEAD(path, m.apply(handler))
}

// Resource creates multiple REST handlers from given interface.
func (m *Middleware) Resource(resourceName string, resource Resource) {
	resourcePath := fmt.Sprintf("%s/{%s}", resourceName, resource.Identifier())
	m.engine.registerRoute(http.MethodGet, resourceName, m.apply(resource.Index))
	m.engine.registerRoute(http.MethodGet, resourcePath, m.apply(resource.Get))
	m.engine.registerRoute(http.MethodPut, resourcePath, m.apply(resource.Put))
	m.engine.registerRoute(http.MethodDelete, resourcePath, m.apply(resource.Delete))
	m.engine.registerRoute(http.MethodPost, resourceName, m.apply(resource.Post))
}
