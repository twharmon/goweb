package main

import (
	"net/http"

	"github.com/twharmon/goweb"
)

func main() {
	app := goweb.New()

	app.GET("/hello", func(c *goweb.Context) goweb.Responder {
		return c.JSON(http.StatusOK, goweb.Map{
			"hello": "world",
		})
	})
	app.GET("/hello/{name}", func(c *goweb.Context) goweb.Responder {
		return c.JSON(http.StatusOK, goweb.Map{
			"hello": c.Param("name"),
		})
	})
	app.GET("/hello/{name}/{age:[0-9]+}", func(c *goweb.Context) goweb.Responder {
		return c.JSON(http.StatusOK, goweb.Map{
			"hello": c.Param("name"),
			"age":   c.Param("age"),
		})
	})

	app.Run(":8080")
}
