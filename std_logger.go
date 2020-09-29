package goweb

import (
	"fmt"
	"log"
)

type stdLogger struct{}

// DefaultLogger logs messages to stdout.
var DefaultLogger = &stdLogger{}

func (l *stdLogger) Log(c *Context, level LogLevel, messages ...interface{}) {
	query := c.Request.URL.Query().Encode()
	if query != "" {
		query = "?" + query
	}
	uri := fmt.Sprintf("%s %s%s%s", c.Request.Method, c.Request.Host, c.Request.URL.Path, query)
	prefix := fmt.Sprintf("[%s] %s -", level, uri)
	messages = append([]interface{}{prefix}, messages...)
	log.Println(messages...)
}
