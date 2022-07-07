package garnish_test

import (
	"context"
	"fmt"
	"github.com/bkielbasa/garnish/garnish"
	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"inet.af/netaddr"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestGarnish_CacheRequest(t *testing.T) {
	var listen pan.IPPortValue
	
	err := listen.Set("18-ffaa:1:fc1,147.28.145.13:5000")
	if err != nil {
		return
	}

	stop := mockServer(listen.Get())
	defer stop() //stop will run when we're finished

	expectedXCacheHeaders := []string{garnish.XcacheMiss, garnish.XcacheHit}
	g := garnish.New(url.URL{Scheme: "http", Host: "localhost:8088"})

	// send the request twice to get one miss and one hit
	for _, expectedHeader := range expectedXCacheHeaders {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8088", nil)
		w := httptest.NewRecorder()
		g.ServeHTTP(w, req, listen.String())

		require.Equal(t, http.StatusOK, w.Code)
		xcache := w.Header().Get("X-Cache")
		assert.Equal(t, expectedHeader, xcache)
	}
}

//change into scion
func mockServer(listen netaddr.IPPort) func() {
	conn, err := pan.ListenUDP(context.Background(), listen, nil)
	if err != nil {
		fmt.Printf("listen error")
	}
	defer conn.Close()
	fmt.Println(conn.LocalAddr())
	buffer := make([]byte, 16*1024)
	_, from, err := conn.ReadFrom(buffer)
	if err != nil {
		fmt.Printf("read error")
	}
	msg := fmt.Sprintf("This is the data in server!")
	_, err = conn.WriteTo([]byte(msg), from)
	if err != nil {
		fmt.Printf("write error")
	}
	fmt.Printf("Wrote %d bytes.\n", 1)
	time.Sleep(time.Millisecond * 10)
	return func() {
		panicOnErr(err)
	}
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
