package goweb_test

import (
	"fmt"
	"net/http"
	"testing"

	"goweb"
)

func TestGET(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.OK()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, http.StatusOK, "")
}

func TestPUT(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.OK()
	}
	app := goweb.New()
	app.PUT("/", handler)
	assert(t, app, "PUT", "/", nil, http.StatusOK, "")
}

func TestPATCH(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.OK()
	}
	app := goweb.New()
	app.PATCH("/", handler)
	assert(t, app, "PATCH", "/", nil, http.StatusOK, "")
}

func TestPOST(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.OK()
	}
	app := goweb.New()
	app.POST("/", handler)
	assert(t, app, "POST", "/", nil, http.StatusOK, "")
}

func TestDELETE(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.OK()
	}
	app := goweb.New()
	app.DELETE("/", handler)
	assert(t, app, "DELETE", "/", nil, http.StatusOK, "")
}

func TestOPTIONS(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.OK()
	}
	app := goweb.New()
	app.OPTIONS("/", handler)
	assert(t, app, "OPTIONS", "/", nil, http.StatusOK, "")
}

func TestHEAD(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.OK()
	}
	app := goweb.New()
	app.HEAD("/", handler)
	assert(t, app, "HEAD", "/", nil, http.StatusOK, "")
}

func TestRouteNotFound(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.OK()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/foo", nil, http.StatusNotFound, "Page not found")
}

func TestCustomNotFound(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.OK()
	}
	app := goweb.New()
	app.GET("/", handler)
	app.NotFound(func(c *goweb.Context) *goweb.Response {
		return c.NotFound()
	})
	assert(t, app, "GET", "/foo", nil, http.StatusNotFound, "")
}

func TestEmptyPath(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.OK()
	}
	app := goweb.New()
	assertPanic(t, func() {
		app.GET("", handler)
	})
}

func TestNonSlashLeadingPath(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.OK()
	}
	app := goweb.New()
	assertPanic(t, func() {
		app.GET("foo", handler)
	})
}

func TestParams(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		res := fmt.Sprintf("%s %s", c.Param("name"), c.Param("age"))
		return c.OK().PlainText(res)
	}
	app := goweb.New()
	app.GET("/hello/{name}/{age}", handler)
	assert(t, app, "GET", "/hello/Gopher/5", nil, http.StatusOK, "Gopher 5")
}
