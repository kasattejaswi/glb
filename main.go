package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/kasattejaswi/glb/internal/config"
	"github.com/kasattejaswi/glb/internal/datastore"
	"github.com/kasattejaswi/glb/internal/kickstarter"
)

func main() {
	log.Println("Starting load balancer ...")
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to read user's home directory")
	}
	log.Printf("Reading default configuration at path %v\n", filepath.Join(homedir, ".glb", "config.yaml"))
	conf := config.GetConfig(filepath.Join(homedir, ".glb", "config.yaml"))
	datastore.LoadRegistry()
	var wg sync.WaitGroup
	for _, v := range conf.BaseConf {
		wg.Add(1)
		go kickstarter.KickStartListener(&wg, v)
	}
	datastore.LoadRegistryEndpoints()
	if err := http.ListenAndServe(":8000", config.LoadMux()); err != nil {
		log.Fatal(err)
	}
	wg.Wait()
}
