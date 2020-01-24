package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/twharmon/goweb"
)

func main() {
	app := goweb.New()

	// uses std lib "log"
	app.AddStdLogger(goweb.LogLevelDebug)

	// create your own logger
	app.AddCustomLogger(&myLogger{goweb.LogLevelNotice})

	app.GET("/divide/{a}/by/{b}", divide)

	app.Run(":8080")
}

func divide(c *goweb.Context) goweb.Responder {
	a, err := strconv.ParseFloat(c.Param("a"), 64)
	if err != nil {
		c.LogNotice(fmt.Sprintf("attempted to divide %s", c.Param("a")))
		return c.BadRequest().Text("only numbers can be divided")
	}
	b, err := strconv.ParseFloat(c.Param("b"), 64)
	if err != nil {
		c.LogNotice(fmt.Sprintf("attempted to divide %s", c.Param("b")))
		return c.BadRequest().Text("only numbers can be divided")
	}

	if b == 0 {
		c.LogEmergency("division by zero")
		return c.BadRequest().Text("can not divide by zero")
	}

	return c.OK().Text(fmt.Sprintf("%f", a/b))
}

type myLogger struct {
	minLevel goweb.LogLevel
}

func (l *myLogger) Log(_ *goweb.Context, logLevel goweb.LogLevel, message interface{}) {
	if logLevel >= l.minLevel {
		log.Printf("Notice - %v", message)
	}
}
