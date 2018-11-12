package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
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
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set(
					"Access-Control-Allow-Headers",
					strings.Join(r.Header["Access-Control-Request-Headers"], ","),
				)
				w.Header().Set(
					"Access-Control-Allow-Method",
					strings.Join(r.Header["Access-Control-Request-Method"], ","),
				)
				w.WriteHeader(http.StatusOK)

				return
			}

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")

			var headerKeys []string

			for k := range r.Header {
				headerKeys = append(headerKeys, k)
			}

			w.Header().Set("Access-Control-Expose-Headers", strings.Join(headerKeys, ","))

			pxy.ServeHTTP(w, r)
		}),
	}

	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("Could not serve: %s", err)
	}
}
