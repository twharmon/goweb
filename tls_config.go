package goweb

// TLSConfig contains TLS configuration information.
type TLSConfig struct {
	CertDir      string
	HostPolicy   func(string) error
	RedirectHTTP bool
}
