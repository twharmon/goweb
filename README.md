# Goweb

![](https://github.com/twharmon/goweb/workflows/Test/badge.svg) [![](https://goreportcard.com/badge/github.com/twharmon/goweb)](https://goreportcard.com/report/github.com/twharmon/goweb) [![](https://gocover.io/_badge/github.com/twharmon/goweb)](https://gocover.io/github.com/twharmon/goweb) [![GoDoc](https://godoc.org/github.com/twharmon/goweb?status.svg)](https://godoc.org/github.com/twharmon/goweb)

Light weight web framework based on net/http.

Includes
- routing
- middleware
- autoTLS
- logging
- websockets

Goweb aims to
1. rely only on the standard library as much as possible
2. be flexible
3. perform well

## Usage
See [examples](https://github.com/twharmon/goweb/tree/master/examples).
```
package main

import (
	"github.com/twharmon/goweb"
)

func main() {
    app := goweb.New()
    app.GET("/hello/{name}", func(c *goweb.Context) goweb.Responder {
        return c.JSON(goweb.Map{
            "hello": c.Param("name"),
        })
    })
    app.Run(":8080")
}
```

## Documentation
For full documentation see [godoc](https://godoc.org/github.com/twharmon/goweb).

## Benchmarks
Goweb is built on Golang's standard library.
Frameworks that significantly beat Goweb usually depend on packages like `fasthttp`/`fasthttprouter`.
These are great packages, but they also have some limitations.
First, `fasthttp` does not support HTTP/2 (yet).
Second, no HTTP/2 support means packages like autocert will also not work with `fasthttp`.
Third, `fasthttprouter` does not allow you to register static routes and parameters for the same path segment.
For example, `/posts/{id}` and `/post/new` can not both be registered.

Goweb avoids these limitations by using the standard library's `net/http` instead of `fasthttp`/`fasthttprouter`.

Plaintext response "hello world" Requests/sec:
```
Goweb       79649.12
Gin         81365.57
Gorilla     77746.85
```

JSON response 100 posts Requests/sec:
```
Goweb       15089.02
Gin         13191.71
Gorilla     12836.61
```


JSON response with path parameter Requests/sec:
```
Goweb       71947.38
Gin         72422.31
Gorilla     57994.28
```

## Contribute
Create a pull request to contribute to Goweb.
