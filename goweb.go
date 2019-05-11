package goweb

// Map is an alias for map[string]interface{}.
type Map map[string]interface{}

// Handler handles HTTP requests.
type Handler func(*Context) Responder

// New returns a new Engine.
func New() *Engine {
	return &Engine{
		notFoundHandler: func(c *Context) Responder {
			return c.NotFound().PlainText("Page not found")
		},
	}
}

// NewMiddleware returns a new MiddlewareChain that can be
// applied to any Handler.
func NewMiddleware() *Middleware {
	return &Middleware{
		chain: []Handler{},
	}
}
