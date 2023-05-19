package datastore

import (
	"reflect"
	"sync"
	"testing"

	"github.com/kasattejaswi/glb/internal/config"
)

func TestLoadRegistry(t *testing.T) {
	t.Run("Test successfully loading registry", func(t *testing.T) {
		// Initializing config with a default value
		conf := config.GetConfig("datastore_test_res/testConfig.yaml")
		registry := LoadRegistry()
		if _, ok := registry.HealthRegistry[conf.BaseConf[0].Hosts[0].UniqueId]; !ok {
			t.Errorf("Expected host to be present in registry with ID: %v", conf.BaseConf[0].Hosts[0].UniqueId)
		}
		if _, ok := registry.HealthRegistry[conf.BaseConf[0].Hosts[1].UniqueId]; !ok {
			t.Errorf("Expected host to be present in registry with ID: %v", conf.BaseConf[0].Hosts[0].UniqueId)
		}
		expectedHost1Health := HealthStatus{
			HostConfig: conf.BaseConf[0].Hosts[0],
		}
		actualHost1Health := registry.HealthRegistry[conf.BaseConf[0].Hosts[0].UniqueId]
		if !reflect.DeepEqual(expectedHost1Health, actualHost1Health) {
			t.Errorf("Unexpected host health state. Expected %v, Actual %v", expectedHost1Health, actualHost1Health)
		}

		expectedHost2Health := HealthStatus{
			HostConfig: conf.BaseConf[0].Hosts[1],
		}
		actualHost2Health := registry.HealthRegistry[conf.BaseConf[0].Hosts[1].UniqueId]
		if !reflect.DeepEqual(expectedHost2Health, actualHost2Health) {
			t.Errorf("Unexpected host health state. Expected %v, Actual %v", expectedHost2Health, actualHost2Health)
		}
	})
}

func TestUpdateHealth(t *testing.T) {
	/* Initializing config with a default value
	For Host 1:
	Min healthy hits: 2 - Host will be marked as healthy after 2 consequtive successful hits at the health endpoint
	Min unhealthy hits: 1 - Host will be marked as unhealthy after 1 failed hit at the health endpoint

	For Host 2:
	Min healthy hits: 5 - Host will be marked as healthy after 5 consequtive successful hits at the health endpoint
	Min unhealthy hits: 6 - Host will be marked as unhealthy after 6 failed hit at the health endpoint
	*/
	conf := config.GetConfig("datastore_test_res/testConfig.yaml")

	t.Run("Host at index 0 should be marked as healthy only after 2 consequtive successful hits", func(t *testing.T) {
		UpdateHealth(conf.BaseConf[0].Hosts[0].UniqueId, true)
		// Host shouldn't be marked as healthy since min healthy hits is 2
		registry := LoadRegistry()
		if registry.HealthRegistry[conf.BaseConf[0].Hosts[0].UniqueId].IsHealthy {
			t.Errorf("Expected host at index 0 to be unhealthy, but found healthy")
		}

		UpdateHealth(conf.BaseConf[0].Hosts[0].UniqueId, true)
		// Host should be marked as healthy now
		registry = LoadRegistry()
		if !registry.HealthRegistry[conf.BaseConf[0].Hosts[0].UniqueId].IsHealthy {
			t.Errorf("Expected host at index 0 to be healthy, but found unhealthy")
		}
	})

	t.Run("Host at index 0 should be marked as unhealthy only after 1 consequtive unsuccessful hit", func(t *testing.T) {
		UpdateHealth(conf.BaseConf[0].Hosts[0].UniqueId, false)
		// Host should be marked as unhealthy
		registry := LoadRegistry()
		if registry.HealthRegistry[conf.BaseConf[0].Hosts[0].UniqueId].IsHealthy {
			t.Errorf("Expected host at index 0 to be unhealthy, but found healthy")
		}
	})

	t.Run("Host at index 1 should be marked as healthy only after 5 consequtive successful hits", func(t *testing.T) {
		UpdateHealth(conf.BaseConf[0].Hosts[1].UniqueId, true)
		// Host shouldn't be marked as healthy since min healthy hits is 2
		registry := LoadRegistry()
		if registry.HealthRegistry[conf.BaseConf[0].Hosts[1].UniqueId].IsHealthy {
			t.Errorf("Expected host at index 1 to be unhealthy, but found healthy")
		}

		UpdateHealth(conf.BaseConf[0].Hosts[1].UniqueId, true)
		// Host shouldn't be marked as healthy since min healthy hits is 5
		registry = LoadRegistry()
		if registry.HealthRegistry[conf.BaseConf[0].Hosts[1].UniqueId].IsHealthy {
			t.Errorf("Expected host at index 1 to be unhealthy, but found healthy")
		}

		UpdateHealth(conf.BaseConf[0].Hosts[1].UniqueId, true)
		// Host shouldn't be marked as healthy since min healthy hits is 5
		registry = LoadRegistry()
		if registry.HealthRegistry[conf.BaseConf[0].Hosts[1].UniqueId].IsHealthy {
			t.Errorf("Expected host at index 1 to be unhealthy, but found healthy")
		}

		UpdateHealth(conf.BaseConf[0].Hosts[1].UniqueId, true)
		// Host shouldn't be marked as healthy since min healthy hits is 5
		registry = LoadRegistry()
		if registry.HealthRegistry[conf.BaseConf[0].Hosts[1].UniqueId].IsHealthy {
			t.Errorf("Expected host at index 1 to be unhealthy, but found healthy")
		}

		UpdateHealth(conf.BaseConf[0].Hosts[1].UniqueId, true)
		// Host should be marked as healthy since min healthy hits is 5
		registry = LoadRegistry()
		if !registry.HealthRegistry[conf.BaseConf[0].Hosts[1].UniqueId].IsHealthy {
			t.Errorf("Expected host at index 1 to be healthy, but found unhealthy")
		}

	})

	t.Run("Host at index 1 should be marked as unhealthy only after 6 consequtive unsuccessful hits", func(t *testing.T) {
		UpdateHealth(conf.BaseConf[0].Hosts[1].UniqueId, false)
		// Host shouldn't be marked as healthy since min healthy hits is 2
		registry := LoadRegistry()
		if !registry.HealthRegistry[conf.BaseConf[0].Hosts[1].UniqueId].IsHealthy {
			t.Errorf("Expected host at index 0 to be healthy, but found unhealthy")
		}

		UpdateHealth(conf.BaseConf[0].Hosts[1].UniqueId, false)
		// Host shouldn't be marked as unhealthy since min unhealthy hits is 6
		registry = LoadRegistry()
		if !registry.HealthRegistry[conf.BaseConf[0].Hosts[1].UniqueId].IsHealthy {
			t.Errorf("Expected host at index 0 to be healthy, but found unhealthy")
		}

		UpdateHealth(conf.BaseConf[0].Hosts[1].UniqueId, false)
		// Host shouldn't be marked as unhealthy since min unhealthy hits is 6
		registry = LoadRegistry()
		if !registry.HealthRegistry[conf.BaseConf[0].Hosts[1].UniqueId].IsHealthy {
			t.Errorf("Expected host at index 0 to be healthy, but found unhealthy")
		}

		UpdateHealth(conf.BaseConf[0].Hosts[1].UniqueId, false)
		// Host shouldn't be marked as unhealthy since min unhealthy hits is 6
		registry = LoadRegistry()
		if !registry.HealthRegistry[conf.BaseConf[0].Hosts[1].UniqueId].IsHealthy {
			t.Errorf("Expected host at index 0 to be healthy, but found unhealthy")
		}

		UpdateHealth(conf.BaseConf[0].Hosts[1].UniqueId, false)
		// Host shouldn't be marked as unhealthy since min unhealthy hits is 6
		registry = LoadRegistry()
		if !registry.HealthRegistry[conf.BaseConf[0].Hosts[1].UniqueId].IsHealthy {
			t.Errorf("Expected host at index 0 to be healthy, but found unhealthy")
		}

		UpdateHealth(conf.BaseConf[0].Hosts[1].UniqueId, false)
		// Host should be marked as unhealthy since min unhealthy hits is 6
		registry = LoadRegistry()
		if registry.HealthRegistry[conf.BaseConf[0].Hosts[1].UniqueId].IsHealthy {
			t.Errorf("Expected host at index 1 to be unhealthy, but found healthy")
		}
	})

	t.Run("Host at index 0 should be marked as healthy only after 2 consequtive successful hits", func(t *testing.T) {
		UpdateHealth(conf.BaseConf[0].Hosts[0].UniqueId, true)
		// Host shouldn't be marked as healthy since min healthy hits is 2
		registry := LoadRegistry()
		if registry.HealthRegistry[conf.BaseConf[0].Hosts[0].UniqueId].IsHealthy {
			t.Errorf("Expected host at index 0 to be unhealthy, but found healthy")
		}

		UpdateHealth(conf.BaseConf[0].Hosts[0].UniqueId, true)
		// Host should be marked as healthy now
		registry = LoadRegistry()
		if !registry.HealthRegistry[conf.BaseConf[0].Hosts[0].UniqueId].IsHealthy {
			t.Errorf("Expected host at index 0 to be healthy, but found unhealthy")
		}
	})

}

func TestUpdateHealthConcurrency(t *testing.T) {
	t.Run("Test registry is thread safe", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(2)
		conf := config.GetConfig("datastore_test_res/testConfig.yaml")
		// Update health status for host at index 1
		go func(wg *sync.WaitGroup, conf *config.Config) {
			UpdateHealth(conf.BaseConf[0].Hosts[0].UniqueId, true)
			wg.Done()
		}(&wg, conf)
		go func(wg *sync.WaitGroup, conf *config.Config) {
			for i := 1; i <= 4; i++ {
				UpdateHealth(conf.BaseConf[0].Hosts[1].UniqueId, true)
			}
			wg.Done()
		}(&wg, conf)
		wg.Wait()
		r := LoadRegistry()
		if r.HealthRegistry[conf.BaseConf[0].Hosts[0].UniqueId].HealthyHitCount != 1 {
			t.Errorf("Expected healthy hit counts to be 1, but got %v", r.HealthRegistry[conf.BaseConf[0].Hosts[0].UniqueId].HealthyHitCount)
		}

		if r.HealthRegistry[conf.BaseConf[0].Hosts[1].UniqueId].HealthyHitCount != 4 {
			t.Errorf("Expected healthy hit counts to be 4, but got %v", r.HealthRegistry[conf.BaseConf[0].Hosts[0].UniqueId].HealthyHitCount)
		}
	})
}

func TestDecideHitEndpoint(t *testing.T) {
	// load config
	conf := config.GetConfig("datastore_test_res/testConfig.yaml")

	t.Run("Should return healthy host at index 1", func(t *testing.T) {
		// Mark endpoint at index 0 as healthy by updating it as healthy twice
		UpdateHealth(conf.BaseConf[0].Hosts[0].UniqueId, true)
		UpdateHealth(conf.BaseConf[0].Hosts[0].UniqueId, true)

		e, ok := DecideHitEndpoint([]string{conf.BaseConf[0].Hosts[0].UniqueId, conf.BaseConf[0].Hosts[1].UniqueId})
		// ok should be true
		if !ok {
			t.Errorf("Expected ok, got not ok")
		}

		if e.UniqueId != conf.BaseConf[0].Hosts[0].UniqueId {
			t.Errorf("Expected host ID be %v, Actual %v", conf.BaseConf[0].Hosts[0].UniqueId, e.UniqueId)
		}
	})

	t.Run("Should return not ok if no host is healthy", func(t *testing.T) {
		// Mark host at index 0 as unhealthy
		UpdateHealth(conf.BaseConf[0].Hosts[0].UniqueId, false)

		_, ok := DecideHitEndpoint([]string{conf.BaseConf[0].Hosts[0].UniqueId, conf.BaseConf[0].Hosts[1].UniqueId})
		// ok should be true
		if ok {
			t.Errorf("Expected not ok, got ok")
		}
	})

	t.Run("Should return second host if first one is unhealthy and second on is healthy", func(t *testing.T) {
		// Mark host at index 1 as healthy
		for i := 1; i <= 5; i++ {
			UpdateHealth(conf.BaseConf[0].Hosts[1].UniqueId, true)
		}
		e, ok := DecideHitEndpoint([]string{conf.BaseConf[0].Hosts[0].UniqueId, conf.BaseConf[0].Hosts[1].UniqueId})
		// ok should be true
		if !ok {
			t.Errorf("Expected ok, got not ok")
		}

		if e.UniqueId != conf.BaseConf[0].Hosts[1].UniqueId {
			t.Errorf("Expected host ID be %v, Actual %v", conf.BaseConf[0].Hosts[1].UniqueId, e.UniqueId)
		}
	})
}
