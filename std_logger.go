package goweb

import "log"

type stdLogger struct {
	level int
}

func (l *stdLogger) ShouldLog(level int) bool {
	return level >= l.level
}

func (l *stdLogger) Log(level int, message interface{}) {
	title, _ := getTitleAndColor(level)
	log.Printf("%s: %s - %s\n", title, caller(), message)
}
