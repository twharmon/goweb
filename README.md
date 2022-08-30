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
    return c.JSON(http.StatusOK, goweb.Map{
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
    return c.JSON(http.StatusOK, goweb.Map{
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
		return c.JSON(http.StatusOK, goweb.Map{
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
BenchmarkGinPlaintext-10         	 2706439	       440.0 ns/op	    1040 B/op	       9 allocs/op
BenchmarkEchoPlaintext-10        	 2549317	       470.7 ns/op	    1024 B/op	      10 allocs/op
BenchmarkGowebPlaintext-10       	 1584044	       756.6 ns/op	    1456 B/op	      16 allocs/op
BenchmarkGorillaPlaintext-10     	 1000000	      1027 ns/op	    1744 B/op	      17 allocs/op
BenchmarkMartiniPlaintext-10     	  223416	      5364 ns/op	    1789 B/op	      39 allocs/op

BenchmarkGowebJSON-10            	   25945	     46359 ns/op	   50905 B/op	      15 allocs/op
BenchmarkEchoJSON-10             	   25664	     46571 ns/op	   50641 B/op	      10 allocs/op
BenchmarkGorillaJSON-10          	   25716	     46857 ns/op	   51115 B/op	      16 allocs/op
BenchmarkGinJSON-10              	   23697	     50697 ns/op	  100836 B/op	      10 allocs/op
BenchmarkMartiniJSON-10          	   22746	     52613 ns/op	   52665 B/op	      41 allocs/op

BenchmarkGinPathParams-10        	  914139	      1273 ns/op	    1849 B/op	      25 allocs/op
BenchmarkEchoPathParams-10       	  889014	      1309 ns/op	    1865 B/op	      25 allocs/op
BenchmarkGowebPathParams-10      	  627306	      1902 ns/op	    2570 B/op	      33 allocs/op
BenchmarkGorillaPathParams-10    	  552852	      2144 ns/op	    2874 B/op	      32 allocs/op
BenchmarkMartiniPathParams-10    	  188500	      6215 ns/op	    2641 B/op	      47 allocs/op
```

## Contribute
Create a pull request to contribute to Goweb.
