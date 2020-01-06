// Package proxy helps make http requests with different proxies and user-agents
package proxy

import (
	"io"
	"net/http"
	"proxy/proxy"
	"proxy/proxysources/getproxylist"
	"proxy/proxysources/gimmeproxy"
)

// DefaultSources are the default proxy sources available
var DefaultSources = []proxy.ProxySource{
	&getproxylist.GetProxyListSource{},
	&gimmeproxy.GimmeProxySource{},
}

// Proxier
type Proxier struct {
	// proxySources are sources of proxies
	proxySources []proxy.ProxySource
	// currentProxies are the currently being used proxies
	currentProxies []proxy.Proxy
	// proxyDB is a store of proxies we currently know about
	proxyDB proxy.ProxyDB
}

// NewBare Creates a new bare proxier with no proxy sources
func NewBare() *Proxier {
	p := &Proxier{}
	return p
}

// New Creates a new proxier with default proxy sources and no proxyDB
func New() *Proxier {
	return NewBare().WithProxySources(DefaultSources...)
}

// WithProxySources Add proxy sources
func (p *Proxier) WithProxySources(sources ...proxy.ProxySource) *Proxier {
	p.proxySources = append(p.proxySources, sources...)
	return p
}

// WithProxyDB Add proxy DB
func (p *Proxier) WithProxyDB(proxyDB proxy.ProxyDB) *Proxier {
	p.proxyDB = proxyDB
	return p
}

func (p *Proxier) NewRequest(method, URL, string, body io.Reader) (*http.Request, error) {
	// TODO: Get proxy
	// TODO: Set user agent

	// Return req
	return nil, nil
}
