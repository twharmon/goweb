package goweb

import "regexp"

type route struct {
	regexp     *regexp.Regexp
	handler    Handler
	paramNames []string
	method     string
	host       string
}
