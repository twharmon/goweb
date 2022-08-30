package goweb_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/twharmon/goweb"
)

func TestRun(t *testing.T) {
	app := goweb.New()
	app.GET("/", func(c *goweb.Context) goweb.Responder {
		return c.Empty(http.StatusOK)
	})
	go app.Run(":9999")
	res, err := http.DefaultClient.Get("http://localhost:9999/")
	if err != nil {
		t.Error(err)
		return
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200; got %d", res.StatusCode)
	}
	app.Shutdown()
}

func TestGET(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty(http.StatusOK)
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusOK, "")
}

func TestPUT(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty(http.StatusOK)
	}
	app := goweb.New()
	app.PUT("/", handler)
	assert(t, app, "PUT", "/", nil, nil, http.StatusOK, "")
}

func TestPATCH(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty(http.StatusOK)
	}
	app := goweb.New()
	app.PATCH("/", handler)
	assert(t, app, "PATCH", "/", nil, nil, http.StatusOK, "")
}

func TestPOST(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty(http.StatusOK)
	}
	app := goweb.New()
	app.POST("/", handler)
	assert(t, app, "POST", "/", nil, nil, http.StatusOK, "")
}

func TestDELETE(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty(http.StatusOK)
	}
	app := goweb.New()
	app.DELETE("/", handler)
	assert(t, app, "DELETE", "/", nil, nil, http.StatusOK, "")
}

func TestHEAD(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty(http.StatusOK)
	}
	app := goweb.New()
	app.HEAD("/", handler)
	assert(t, app, "HEAD", "/", nil, nil, http.StatusOK, "")
}

func TestOPTIONS(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty(http.StatusOK)
	}
	app := goweb.New()
	app.OPTIONS("/", handler)
	assert(t, app, "OPTIONS", "/", nil, nil, http.StatusOK, "")
}

func TestPostParamRoute(t *testing.T) {
	wrongHandler := func(c *goweb.Context) goweb.Responder {
		return c.Empty(http.StatusBadRequest)
	}
	correctHandler := func(c *goweb.Context) goweb.Responder {
		return c.Empty(http.StatusOK)
	}
	app := goweb.New()
	app.GET("/a/{b}", wrongHandler)
	app.GET("/a/{b}/c", correctHandler)
	assert(t, app, "GET", "/a/b/c", nil, nil, http.StatusOK, "")
}

func TestMultiParamRoute(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty(http.StatusOK)
	}
	app := goweb.New()
	app.GET("/a/{b}/c/{d}", handler)
	assert(t, app, "GET", "/a/b/c/d", nil, nil, http.StatusOK, "")
}

func TestRouteNotFound(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty(http.StatusOK)
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/foo", nil, nil, http.StatusNotFound, "{\"message\":\"Page Not Found\"}")
}

func TestCustomNotFound(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty(http.StatusOK)
	}
	app := goweb.New()
	app.GET("/", handler)
	app.NotFound(func(c *goweb.Context) goweb.Responder {
		return c.Empty(http.StatusNotFound)
	})
	assert(t, app, "GET", "/foo", nil, nil, http.StatusNotFound, "")
}

func TestEmptyPath(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty(http.StatusOK)
	}
	app := goweb.New()
	assertPanic(t, func() {
		app.GET("", handler)
	})
}

func TestNonSlashLeadingPath(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty(http.StatusOK)
	}
	app := goweb.New()
	assertPanic(t, func() {
		app.GET("foo", handler)
	})
}

func TestParams(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		res := fmt.Sprintf("%s %s", c.Param("name"), c.Param("age"))
		return c.Text(http.StatusOK, res)
	}
	app := goweb.New()
	app.GET("/hello/{name}/{age}", handler)
	assert(t, app, "GET", "/hello/Gopher/5", nil, nil, http.StatusOK, "Gopher 5")
}

func TestParamsNotFound(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		res := fmt.Sprintf("%s %s", c.Param("name_"), c.Param("age"))
		return c.Text(http.StatusOK, res)
	}
	app := goweb.New()
	app.GET("/hello/{name}/{age}", handler)
	assert(t, app, "GET", "/hello/Gopher/5", nil, nil, http.StatusOK, " 5")
}
