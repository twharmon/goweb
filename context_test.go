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
		return c.Empty(http.StatusOK)
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusOK, "")
}

func TestNil(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Nil()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusOK, "")
}

func TestQuery(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Text(http.StatusOK, c.Query("foo"))
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
		return c.JSON(http.StatusOK, &Msg{
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
			return c.Empty(http.StatusBadRequest)
		}
		return c.JSON(http.StatusOK, &msg)
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
		return c.Empty(http.StatusOK)
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
		return c.Empty(http.StatusOK)
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
		return c.Empty(http.StatusOK)
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
		return c.Empty(http.StatusOK)
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
		return c.Empty(http.StatusOK)
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
		return c.Empty(http.StatusOK)
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
		return c.Empty(http.StatusOK)
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
		return c.Empty(http.StatusOK)
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
		return c.Empty(http.StatusOK)
	}
	l := newLogger()
	app := goweb.New()
	app.RegisterLogger(l)
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", l, logMsg)
}

func TestRedirect(t *testing.T) {
	app := goweb.New()
	app.GET("/", func(c *goweb.Context) goweb.Responder {
		return c.Redirect(http.StatusTemporaryRedirect, "/foo")
	})
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)
	if rr.Result().StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("handler returned unexpected status code: got '%v' want '%v'", rr.Result().StatusCode, http.StatusTemporaryRedirect)
	}
	if rr.Result().Header.Get("Location") != "/foo" {
		t.Errorf("handler returned unexpected location header: got '%v' want '%v'", rr.Result().Header.Get("Location"), "/foo")
	}
}
