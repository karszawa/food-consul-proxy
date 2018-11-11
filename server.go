// package main

// import (
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strings"
// )

// func main() {
// 	fooLogHost := os.Getenv("FOO_LOG_HOST")
// 	client := &http.Client{}
// 	mux := http.NewServeMux()

// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		r.URL.Scheme = "https"
// 		r.URL.Host = fooLogHost

// 		req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)

// 		if err != nil {
// 			log.Printf("Failed to new request: %v", err)
// 			return
// 		}

// 		req.Header = r.Header
// 		resp, err := client.Do(req)

// 		if err != nil {
// 			log.Printf("Failed to do request: %v", err)
// 			return
// 		}

// 		w.WriteHeader(resp.StatusCode)

// 		for k, v := range resp.Header {
// 			w.Header().Set(k, strings.Join(v, ","))
// 		}

// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")

// 		body, err := ioutil.ReadAll(resp.Body)

// 		if err != nil {
// 			log.Printf("Failed to read body: %v", err)
// 			return
// 		}

// 		w.Write(body)
// 	})

// 	http.ListenAndServe(":"+os.Getenv("PORT"), mux)
// }

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
			pxy.ServeHTTP(w, r)
		}),
	}

	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("Could not serve: %s", err)
	}
}
