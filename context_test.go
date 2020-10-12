package goweb_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/twharmon/goweb"
)

func TestOKEmpty(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusOK, "")
}

func TestStatus(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Status(http.StatusTeapot).Empty()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusTeapot, "")
}

func TestCreated(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Created().Empty()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusCreated, "")
}

func TestBadRequest(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.BadRequest().Empty()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusBadRequest, "")
}

func TestUnauthorized(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Unauthorized().Empty()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusUnauthorized, "")
}

func TestForbidden(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Forbidden().Empty()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusForbidden, "")
}

func TestNotFound(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.NotFound().Empty()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusNotFound, "")
}

func TestConflict(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Conflict().Empty()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusConflict, "")
}

func TestUnprocessableEntity(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.UnprocessableEntity().Empty()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusUnprocessableEntity, "")
}

func TestInternalServerError(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.InternalServerError().Empty()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusInternalServerError, "")
}

func TestQuery(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Text(c.Query("foo"))
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/?foo=bar", nil, nil, http.StatusOK, "bar")
}

func TestJSON(t *testing.T) {
	type Msg struct {
		Hello string `json:"hello"`
	}
	handler := func(c *goweb.Context) goweb.Responder {
		return c.JSON(&Msg{
			Hello: "world",
		})
	}
	app := goweb.New()
	app.GET("/", handler)
	resBody := "{\"hello\":\"world\"}\n"
	assert(t, app, "GET", "/", nil, nil, http.StatusOK, resBody)
}

func TestParseJSON(t *testing.T) {
	type Msg struct {
		Hello string `json:"hello"`
	}
	body := "{\"hello\":\"world\"}\n"
	handler := func(c *goweb.Context) goweb.Responder {
		var msg Msg
		if err := c.ParseJSON(&msg); err != nil {
			return c.BadRequest().Empty()
		}
		return c.JSON(&msg)
	}
	app := goweb.New()
	app.POST("/", handler)
	assert(t, app, "POST", "/", strings.NewReader(body), nil, http.StatusOK, body)
}

func TestSetCookie(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		c.SetCookie(&http.Cookie{
			Name:  "foo",
			Value: "bar",
		})
		return c.Empty()
	}
	app := goweb.New()
	app.GET("/", handler)
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", "foo=bar")
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)
	got := rr.Header().Get("Set-Cookie")
	if got != "foo=bar" {
		t.Errorf("handler returned unexpected cookie header: got '%v' want '%v'", got, "foo=bar")
	}
}

func TestLogDebug(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogDebug(logMsg)
		return c.Empty()
	}
	l := newLogger()
	app := goweb.New()
	app.RegisterLogger(l)
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", l, logMsg)
}

func TestLogInfo(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogInfo(logMsg)
		return c.Empty()
	}
	l := newLogger()
	app := goweb.New()
	app.RegisterLogger(l)
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", l, logMsg)
}

func TestLogNotice(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogNotice(logMsg)
		return c.Empty()
	}
	l := newLogger()
	app := goweb.New()
	app.RegisterLogger(l)
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", l, logMsg)
}

func TestLogWarning(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogWarning(logMsg)
		return c.Empty()
	}
	l := newLogger()
	app := goweb.New()
	app.RegisterLogger(l)
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", l, logMsg)
}

func TestLogError(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogError(logMsg)
		return c.Empty()
	}
	l := newLogger()
	app := goweb.New()
	app.RegisterLogger(l)
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", l, logMsg)
}

func TestLogAlert(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogAlert(logMsg)
		return c.Empty()
	}
	l := newLogger()
	app := goweb.New()
	app.RegisterLogger(l)
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", l, logMsg)
}

func TestLogCritical(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogCritical(logMsg)
		return c.Empty()
	}
	l := newLogger()
	app := goweb.New()
	app.RegisterLogger(l)
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", l, logMsg)
}

func TestLogEmergency(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogEmergency(logMsg)
		return c.Empty()
	}
	l := newLogger()
	app := goweb.New()
	app.RegisterLogger(l)
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", l, logMsg)
}
