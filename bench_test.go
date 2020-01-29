package goweb_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/twharmon/goweb"
)

func BenchmarkGoweb(b *testing.B) {
	app := goweb.New()
	app.GET("/", func(c *goweb.Context) goweb.Responder {
		return c.Text("Hello, world!")
	})
	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			res, err := http.Get("http://localhost:8080/")
			if err != nil {
				b.Fatal(err)
			}
			if res.StatusCode != http.StatusOK {
				b.Fatalf("expected response status code to be 200; got %d", res.StatusCode)
			}
		}
		if err := app.Shutdown(); err != nil {
			b.Fatal(err)
		}
	}()
	app.Run(":8080")
}

func BenchmarkGin(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, world!")
	})
	srv := &http.Server{
		Addr:    ":8081",
		Handler: app,
	}
	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			res, err := http.Get("http://localhost:8081/")
			if err != nil {
				b.Fatal(err)
			}
			if res.StatusCode != http.StatusOK {
				b.Fatalf("expected response status code to be 200; got %d", res.StatusCode)
			}
		}
		if err := srv.Shutdown(context.TODO()); err != nil {
			b.Fatal(err)
		}
	}()
	srv.ListenAndServe()
}
