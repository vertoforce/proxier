package proxydbs

import (
	"context"
	"proxy/proxy"
	"proxy/proxy/proxyDBs/inmemory"
	"proxy/proxy/proxyDBs/mongodb"
	"reflect"
	"testing"
)

const (
	MongoDBURL      = "mongodb://root:pass@localhost"
	MongoDB         = "proxies_test"
	MongoCollection = "proxies"
)

var testingProxies = []*proxy.Proxy{
	&proxy.Proxy{
		IP:       "1.1.1.1",
		Port:     1000,
		Protocol: proxy.Socks5Protocol,
	},
	&proxy.Proxy{
		IP:       "1.1.1.2",
		Port:     1000,
		Protocol: proxy.Socks5Protocol,
	},
	&proxy.Proxy{
		IP:       "1.1.1.3",
		Port:     1000,
		Protocol: proxy.Socks5Protocol,
	},
}

func TestDBs(t *testing.T) {
	// Init DBs
	var dbs = []proxy.ProxyDB{}
	mongodb, err := mongodb.New(context.Background(), MongoDBURL, MongoDB, MongoCollection)
	if err != nil {
		t.Errorf("Error creating mongodb")
		return
	}
	dbs = append(dbs, inmemory.New())
	dbs = append(dbs, mongodb)

	// Test DBs
	for _, proxydb := range dbs {
		ctx := context.Background()
		// Clear Database
		err = proxydb.Clear(ctx)
		if err != nil {
			t.Error(err)
		}

		// Test GetProxies
		if proxies, err := proxydb.GetProxies(ctx); err == nil {
			if len(proxies) != 0 {
				t.Errorf("Database has entries")
			}
		} else {
			t.Errorf(err.Error())
			continue
		}

		// Store some entries
		for _, proxy := range testingProxies {
			if err := proxydb.StoreProxy(ctx, proxy); err != nil {
				t.Errorf(err.Error())
			}
		}

		// Make sure it stored
		if proxies, err := proxydb.GetProxies(ctx); err == nil {
			if len(proxies) != 3 {
				t.Errorf("Database did not store proxies")
			}
			if !reflect.DeepEqual(proxies, testingProxies) {
				t.Errorf("Incorrect proxies stored")
			}
		} else {
			t.Errorf(err.Error())
			continue
		}

		// Test delete
		if err := proxydb.DelProxy(ctx, testingProxies[0]); err != nil {
			t.Errorf(err.Error())
		}

		// Make sure it removed that specific one
		if proxies, err := proxydb.GetProxies(ctx); err == nil {
			if len(proxies) != 2 {
				t.Errorf("Database has wrong number of entries")
			}
			// Check other entries to see if it has the deleted one
			for _, proxy := range proxies {
				if proxy.IP == testingProxies[0].IP {
					t.Errorf("did not delete proxy")
				}
			}
		} else {
			t.Errorf(err.Error())
			continue
		}
	}
}
