package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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
	r := mux.NewRouter()
	r.HandleFunc("/plaintext", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}).Methods("GET")
	r.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&posts)
	}).Methods("GET")
	r.HandleFunc("/params/{val}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"val": params["val"],
		})
	})
	http.ListenAndServe(":8082", r)
}
