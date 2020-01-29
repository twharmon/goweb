package goweb_test

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"golang.org/x/net/websocket"

	"github.com/twharmon/goweb"
)

func TestGET(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/", nil, nil, http.StatusOK, "")
}

func TestWS(t *testing.T) {
	app := goweb.New()
	app.WebSocket("/", func(c *websocket.Conn) {
		io.Copy(c, c)
	})
	go func() {
		conn, err := websocket.Dial("ws://localhost:8080/", "", "http://localhost/")
		if err != nil {
			t.Error(err)
		}
		want := []byte("hi")
		conn.Write(want)
		var got = make([]byte, 512)
		_, err = conn.Read(got)
		if err != nil {
			t.Error(err)
		}
		if !strings.HasPrefix(string(got), string(want)) {
			t.Errorf("expected %s to equal %s", string(got), string(want))
		}
		if err := app.Shutdown(); err != nil {
			t.Error(err)
		}
	}()
	if err := app.Run(":8080"); err != http.ErrServerClosed {
		t.Errorf("expected %v to equal %v", err, http.ErrServerClosed)
	}
}

func TestPUT(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty()
	}
	app := goweb.New()
	app.PUT("/", handler)
	assert(t, app, "PUT", "/", nil, nil, http.StatusOK, "")
}

func TestPATCH(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty()
	}
	app := goweb.New()
	app.PATCH("/", handler)
	assert(t, app, "PATCH", "/", nil, nil, http.StatusOK, "")
}

func TestPOST(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty()
	}
	app := goweb.New()
	app.POST("/", handler)
	assert(t, app, "POST", "/", nil, nil, http.StatusOK, "")
}

func TestDELETE(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty()
	}
	app := goweb.New()
	app.DELETE("/", handler)
	assert(t, app, "DELETE", "/", nil, nil, http.StatusOK, "")
}

func TestOPTIONS(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty()
	}
	app := goweb.New()
	app.OPTIONS("/", handler)
	assert(t, app, "OPTIONS", "/", nil, nil, http.StatusOK, "")
}

func TestHEAD(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty()
	}
	app := goweb.New()
	app.HEAD("/", handler)
	assert(t, app, "HEAD", "/", nil, nil, http.StatusOK, "")
}

func TestPostParamRoute(t *testing.T) {
	wrongHandler := func(c *goweb.Context) goweb.Responder {
		return c.BadRequest().Empty()
	}
	correctHandler := func(c *goweb.Context) goweb.Responder {
		return c.Empty()
	}
	app := goweb.New()
	app.GET("/a/{b}", wrongHandler)
	app.GET("/a/{b}/c", correctHandler)
	assert(t, app, "GET", "/a/b/c", nil, nil, http.StatusOK, "")
}

func TestMultiParamRoute(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty()
	}
	app := goweb.New()
	app.GET("/a/{b}/c/{d}", handler)
	assert(t, app, "GET", "/a/b/c/d", nil, nil, http.StatusOK, "")
}

func TestServeFiles(t *testing.T) {
	content := "test file contents"
	data := []byte(content)
	if err := ioutil.WriteFile("./test.txt", data, 0700); err != nil {
		t.Error(err)
	}
	app := goweb.New()
	app.ServeFiles("/", ".")
	assert(t, app, "GET", "/test.txt", nil, nil, http.StatusOK, content)
	os.Remove("./test.txt")
}

func TestServeFilesIndex(t *testing.T) {
	content := "<html>hello world</html>"
	data := []byte(content)
	if err := ioutil.WriteFile("./index.html", data, 0700); err != nil {
		t.Error(err)
	}
	app := goweb.New()
	app.ServeFiles("/", ".")
	assert(t, app, "GET", "/", nil, nil, http.StatusOK, content)
	os.Remove("./index.html")
}

func TestGzipAndServeFiles(t *testing.T) {
	content := "test file contents"
	data := []byte(content)
	if err := ioutil.WriteFile("./test.txt", data, 0700); err != nil {
		t.Error(err)
	}
	app := goweb.New()
	app.GzipAndServeFiles("/", ".", 10)
	assertOK(t, app, "GET", "/test.txt", nil, nil, http.StatusOK)
	os.Remove("./test.txt")
	os.Remove("./test.txt.gz")
}

func TestGzipAndServeFilesIndex(t *testing.T) {
	content := "<html>hello world</html>"
	data := []byte(content)
	if err := ioutil.WriteFile("./index.html", data, 0700); err != nil {
		t.Error(err)
	}
	app := goweb.New()
	app.GzipAndServeFiles("/", ".", 10)
	transformer := func(r *http.Request) {
		r.Header.Set("Accept-Encoding", "gzip")
	}
	assertOK(t, app, "GET", "/", nil, transformer, http.StatusOK)
	os.Remove("./index.html")
	os.Remove("./index.html.gz")
}

func TestRouteNotFound(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty()
	}
	app := goweb.New()
	app.GET("/", handler)
	assert(t, app, "GET", "/foo", nil, nil, http.StatusNotFound, "Page not found")
}

func TestCustomNotFound(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty()
	}
	app := goweb.New()
	app.GET("/", handler)
	app.NotFound(func(c *goweb.Context) goweb.Responder {
		return c.NotFound().Empty()
	})
	assert(t, app, "GET", "/foo", nil, nil, http.StatusNotFound, "")
}

func TestEmptyPath(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty()
	}
	app := goweb.New()
	assertPanic(t, func() {
		app.GET("", handler)
	})
}

func TestNonSlashLeadingPath(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty()
	}
	app := goweb.New()
	assertPanic(t, func() {
		app.GET("foo", handler)
	})
}

func TestParams(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		res := fmt.Sprintf("%s %s", c.Param("name"), c.Param("age"))
		return c.Text(res)
	}
	app := goweb.New()
	app.GET("/hello/{name}/{age}", handler)
	assert(t, app, "GET", "/hello/Gopher/5", nil, nil, http.StatusOK, "Gopher 5")
}

func TestParamsNotFound(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		res := fmt.Sprintf("%s %s", c.Param("name_"), c.Param("age"))
		return c.Text(res)
	}
	app := goweb.New()
	app.GET("/hello/{name}/{age}", handler)
	assert(t, app, "GET", "/hello/Gopher/5", nil, nil, http.StatusOK, " 5")
}

func TestNoWWW(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty()
	}
	app := goweb.New()
	app.RedirectWWW()
	app.GET("/", handler)
	transformer := func(r *http.Request) {
		r.Host = "www.example.com"
	}
	assert(t, app, "GET", "/", nil, transformer, http.StatusMovedPermanently, "<a href=\"http://example.com/\">Moved Permanently</a>.\n\n")
}

func TestRedirectWWWTLS(t *testing.T) {
	handler := func(c *goweb.Context) goweb.Responder {
		return c.Empty()
	}
	app := goweb.New()
	app.RedirectWWW()
	app.GET("/", handler)
	transformer := func(r *http.Request) {
		r.Host = "www.example.com"
		r.TLS = &tls.ConnectionState{}
	}
	assert(t, app, "GET", "/", nil, transformer, http.StatusMovedPermanently, "<a href=\"https://example.com/\">Moved Permanently</a>.\n\n")
}
