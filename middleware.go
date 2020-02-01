package goweb

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

// OPTIONS registers a route for method OPTIONS.
func (m *Middleware) OPTIONS(path string, handler Handler) {
	m.engine.OPTIONS(path, m.apply(handler))
}

// HEAD registers a route for method HEAD.
func (m *Middleware) HEAD(path string, handler Handler) {
	m.engine.HEAD(path, m.apply(handler))
}
