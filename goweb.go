package goweb

// Map is an alias for map[string]interface{}.
type Map map[string]interface{}

// Handler handles HTTP requests.
type Handler func(*Context) Responder

// Config configures the app engine.
type Config struct {
	Logger Logger
}

// New returns a new Engine.
func New(config *Config) *Engine {
	e := &Engine{
		notFoundHandler: func(c *Context) Responder {
			return c.NotFound().Text("Page not found")
		},
	}

	if config != nil {
		e.logger = config.Logger
	}

	return e
}
