package main

import (
	"github.com/twharmon/goweb"
)

func main() {
	app := goweb.New()

	app.GET("/", hello)

	app.Run(":8080")
}

func hello(c *goweb.Context) goweb.Responder {
	return c.JSON(goweb.Map{
		"hello": "world",
	})
}
