# Goweb

![](https://github.com/twharmon/goweb/workflows/Test/badge.svg) [![](https://goreportcard.com/badge/github.com/twharmon/goweb)](https://goreportcard.com/report/github.com/twharmon/goweb) [![](https://gocover.io/_badge/github.com/twharmon/goweb)](https://gocover.io/github.com/twharmon/goweb)

Light weight web framework based on net/http.

Includes
- routing
- middleware
- logging

Goweb aims to
1. rely only on the standard library as much as possible
2. be flexible
3. perform well

## Usage
See [examples](https://github.com/twharmon/goweb/tree/master/examples).

### Basic
```go
package main

import (
	"github.com/twharmon/goweb"
)

func main() {
    app := goweb.New()
    app.GET("/hello/{name}", hello)
    app.Run(":8080")
}

func hello(c *goweb.Context) goweb.Responder {
    return c.JSON(goweb.Map{
        "hello": c.Param("name"),
    })
}
```

### Logging
```go
package main

import (
	"github.com/twharmon/goweb"
)

func main() {
    app := goweb.New()
	app.RegisterLogger(newLogger(goweb.LogLevelInfo))
    app.GET("/hello/{name}", hello)
    app.Run(":8080")
}

func hello(c *goweb.Context) goweb.Responder {
    c.LogInfo("param name:", c.Param("name"))
    // logs "[INFO] /hello/Gopher param name: Gopher"
    return c.JSON(goweb.Map{
        "hello": c.Param("name"),
    })
}

type logger struct{
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
```

### Auto TLS
```go
package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/twharmon/goweb"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	app := goweb.New()
	app.GET("/", func(c *goweb.Context) goweb.Responder {
		return c.JSON(goweb.Map{
			"hello": "world",
		})
	})
	serveTLS(app)
}

func serveTLS(app *goweb.Engine) {
	m := &autocert.Manager{
		Cache:  autocert.DirCache(".certs"),
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(_ context.Context, host string) error {
			if host == "example.com" {
				return nil
			}
			return errors.New("host not configured")
		},
	}
	go http.ListenAndServe(":http", m.HTTPHandler(nil))
	s := &http.Server{
		Addr:      ":https",
		TLSConfig: m.TLSConfig(),
		Handler:   app,
	}
	log.Fatalln(s.ListenAndServeTLS("", ""))
}
```

### Easily extendable
See [serving files](https://github.com/twharmon/goweb/tree/master/examples/files), [template rendering](https://github.com/twharmon/goweb/tree/master/examples/templates), [tls](https://github.com/twharmon/goweb/tree/master/examples/tls), and [logging](https://github.com/twharmon/goweb/tree/master/examples/logging) for examples.

## Documentation
For full documentation see [pkg.go.dev](https://pkg.go.dev/github.com/twharmon/goweb).

## Benchmarks
```
BenchmarkGinPlaintext        	 	       780 ns/op	    1040 B/op	       9 allocs/op
BenchmarkEchoPlaintext       	 	       817 ns/op	    1024 B/op	      10 allocs/op
BenchmarkGowebPlaintext      	  	      1241 ns/op	    1456 B/op	      16 allocs/op
BenchmarkGorillaPlaintext    	  	      1916 ns/op	    2032 B/op	      19 allocs/op
BenchmarkMartiniPlaintext    	   	     14448 ns/op	    1779 B/op	      36 allocs/op

BenchmarkGowebJSON           	   	     60042 ns/op	   50798 B/op	      15 allocs/op
BenchmarkGorillaJSON         	   	     61086 ns/op	   51330 B/op	      18 allocs/op
BenchmarkEchoJSON            	   	     61115 ns/op	   50280 B/op	      10 allocs/op
BenchmarkGinJSON             	   	     68322 ns/op	  100116 B/op	      10 allocs/op
BenchmarkMartiniJSON         	   	     96365 ns/op	  144335 B/op	      38 allocs/op

BenchmarkGinPathParams       	  	      2464 ns/op	    1952 B/op	      27 allocs/op
BenchmarkEchoPathParams      	  	      2600 ns/op	    1968 B/op	      27 allocs/op
BenchmarkGowebPathParams     	  	      3591 ns/op	    2673 B/op	      35 allocs/op
BenchmarkGorillaPathParams   	  	      4220 ns/op	    3265 B/op	      36 allocs/op
BenchmarkMartiniPathParams   	   	     15211 ns/op	    2657 B/op	      45 allocs/op
```

## Contribute
Create a pull request to contribute to Goweb.
