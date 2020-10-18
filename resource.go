package goweb

import (
	"fmt"
	"net/http"
)

// Resource creates multiple REST handlers from given interface.
func (e *Engine) Resource(resourceName string, resource iResource) {
	resourcePath := fmt.Sprintf("/%s/{%s}", resourceName, resource.Identifier())
	e.registerRoute(http.MethodGet, resourceName, resource.Index)
	e.registerRoute(http.MethodGet, resourcePath, resource.Get)
	e.registerRoute(http.MethodPut, resourcePath, resource.Put)
	e.registerRoute(http.MethodDelete, resourcePath, resource.Delete)
	e.registerRoute(http.MethodPost, resourceName, resource.Post)
}

type iResource interface {
	Index(c *Context) Responder
	Get(c *Context) Responder
	Put(c *Context) Responder
	Delete(c *Context) Responder
	Post(c *Context) Responder
	Identifier() string
}
