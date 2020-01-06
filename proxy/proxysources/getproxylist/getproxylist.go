package getproxylist

import (
	"context"
	"fmt"
	"proxy/proxy"
	"proxy/proxy/proxysources/help"
)

const (
	GetProxyListURL = "https://api.getproxylist.com/proxy"
)

type GetProxyListSource struct {
}

func (p *GetProxyListSource) GetProxy(ctx context.Context) (*proxy.Proxy, error) {
	getProxyListProxy := GetProxyListProxy{}
	resp, err := help.DoRequestObj(ctx, "GET", GetProxyListURL, nil, &getProxyListProxy)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid response code: %d", resp.StatusCode)
	}

	return getProxyListProxy.Standardize(), nil
}

func (p *GetProxyListProxy) Standardize() *proxy.Proxy {
	ret := &proxy.Proxy{
		IP:       p.IP,
		Port:     uint16(p.Port),
		Protocol: proxy.Protocol(p.Protocol),
	}
	return ret
}
