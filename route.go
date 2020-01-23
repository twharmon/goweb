package goweb

import (
	"regexp"

	"golang.org/x/net/websocket"
)

type route struct {
	regexp           *regexp.Regexp
	handler          Handler
	webSocketHandler websocket.Handler
	paramNames       []string
	method           string
	host             string
}
