package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/kasattejaswi/glb/internal/config"
	"github.com/kasattejaswi/glb/internal/kickstarter"
)

func main() {
	log.Println("Starting load balancer ...")
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to read user's home directory")
	}
	log.Printf("Reading default configuration at path %v\n", filepath.Join(homedir, ".glb", "config.yaml"))
	config := config.GetConfig(filepath.Join(homedir, ".glb", "config.yaml"))
	var wg sync.WaitGroup
	for _, v := range *config {
		wg.Add(1)
		go kickstarter.KickStartListener(&wg, v)
	}
	wg.Wait()
}
