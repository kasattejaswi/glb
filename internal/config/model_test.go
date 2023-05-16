package config

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestReadConfig(t *testing.T) {
	t.Run("Should successfully read config file and store in struct", func(t *testing.T) {
		expectedConfig := Config{
			BaseConf: []BaseConf{
				{
					UniqueId:  reflect.Array.String(),
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
			},
		}
		actualConfig := GetConfig("model_test_res/testConfig.yaml")
		expectedConfig.BaseConf[0].UniqueId = actualConfig.BaseConf[0].UniqueId
		expectedConfig.BaseConf[0].Hosts[0].UniqueId = actualConfig.BaseConf[0].Hosts[0].UniqueId
		expectedConfig.BaseConf[0].Hosts[1].UniqueId = actualConfig.BaseConf[0].Hosts[1].UniqueId
		if !reflect.DeepEqual(&expectedConfig, actualConfig) {
			t.Errorf("Expected and actual config didn't match. Expected: %v, Actual %v", expectedConfig, actualConfig)
		}
		for _, v := range actualConfig.BaseConf {
			if v.UniqueId == "" {
				t.Errorf("Expected an assigned unique ID.")
				for _, h := range v.Hosts {
					if h.UniqueId == "" {
						t.Errorf("Expected an assigned unique ID")
					}
				}
			}
		}
	})

	t.Run("Should panic if file is not present", func(t *testing.T) {

		once = sync.Once{}

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected function to panic")
			}
		}()
		gen := GetConfig("nonexistingfile.yaml")
		fmt.Println(gen)
	})

	t.Run("Should panic if invalid config is provided", func(t *testing.T) {

		once = sync.Once{}

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected function to panic")
			}
		}()
		GetConfig("model_test_res/invalidConfig.yaml")
	})
}

func TestGetUniqueIDByPath(t *testing.T) {
	expectedUniqueID := GetConfig("model_test_res/testConfig.yaml").BaseConf[0].UniqueId
	actualUniqueId := GetUniqueIDByPath("/a")
	if expectedUniqueID != actualUniqueId {
		t.Errorf("Expected Unique ID to be %v, Actual %v\n", expectedUniqueID, actualUniqueId)
	}

	// should return empty string if path is not present
	if GetUniqueIDByPath("/b") != "" {
		t.Errorf("Expected Unique ID to be %v, Actual %v\n", expectedUniqueID, actualUniqueId)
	}

}

func TestGetUniqueIdsOfHostsByPath(t *testing.T) {
	baseConf := GetConfig("model_test_res/testConfig.yaml").BaseConf[0]
	expectedHostsUniqueIds := []string{baseConf.Hosts[0].UniqueId, baseConf.Hosts[1].UniqueId}
	actualHostsUniqueIds := GetUniqueIdsOfHostsByPath("/a")
	if !reflect.DeepEqual(expectedHostsUniqueIds, actualHostsUniqueIds) {
		t.Errorf("Expected Unique ID to be %v, Actual %v\n", expectedHostsUniqueIds, actualHostsUniqueIds)
	}
}
