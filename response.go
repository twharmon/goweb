package goweb

import (
	"encoding/json"
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
}

// TextResponse implements Responder interface.
type TextResponse struct {
	context *Context
	body    string
}

// EmptyResponse implements Responder interface.
type EmptyResponse struct {
	context *Context
}

// Responder is the Responder interface that responds to
// HTTP requests.
type Responder interface {
	Respond()
}

// Respond sends a JSON response.
func (r *JSONResponse) Respond() {
	r.context.ResponseWriter.Header().Set(contentTypeHeader, contentTypeApplicationJSON)
	r.context.ResponseWriter.WriteHeader(r.context.status)
	json.NewEncoder(r.context.ResponseWriter).Encode(r.body)
}

// Respond sends a JSON response.
func (r *EmptyResponse) Respond() {
	r.context.ResponseWriter.WriteHeader(r.context.status)
}

// Respond sends a plain text response.
func (r *TextResponse) Respond() {
	r.context.ResponseWriter.Header().Set(contentTypeHeader, contentTypeTextPlain)
	r.context.ResponseWriter.WriteHeader(r.context.status)
	r.context.ResponseWriter.Write([]byte(r.body))
}
