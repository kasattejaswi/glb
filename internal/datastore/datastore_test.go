package datastore

import (
	"reflect"
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
	type args struct {
		id        string
		isHealthy bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UpdateHealth(tt.args.id, tt.args.isHealthy)
		})
	}
}

func TestDecideHitEndpoint(t *testing.T) {
	type args struct {
		hostUniqueIds []string
	}
	tests := []struct {
		name           string
		args           args
		wantHostConfig config.Hosts
		wantOk         bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHostConfig, gotOk := DecideHitEndpoint(tt.args.hostUniqueIds)
			if !reflect.DeepEqual(gotHostConfig, tt.wantHostConfig) {
				t.Errorf("DecideHitEndpoint() gotHostConfig = %v, want %v", gotHostConfig, tt.wantHostConfig)
			}
			if gotOk != tt.wantOk {
				t.Errorf("DecideHitEndpoint() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
