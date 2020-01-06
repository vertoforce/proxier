package proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/proxy"
)

type Protocol string

const (
	Socks5 Protocol = "socks5"
	HTTP   Protocol = "HTTP"
)

type Proxy struct {
	IP       string
	Port     int16
	Protocol Protocol
}

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
}

func (p *Proxy) Address() string {
	return fmt.Sprintf("%s:%d", p.IP, p.Port)
}

func (p *Proxy) DoRequest(ctx context.Context, method, URL string, body io.Reader) (*http.Response, error) {
	switch p.Protocol {
	case "socks5":
		return p.doRequestSocks5(ctx, method, URL, body)
	default:
		return nil, fmt.Errorf("No function to use this protocol")
	}
}

func (p *Proxy) doRequestSocks5(ctx context.Context, method, URL string, body io.Reader) (*http.Response, error) {
	// Create socks5 client
	dialer, err := proxy.SOCKS5("tcp", p.Address(), nil, proxy.Direct)
	if err != nil {
		return nil, err
	}
	httpTransport := &http.Transport{}
	httpTransport.Dial = dialer.Dial
	httpClient := &http.Client{Transport: httpTransport}

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
