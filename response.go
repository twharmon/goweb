package goweb

import (
	"compress/gzip"
	"errors"
	"fmt"
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
	if r.gzip == 0 {
		http.ServeFile(r.context.ResponseWriter, r.context.Request, r.path)
		return
	}

	if !strings.Contains(r.context.Request.Header.Get(acceptEncodingHeader), contentEncodingGzip) {
		http.ServeFile(r.context.ResponseWriter, r.context.Request, r.path)
		return
	}

	fi, err := os.Stat(r.path)
	if err != nil {
		r.context.LogError(fmt.Errorf("unable to stat file: %w", err))
		http.ServeFile(r.context.ResponseWriter, r.context.Request, r.path)
		return
	}
	if fi.Size() < r.gzip {
		http.ServeFile(r.context.ResponseWriter, r.context.Request, r.path)
		return
	}

	gzipPath := r.path + ".gz"
	if err := r.gzipIfNeeded(r.path, gzipPath, fi); err != nil {
		r.context.LogError(fmt.Errorf("gzip error: %w", err))
		http.ServeFile(r.context.ResponseWriter, r.context.Request, r.path)
		return
	}

	r.context.ResponseWriter.Header().Set(contentEncodingHeader, contentEncodingGzip)
	r.context.ResponseWriter.Header().Set(contentTypeHeader, contentTypeWasm)
	http.ServeFile(r.context.ResponseWriter, r.context.Request, gzipPath)
}

func (r *FileResponse) gzipIfNeeded(path string, gzipPath string, fi os.FileInfo) error {
	if gfi, err := os.Stat(gzipPath); errors.Is(err, os.ErrNotExist) || gfi.ModTime().Before(fi.ModTime()) {
		f, err := os.Open(r.path)
		if err != nil {
			return fmt.Errorf("unable to open file: %w", err)
		}
		defer f.Close()

		fgz, err := os.Create(gzipPath)
		if err != nil {
			return fmt.Errorf("unable to create gzip file: %w", err)
		}
		defer fgz.Close()

		gw := gzip.NewWriter(fgz)
		_, err = io.Copy(gw, f)
		if err != nil {
			return fmt.Errorf("unable to write to file: %w", err)
		}
		err = gw.Close()
		if err != nil {
			return fmt.Errorf("unable to close gzipper: %w", err)
		}
	}
	return nil
}

// Respond sends a plain text response.
func (r *PlainTextResponse) Respond() {
	r.context.ResponseWriter.Header().Set(contentTypeHeader, contentTypeTextPlain)
	r.context.ResponseWriter.WriteHeader(r.context.status)
	r.context.ResponseWriter.Write([]byte(r.body))
}
