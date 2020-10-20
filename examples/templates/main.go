package main

import (
	"html/template"
	"net/http"

	"github.com/twharmon/goweb"
)

func main() {
	app := goweb.New()

	app.GET("/", func(c *goweb.Context) goweb.Responder {
		return newTemplateResponse(c, "index.html", goweb.Map{
			"title": "Hello world!",
			"body":  "Lorem ipsum",
		})
	})

	app.Run(":8080")
}

func newTemplateResponse(c *goweb.Context, path string, data goweb.Map) *templateResponse {
	return &templateResponse{
		context: c,
		path:    path,
		data:    data,
	}
}

type templateResponse struct {
	context *goweb.Context
	path    string
	data    goweb.Map
}

func (r *templateResponse) Respond() {
	t, err := template.New(r.path).ParseFiles("html/" + r.path)
	if err != nil {
		r.context.Status(http.StatusInternalServerError).Text("Unable to parse templates: " + err.Error()).Respond()
		return
	}
	if err := t.Execute(r.context.ResponseWriter, r.data); err != nil {
		r.context.Status(http.StatusInternalServerError).Text("Unable to execute templates: " + err.Error()).Respond()
	}
}
