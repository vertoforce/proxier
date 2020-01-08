# Proxier

[![Go Report Card](https://goreportcard.com/badge/github.com/vertoforce/proxier)](https://goreportcard.com/report/github.com/vertoforce/proxier)
[![Documentation](https://godoc.org/github.com/vertoforce/proxier?status.svg)](https://godoc.org/github.com/vertoforce/proxier)

Proxier finds proxies and makes requests through them for web crawling.

Proxier helps avoid

- Too many requests (HTTP 429)
- Crawler detection (HTTP 403 or otherwise)

## Usage

### Super Simple Usage

This will make a new proxier automatically finding proxies to use, and then making a request through the proxy.

```go
proxier := New()
resp, _ := proxier.DoRequest(context.Background(), "GET", "https://google.com", nil)
// Do something with resp
```

### Usage with DB

The APIs used to find proxies only allow about 10 requests per day.  So if you want to store the found proxies, you can use a DB.

By default proxier uses an in memory cache. However for my usage I wanted this to work across a horizontally scalable api.  That way each microservice can have access to valid proxies, instead of each instance requesting proxies from the API and quickly using up the quota.
Also if each API was storing valid proxies in memory that would make each API stateful, and I prefer my APIs to be stateless.

You can use your own proxy cache by implementing the `proxy.ProxyDB` interface

This example is with MongoDB

```go
import "github.com/vertoforce/proxier/proxy/proxyDBs/mongodb"

ProxyDB, _ := mongodb.New(context.Background(), MongoDBURL, MongoDB, MongoCollection)
proxier := New().WithProxyDB(ProxyDB)
// Use proxier...
```

## Proxy Sources

- [Gimme Proxy](https://gimmeproxy.com/api/getProxy)
- [GetProxyList](https://api.getproxylist.com/proxy)

## How it works

1. When you create a new Proxier object it contains no proxies.
2. When you make your first request, it will check the ProxyDB (in memory by default) for available proxies.
    - It will loop through the proxies _randomly_ and try each until it finds one that works.
    - If one "succeeds", it will return the `*http.Response` from that proxy
3. If no more proxies exist in the DB it will try to fetch one from a ProxySource
    - It will loop over proxy sources randomly until it finds a new proxy
    - It will try the proxy to see if it "succeeds"
    - If it does, it will save that to the DB and return the `*http.Response`

"succeeds" is in quotes because for a request to succeed it depends on what you're expecting.  By default a proxier object checks for http `200`.  You can pass in your own function to check a response by calling `DoRequestExtra` with your own `CheckResponseFunc`.

## Notes

- Depending on if you cached proxies or how long it takes to find a socks proxy, it could take around 15 seconds to make a request
- This only saves socks5, socks5h, socks4 and socks4a proxies as http proxies don't allow for https

## Known bugs

- If the target server times out, there is undefined behavior, not sure how to handle that yet
