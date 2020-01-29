package main

import "github.com/twharmon/goweb"

func main() {
	app := goweb.New()
	app.AddStdLogger(goweb.LogLevelInfo)
	app.GET("/", func(c *goweb.Context) goweb.Responder {
		if err := c.Push("main.css"); err != nil {
			c.LogError(err)
		}
		return c.File("index.html")
	})
	app.ServeFiles("/", ".")
	app.RunTLS(&goweb.TLSConfig{
		HostPolicy: func(host string) error { return nil },
	})
}
