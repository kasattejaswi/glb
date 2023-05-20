package kickstarter

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/kasattejaswi/glb/internal/config"
	"github.com/kasattejaswi/glb/internal/datastore"
	"github.com/kasattejaswi/glb/internal/proxy"
)

// Starts listeners on incoming requests
// Starts go routines for health checks

func KickStartListener(wg *sync.WaitGroup, conf config.BaseConf) {
	defer wg.Done()
	var hcwg sync.WaitGroup
	hcwg.Add(len(conf.Hosts))
	var hostUniqueIds = make([]string, len(conf.Hosts))
	for _, v := range conf.Hosts {
		hostUniqueIds = append(hostUniqueIds, v.UniqueId)
		go KickStartHealthChecks(&hcwg, v)
	}
	log.Printf("Listening for requests at endpoint %v\n", conf.Path)
	config.LoadMux().HandleFunc(conf.Path+"/", proxy.ProxyToServerHandler(conf.Path, hostUniqueIds))
	hcwg.Wait()
}

func KickStartHealthChecks(wg *sync.WaitGroup, conf config.Hosts) {
	defer wg.Done()
	for {
		log.Printf("Running health checks for Hostname: %v, Port: %v, Endpoint: %v\n", conf.Hostname, conf.Port, conf.Health.Endpoint)
		// Create an HTTP request with the specified method and URL
		req, err := http.NewRequest(strings.ToUpper(conf.Health.Method), fmt.Sprintf("%v://%v:%v%v", conf.Protocol, conf.Hostname, conf.Port, conf.Health.Endpoint), nil)
		if err != nil {
			log.Printf("Error creating request: %v", err)
			datastore.UpdateHealth(conf.UniqueId, false)
			time.Sleep(time.Duration(conf.HitFrequencyInSeconds) * time.Second)
			continue
		}

		// Send the HTTP request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Error sending request: %v", err)
			datastore.UpdateHealth(conf.UniqueId, false)
			time.Sleep(time.Duration(conf.HitFrequencyInSeconds) * time.Second)
			continue
		}
		// Retrieve the HTTP status code from the response
		statusCode := resp.StatusCode
		log.Printf("Health response code for Hostname: %v, Port: %v, Endpoint: %v: %v \n", conf.Hostname, conf.Port, conf.Health.Endpoint, statusCode)
		datastore.UpdateHealth(conf.UniqueId, statusCode == conf.Health.SuccessCode)
		time.Sleep(time.Duration(conf.HitFrequencyInSeconds) * time.Second)
	}
}
