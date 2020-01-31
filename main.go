package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bkielbasa/garnish/garnish"
)

func main() {
	u, err := url.Parse("http://localhost:8080")
	if err != nil {
		panic(err)
	}
	port := flag.Int("p", 80, "port")
	flag.Parse()

	g := garnish.New(*u)

	if *port == 443 {
		http.ListenAndServeTLS(fmt.Sprintf(":%d", *port), "localhost.pem", "localhost-key.pem", g)
	} else {
		http.ListenAndServe(fmt.Sprintf(":%d", *port), g)
	}
}
