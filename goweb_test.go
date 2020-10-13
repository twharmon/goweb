package goweb_test

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/twharmon/goweb"
)

type testLogger struct {
	out *bytes.Buffer
	log *log.Logger
}

func (l *testLogger) Log(_ *goweb.Context, level goweb.LogLevel, msgs ...interface{}) {
	l.log.Print(msgs...)
}

var stdLoggerOut bytes.Buffer

func init() {
	log.SetOutput(&stdLoggerOut)
	log.SetFlags(0)
}

func newLogger() *testLogger {
	l := new(testLogger)
	l.out = bytes.NewBuffer(nil)
	l.log = log.New(l.out, "", 0)
	return l
}

func assert(t *testing.T, app *goweb.Engine, method string, path string, reqBody io.Reader, reqTransformer func(*http.Request), status int, resBody string) {
	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		t.Fatal(err)
	}
	if reqTransformer != nil {
		reqTransformer(req)
	}
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)
	if rr.Code != status {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, status)
	}
	if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(resBody) {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'", rr.Body.String(), resBody)
	}
}

func assertOK(t *testing.T, app *goweb.Engine, method string, path string, reqBody io.Reader, reqTransformer func(*http.Request), status int) {
	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		t.Fatal(err)
	}
	if reqTransformer != nil {
		reqTransformer(req)
	}
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)
	if rr.Code != status {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, status)
	}
}

func assertLog(t *testing.T, app *goweb.Engine, method string, path string, logger *testLogger, logMsg string) {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)
	got := strings.TrimSuffix(logger.out.String(), "\n")
	if got != logMsg {
		t.Errorf("logged wrong message: got '%v' want '%v'", got, logMsg)
	}
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("assert panic: (no panic)")
		}
	}()
	f()
}
