// Package proxier finds proxies and makes requests through them for web crawling
package proxier

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vertoforce/proxier/proxy"
	"github.com/vertoforce/proxier/proxy/proxyDBs/inmemory"
	"github.com/vertoforce/proxier/proxy/proxysources/getproxylist"
	"github.com/vertoforce/proxier/proxy/proxysources/gimmeproxy"
)

// Defaults
const (
	DefaultProxyDBTimeout = time.Second * 5
	DefaultProxyTimeout   = time.Second * 4
)

// DefaultSources are the default proxy sources available
var (
	DefaultProxySources = []proxy.ProxySource{
		&getproxylist.GetProxyListSource{},
		&gimmeproxy.GimmeProxySource{},
	}

	AllowedProxyProtocols = []proxy.Protocol{proxy.Socks4Protocol, proxy.Socks4aProtocol, proxy.Socks5Protocol, proxy.Socks5hProtocol, proxy.SocksProtocol}
)

// CheckResponseFunc is a func given a http response it returns true if the request succeeded
// By default proxier looks for a HTTP 200
type CheckResponseFunc func(*http.Response) bool

// Proxier A proxier object
type Proxier struct {
	// proxySources are sources of proxies, using a map so we randomize our use of each
	proxySources map[proxy.ProxySource]bool
	// proxyDB is where we store the proxies we know about
	proxyDB proxy.ProxyDB
	// ProxyTimeout is how long to try a proxy before giving up
	ProxyTimeout time.Duration
	// True to check no proxy first by default
	TryNoProxyFirst bool
}

// NewBare Creates a new bare proxier with no proxy sources
func NewBare() *Proxier {
	p := &Proxier{}
	p.proxySources = map[proxy.ProxySource]bool{}
	p.ProxyTimeout = DefaultProxyDBTimeout
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
	for proxySource := range p.proxySources {
		// Try to find a valid proxy from this source
		for {
			var err error
			proxy, err = proxySource.GetProxy(ctx)
			if err != nil {
				// This proxy source has no more proxies (for now)
				break
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

// DoRequest Do a request using a random proxy in our DB and keep cycling through proxies until we find one that passes DefaultCheckResponseFunc
func (p *Proxier) DoRequest(ctx context.Context, method, URL string, body io.Reader) (*http.Response, error) {
	return p.DoRequestExtra(ctx, method, URL, body, p.TryNoProxyFirst, DefaultCheckResponseFunc)
}

// DoRequestExtra Same as DoRequest with additional
/// tryNoProxyFirst to try a normal request first
// and also a checkResponeFunc to check if the response was successful
func (p *Proxier) DoRequestExtra(ctx context.Context, method, URL string, body io.Reader, tryNoProxyFirst bool, checkResponseFunc CheckResponseFunc) (*http.Response, error) {
	// TODO: Add max request count to avoid endless requesting for a down server, or a server that always returns 403

	// -- Try default request --
	if tryNoProxyFirst {
		req, err := http.NewRequestWithContext(ctx, method, URL, body)
		if err == nil {
			resp, err := http.DefaultClient.Do(req)
			if err == nil && checkResponseFunc(resp) {
				return resp, nil
			}
		}
	}

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

	for proxy := range proxiesMap {
		resp, err := p.makeProxyRequest(ctx, proxy, method, URL, body)

		// Check if this was a success
		if err == nil && checkResponseFunc(resp) {
			return resp, nil
		}

		// This wasn't a success, we should ditch this proxy from the database
		// TODO: Change this to delete after 3 failures or something
		p.proxyDB.DelProxy(ctx, proxy)
	}

	// -- Get new proxies --

	// If we are here, there are no valid proxies available in the proxyDB
	// Keep trying to get new proxies forever (until we run out of proxies from our proxy sources)
	for {
		proxy, err := p.GetProxyFromSources(ctx)
		if err != nil {
			// No more proxies to try
			return nil, fmt.Errorf("no proxies available")
		}

		// Try this proxy
		resp, err := p.makeProxyRequest(ctx, proxy, method, URL, body)
		// Check if this was a success
		if err != nil || !checkResponseFunc(resp) {
			continue
		}

		// It worked!  Add this to our database
		p.proxyDB.StoreProxy(ctx, proxy)
		// Return response
		return resp, nil
	}
}

func (p *Proxier) makeProxyRequest(ctx context.Context, proxy *proxy.Proxy, method, URL string, body io.Reader) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, p.ProxyTimeout)
	defer cancel()
	return proxy.DoRequest(ctx, method, URL, body)
}

// DefaultCheckResponseFunc Returns true if status code is 200 OR 500-599, false if status code is 429 OR 403, true otherwise
func DefaultCheckResponseFunc(resp *http.Response) bool {
	// TODO: Change this from checking 200
	if resp.StatusCode == 200 {
		return true
	}

	// If there is a server error, it's not the proxies fault, the request succeeded
	if resp.StatusCode >= 500 && resp.StatusCode <= 599 {
		return true
	}

	// These typically indicate a failure
	if resp.StatusCode == 429 || resp.StatusCode == 403 {
		return false
	}

	// Return true otherwise
	return true
}
