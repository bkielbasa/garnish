package garnish

import (
	"context"
	"fmt"
	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"inet.af/netaddr"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const Xcache = "X-Cache"
const XcacheHit = "HIT"
const XcacheMiss = "MISS"

/**
  Client .-----request---->> Garnish (Cache GET requests) .-------request------>> original Server
      <<-----response-----.                             <<-------response------.
*/

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

func (g *garnish) ServeHTTP(rw http.ResponseWriter, r *http.Request, serverAddress string) {
	// only GET requests should be cached
	// send response back to the client
	// do not need to change
	if r.Method != http.MethodGet {
		rw.Header().Set(Xcache, XcacheMiss)
		g.proxy.ServeHTTP(rw, r)
		return
	}

	u := r.URL.String()
	cached := g.c.get(u)

	//if cached, return the cached data
	// send response back to the client
	// do not need to change
	if cached != nil {
		rw.Header().Set(Xcache, XcacheHit)
		_, _ = rw.Write(cached)
		return
	}

	// if not cached, should ask for the original server
	rw.Header().Set(Xcache, XcacheMiss)

	// change from http connection to scion
	// for a server, the proxy is like a client, so:
	addr, err := pan.ResolveUDPAddr(serverAddress)
	if err != nil {
		fmt.Println("server address error")
		return
	}
	// garnish connect to the server
	conn, err := pan.DialUDP(context.Background(), netaddr.IPPort{}, addr, nil, nil)
	if err != nil {
		fmt.Println("connect to server error")
		return
	}
	defer conn.Close()
	buffer := make([]byte, 16*1024)
	if err = conn.SetReadDeadline(time.Now().Add(1 * time.Second)); err != nil {
		fmt.Println("SetReadDeadline error")
		return
	}
	n, err := conn.Read(buffer) //read to this buffer
	if err != nil {
		fmt.Println("read from server error")
		return
	}
	data := buffer[:n]
	cc := rw.Header().Get(cacheControl)
	toCache, duration := parseCacheControl(cc)
	//check if it needs cache
	if toCache {
		g.c.store(u, data, duration)
	}

}
