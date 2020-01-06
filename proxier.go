// Package proxy helps make http requests with different proxies and user-agents
package proxier

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"proxy/proxy"
	"proxy/proxy/proxysources/getproxylist"
	"proxy/proxy/proxysources/gimmeproxy"
	"time"
)

const (
	DefaultProxyDBTimeout = time.Second * 5
)

// DefaultSources are the default proxy sources available
var DefaultSources = []proxy.ProxySource{
	&getproxylist.GetProxyListSource{},
	&gimmeproxy.GimmeProxySource{},
}

// Proxier
type Proxier struct {
	// proxySources are sources of proxies, using a map so we randomize our use of each
	proxySources map[proxy.ProxySource]bool
	// proxyDB is where we store the proxies we know about
	proxyDB proxy.ProxyDB
}

// NewBare Creates a new bare proxier with no proxy sources
func NewBare() *Proxier {
	p := &Proxier{}
	p.proxySources = map[proxy.ProxySource]bool{}
	return p
}

// New Creates a new proxier with default proxy sources and no proxyDB
func New() *Proxier {
	return NewBare().WithProxySources(DefaultSources...)
}

// WithProxySources Add proxy sources
func (p *Proxier) WithProxySources(sources ...proxy.ProxySource) *Proxier {
	for _, proxySource := range sources {
		p.proxySources[proxySource] = true
	}
	return p
}

// WithProxyDB Add proxy DB, there can only be one proxy db
func (p *Proxier) WithProxyDB(proxyDB proxy.ProxyDB) *Proxier {
	p.proxyDB = proxyDB
	return p
}

// -- functionality --

// GetProxyFromSources Get a ProxySource from one of our proxySources
func (p *Proxier) GetProxyFromSources(ctx context.Context) (*proxy.Proxy, error) {
	var proxy *proxy.Proxy
	for proxySource, _ := range p.proxySources {
		var err error
		proxy, err = proxySource.GetProxy(ctx)
		if err != nil {
			continue
		}
		// We found a proxy!
		return proxy, nil
	}

	// No proxies to be found
	return nil, fmt.Errorf("no new proxies available")
}

// CacheProxies Get count proxies from our sources and put it in the database for later use
func (p *Proxier) CacheProxies(count int64) error {
	// TODO:
	return nil
}

func (p *Proxier) DoRequest(ctx context.Context, method, URL string, body io.Reader) (*http.Response, error) {
	var resp *http.Response
	var err error
	// Try our current proxies in the database
	if p.proxyDB == nil {
		return nil, fmt.Errorf("no proxydb set and is required")
	}
	proxies, err := p.proxyDB.GetProxies(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: Randomize which proxies we try
	for _, proxy := range proxies {
		resp, err = proxy.DoRequest(ctx, method, URL, body)

		// Check if this was a success
		if err == nil && resp.StatusCode == 200 {
			return resp, nil
		}

		// This wasn't a success, we should ditch this proxy from the database
		// p.proxyDB.DelProxy(ctx, proxy)
	}

	// If we are here, there are no valid proxies available in the proxyDB
	// Keep trying to get new proxies
	for {
		proxy, err := p.GetProxyFromSources(ctx)
		if err != nil {
			// No more proxies to try
			return nil, fmt.Errorf("no proxies available")
		}

		// Try this proxy
		resp, err := proxy.DoRequest(ctx, method, URL, body)
		if err != nil {
			continue
		}

		// It worked!  Add this to our database
		p.proxyDB.StoreProxy(ctx, proxy)

		return resp, nil
	}
}
