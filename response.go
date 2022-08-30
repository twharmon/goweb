package goweb

import (
	"encoding/json"
	"net/http"
)

const (
	contentTypeHeader          = "Content-Type"
	contentTypeApplicationJSON = "application/json; charset=utf-8"
	contentTypeTextPlain       = "text/plain; charset=utf-8"
)

// JSONResponse implements Responder interface.
type JSONResponse struct {
	context *Context
	body    interface{}
	status  int
}

// TextResponse implements Responder interface.
type TextResponse struct {
	context *Context
	body    string
	status  int
}

// EmptyResponse implements Responder interface. It sends
// as response without a body.
type EmptyResponse struct {
	context *Context
	status  int
}

// NilResponse implements Responder interface. It does
// not send any response. Useful for handlers that upgrade
// to WebSockets.
type NilResponse struct {
}

// RedirectResponse implements Responder interface.
type RedirectResponse struct {
	context *Context
	url     string
	status  int
}

// Responder is the Responder interface that responds to
// HTTP requests.
type Responder interface {
	Respond()
}

// Respond sends a JSON response.
func (r *JSONResponse) Respond() {
	r.context.ResponseWriter.Header().Set(contentTypeHeader, contentTypeApplicationJSON)
	r.context.ResponseWriter.WriteHeader(r.status)
	json.NewEncoder(r.context.ResponseWriter).Encode(r.body)
}

// Respond sends a JSON response.
func (r *EmptyResponse) Respond() {
	r.context.ResponseWriter.WriteHeader(r.status)
}

// Respond does not do anything.
func (r *NilResponse) Respond() {}

// Respond sends a plain text response.
func (r *TextResponse) Respond() {
	r.context.ResponseWriter.Header().Set(contentTypeHeader, contentTypeTextPlain)
	r.context.ResponseWriter.WriteHeader(r.status)
	r.context.ResponseWriter.Write([]byte(r.body))
}

// Respond sends a plain text response.
func (r *RedirectResponse) Respond() {
	http.Redirect(r.context.ResponseWriter, r.context.Request, r.url, r.status)
}
