package inmemory

import (
	"context"

	"github.com/vertoforce/proxier/proxy"
)

type InMemoryProxyDB struct {
	proxies []*proxy.Proxy
}

func New() *InMemoryProxyDB {
	db := &InMemoryProxyDB{}
	return db
}

func (db *InMemoryProxyDB) GetProxies(ctx context.Context) ([]*proxy.Proxy, error) {
	return db.proxies, nil
}

func (db *InMemoryProxyDB) StoreProxy(ctx context.Context, proxy *proxy.Proxy) error {
	db.proxies = append(db.proxies, proxy)
	return nil
}

func (db *InMemoryProxyDB) DelProxy(ctx context.Context, proxy *proxy.Proxy) error {
	// Find proxy
	removeIndex := -1
	for i, proxyI := range db.proxies {
		if proxyI == proxy {
			removeIndex = i
			break
		}
	}

	// Remove this element
	if removeIndex != -1 {
		db.proxies[removeIndex] = db.proxies[len(db.proxies)-1] // copy last element to this index
		db.proxies[len(db.proxies)-1] = nil                     // clear last entry
		db.proxies = db.proxies[:len(db.proxies)-1]             // Remove last element
	}

	return nil
}

func (db *InMemoryProxyDB) Clear(ctx context.Context) error {
	db.proxies = []*proxy.Proxy{}
	return nil
}
