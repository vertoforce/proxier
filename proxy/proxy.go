package proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"h12.io/socks"
)

// Protocol Proxy protocol
type Protocol string

// Protocols
const (
	Socks5hProtocol Protocol = "socks5h"
	Socks5Protocol  Protocol = "socks5"
	Socks4aProtocol Protocol = "socks4a"
	Socks4Protocol  Protocol = "socks4"
	SocksProtocol   Protocol = "socks"
	HTTPProtocol    Protocol = "http"
)

// Proxy Proxy
type Proxy struct {
	IP       string
	Port     uint16
	Protocol Protocol
}

// ProxySource is a source of proxies
type ProxySource interface {
	GetProxy(ctx context.Context) (*Proxy, error)
}

// ProxyDB is an interface to store and retrieve proxies.
// Most proxy sources cost money to keep getting proxies, so
// this interface allows for the storing of proxies (in a database or something else)
// Note that Proxier still stores a local slice of proxies in use which is a cache of this DB
type ProxyDB interface {
	GetProxies(context.Context) ([]*Proxy, error)
	StoreProxy(context.Context, *Proxy) error
	DelProxy(context.Context, *Proxy) error
	Clear(context.Context) error
}

// Address gets address of proxy
func (p *Proxy) Address() string {
	return fmt.Sprintf("%s:%d", p.IP, p.Port)
}

// DoRequest makes a request to this proxy
func (p *Proxy) DoRequest(ctx context.Context, method, URL string, body io.Reader) (*http.Response, error) {
	switch p.Protocol {
	case SocksProtocol, Socks4Protocol, Socks4aProtocol, Socks5Protocol:
		return p.doRequestSocks(ctx, method, URL, body)
	case HTTPProtocol:
		return p.doRequestHTTP(ctx, method, URL, body)
	default:
		return nil, fmt.Errorf("No function to use this protocol")
	}
}

func (p *Proxy) doRequestSocks(ctx context.Context, method, URL string, body io.Reader) (*http.Response, error) {
	// Create socks proxy
	socksType := socks.SOCKS4
	switch p.Protocol {
	case Socks4Protocol:
		socksType = socks.SOCKS4
	case Socks4aProtocol:
		socksType = socks.SOCKS4A
	case Socks5Protocol:
		socksType = socks.SOCKS5
	}
	dialSocksProxy := socks.DialSocksProxy(socksType, p.Address())
	tr := &http.Transport{Dial: dialSocksProxy}
	httpClient := &http.Client{Transport: tr}

	// Create request
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	// Make request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *Proxy) doRequestHTTP(ctx context.Context, method, URL string, body io.Reader) (*http.Response, error) {
	proxyURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}

	// Create request
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	// Make request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
