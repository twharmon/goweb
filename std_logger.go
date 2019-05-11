package goweb

import "log"

type stdLogger struct {
	minLevel int
}

func (l *stdLogger) Log(level int, message interface{}) {
	if level >= l.minLevel {
		title, _ := getTitleAndColor(level)
		log.Printf("%s: %s - %s\n", title, caller(), message)
	}
}
