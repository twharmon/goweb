package goweb

import "log"

type stdLogger struct {
	minLevel int
}

func (l *stdLogger) Log(c *Context, level int, message interface{}) {
	if level >= l.minLevel {
		title, _ := getTitleAndColor(level)
		log.Printf("%s: %s - %s\n", title, c.Request.URL.Path, message)
	}
}
