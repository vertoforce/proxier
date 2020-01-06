package gimmeproxy

type GimmeProxyProxy struct {
	SupportsHTTPS      bool        `json:"supportsHttps"`
	Protocol           string      `json:"protocol"`
	IP                 string      `json:"ip"`
	Port               string      `json:"port"`
	Get                bool        `json:"get"`
	Post               bool        `json:"post"`
	Cookies            bool        `json:"cookies"`
	Referer            bool        `json:"referer"`
	UserAgent          bool        `json:"user-agent"`
	AnonymityLevel     int64       `json:"anonymityLevel"`
	Websites           interface{} `json:"websites"`
	Country            string      `json:"country"`
	UnixTimestampMS    int64       `json:"unixTimestampMs"`
	TsChecked          int64       `json:"tsChecked"`
	UnixTimestamp      int64       `json:"unixTimestamp"`
	Curl               string      `json:"curl"`
	IPPort             string      `json:"ipPort"`
	Type               string      `json:"type"`
	Speed              float64     `json:"speed"`
	OtherProtocols     interface{} `json:"otherProtocols"`
	VerifiedSecondsAgo int64       `json:"verifiedSecondsAgo"`
}

// type Websites struct {
// 	Example    bool `json:"example"`
// 	Google     bool `json:"google"`
// 	Amazon     bool `json:"amazon"`
// 	Yelp       bool `json:"yelp"`
// 	GoogleMaps bool `json:"google_maps"`
// }
