package goweb_test

import (
	"net/http"
	"testing"

	"github.com/twharmon/goweb"
)

type todoResource struct {
}

func (t todoResource) Put(c *goweb.Context) goweb.Responder {
	return c.Text(http.StatusOK, c.Param(t.Identifier()))
}

func (t todoResource) Delete(c *goweb.Context) goweb.Responder {
	return c.Text(http.StatusOK, c.Param(t.Identifier()))
}

func (t todoResource) Post(c *goweb.Context) goweb.Responder {
	return c.Text(http.StatusOK, "4")
}

func (t todoResource) Index(c *goweb.Context) goweb.Responder {
	return c.Text(http.StatusOK, "index")
}

func (t todoResource) Get(c *goweb.Context) goweb.Responder {
	return c.Text(http.StatusOK, c.Param(t.Identifier()))
}

func (t todoResource) Identifier() string {
	return "id"
}

func TestResource(t *testing.T) {
	app := goweb.New()

	todos := todoResource{}
	app.Resource("/todo", todos)

	assert(t, app, "GET", "/todo", nil, nil, http.StatusOK, "index")
	assert(t, app, "GET", "/todo/1", nil, nil, http.StatusOK, "1")
	assert(t, app, "PUT", "/todo/2", nil, nil, http.StatusOK, "2")
	assert(t, app, "DELETE", "/todo/3", nil, nil, http.StatusOK, "3")
	assert(t, app, "POST", "/todo", nil, nil, http.StatusOK, "4")
}

func TestResourceMiddleware(t *testing.T) {
	app := goweb.New()

	m := app.Middleware(func(c *goweb.Context) goweb.Responder {
		if c.Param("id") == "1" {
			return c.Text(http.StatusOK, "param was one")
		}
		return nil
	})
	m.Resource("/todo", todoResource{})

	assert(t, app, "GET", "/todo/1", nil, nil, http.StatusOK, "param was one")
}
