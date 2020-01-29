package main

import (
	"github.com/twharmon/goweb"
)

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func main() {
	var posts []Post
	for i := 0; i < 100; i++ {
		posts = append(posts, Post{
			ID:    123,
			Title: "Lorem Ipsum",
			Body:  "Veniam ipsum officia consequat minim veniam cillum incididunt laborum aliqua ad do magna sed aliquip fugiat. Cillum et aliqua commodo, velit minim anim, pariatur, magna culpa officia dolor quis consectetur. Proident commodo laboris eu eu quis esse ea exercitation irure pariatur duis nulla deserunt dolor sed. Nulla qui laboris ut ea non consectetur amet culpa, pariatur commodo magna deserunt nostrud in.",
		})
	}
	app := goweb.New()
	app.GET("/plaintext", func(c *goweb.Context) goweb.Responder {
		return c.Text("hello world")
	})
	app.GET("/json", func(c *goweb.Context) goweb.Responder {
		return c.JSON(&posts)
	})
	app.GET("/params/{val}", func(c *goweb.Context) goweb.Responder {
		return c.JSON(goweb.Map{
			"val": c.Param("val"),
		})
	})
	app.Run(":8080")
}
