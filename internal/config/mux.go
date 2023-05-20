package config

import (
	"net/http"
	"sync"
)

var mux *http.ServeMux
var onceMux sync.Once

func LoadMux() *http.ServeMux {
	onceMux.Do(func() {
		mux = http.NewServeMux()
	})
	return mux
}
