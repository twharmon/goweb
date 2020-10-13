package goweb

import "golang.org/x/crypto/acme/autocert"

// TLSConfig contains TLS configuration information.
type TLSConfig struct {
	HostPolicy func(string) error
	AllowHTTP  bool
	Cache      autocert.Cache
}
