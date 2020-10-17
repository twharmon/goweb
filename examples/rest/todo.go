package main

import (
	"fmt"
	"github.com/twharmon/goweb"
	"strconv"
)

type Todo struct {
	Text string
}
type todoResource struct {
	Todos   map[string]Todo
	counter int
}

func New() todoResource {
	return todoResource{
		Todos:   make(map[string]Todo),
		counter: 0,
	}
}
func (t todoResource) Index(c *goweb.Context) goweb.Responder {
	fmt.Println(t.Todos)
	return c.JSON(t.Todos)
}

func (t todoResource) Get(c *goweb.Context) goweb.Responder {
	id := c.Param(t.Identifier())
	return c.JSON(t.Todos[id])
}

func (t *todoResource) Put(c *goweb.Context) goweb.Responder {
	id := c.Param(t.Identifier())
	t.Todos[id] = Todo{"put called"}
	return c.JSON(t.Todos[id])
}

func (t *todoResource) Delete(c *goweb.Context) goweb.Responder {
	delete(t.Todos, c.Param(t.Identifier()))
	return c.JSON(len(t.Todos))
}

func (t *todoResource) Post(c *goweb.Context) goweb.Responder {
	id := t.counter
	t.counter++
	t.Todos[strconv.Itoa(id)] = Todo{Text: "new from post"}
	return c.JSON(id)
}

func (t todoResource) Identifier() string {
	return "id"
}
