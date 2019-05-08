package goweb

import jsoniter "github.com/json-iterator/go"

const (
	contentTypeHeader          = "Content-Type"
	contentTypeApplicationJSON = "application/json; charset=utf-8"
	contentTypeTextPlain       = "text/plain; charset=utf-8"
)

// Response containts the Context, HTTP status code, a
// response body, and content type.
type Response struct {
	context *Context
	body    interface{}
}

// JSON returns a Response with JSON body.
func (r *Response) JSON(value interface{}) *Response {
	r.context.writer.Header().Set(contentTypeHeader, contentTypeApplicationJSON)
	r.body = value
	return r
}

// PlainText returns a Response with plain text body.
func (r *Response) PlainText(text string) *Response {
	r.context.writer.Header().Set(contentTypeHeader, contentTypeTextPlain)
	r.body = text
	return r
}

func (r *Response) send() {
	if r.body != nil {
		switch r.context.writer.Header().Get(contentTypeHeader) {
		case contentTypeApplicationJSON:
			jsoniter.NewEncoder(r.context.writer).Encode(r.body)
		case contentTypeTextPlain:
			r.context.writer.Write([]byte(r.body.(string)))
		}
	}
}
