package proxier

import (
	"context"
	"fmt"
	"proxy/proxy/proxyDBs/mongodb"
	"testing"
)

const (
	TestURL         = "https://www.meetup.com/Huey-Spheres-by-GoHuey-com/events/kjwzvqyzqbgc"
	MongoDBURL      = "mongodb://root:pass@localhost"
	MongoDB         = "proxies"
	MongoCollection = "proxies"
)

func TestDoRequest(t *testing.T) {
	mongodbProxyDB, err := mongodb.New(context.Background(), MongoDBURL, MongoDB, MongoCollection)
	if err != nil {
		t.Error(err)
		return
	}
	proxier := New().WithProxyDB(mongodbProxyDB)
	resp, err := proxier.DoRequest(context.Background(), "GET", TestURL, nil)
	if err != nil {
		t.Error(err)
		return
	}
	// bytes, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// ioutil.WriteFile("out.htm", bytes, 0644)
	fmt.Println(resp)
}

func TestCacheProxies(t *testing.T) {
	mongodbProxyDB, err := mongodb.New(context.Background(), MongoDBURL, MongoDB, MongoCollection)
	if err != nil {
		t.Error(err)
		return
	}
	proxier := New().WithProxyDB(mongodbProxyDB)
	added, err := proxier.CacheProxies(context.Background(), 10)
	if err != nil {
		t.Error(err)
	}

	if added == 0 {
		t.Errorf("Didn't add any proxies")
	}
}
