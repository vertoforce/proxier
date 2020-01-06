package proxy

import "context"

type Proxy struct {
	IP   string
	Port int16
}

type ProxySource interface {
	GetProxy(ctx context.Context) (*Proxy, error)
}

// ProxyDB is an interface to store and retrieve proxies.
// Most proxy sources cost money to keep getting proxies, so
// this interface allows for the storing of proxies (in a database or something else)
// Note that Proxier still stores a local slice of proxies in use which is a cache of this DB
type ProxyDB interface {
	GetProxies() []Proxy
	DelProxy(Proxy) error
	StoreProxy(Proxy) error
}
