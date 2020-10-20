# Goweb

![](https://github.com/twharmon/goweb/workflows/Test/badge.svg) [![](https://goreportcard.com/badge/github.com/twharmon/goweb)](https://goreportcard.com/report/github.com/twharmon/goweb) [![](https://gocover.io/_badge/github.com/twharmon/goweb)](https://gocover.io/github.com/twharmon/goweb)

Light weight web framework based on net/http.

Includes
- routing
- middleware
- logging
- auto TLS (https)
- easy CORS

Goweb aims to
1. rely only on the standard library as much as possible
2. be flexible
3. perform well

## Usage
See [examples](https://github.com/twharmon/goweb/tree/master/examples).

### Basic usage
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

### Easily extendable
See [serving files](https://github.com/twharmon/goweb/tree/master/examples/files) or [template rendering](https://github.com/twharmon/goweb/tree/master/examples/templates) for examples.

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
