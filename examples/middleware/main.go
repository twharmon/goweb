package main

import (
	"strings"

	"github.com/twharmon/goweb"
)

// User contains user information.
type User struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func main() {
	app := goweb.New()

	auth := goweb.NewMiddleware()
	auth.Use(isAuthUser, isNotGopher)

	app.GET("/", hello)
	app.GET("/user/{name}", auth.Apply(getUser))

	app.Run(":8080")
}

func isAuthUser(c *goweb.Context) goweb.Responder {
	c.Set("user", &User{
		ID:   45,
		Name: c.Param("name"),
	})

	// return nil to continue the middleware chain
	return nil
}

func isNotGopher(c *goweb.Context) goweb.Responder {
	user := c.Get("user").(*User)
	if strings.ToLower(user.Name) == "gopher" {
		// return non nil Responder to terminate middleware chain
		return c.Forbidden().Text("no gophers allowed")
	}

	// return nil to continue the middleware chain
	return nil
}

func hello(c *goweb.Context) goweb.Responder {
	return c.JSON(goweb.Map{
		"hello": "world",
	})
}

func getUser(c *goweb.Context) goweb.Responder {
	return c.JSON(c.Get("user").(*User))
}
