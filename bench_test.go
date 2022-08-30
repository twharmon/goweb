package goweb_test

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-martini/martini"
	"github.com/labstack/echo/v4"

	"github.com/gorilla/mux"

	"github.com/gin-gonic/gin"
	"github.com/twharmon/goweb"
)

var plainTextBody = "Hello, World!"

var gowebApp *goweb.Engine
var ginApp *gin.Engine
var echoApp *echo.Echo
var gorillaApp *mux.Router
var martiniApp *martini.ClassicMartini

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

var posts []Post

func init() {
	for i := 0; i < 100; i++ {
		posts = append(posts, Post{
			ID:    1234,
			Title: "Lorem Ipsum",
			Body:  "Veniam ipsum officia consequat minim veniam cillum incididunt laborum aliqua ad do magna sed aliquip fugiat. Cillum et aliqua commodo, velit minim anim, pariatur, magna culpa officia dolor quis consectetur. Proident commodo laboris eu eu quis esse ea exercitation irure pariatur duis nulla deserunt dolor sed. Nulla qui laboris ut ea non consectetur amet culpa, pariatur commodo magna deserunt nostrud in.",
		})
	}

	gowebApp = goweb.New()
	gowebApp.GET("/plaintext", func(c *goweb.Context) goweb.Responder {
		return c.Text(http.StatusOK, plainTextBody)
	})
	gowebApp.GET("/json", func(c *goweb.Context) goweb.Responder {
		return c.JSON(http.StatusOK, &posts)
	})
	gowebApp.GET("/params/{foo}/{bar}/{baz}", func(c *goweb.Context) goweb.Responder {
		return c.JSON(http.StatusOK, goweb.Map{
			"foo": c.Param("foo"),
			"bar": c.Param("bar"),
			"baz": c.Param("baz"),
		})
	})

	gin.SetMode(gin.ReleaseMode)
	ginApp = gin.New()
	ginApp.GET("/plaintext", func(c *gin.Context) {
		c.String(http.StatusOK, plainTextBody)
	})
	ginApp.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, &posts)
	})
	ginApp.GET("/params/:foo/:bar/:baz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"foo": c.Param("foo"),
			"bar": c.Param("bar"),
			"baz": c.Param("baz"),
		})
	})

	gorillaApp = mux.NewRouter()
	gorillaApp.HandleFunc("/plaintext", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(plainTextBody))
	}).Methods("GET")
	gorillaApp.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(&posts)
	}).Methods("GET")
	gorillaApp.HandleFunc("/params/{foo}/{bar}/{baz}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"foo": params["foo"],
			"bar": params["bar"],
			"baz": params["baz"],
		})
	}).Methods("GET")

	echoApp = echo.New()
	echoApp.GET("/plaintext", func(c echo.Context) error {
		return c.String(http.StatusOK, plainTextBody)
	})
	echoApp.GET("/json", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &posts)
	})
	echoApp.GET("/params/:foo/:bar/:baz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"foo": c.Param("foo"),
			"bar": c.Param("bar"),
			"baz": c.Param("baz"),
		})
	})

	martiniApp = martini.Classic()
	martiniApp.Logger(log.New(io.Discard, "", 0))
	martiniApp.Get("/plaintext", func() string {
		return plainTextBody
	})
	martiniApp.Get("/json", func(w http.ResponseWriter) error {
		return json.NewEncoder(w).Encode(&posts)
	})
	martiniApp.Get("/params/:foo/:bar/:baz", func(params martini.Params) map[string]interface{} {
		return map[string]interface{}{
			"foo": params["foo"],
			"bar": params["bar"],
			"baz": params["baz"],
		}
	})
}

type Fatalfer interface {
	Fatalf(string, ...interface{})
}

func equals(f Fatalfer, a interface{}, b interface{}) {
	if a != b {
		f.Fatalf("expected %v to equal %v", a, b)
	}
}

func BenchmarkGowebPlaintext(b *testing.B) {
	req, err := http.NewRequest("GET", "/plaintext", nil)
	equals(b, err, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		gowebApp.ServeHTTP(rr, req)
		equals(b, rr.Code, http.StatusOK)
	}
}

func BenchmarkGinPlaintext(b *testing.B) {
	req, err := http.NewRequest("GET", "/plaintext", nil)
	equals(b, err, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ginApp.ServeHTTP(rr, req)
		equals(b, rr.Code, http.StatusOK)
	}
}

func BenchmarkGorillaPlaintext(b *testing.B) {
	req, err := http.NewRequest("GET", "/plaintext", nil)
	equals(b, err, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		gorillaApp.ServeHTTP(rr, req)
		equals(b, rr.Code, http.StatusOK)
	}
}

func BenchmarkEchoPlaintext(b *testing.B) {
	req, err := http.NewRequest("GET", "/plaintext", nil)
	equals(b, err, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		echoApp.ServeHTTP(rr, req)
		equals(b, rr.Code, http.StatusOK)
	}
}

func BenchmarkMartiniPlaintext(b *testing.B) {
	req, err := http.NewRequest("GET", "/plaintext", nil)
	equals(b, err, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		martiniApp.ServeHTTP(rr, req)
		equals(b, rr.Code, http.StatusOK)
	}
}

func BenchmarkGowebJSON(b *testing.B) {
	req, err := http.NewRequest("GET", "/json", nil)
	equals(b, err, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		gowebApp.ServeHTTP(rr, req)
		equals(b, rr.Code, http.StatusOK)
	}
}

func BenchmarkGinJSON(b *testing.B) {
	req, err := http.NewRequest("GET", "/json", nil)
	equals(b, err, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ginApp.ServeHTTP(rr, req)
		equals(b, rr.Code, http.StatusOK)
	}
}

func BenchmarkGorillaJSON(b *testing.B) {
	req, err := http.NewRequest("GET", "/json", nil)
	equals(b, err, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		gorillaApp.ServeHTTP(rr, req)
		equals(b, rr.Code, http.StatusOK)
	}
}

func BenchmarkEchoJSON(b *testing.B) {
	req, err := http.NewRequest("GET", "/json", nil)
	equals(b, err, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		echoApp.ServeHTTP(rr, req)
		equals(b, rr.Code, http.StatusOK)
	}
}

func BenchmarkMartiniJSON(b *testing.B) {
	req, err := http.NewRequest("GET", "/json", nil)
	equals(b, err, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		martiniApp.ServeHTTP(rr, req)
		equals(b, rr.Code, http.StatusOK)
	}
}

func BenchmarkGowebPathParams(b *testing.B) {
	req, err := http.NewRequest("GET", "/params/a/b/c", nil)
	equals(b, err, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		gowebApp.ServeHTTP(rr, req)
		equals(b, rr.Code, http.StatusOK)
	}
}

func BenchmarkGinPathParams(b *testing.B) {
	req, err := http.NewRequest("GET", "/params/a/b/c", nil)
	equals(b, err, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ginApp.ServeHTTP(rr, req)
		equals(b, rr.Code, http.StatusOK)
	}
}

func BenchmarkGorillaPathParams(b *testing.B) {
	req, err := http.NewRequest("GET", "/params/a/b/c", nil)
	equals(b, err, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		gorillaApp.ServeHTTP(rr, req)
		equals(b, rr.Code, http.StatusOK)
	}
}

func BenchmarkEchoPathParams(b *testing.B) {
	req, err := http.NewRequest("GET", "/params/a/b/c", nil)
	equals(b, err, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		echoApp.ServeHTTP(rr, req)
		equals(b, rr.Code, http.StatusOK)
	}
}

func BenchmarkMartiniPathParams(b *testing.B) {
	req, err := http.NewRequest("GET", "/params/a/b/c", nil)
	equals(b, err, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		martiniApp.ServeHTTP(rr, req)
		equals(b, rr.Code, http.StatusOK)
	}
}
