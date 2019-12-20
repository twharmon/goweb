package goweb

import (
	"fmt"
	"log"
)

type stdLogger struct {
	minLevel LogLevel
}

func (l *stdLogger) Log(c *Context, level LogLevel, message interface{}) {
	if level >= l.minLevel {
		query := c.Request.URL.Query().Encode()
		if query != "" {
			query = "?" + query
		}
		scheme := "http://"
		if c.Request.TLS != nil {
			scheme = "https://"
		}
		uri := fmt.Sprintf("%s %s%s%s%s", c.Request.Method, scheme, c.Request.Host, c.Request.URL.Path, query)
		log.Printf("%s: %s - %s\n", level.String(), uri, message)
	}
}
