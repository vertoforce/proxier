package gimmeproxy

import (
	"context"
	"proxy/proxy"
)

const (
	GimmeProxyURL = "https://gimmeproxy.com/api/getProxy"
)

type GimmeProxySource struct {
}

func (p *GimmeProxySource) GetProxy(ctx context.Context) (*proxy.Proxy, error) {
	return &proxy.Proxy{}, nil
}
