package main

import "github.com/twharmon/goweb"

func main() {
	app := goweb.New()
	app.ServeFiles("/", ".")
	app.Run(":8080")
}
