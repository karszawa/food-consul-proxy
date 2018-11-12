package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	target, err := url.Parse(os.Getenv("FOO_LOG_URL"))

	if err != nil {
		log.Fatalf("Could not parse URL: %s", err)
	}

	pxy := httputil.NewSingleHostReverseProxy(target)
	server := http.Server{
		Addr: ":" + os.Getenv("PORT"),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")

			pxy.ServeHTTP(w, r)
		}),
	}

	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("Could not serve: %s", err)
	}
}
