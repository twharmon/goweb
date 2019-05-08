package goweb_test

import (
	"net/http"
	"testing"

	"goweb"
)

func TestPassThroughMiddleware(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.OK().PlainText(c.Get("foo").(string))
	}
	app := goweb.New()
	mw := goweb.NewMiddleware()
	mw.Use(func(c *goweb.Context) *goweb.Response {
		c.Set("foo", "bar")
		return nil
	})
	app.GET("/", mw.Apply(handler))
	assert(t, app, "GET", "/", nil, nil, http.StatusOK, "bar")
}

func TestInterruptingMiddleware(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.OK()
	}
	app := goweb.New()
	mw := goweb.NewMiddleware()
	mw.Use(func(c *goweb.Context) *goweb.Response {
		return c.BadRequest()
	})
	app.GET("/", mw.Apply(handler))
	assert(t, app, "GET", "/", nil, nil, http.StatusBadRequest, "")
}
