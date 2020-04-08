package proxy

import (
	"fmt"
	"net/http"
)

// Proxy server protocol.
type Protocol int

const (
	HTTP Protocol = iota + 1
	HTTPS
)

// Returns string representation of the protocol.
func (p Protocol) String() string {
	switch p {
	case HTTP:
		return "http"
	case HTTPS:
		return "https"
	default:
		return "unsupported protocol"
	}
}

// Proxy server representation.
type Proxy struct {
	Host  string
	Port  string
	Proto Protocol
}

// Checks if the proxy available for usage.
func (p Proxy) isAvailable(client *http.Client) bool {
	_, err := client.Head(p.String())
	return err == nil
}

// Returns proxy address.
func (p Proxy) Address() string {
	return fmt.Sprintf("%s://%s:%s", p.Proto, p.Host, p.Port)
}

// Returns string representation of the proxy server.
func (p Proxy) String() string {
	return p.Address()
}

// Returns true if proxy is correctly specified.
func (p Proxy) IsValid() bool {
	return len(p.Host) > 0 && len(p.Port) > 0 && p.Proto != Protocol(0)
}
