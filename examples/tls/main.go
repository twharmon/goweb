package main

import (
	"fmt"

	"github.com/twharmon/goweb"
)

func main() {
	app := goweb.New()

	app.GET("/", func(c *goweb.Context) goweb.Responder {
		return c.OK().PlainText("hello world")
	})

	app.RunTLS(&goweb.TLSConfig{
		// directory to store certificates
		CertDir: ".certs",

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

		// redirect http to https
		RedirectHTTP: true,
	})
}
