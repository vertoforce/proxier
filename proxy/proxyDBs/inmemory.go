// Pacakge inmemory is a simple in memory proxy DB for testing
package inmemory

import "proxy/proxy"

type InMemoryProxyDB struct {
	Proxies []proxy.Proxy
}

func (db *InMemoryProxyDB) GetProxies() []proxy.Proxy {
	return db.Proxies
}
