package proxydbs

import (
	"context"
	"proxy/proxy"
	"proxy/proxy/proxyDBs/mongodb"
	"reflect"
	"testing"
)

var testingProxy = proxy.Proxy{
	IP:       "1.1.1.1",
	Port:     1000,
	Protocol: proxy.Socks5,
}

func TestDBs(t *testing.T) {
	// Init DBs
	var dbs = []proxy.ProxyDB{}
	mongodb, err := mongodb.New(context.Background(), "mongodb://root:pass@localhost:27017", "proxies", "proxies")
	if err != nil {
		t.Errorf("Error creating mongodb")
	}
	dbs = append(dbs, mongodb)

	for _, proxydb := range dbs {
		// Test GetProxies
		if proxies, err := proxydb.GetProxies(context.Background()); err == nil {
			if len(proxies) != 0 {
				t.Errorf("Database has entries")
			}
		} else {
			t.Errorf(err.Error())
			continue
		}

		// Store
		if err := proxydb.StoreProxy(context.Background(), testingProxy); err != nil {
			t.Errorf(err.Error())
		}

		// Make sure it stored
		if proxies, err := proxydb.GetProxies(context.Background()); err == nil {
			if len(proxies) != 1 {
				t.Errorf("Database did not store proxy")
			}
			if !reflect.DeepEqual(proxies[0], testingProxy) {
				t.Errorf("Incorrect proxy stored")
			}
		} else {
			t.Errorf(err.Error())
			continue
		}

		// Test delete
		if err := proxydb.DelProxy(context.Background(), testingProxy); err != nil {
			t.Errorf(err.Error())
		}

		// Make sure its empty
		if proxies, err := proxydb.GetProxies(context.Background()); err == nil {
			if len(proxies) != 0 {
				t.Errorf("Database has entries")
			}
		} else {
			t.Errorf(err.Error())
			continue
		}
	}
}
