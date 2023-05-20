package datastore

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/kasattejaswi/glb/internal/config"
)

// Datastore is responsible for in-memory data
// All data is retained until the service is running. On shutdown, the data will be lost.
// The data that load balancer stores does not need persistence.

// Registry stores live health data of available services. The key is the unique ID of each host generated at load balancer startup.
type Registry struct {
	// RWMutex makes it thread safe since multiple Go Routines will be accessing it at the same time
	sync.RWMutex
	HealthRegistry map[string]HealthStatus `json:"registry"`
}
type HealthStatus struct {
	HostConfig        config.Hosts `json:"hostConfig"`
	LastChecked       time.Time    `json:"lastChecked"`
	IsHealthy         bool         `json:"isHealthy"`
	HealthyHitCount   int          `json:"healthyHitCount"`
	UnhealthyHitCount int          `json:"unhealthyHitCount"`
	LastHitAt         time.Time    `json:"lastHitAt"`
}

var registry Registry
var once sync.Once

func LoadRegistry() *Registry {
	once.Do(func() {
		registry = Registry{
			HealthRegistry: make(map[string]HealthStatus),
		}
		conf := config.GetConfig("")
		for _, v := range conf.BaseConf {
			for _, h := range v.Hosts {
				registry.HealthRegistry[h.UniqueId] = HealthStatus{
					HostConfig: h,
				}
			}
		}
	})
	return &registry
}

func UpdateHealth(id string, isHealthy bool) {
	r := LoadRegistry()
	r.Lock()
	v := r.HealthRegistry[id]

	v.LastChecked = time.Now()
	if isHealthy {
		if v.HealthyHitCount+1 == v.HostConfig.MinHealthyHits {
			v.IsHealthy = true
			v.HealthyHitCount = 0
			v.UnhealthyHitCount = 0

		}
		v.HealthyHitCount++
	} else {
		if v.UnhealthyHitCount+1 == v.HostConfig.MinUnhealthyHits {
			v.IsHealthy = false
			v.HealthyHitCount = 0
			v.UnhealthyHitCount = 0
		}
		v.UnhealthyHitCount++
	}
	r.HealthRegistry[id] = v
	r.Unlock()
}

func DecideHitEndpoint(hostUniqueIds []string) (hostConfig config.Hosts, ok bool) {
	r := LoadRegistry()
	r.RLock()
	defer r.RUnlock()
	for _, v := range hostUniqueIds {
		if r.HealthRegistry[v].IsHealthy {
			return r.HealthRegistry[v].HostConfig, true
		}
	}
	return config.Hosts{}, false
}

func LoadRegistryEndpoints() {
	config.LoadMux().Handle("/glb/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		registry := LoadRegistry()
		b, e := json.Marshal(registry)
		if e != nil {
			fmt.Fprintf(w, "error loading registry: %v", e)
		}
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(b))
	}))
}
