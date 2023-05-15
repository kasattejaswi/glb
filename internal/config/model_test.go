package config

import (
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	t.Run("Should successfully read config file and store in struct", func(t *testing.T) {
		expectedConfig := Config{
			{
				Path:      "/a",
				Algorithm: "roundRobin",
				Sticky:    false,
				Hosts: []Hosts{
					{
						Protocol:              "http",
						Hostname:              "localhost",
						Port:                  8080,
						MinHealthyHits:        5,
						MinUnhealthyHits:      6,
						HitFrequencyInSeconds: 20,
						Health: Health{
							Endpoint:    "/health",
							SuccessCode: 200,
							Method:      "GET",
						},
					},
					{
						Protocol:              "http",
						Hostname:              "localhost",
						Port:                  9090,
						MinHealthyHits:        5,
						MinUnhealthyHits:      6,
						HitFrequencyInSeconds: 20,
						Health: Health{
							Endpoint:    "/health",
							SuccessCode: 200,
							Method:      "GET",
						},
					},
				},
			},
		}
		actualConfig := ReadConfig("model_test_res/testConfig.yaml")
		if !reflect.DeepEqual(expectedConfig, actualConfig) {
			t.Errorf("Expected and actual config didn't match. Expected: %v, Actual %v", expectedConfig, actualConfig)
		}
	})
	t.Run("Should panic if file is not present", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected function to panic")
			}
		}()
		ReadConfig("nonexistingfile.yaml")
	})
	t.Run("Should panic if invalid config is provided", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected function to panic")
			}
		}()
		ReadConfig("model_test_res/invalidConfig.yaml")
	})
}
