package goweb_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"testing"

	"goweb"
)

func TestStatus(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.Status(-1)
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, -1, "")
}

func TestOK(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.OK()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, http.StatusOK, "")
}

func TestBadRequest(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.BadRequest()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, http.StatusBadRequest, "")
}

func TestUnauthorized(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.Unauthorized()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, http.StatusUnauthorized, "")
}

func TestForbidden(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.Forbidden()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, http.StatusForbidden, "")
}

func TestNotFound(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.NotFound()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, http.StatusNotFound, "")
}

func TestConflict(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.Conflict()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, http.StatusConflict, "")
}

func TestUnprocessableEntity(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.UnprocessableEntity()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, http.StatusUnprocessableEntity, "")
}

func TestInternalServerError(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.InternalServerError()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, http.StatusInternalServerError, "")
}

func TestQuery(t *testing.T) {
	handler := func(c *goweb.Context) *goweb.Response {
		return c.OK().PlainText(c.Query("foo"))
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/?foo=bar", nil, http.StatusOK, "bar")
}

type myLogger struct {
	level int
}

var logOut bytes.Buffer

func init() {
	log.SetOutput(&logOut)
}

func (l *myLogger) ShouldLog(level int) bool {
	return level >= l.level
}

func (l *myLogger) Log(level int, msg interface{}) {
	log.Print(msg)
}

func TestLog(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) *goweb.Response {
		c.LogInfo(logMsg)
		return c.OK()
	}
	fmt.Println("logOut", logOut.String())
	app := goweb.New()

	app.AddCustomLogger(&myLogger{goweb.LogLevelInfo})
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", logMsg)
}

func TestJSON(t *testing.T) {
	type Msg struct {
		Hello string `json:"hello"`
	}
	handler := func(c *goweb.Context) *goweb.Response {
		return c.OK().JSON(&Msg{
			Hello: "world",
		})
	}
	app := goweb.New()
	app.GET("/", handler)
	resBody := "{\"hello\":\"world\"}\n"
	assert(t, app, "GET", "/", nil, http.StatusOK, resBody)
}
