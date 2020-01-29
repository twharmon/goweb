package main

import (
	"fmt"

	"github.com/twharmon/goweb"
)

func main() {
	app := goweb.New()

	app.GET("/", func(c *goweb.Context) goweb.Responder {
		return c.Text("hello world")
	})

	app.RedirectWWW()

	app.RunTLS(&goweb.TLSConfig{
		// cache for Let's Encrypt certificates
		Cache: goweb.DirCache("certs"), // this is the default

		// function to determine which hosts are allowed
		HostPolicy: func(host string) error {
			allowedHosts := []string{"example.com", "www.example.com"}
			for _, allowedHost := range allowedHosts {
				if host == allowedHost {
					return nil
				}
			}
			return fmt.Errorf("host %s not allowed", host)
		},

		// When set to false (default), http is redirected to https.
		AllowHTTP: true,
	})
}
