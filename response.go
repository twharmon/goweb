package goweb

import (
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

const (
	acceptEncodingHeader       = "Accept-Encoding"
	contentEncodingHeader      = "Content-Encoding"
	contentEncodingGzip        = "gzip"
	contentTypeHeader          = "Content-Type"
	contentTypeApplicationJSON = "application/json; charset=utf-8"
	contentTypeTextPlain       = "text/plain; charset=utf-8"
	contentTypeWasm            = "application/wasm"
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
	gzip    int64
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
	if r.gzip > 0 {
		if !strings.Contains(r.context.Request.Header.Get(acceptEncodingHeader), contentEncodingGzip) {
			http.ServeFile(r.context.ResponseWriter, r.context.Request, r.path)
			return
		}

		fi, err := os.Stat(r.path)
		if err != nil {
			r.context.LogError(err)
			http.ServeFile(r.context.ResponseWriter, r.context.Request, r.path)
			return
		}
		if fi.Size() < r.gzip {
			http.ServeFile(r.context.ResponseWriter, r.context.Request, r.path)
			return
		}

		gzipPath := r.path + ".gz"
		if _, err := os.Stat(gzipPath); errors.Is(err, os.ErrNotExist) {
			f, err := os.Open(r.path)
			if err != nil {
				r.context.LogError(err)
				http.ServeFile(r.context.ResponseWriter, r.context.Request, r.path)
				return
			}
			defer f.Close()

			fgz, err := os.Create(gzipPath)
			if err != nil {
				r.context.LogError(err)
				http.ServeFile(r.context.ResponseWriter, r.context.Request, r.path)
				return
			}
			defer fgz.Close()

			gw := gzip.NewWriter(fgz)
			_, err = io.Copy(gw, f)
			if err != nil {
				r.context.LogError(err)
				http.ServeFile(r.context.ResponseWriter, r.context.Request, r.path)
				return
			}
			err = gw.Close()
			if err != nil {
				r.context.LogError(err)
				http.ServeFile(r.context.ResponseWriter, r.context.Request, r.path)
				return
			}
		}
		r.context.ResponseWriter.Header().Set(contentEncodingHeader, contentEncodingGzip)
		r.context.ResponseWriter.Header().Set(contentTypeHeader, contentTypeWasm)
		http.ServeFile(r.context.ResponseWriter, r.context.Request, gzipPath)
	} else {
		http.ServeFile(r.context.ResponseWriter, r.context.Request, r.path)
	}
}

// Respond sends a plain text response.
func (r *PlainTextResponse) Respond() {
	r.context.ResponseWriter.Header().Set(contentTypeHeader, contentTypeTextPlain)
	r.context.ResponseWriter.WriteHeader(r.context.status)
	r.context.ResponseWriter.Write([]byte(r.body))
}
