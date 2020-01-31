package garnish

import "net/http"

type responseWriter struct {
	proxied http.ResponseWriter
	body    []byte
}

func (r *responseWriter) Header() http.Header {
	return r.proxied.Header()
}

func (r *responseWriter) Write(data []byte) (int, error) {
	r.body = data
	return r.proxied.Write(data)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.proxied.WriteHeader(statusCode)
}
