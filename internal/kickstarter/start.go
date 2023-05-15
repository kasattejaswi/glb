package kickstarter

import (
	"log"
	"sync"
	"time"

	"github.com/kasattejaswi/glb/internal/config"
)

// Starts listeners on incoming requests
// Starts go routines for health checks

func KickStartListener(wg *sync.WaitGroup, conf config.BaseConf) {
	defer wg.Done()
	var hcwg sync.WaitGroup
	hcwg.Add(len(conf.Hosts))
	for _, v := range conf.Hosts {
		go KickStartHealthChecks(&hcwg, v)
	}
	log.Printf("Listening for requests at endpoint %v\n", conf.Path)
	hcwg.Wait()
}

func KickStartHealthChecks(wg *sync.WaitGroup, conf config.Hosts) {
	defer wg.Done()
	for {
		log.Printf("Running health checks for Hostname: %v, Port: %v, Endpoint: %v\n", conf.Hostname, conf.Port, conf.Health.Endpoint)
		time.Sleep(time.Duration(conf.HitFrequencyInSeconds) * time.Second)
	}
}
