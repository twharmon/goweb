package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/twharmon/goweb"
)

// Todo represents a todo item.
type Todo struct {
	Text string
}

// TodoResource implements the resource interface.
type TodoResource struct {
	Todos   map[string]Todo
	counter int
}

// New returns a new todo resource instance.
func New() *TodoResource {
	return &TodoResource{
		Todos:   make(map[string]Todo),
		counter: 0,
	}
}

// Index returns a list of todos.
func (t *TodoResource) Index(c *goweb.Context) goweb.Responder {
	fmt.Println(t.Todos)
	return c.JSON(http.StatusOK, t.Todos)
}

// Get returns one todo.
func (t *TodoResource) Get(c *goweb.Context) goweb.Responder {
	id := c.Param(t.Identifier())
	return c.JSON(http.StatusOK, t.Todos[id])
}

// Put updates a todo.
func (t *TodoResource) Put(c *goweb.Context) goweb.Responder {
	id := c.Param(t.Identifier())
	t.Todos[id] = Todo{"put called"}
	return c.JSON(http.StatusOK, t.Todos[id])
}

// Delete removes a todo.
func (t *TodoResource) Delete(c *goweb.Context) goweb.Responder {
	delete(t.Todos, c.Param(t.Identifier()))
	return c.JSON(http.StatusOK, len(t.Todos))
}

// Post creates a new todo.
func (t *TodoResource) Post(c *goweb.Context) goweb.Responder {
	id := t.counter
	t.counter++
	t.Todos[strconv.Itoa(id)] = Todo{Text: "new from post"}
	return c.JSON(http.StatusOK, id)
}

// Identifier returns the identifying property.
func (t *TodoResource) Identifier() string {
	return "id"
}
