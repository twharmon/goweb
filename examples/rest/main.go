package main

import (
	"github.com/twharmon/goweb"
)

func main() {
	app := goweb.New()
	t := New()
	app.Resource("/todos", t)
	app.Run(":8080")
}
