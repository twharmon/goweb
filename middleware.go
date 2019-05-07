package goweb

// Middleware contains a set of Handler functions that will
// be applied in the same order in which they were
// registerend.
type Middleware struct {
	chain []Handler
}

// Use appends a Handler to the Middleware.
func (m *Middleware) Use(handlers ...Handler) {
	m.chain = append(m.chain, handlers...)
}

// Apply wraps all of a Middleware's Handlers to a Handler.
func (m *Middleware) Apply(handler Handler) Handler {
	return func(c *Context) *Response {
		for _, mw := range m.chain {
			if res := mw(c); res != nil {
				return res
			}
		}
		return handler(c)
	}
}
