package garnish_test

import (
	"net/http"
	"time"
)

//change into scion
func mockServer() func() {
	m := http.NewServeMux()
	s := http.Server{Addr: ":8088", Handler: m}
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=100")
		_, _ = w.Write([]byte("OK"))
	})

	go func() {
		_ = s.ListenAndServe()
	}()

	time.Sleep(time.Millisecond * 10)

	return func() {
		panicOnErr(s.Close())
	}
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
