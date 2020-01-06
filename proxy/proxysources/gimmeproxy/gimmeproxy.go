package gimmeproxy

import (
	"context"
	"fmt"
	"proxy/proxy"
	"proxy/proxy/proxysources/help"
	"strconv"
)

const (
	GimmeProxyURL = "https://gimmeproxy.com/api/getProxy"
)

type GimmeProxySource struct {
}

func (p *GimmeProxySource) GetProxy(ctx context.Context) (*proxy.Proxy, error) {
	gimmeProxyProxy := GimmeProxyProxy{}
	resp, err := help.DoRequestObj(ctx, "GET", GimmeProxyURL, nil, &gimmeProxyProxy)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid response code: %d", resp.StatusCode)
	}

	return gimmeProxyProxy.Standardize(), nil
}

func (p *GimmeProxyProxy) Standardize() *proxy.Proxy {
	port, err := strconv.ParseUint(p.Port, 10, 64)
	if err != nil {
		port = 0
	}
	ret := &proxy.Proxy{
		IP:       p.IP,
		Port:     uint16(port),
		Protocol: proxy.Protocol(p.Type),
	}

	return ret
}
