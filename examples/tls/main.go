package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/twharmon/goweb"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	app := goweb.New()

	app.GET("/", func(c *goweb.Context) goweb.Responder {
		return c.JSON(goweb.Map{
			"hello": "world",
		})
	})

	serveTLS(app)
}

func serveTLS(app *goweb.Engine) {
	m := &autocert.Manager{
		Cache:  autocert.DirCache(".certs"),
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(_ context.Context, host string) error {
			if host == "example.com" {
				return nil
			}
			return errors.New("host not configured")
		},
	}
	go http.ListenAndServe(":http", m.HTTPHandler(nil))
	s := &http.Server{
		Addr:      ":https",
		TLSConfig: m.TLSConfig(),
		Handler:   app,
	}

	log.Fatalln(s.ListenAndServeTLS("", ""))
}
