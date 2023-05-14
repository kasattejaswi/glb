package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config []struct {
	Path      string  `yaml:"path"`
	Algorithm string  `yaml:"algorithm"`
	Sticky    bool    `yaml:"sticky"`
	Hosts     []Hosts `yaml:"hosts"`
}

type Hosts struct {
	Protocol              string `yaml:"protocol"`
	Hostname              string `yaml:"hostname"`
	Port                  int    `yaml:"port"`
	Health                Health `yaml:"health"`
	MinHealthyHits        int    `yaml:"minHealthyHits"`
	MinUnhealthyHits      int    `yaml:"minUnhealthyHits"`
	HitFrequencyInSeconds int    `yaml:"hitFrequencyInSeconds"`
}

type Health struct {
	Endpoint    string `yaml:"endpoint"`
	SuccessCode int    `yaml:"successCode"`
	Method      string `yaml:"method"`
}

func ReadConfig(path string) Config {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}
	return config
}
