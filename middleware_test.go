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

func TestPUTMiddleware(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Text(c.Get("foo").(string))
	}
	app := goweb.New()
	mw := app.Middleware(func(c *goweb.Context) goweb.Responder {
		c.Set("foo", "bar")
		return nil
	})
	mw.PUT("/", handler)
	assert(t, app, "PUT", "/", nil, nil, http.StatusOK, "bar")
}

func TestPOSTMiddleware(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Text(c.Get("foo").(string))
	}
	app := goweb.New()
	mw := app.Middleware(func(c *goweb.Context) goweb.Responder {
		c.Set("foo", "bar")
		return nil
	})
	mw.POST("/", handler)
	assert(t, app, "POST", "/", nil, nil, http.StatusOK, "bar")
}

func TestPATCHMiddleware(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Text(c.Get("foo").(string))
	}
	app := goweb.New()
	mw := app.Middleware(func(c *goweb.Context) goweb.Responder {
		c.Set("foo", "bar")
		return nil
	})
	mw.PATCH("/", handler)
	assert(t, app, "PATCH", "/", nil, nil, http.StatusOK, "bar")
}

func TestDELETEMiddleware(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Text(c.Get("foo").(string))
	}
	app := goweb.New()
	mw := app.Middleware(func(c *goweb.Context) goweb.Responder {
		c.Set("foo", "bar")
		return nil
	})
	mw.DELETE("/", handler)
	assert(t, app, "DELETE", "/", nil, nil, http.StatusOK, "bar")
}

func TestOPTIONSMiddleware(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Text(c.Get("foo").(string))
	}
	app := goweb.New()
	mw := app.Middleware(func(c *goweb.Context) goweb.Responder {
		c.Set("foo", "bar")
		return nil
	})
	mw.OPTIONS("/", handler)
	assert(t, app, "OPTIONS", "/", nil, nil, http.StatusOK, "bar")
}

func TestHEADMiddleware(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Text(c.Get("foo").(string))
	}
	app := goweb.New()
	mw := app.Middleware(func(c *goweb.Context) goweb.Responder {
		c.Set("foo", "bar")
		return nil
	})
	mw.HEAD("/", handler)
	assert(t, app, "HEAD", "/", nil, nil, http.StatusOK, "bar")
}

func TestChainMiddleware(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Text(c.Get("foo").(string) + c.Get("bar").(string))
	}
	app := goweb.New()
	mw1 := app.Middleware(func(c *goweb.Context) goweb.Responder {
		c.Set("foo", "bar")
		return nil
	})
	mw2 := mw1.Middleware(func(c *goweb.Context) goweb.Responder {
		c.Set("bar", "baz")
		return nil
	})
	mw2.PUT("/", handler)
	assert(t, app, "PUT", "/", nil, nil, http.StatusOK, "barbaz")
}
