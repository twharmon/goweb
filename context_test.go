package goweb_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/twharmon/goweb"
)

func TestStatus(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Status(-1)
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, -1, "")
}

func TestOK(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.OK()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusOK, "")
}

func TestBadRequest(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.BadRequest()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusBadRequest, "")
}

func TestUnauthorized(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Unauthorized()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusUnauthorized, "")
}

func TestForbidden(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Forbidden()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusForbidden, "")
}

func TestNotFound(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.NotFound()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusNotFound, "")
}

func TestConflict(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Conflict()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusConflict, "")
}

func TestUnprocessableEntity(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.UnprocessableEntity()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusUnprocessableEntity, "")
}

func TestInternalServerError(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.InternalServerError()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusInternalServerError, "")
}

func TestQuery(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.OK().PlainText(c.Query("foo"))
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
		return c.OK().JSON(&Msg{
			Hello: "world",
		})
	}
	app := goweb.New()
	app.GET("/", handler)
	resBody := "{\"hello\":\"world\"}\n"
	assert(t, app, "GET", "/", nil, nil, http.StatusOK, resBody)
}

func TestFile(t *testing.T) {
	content := "test file contents"
	data := []byte(content)
	if err := ioutil.WriteFile("./test.txt", data, 0700); err != nil {
		t.Error(err)
	}
	app := goweb.New()
	app.GET("/", func(c *goweb.Context) goweb.Responder {
		return c.File("test.txt")
	})
	assert(t, app, "GET", "/", nil, nil, http.StatusOK, content)
	os.Remove("./test.txt")
}

func TestWrite(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		c.Write([]byte("Hello, world!"))
		return c.OK()
	}
	app := goweb.New()
	app.GET("/", handler)
	resBody := "Hello, world!"
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
			return c.BadRequest()
		}
		return c.OK().JSON(&msg)
	}
	app := goweb.New()
	app.POST("/", handler)
	assert(t, app, "POST", "/", strings.NewReader(body), nil, http.StatusOK, body)
}

func TestHost(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.OK().PlainText(c.Host())
	}
	app := goweb.New()
	app.GET("/", handler)
	host := "example.com"
	transformer := func(r *http.Request) {
		r.Host = host
	}
	assert(t, app, "GET", "/", nil, transformer, http.StatusOK, host)
}

func TestRequestHeader(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.OK().PlainText(c.RequestHeader().Get("foo"))
	}
	app := goweb.New()
	app.GET("/", handler)
	transformer := func(r *http.Request) {
		r.Header.Set("foo", "bar")
	}
	assert(t, app, "GET", "/", nil, transformer, http.StatusOK, "bar")
}

func TestResponseHeader(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		c.ResponseHeader().Set("foo", "bar")
		return c.OK().PlainText(c.ResponseHeader().Get("foo"))
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusOK, "bar")
}

func TestCookie(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		cookie, err := c.Cookie("foo")
		if err != nil {
			return c.BadRequest()
		}
		return c.OK().PlainText(cookie.Value)
	}
	app := goweb.New()
	app.GET("/", handler)
	transformer := func(r *http.Request) {
		r.Header.Set("Cookie", "foo=bar")
	}
	assert(t, app, "GET", "/", nil, transformer, http.StatusOK, "bar")
}

func TestSetCookie(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		c.SetCookie(&http.Cookie{
			Name:  "foo",
			Value: "bar",
		})
		return c.OK()
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
		return c.OK()
	}
	app := goweb.New()
	l := newLogger(goweb.LogLevelDebug)
	app.AddCustomLogger(l)
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", l, logMsg)
}

func TestLogInfo(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogInfo(logMsg)
		return c.OK()
	}
	app := goweb.New()
	l := newLogger(goweb.LogLevelInfo)
	app.AddCustomLogger(l)
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", l, logMsg)
}

func TestLogNotice(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogNotice(logMsg)
		return c.OK()
	}
	app := goweb.New()
	l := newLogger(goweb.LogLevelNotice)
	app.AddCustomLogger(l)
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", l, logMsg)
}

func TestLogWarning(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogWarning(logMsg)
		return c.OK()
	}
	app := goweb.New()
	l := newLogger(goweb.LogLevelWarning)
	app.AddCustomLogger(l)
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", l, logMsg)
}

func TestLogError(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogError(logMsg)
		return c.OK()
	}
	app := goweb.New()
	l := newLogger(goweb.LogLevelError)
	app.AddCustomLogger(l)
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", l, logMsg)
}

func TestLogAlert(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogAlert(logMsg)
		return c.OK()
	}
	app := goweb.New()
	l := newLogger(goweb.LogLevelAlert)
	app.AddCustomLogger(l)
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", l, logMsg)
}

func TestLogCritical(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogCritical(logMsg)
		return c.OK()
	}
	app := goweb.New()
	l := newLogger(goweb.LogLevelCritical)
	app.AddCustomLogger(l)
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", l, logMsg)
}

func TestLogEmergency(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogEmergency(logMsg)
		return c.OK()
	}
	app := goweb.New()
	l := newLogger(goweb.LogLevelEmergency)
	app.AddCustomLogger(l)
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", l, logMsg)
}

func TestLogPassthrough(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogDebug("this should not get logged")
		return c.OK()
	}
	app := goweb.New()
	l := newLogger(goweb.LogLevelInfo)
	app.AddCustomLogger(l)
	app.GET("/", handler)
	assertLog(t, app, "GET", "/", l, "")
}

func TestStdLoggerDebug(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogDebug(logMsg)
		return c.OK()
	}
	app := goweb.New()
	app.AddStdLogger(goweb.LogLevelDebug)
	app.GET("/", handler)
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	app.ServeHTTP(httptest.NewRecorder(), req)
	got := stdLoggerOut.String()
	if !strings.Contains(got, logMsg) {
		t.Errorf("logged wrong message: got '%v' want '%v'", got, logMsg)
	}
}

func TestStdLoggerInfo(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogInfo(logMsg)
		return c.OK()
	}
	app := goweb.New()
	app.AddStdLogger(goweb.LogLevelInfo)
	app.GET("/", handler)
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	app.ServeHTTP(httptest.NewRecorder(), req)
	got := stdLoggerOut.String()
	if !strings.Contains(got, logMsg) {
		t.Errorf("logged wrong message: got '%v' want '%v'", got, logMsg)
	}
}

func TestStdLoggerNotice(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogNotice(logMsg)
		return c.OK()
	}
	app := goweb.New()
	app.AddStdLogger(goweb.LogLevelNotice)
	app.GET("/", handler)
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	app.ServeHTTP(httptest.NewRecorder(), req)
	got := stdLoggerOut.String()
	if !strings.Contains(got, logMsg) {
		t.Errorf("logged wrong message: got '%v' want '%v'", got, logMsg)
	}
}

func TestStdLoggerWarning(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogWarning(logMsg)
		return c.OK()
	}
	app := goweb.New()
	app.AddStdLogger(goweb.LogLevelWarning)
	app.GET("/", handler)
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	app.ServeHTTP(httptest.NewRecorder(), req)
	got := stdLoggerOut.String()
	if !strings.Contains(got, logMsg) {
		t.Errorf("logged wrong message: got '%v' want '%v'", got, logMsg)
	}
}

func TestStdLoggerError(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogError(logMsg)
		return c.OK()
	}
	app := goweb.New()
	app.AddStdLogger(goweb.LogLevelError)
	app.GET("/", handler)
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	app.ServeHTTP(httptest.NewRecorder(), req)
	got := stdLoggerOut.String()
	if !strings.Contains(got, logMsg) {
		t.Errorf("logged wrong message: got '%v' want '%v'", got, logMsg)
	}
}

func TestStdLoggerAlert(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogAlert(logMsg)
		return c.OK()
	}
	app := goweb.New()
	app.AddStdLogger(goweb.LogLevelAlert)
	app.GET("/", handler)
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	app.ServeHTTP(httptest.NewRecorder(), req)
	got := stdLoggerOut.String()
	if !strings.Contains(got, logMsg) {
		t.Errorf("logged wrong message: got '%v' want '%v'", got, logMsg)
	}
}

func TestStdLoggerCritical(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogCritical(logMsg)
		return c.OK()
	}
	app := goweb.New()
	app.AddStdLogger(goweb.LogLevelCritical)
	app.GET("/", handler)
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	app.ServeHTTP(httptest.NewRecorder(), req)
	got := stdLoggerOut.String()
	if !strings.Contains(got, logMsg) {
		t.Errorf("logged wrong message: got '%v' want '%v'", got, logMsg)
	}
}

func TestStdLoggerEmergency(t *testing.T) {
	logMsg := "test log message"
	handler := func(c *goweb.Context) goweb.Responder {
		c.LogEmergency(logMsg)
		return c.OK()
	}
	app := goweb.New()
	app.AddStdLogger(goweb.LogLevelEmergency)
	app.GET("/", handler)
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	app.ServeHTTP(httptest.NewRecorder(), req)
	got := stdLoggerOut.String()
	if !strings.Contains(got, logMsg) {
		t.Errorf("logged wrong message: got '%v' want '%v'", got, logMsg)
	}
}
