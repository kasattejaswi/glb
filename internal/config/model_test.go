package config

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	expectedConfig := Config{}
	actualConfig := ReadConfig("testconfig.yaml")

}
