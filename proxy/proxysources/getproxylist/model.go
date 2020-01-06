package getproxylist

type Protocol string

// Protocols
const (
	ProtocolHTTP    = "http"
	ProtocolSocks4  = "socks4"
	ProtocolSocks4a = "socks4a"
	ProtocolSocks5  = "socks5"
	ProtocolSocks5h = "socks5h"
)

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
