package main

import (
	"time"

	"github.com/twharmon/goweb"
)

func main() {
	app := goweb.New()
	app.AutoCors(&goweb.CorsConfig{
		Headers: []string{"Authorization", "Content-Type"},
		MaxAge:  time.Hour,
		Origin:  "*",
	})

	app.GET("/hello", hello)

	app.Run(":8080")
}

func hello(c *goweb.Context) goweb.Responder {
	return c.JSON(goweb.Map{
		"hello": "world",
	})
}
