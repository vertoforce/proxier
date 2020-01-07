# Proxier

Proxier finds proxies and makes requests through them for web crawling.

Proxier helps avoid

- Too many requests (HTTP 429)
- Crawler detection (HTTP 403 or otherwise)

## Usage

### Super Simple Usage

```go
proxier := New()
resp, _ := proxier.DoRequest(context.Background(), "GET", "https://google.com", nil)
// Do something with resp
```

### Usage with DB

The APIs used to find proxies only allow about 10 requests per day.  So if you want to store the found proxies, you can use a DB.

By default it uses an in memory cache. However for my usage I wanted this to work across a horizontally scalable api.  That way each microservice can have access to valid proxies, instead of each instance requesting proxies from the API and quickly using up the quota.
Also if each API was storing valid proxies in memory that would make each API stateful, and I prefer my APIs to be stateless.

You can use your own proxy cache by implementing the `proxy.ProxyDB` interface

This example is with MongoDB

```go
import "github.com/vertoforce/proxier/proxy/proxyDBs/mongodb"

mongoDBProxyDB, _ := mongodb.New(context.Background(), MongoDBURL, MongoDB, MongoCollection)
proxier := New().WithProxyDB(mongoDBProxyDB)
// Use proxier...
```

## Notes

- This only saves socks5, socks5h, socks4 and socks4a proxies as http proxies don't allow for https
