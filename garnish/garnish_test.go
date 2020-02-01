package garnish_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bkielbasa/garnish/garnish"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGarnish_CacheRequest(t *testing.T) {
	stop := mockServer()
	defer stop()

	expectedXCacheHeaders := []string{garnish.XcacheMiss, garnish.XcacheHit}
	g := garnish.New(url.URL{Scheme: "http", Host: "localhost:8080"})

	for _, expectedHeader := range expectedXCacheHeaders {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		w := httptest.NewRecorder()
		g.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		xcache := w.Header().Get("X-Cache")
		assert.Equal(t, expectedHeader, xcache)
	}
}

func BenchmarkGarnish_ServeHTTP(b *testing.B) {
	stop := mockServer()
	defer stop()
	g := garnish.New(url.URL{Scheme: "http", Host: "localhost:8080"})
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	w := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
		g.ServeHTTP(w, req)
	}
}

func mockServer() func() {
	m := http.NewServeMux()
	s := http.Server{Addr: ":8080", Handler: m}
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	go func() {
		_ = s.ListenAndServe()
	}()

	return func() {
		_ = s.Close()
	}
}
