package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/twharmon/goweb"
)

func main() {
	app := goweb.New()

	app.RegisterLogger(newLogger(goweb.LogLevelInfo))

	app.GET("/divide/{a}/by/{b}", divide)

	app.Run(":8080")
}

func divide(c *goweb.Context) goweb.Responder {
	a, err := strconv.ParseFloat(c.Param("a"), 64)
	if err != nil {
		c.LogNotice(fmt.Sprintf("attempted to divide %s", c.Param("a")))
		return c.Text(http.StatusBadRequest, "only numbers can be divided")
	}
	b, err := strconv.ParseFloat(c.Param("b"), 64)
	if err != nil {
		c.LogNotice(fmt.Sprintf("attempted to divide %s", c.Param("b")))
		return c.Text(http.StatusBadRequest, "only numbers can be divided")
	}

	if b == 0 {
		c.LogEmergency("division by zero", 5, "asdf")
		return c.Text(http.StatusBadRequest, "can not divide by zero")
	}

	return c.Text(http.StatusOK, fmt.Sprintf("%f", a/b))
}

type logger struct {
	level goweb.LogLevel
}

func newLogger(level goweb.LogLevel) goweb.Logger {
	return &logger{level: level}
}

func (l *logger) Log(c *goweb.Context, logLevel goweb.LogLevel, messages ...interface{}) {
	if l.level > logLevel {
		return
	}
	prefix := fmt.Sprintf("[%s] %s", logLevel, c.Request.URL.Path)
	messages = append([]interface{}{prefix}, messages...)
	log.Println(messages...)
}
