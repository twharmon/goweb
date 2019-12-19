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
	r.context.ResponseWriter.Header().Set(contentTypeHeader, contentTypeApplicationJSON)
	r.context.ResponseWriter.WriteHeader(r.context.status)
	jsoniter.NewEncoder(r.context.ResponseWriter).Encode(r.body)
}

// Respond sends a JSON response.
func (r *FileResponse) Respond() {
	http.ServeFile(r.context.ResponseWriter, r.context.Request, r.path)
}

// Respond sends a plain text response.
func (r *PlainTextResponse) Respond() {
	r.context.ResponseWriter.Header().Set(contentTypeHeader, contentTypeTextPlain)
	r.context.ResponseWriter.WriteHeader(r.context.status)
	r.context.ResponseWriter.Write([]byte(r.body))
}
