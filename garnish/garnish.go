package garnish

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

const Xcache = "X-Cache"
const XcacheHit = "HIT"
const XcacheMiss = "MISS"

type garnish struct {
	c     *cache
	proxy *httputil.ReverseProxy
}

func New(url url.URL) *garnish {
	director := func(req *http.Request) {
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
	}

	reverseProxy := &httputil.ReverseProxy{Director: director}
	return &garnish{c: newCache(), proxy: reverseProxy}
}

func (g *garnish) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// only GET requests should be cached
	if r.Method != http.MethodGet {
		rw.Header().Set(Xcache, XcacheMiss)
		g.proxy.ServeHTTP(rw, r)
		return
	}

	u := r.URL.String()
	cached := g.c.get(u)
	if cached != nil {
		rw.Header().Set(Xcache, XcacheHit)
		_, _ = rw.Write(cached)
		return
	}

	proxyRW := &responseWriter{
		proxied: rw,
	}
	proxyRW.Header().Set(Xcache, XcacheMiss)
	g.proxy.ServeHTTP(proxyRW, r)

	cc := rw.Header().Get(cacheControl)
	toCache, duration := parseCacheControl(cc)
	if toCache {
		g.c.store(u, proxyRW.body, duration)
	}
}
