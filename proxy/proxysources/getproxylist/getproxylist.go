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

type Protocol string

// Protocols
const (
	ProtocolHTTP    = "http"
	ProtocolSocks4  = "socks4"
	ProtocolSocks4a = "socks4a"
	ProtocolSocks5  = "socks5"
	ProtocolSocks5h = "socks5h"
)

type GetProxyListSource struct {
}

type GetProxyListProxy struct {
	Links                 Links    `json:"_links"`
	IP                    string   `json:"ip"`
	Port                  int64    `json:"port"`
	Protocol              Protocol `json:"protocol"`
	Anonymity             string   `json:"anonymity"`
	LastTested            string   `json:"lastTested"`
	AllowsRefererHeader   bool     `json:"allowsRefererHeader"`
	AllowsUserAgentHeader bool     `json:"allowsUserAgentHeader"`
	AllowsCustomHeaders   bool     `json:"allowsCustomHeaders"`
	AllowsCookies         bool     `json:"allowsCookies"`
	AllowsPost            bool     `json:"allowsPost"`
	AllowsHTTPS           bool     `json:"allowsHttps"`
	Country               string   `json:"country"`
	ConnectTime           string   `json:"connectTime"`
	DownloadSpeed         string   `json:"downloadSpeed"`
	SecondsToFirstByte    string   `json:"secondsToFirstByte"`
	Uptime                string   `json:"uptime"`
}

type Links struct {
	Self   string `json:"_self"`
	Parent string `json:"_parent"`
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
		Port:     int16(p.Port),
		Protocol: proxy.Protocol(p.Protocol),
	}
	return ret
}
