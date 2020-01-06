package getproxylist

import (
	"context"
	"encoding/json"
	"net/http"
	"proxy/proxy"
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
	req, err := http.NewRequest("GET", GetProxyListURL, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	getProxyListProxy := &GetProxyListProxy{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(getProxyListProxy)
	if err != nil {
		return nil, err
	}

	return getProxyListProxy.Standardize(), nil
}

func (p *GetProxyListProxy) Standardize() *proxy.Proxy {
	// TODO: Fill
	ret := &proxy.Proxy{}
	return ret
}
