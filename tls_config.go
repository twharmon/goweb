package goweb

// TLSConfig contains TLS information.
type TLSConfig struct {
	CertDir      string
	HostPolicy   func(string) error
	RedirectHTTP bool
}
