package goweb

import (
	"fmt"
	"net/http"
)

//Resource creates multiple REST handlers from given interface.
func (e *Engine) Resource(resourceName string, r iResource) {
	resourcePath := fmt.Sprintf("%s/{%s}", resourceName, r.Identifier())
	e.registerRoute(http.MethodGet, resourceName, r.Index)
	e.registerRoute(http.MethodGet, resourcePath, r.Get)
	e.registerRoute(http.MethodPut, resourcePath, r.Put)
	e.registerRoute(http.MethodDelete, resourcePath, r.Delete)
	e.registerRoute(http.MethodPost, resourceName, r.Post)
}

type iResource interface {
	Index(c *Context) Responder
	Get(c *Context) Responder
	Put(c *Context) Responder
	Delete(c *Context) Responder
	Post(c *Context) Responder
	Identifier() string
}
