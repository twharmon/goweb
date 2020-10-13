package goweb

// Map is an alias for map[string]interface{}.
type Map map[string]interface{}

// Handler handles HTTP requests.
type Handler func(*Context) Responder

// New returns a new Engine.
func New() *Engine {
	e := &Engine{
		notFoundHandler: func(c *Context) Responder {
			return c.NotFound().JSON(Map{
				"message": "Page Not Found",
			})
		},
	}

	return e
}
