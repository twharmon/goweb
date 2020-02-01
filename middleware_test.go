package goweb_test

import (
	"net/http"
	"testing"

	"github.com/twharmon/goweb"
)

func TestPassThroughMiddleware(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Text(c.Get("foo").(string))
	}
	app := goweb.New()
	mw := app.Middleware(func(c *goweb.Context) goweb.Responder {
		c.Set("foo", "bar")
		return nil
	})
	mw.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusOK, "bar")
}

func TestInterruptingMiddleware(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty()
	}
	app := goweb.New()
	mw := app.Middleware(func(c *goweb.Context) goweb.Responder {
		return c.BadRequest().Empty()
	})
	mw.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusBadRequest, "")
}
