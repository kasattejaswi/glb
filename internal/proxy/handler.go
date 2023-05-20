package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kasattejaswi/glb/internal/datastore"
)

// ProxyToServerHandler calls target service and writes response to response writer.
func ProxyToServerHandler(rootEndpoint string, destinationHostsUniqueIds []string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Received request at: %s\n", time.Now())
		destinationHost, ok := datastore.DecideHitEndpoint(destinationHostsUniqueIds)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, "No healthy hosts")
			return
		}
		u, err := url.Parse(fmt.Sprintf("%v://%v:%v%v", destinationHost.Protocol, destinationHost.Hostname, destinationHost.Port, strings.Replace(req.URL.Path, rootEndpoint, "", 1)))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, err)
			return
		}
		req.URL = u
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
