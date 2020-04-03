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
```go
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

```
BenchmarkGowebPlaintext   	      1495 ns/op	    1528 B/op	      20 allocs/op
BenchmarkGinPlaintext     	       994 ns/op	    1096 B/op	      13 allocs/op
BenchmarkGorillaPlaintext 	      2209 ns/op	    2056 B/op	      23 allocs/op
BenchmarkEchoPlaintext    	      1032 ns/op	    1080 B/op	      14 allocs/op
BenchmarkMartiniPlaintext 	     14513 ns/op	    1889 B/op	      41 allocs/op

BenchmarkGowebJSON        	     96569 ns/op	   51026 B/op	      16 allocs/op
BenchmarkGinJSON          	     96568 ns/op	   50485 B/op	      12 allocs/op
BenchmarkGorillaJSON      	     98893 ns/op	   51428 B/op	      19 allocs/op
BenchmarkEchoJSON         	     98160 ns/op	   50482 B/op	      11 allocs/op
BenchmarkMartiniJSON                216111 ns/op	  144936 B/op	      41 allocs/op

BenchmarkGowebPathParams  	      4061 ns/op	    2697 B/op	      36 allocs/op
BenchmarkGinPathParams    	      2974 ns/op	    2024 B/op	      29 allocs/op
BenchmarkGorillaPathParams	      4600 ns/op	    3241 B/op	      37 allocs/op
BenchmarkEchoPathParams   	      2833 ns/op	    1976 B/op	      28 allocs/op
BenchmarkMartiniPathParams	     16084 ns/op	    2734 B/op	      47 allocs/op
```

## Contribute
Create a pull request to contribute to Goweb.
