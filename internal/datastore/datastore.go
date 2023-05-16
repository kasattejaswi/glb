package datastore

import (
	"sync"
	"time"

	"github.com/kasattejaswi/glb/internal/config"
)

// Datastore is responsible for in-memory data
// All data is retained until the service is running. On shutdown, the data will be lost.
// The data that load balancer stores does not need persistence.

// IsHealthyRegistry stores live health data of available services. The key is the unique ID of each host generated at load balancer startup.
type Registry struct {
	sync.RWMutex
	HealthRegistry map[string]HealthStatus
}
type HealthStatus struct {
	HostConfig        config.Hosts
	LastChecked       time.Time
	IsHealthy         bool
	HealthyHitCount   int
	UnhealthyHitCount int
	LastHitAt         time.Time
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
