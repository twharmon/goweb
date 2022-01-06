package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/twharmon/goweb"
)

func main() {
	app := goweb.New()

	app.GET("/{path:.*}", func(c *goweb.Context) goweb.Responder {
		path := "assets/" + c.Param("path")
		if path == "assets/" {
			path = "assets/index.html"
		}
		return newFileResponse(c, strings.TrimRight(path, "/"))
	})

	app.Run("localhost:8080")
}

func newFileResponse(c *goweb.Context, path string) *fileResponse {
	return &fileResponse{
		context: c,
		path:    path,
	}
}

type fileResponse struct {
	context *goweb.Context
	path    string
}

func (r *fileResponse) Respond() {
	fi, err := os.Stat(r.path)
	if err != nil {
		r.context.LogError(fmt.Errorf("unable to stat file: %w", err))
		http.ServeFile(r.context.ResponseWriter, r.context.Request, r.path)
		return
	}
	if fi.IsDir() {
		r.path += "/index.html"
	}
	f, err := os.Open(r.path)
	if err != nil {
		r.context.LogError(fmt.Errorf("unable to open file: %w", err))
		http.ServeFile(r.context.ResponseWriter, r.context.Request, r.path)
		return
	}
	data := make([]byte, 512)
	if _, err := f.Read(data); err != nil {
		r.context.LogError(fmt.Errorf("file read error: %w", err))
		http.ServeFile(r.context.ResponseWriter, r.context.Request, r.path)
		return
	}
	r.context.ResponseWriter.Header().Set("Content-Type", http.DetectContentType(data))
	http.ServeFile(r.context.ResponseWriter, r.context.Request, r.path)
}
