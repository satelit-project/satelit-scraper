package proxy

import "fmt"

type Protocol int

const (
	HTTP Protocol = iota
	HTTPS
	SOCKS5
)

func (p Protocol) String() string {
	switch p {
	case HTTP:
		return "http"
	case HTTPS:
		return "https"
	case SOCKS5:
		return "socks5"
	default:
		panic("unsupported protocol")
	}
}

type Proxy struct {
	Host  string
	Port  string
	Proto Protocol
}

func (p *Proxy) String() string {
	return fmt.Sprintf("%s://%s:%s", p.Proto, p.Host, p.Port)
}
