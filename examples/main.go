package main

import (
	"log"
	"net/http"

	"github.com/twharmon/gorouter"
)

func main() {
	r := gorouter.New()
	r.Get("/home", home)

	// Regurlar expressions are required for all route parameters
	r.Get("/hello/{name:[a-zA-Z]+}", hello)

	// If none of the registered routes are matched, serve static
	// files from this directory
	r.StaticFiles("static")

	log.Fatalln(http.ListenAndServe(":8000", r))
}

func home(w http.ResponseWriter, r *http.Request, params gorouter.Params) {
	w.Write([]byte("home"))
}

func hello(w http.ResponseWriter, r *http.Request, params gorouter.Params) {
	w.Write([]byte("hello " + params.Get("name")))
}
