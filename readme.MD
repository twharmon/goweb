# Goweb
Light weight web framework based on net/http.

## Usage
Basic example
```
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

func isAuthUser(c *goweb.Context) *goweb.Response {
	c.Set("user", &User{
		ID:   45,
		Name: c.Param("name"),
	})
	// return nil to allow middleware to continue.
	return nil
}

func isNotGopher(c *goweb.Context) *goweb.Response {
	user := c.Get("user").(*User)
	if strings.ToLower(user.Name) == "gopher" {
		// return non nil Response to stop middleware.
		return c.Forbidden().PlainText("no gophers allowed")
	}
	return nil
}

func hello(c *goweb.Context) *goweb.Response {
	return c.OK().JSON(goweb.Map{
		"hello": "world",
	})
}

func getUser(c *goweb.Context) *goweb.Response {
	return c.OK().JSON(c.Get("user").(*User))
}
```
