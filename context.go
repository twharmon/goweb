package goweb

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

// Context provides helper methods to read the request, get
// and set values in a data store, and send a response to
// the client. A new Context is constructed for each
// request, and is dropped when the response is sent.
type Context struct {
	writer  http.ResponseWriter
	request *http.Request
	params  params
	store   Map
	engine  *Engine
}

// Status returns a Response with the given status.
func (c *Context) Status(status int) *Response {
	c.writer.WriteHeader(status)
	return &Response{
		context: c,
	}
}

// OK returns a Response with status 200.
func (c *Context) OK() *Response {
	c.writer.WriteHeader(http.StatusOK)
	return &Response{
		context: c,
	}
}

// BadRequest returns a Response with status 400.
func (c *Context) BadRequest() *Response {
	c.writer.WriteHeader(http.StatusBadRequest)
	return &Response{
		context: c,
	}
}

// Unauthorized returns a Response with status 401.
func (c *Context) Unauthorized() *Response {
	c.writer.WriteHeader(http.StatusUnauthorized)
	return &Response{
		context: c,
	}
}

// Forbidden returns a Response with status 403.
func (c *Context) Forbidden() *Response {
	c.writer.WriteHeader(http.StatusForbidden)
	return &Response{
		context: c,
	}
}

// NotFound returns a Response with status 404.
func (c *Context) NotFound() *Response {
	c.writer.WriteHeader(http.StatusNotFound)
	return &Response{
		context: c,
	}
}

// Conflict returns a Response with status 409.
func (c *Context) Conflict() *Response {
	c.writer.WriteHeader(http.StatusConflict)
	return &Response{
		context: c,
	}
}

// UnprocessableEntity returns a Response with status 422.
func (c *Context) UnprocessableEntity() *Response {
	c.writer.WriteHeader(http.StatusUnprocessableEntity)
	return &Response{
		context: c,
	}
}

// InternalServerError returns a Response with status 500.
func (c *Context) InternalServerError() *Response {
	c.writer.WriteHeader(http.StatusInternalServerError)
	return &Response{
		context: c,
	}
}

// Param gets a path parameter by the given name. An Empty
// string is returned if a parameter by the given name
// doesn't exist.
func (c *Context) Param(name string) string {
	return c.params.get(name)
}

// Query gets a query value by the given name. An empty
// string is returned if a value by the given name
// doesn't exist.
func (c *Context) Query(name string) string {
	return c.request.URL.Query().Get(name)
}

// Set sets a value in the Context data store.
func (c *Context) Set(key string, value interface{}) {
	c.store[key] = value
}

// Get gets a value from the Context data store.
func (c *Context) Get(key string) interface{} {
	return c.store[key]
}

// ParseJSON parses the request body into the given target.
func (c *Context) ParseJSON(target interface{}) error {
	return jsoniter.NewDecoder(c.request.Body).Decode(target)
}

// Host returns the requested host.
func (c *Context) Host() string {
	return c.request.Host
}

// RequestHeader returns the request header.
func (c *Context) RequestHeader() http.Header {
	return c.request.Header
}

// ResponseHeader returns the Response header.
func (c *Context) ResponseHeader() http.Header {
	return c.writer.Header()
}

// SetCookie adds a Set-Cookie header to response.
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.writer, cookie)
}

// Cookie returns the cookie by the given name provided in
// the request or http.ErrNoCookie if not found.
func (c *Context) Cookie(name string) (*http.Cookie, error) {
	return c.request.Cookie(name)
}

// LogDebug logs the given messages for all loggers where
// ShouldLog(LogLevelDebug) method returns true.
func (c *Context) LogDebug(message interface{}) {
	for _, l := range c.engine.loggers {
		if l.ShouldLog(LogLevelDebug) {
			l.Log(LogLevelDebug, message)
		}
	}
}

// LogInfo logs the given messages for all loggers where
// ShouldLog(LogLevelInfo) method returns true.
func (c *Context) LogInfo(message interface{}) {
	for _, l := range c.engine.loggers {
		if l.ShouldLog(LogLevelInfo) {
			l.Log(LogLevelInfo, message)
		}
	}
}

// LogNotice logs the given messages for all loggers where
// ShouldLog(LogLevelNotice) method returns true.
func (c *Context) LogNotice(message interface{}) {
	for _, l := range c.engine.loggers {
		if l.ShouldLog(LogLevelNotice) {
			l.Log(LogLevelNotice, message)
		}
	}
}

// LogWarning logs the given messages for all loggers where
// ShouldLog(LogLevelWarning) method returns true.
func (c *Context) LogWarning(message interface{}) {
	for _, l := range c.engine.loggers {
		if l.ShouldLog(LogLevelWarning) {
			l.Log(LogLevelWarning, message)
		}
	}
}

// LogError logs the given messages for all loggers where
// ShouldLog(LogLevelError) method returns true.
func (c *Context) LogError(message interface{}) {
	for _, l := range c.engine.loggers {
		if l.ShouldLog(LogLevelError) {
			l.Log(LogLevelError, message)
		}
	}
}

// LogAlert logs the given messages for all loggers where
// ShouldLog(LogLevelAlert) method returns true.
func (c *Context) LogAlert(message interface{}) {
	for _, l := range c.engine.loggers {
		if l.ShouldLog(LogLevelAlert) {
			l.Log(LogLevelAlert, message)
		}
	}
}

// LogCritical logs the given messages for all loggers where
// ShouldLog(LogLevelCritical) method returns true.
func (c *Context) LogCritical(message interface{}) {
	for _, l := range c.engine.loggers {
		if l.ShouldLog(LogLevelCritical) {
			l.Log(LogLevelCritical, message)
		}
	}
}

// LogEmergency logs the given messages for all loggers where
// ShouldLog(LogLevelEmergency) method returns true.
func (c *Context) LogEmergency(message interface{}) {
	for _, l := range c.engine.loggers {
		if l.ShouldLog(LogLevelEmergency) {
			l.Log(LogLevelEmergency, message)
		}
	}
}
