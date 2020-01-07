package proxysources

import (
	"context"
	"testing"

	"github.com/vertoforce/proxier/proxy"
	"github.com/vertoforce/proxier/proxy/proxysources/getproxylist"
	"github.com/vertoforce/proxier/proxy/proxysources/gimmeproxy"
)

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
