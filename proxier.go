// Package proxy helps make http requests with different proxies and user-agents
package proxier

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"proxy/proxy"
	"proxy/proxy/proxyDBs/inmemory"
	"proxy/proxy/proxysources/getproxylist"
	"proxy/proxy/proxysources/gimmeproxy"
	"time"
)

const (
	DefaultProxyDBTimeout = time.Second * 5
)

// DefaultSources are the default proxy sources available
var (
	DefaultProxySources = []proxy.ProxySource{
		&getproxylist.GetProxyListSource{},
		&gimmeproxy.GimmeProxySource{},
	}

	AllowedProxyProtocols = []proxy.Protocol{proxy.Socks4Protocol, proxy.Socks4aProtocol, proxy.Socks5Protocol, proxy.Socks5hProtocol, proxy.SocksProtocol}
)

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

// New Creates a new proxier with default proxy sources and in memory proxyDB
func New() *Proxier {
	return NewBare().WithProxySources(DefaultProxySources...).WithProxyDB(inmemory.New())
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

// WithProxies Adds proxies to our DB
func (p *Proxier) WithProxies(ctx context.Context, proxies ...*proxy.Proxy) *Proxier {
	for _, proxy := range proxies {
		p.proxyDB.StoreProxy(ctx, proxy)
	}
	return p
}

// -- functionality --

// GetProxyFromSources Get a ProxySource from one of our proxySources
// This will continue to try and get proxies from each source until it finds a SOCKS proxy
func (p *Proxier) GetProxyFromSources(ctx context.Context) (*proxy.Proxy, error) {
	var proxy *proxy.Proxy
proxySourceLoop:
	for proxySource, _ := range p.proxySources {
		// Try to find a valid proxy from this source
		for {
			var err error
			proxy, err = proxySource.GetProxy(ctx)
			if err != nil {
				continue proxySourceLoop
			}

			// Check if it's our allowed protocols
			for _, protocol := range AllowedProxyProtocols {
				if proxy.Protocol == protocol {
					// We found a proxy!
					return proxy, nil
				}
			}
		}
	}

	// No proxies to be found
	return nil, fmt.Errorf("no new proxies available")
}

// CacheProxies Get "count" proxies from our sources and put each in the database for later use
func (p *Proxier) CacheProxies(ctx context.Context, count int) (added int, err error) {
	added = 0
	for i := 0; i < count; i++ {
		// Get proxy
		proxy, err := p.GetProxyFromSources(ctx)
		if err != nil {
			// No more proxies available
			break
		}
		// Store proxy
		err = p.proxyDB.StoreProxy(ctx, proxy)
		if err != nil {
			return added, err
		}
		added++
	}

	return added, nil
}

// DoRequest Do a request using a random proxy in our DB and keep cycling through proxies until we find one that returns 200 OK
// TODO: Do not reply just on 200 OK as indication the proxy "worked"
func (p *Proxier) DoRequest(ctx context.Context, method, URL string, body io.Reader) (*http.Response, error) {
	var resp *http.Response
	var err error

	// -- Try our DB Proxies --
	if p.proxyDB == nil {
		return nil, fmt.Errorf("no proxydb set and is required")
	}
	proxies, err := p.proxyDB.GetProxies(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to map so we use randomly
	proxiesMap := map[*proxy.Proxy]bool{}
	for _, proxy := range proxies {
		proxiesMap[proxy] = true
	}

	for proxy, _ := range proxiesMap {
		resp, err = proxy.DoRequest(ctx, method, URL, body)

		// TODO: Change this
		// Check if this was a success
		if err == nil && resp.StatusCode == 200 {
			return resp, nil
		}

		// This wasn't a success, we should ditch this proxy from the database
		p.proxyDB.DelProxy(ctx, proxy)
	}

	// -- We need new proxies --
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
		// TODO: Change this from checking 200
		// Check if this was a success
		if err != nil || resp.StatusCode != 200 {
			continue
		}

		// It worked!  Add this to our database
		p.proxyDB.StoreProxy(ctx, proxy)

		return resp, nil
	}
}
