package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// ProxyToServerHandler calls target service and writes response to response writer.
func ProxyToServerHandler(destinationHost string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Received request at: %s\n", time.Now())

		d, _ := url.Parse(destinationHost)
		req.Host = d.Host
		req.URL.Host = d.Host
		req.URL.Scheme = d.Scheme
		req.RequestURI = ""

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, err)
			return
		}

		for k, v := range res.Header {
			for _, hdrs := range v {
				w.Header().Add(k, hdrs)
			}
		}
		w.WriteHeader(res.StatusCode)
		io.Copy(w, res.Body)
	}
}
