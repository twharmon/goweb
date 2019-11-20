package goweb

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
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

// FileResponse implements Responder interface.
type FileResponse struct {
	context *Context
	path    string
}

// PlainTextResponse implements Responder interface.
type PlainTextResponse struct {
	context *Context
	body    string
}

// Responder is the Responder interface that responds to
// HTTP requests.
type Responder interface {
	Respond()
}

// Respond sends a JSON response.
func (r *JSONResponse) Respond() {
	r.context.writer.Header().Set(contentTypeHeader, contentTypeApplicationJSON)
	r.context.writer.WriteHeader(r.context.status)
	jsoniter.NewEncoder(r.context.writer).Encode(r.body)
}

// Respond sends a JSON response.
func (r *FileResponse) Respond() {
	http.ServeFile(r.context.writer, r.context.request, r.path)
}

// Respond sends a plain text response.
func (r *PlainTextResponse) Respond() {
	r.context.writer.Header().Set(contentTypeHeader, contentTypeTextPlain)
	r.context.writer.WriteHeader(r.context.status)
	r.context.writer.Write([]byte(r.body))
}
