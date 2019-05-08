package goweb_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"goweb"
)

func assert(t *testing.T, app *goweb.Engine, method string, path string, reqBody io.Reader, status int, resBody string) {
	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	app.ServeHTTP(rr, req)

	if rr.Code != status {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, status)
	}

	if rr.Body.String() != resBody {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'", rr.Body.String(), resBody)
	}
}

func assertLog(t *testing.T, app *goweb.Engine, method string, path string, logMsg string) {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	app.ServeHTTP(rr, req)

	if !strings.Contains(logOut.String(), logMsg) {
		t.Errorf("logged wrong message: got '%v' want '%v'", logOut.String(), logMsg)
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
