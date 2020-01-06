package proxysources

import "testing"

import "proxy/proxy"

import "proxy/proxy/proxysources/gimmeproxy"

import "proxy/proxy/proxysources/getproxylist"

import "context"

var sources = []proxy.ProxySource{
	&gimmeproxy.GimmeProxySource{},
	&getproxylist.GetProxyListSource{},
}

func TestProxySources(t *testing.T) {
	for _, source := range sources {
		proxy, err := source.GetProxy(context.Background())
		if err != nil {
			t.Errorf(err.Error())
			continue
		}

		// Make sure proxy has valid stuff
		if proxy.IP == "" || proxy.Port <= 0 {
			t.Errorf("Invalid proxy got")
		}
	}
}
