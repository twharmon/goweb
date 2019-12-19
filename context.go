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
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	params         params
	store          Map
	engine         *Engine
	status         int
}

// Status sets the response status.
func (c *Context) Status(status int) *Context {
	c.status = status
	return c
}

// OK sets the response status to 200.
func (c *Context) OK() *Context {
	c.status = http.StatusOK
	return c
}

// BadRequest sets the response status to 400.
func (c *Context) BadRequest() *Context {
	c.status = http.StatusBadRequest
	return c
}

// Unauthorized sets the response status to 401.
func (c *Context) Unauthorized() *Context {
	c.status = http.StatusUnauthorized
	return c
}

// Forbidden sets the response status to 403.
func (c *Context) Forbidden() *Context {
	c.status = http.StatusForbidden
	return c
}

// NotFound sets the response status to 404.
func (c *Context) NotFound() *Context {
	c.status = http.StatusNotFound
	return c
}

// Conflict sets the response status to 409.
func (c *Context) Conflict() *Context {
	c.status = http.StatusConflict
	return c
}

// UnprocessableEntity sets the response status to 422.
func (c *Context) UnprocessableEntity() *Context {
	c.status = http.StatusUnprocessableEntity
	return c
}

// InternalServerError sets the response status to 500.
func (c *Context) InternalServerError() *Context {
	c.status = http.StatusInternalServerError
	return c
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
	return c.Request.URL.Query().Get(name)
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
	return jsoniter.NewDecoder(c.Request.Body).Decode(target)
}

// SetCookie adds a Set-Cookie header to response.
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.ResponseWriter, cookie)
}

// LogDebug logs the given messages for all loggers where
// ShouldLog(LogLevelDebug) method returns true.
func (c *Context) LogDebug(message interface{}) {
	for _, l := range c.engine.loggers {
		l.Log(LogLevelDebug, message)
	}
}

// LogInfo logs the given messages for all loggers where
// ShouldLog(LogLevelInfo) method returns true.
func (c *Context) LogInfo(message interface{}) {
	for _, l := range c.engine.loggers {
		l.Log(LogLevelInfo, message)
	}
}

// LogNotice logs the given messages for all loggers where
// ShouldLog(LogLevelNotice) method returns true.
func (c *Context) LogNotice(message interface{}) {
	for _, l := range c.engine.loggers {
		l.Log(LogLevelNotice, message)
	}
}

// LogWarning logs the given messages for all loggers where
// ShouldLog(LogLevelWarning) method returns true.
func (c *Context) LogWarning(message interface{}) {
	for _, l := range c.engine.loggers {
		l.Log(LogLevelWarning, message)
	}
}

// LogError logs the given messages for all loggers where
// ShouldLog(LogLevelError) method returns true.
func (c *Context) LogError(message interface{}) {
	for _, l := range c.engine.loggers {
		l.Log(LogLevelError, message)
	}
}

// LogAlert logs the given messages for all loggers where
// ShouldLog(LogLevelAlert) method returns true.
func (c *Context) LogAlert(message interface{}) {
	for _, l := range c.engine.loggers {
		l.Log(LogLevelAlert, message)
	}
}

// LogCritical logs the given messages for all loggers where
// ShouldLog(LogLevelCritical) method returns true.
func (c *Context) LogCritical(message interface{}) {
	for _, l := range c.engine.loggers {
		l.Log(LogLevelCritical, message)
	}
}

// LogEmergency logs the given messages for all loggers where
// ShouldLog(LogLevelEmergency) method returns true.
func (c *Context) LogEmergency(message interface{}) {
	for _, l := range c.engine.loggers {
		l.Log(LogLevelEmergency, message)
	}
}

// JSON returns a JSONResponse.
func (c *Context) JSON(value interface{}) *JSONResponse {
	return &JSONResponse{
		context: c,
		body:    value,
	}
}

// PlainText returns a PlainTextResponse.
func (c *Context) PlainText(text string) *PlainTextResponse {
	return &PlainTextResponse{
		context: c,
		body:    text,
	}
}

// File serves a file.
func (c *Context) File(path string) *FileResponse {
	return &FileResponse{
		context: c,
		path:    path,
	}
}

// Respond returns an empty response.
func (c *Context) Respond() {
	c.ResponseWriter.WriteHeader(c.status)
}
