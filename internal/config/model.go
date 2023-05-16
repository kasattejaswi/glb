package config

import (
	"io/ioutil"
	"sync"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type Config struct {
	BaseConf []BaseConf
}

type BaseConf struct {
	UniqueId  string
	Path      string  `yaml:"path"`
	Algorithm string  `yaml:"algorithm"`
	Sticky    bool    `yaml:"sticky"`
	Hosts     []Hosts `yaml:"hosts"`
}

type Hosts struct {
	UniqueId              string
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

var conf Config
var once sync.Once

// ReadConfig reads load balancer configuration at the specified path
func GetConfig(path string) *Config {
	once.Do(func() {
		file, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		var baseConf []BaseConf
		err = yaml.Unmarshal(file, &baseConf)
		if err != nil {
			panic(err)
		}
		for i, v := range baseConf {

			baseConf[i].UniqueId = uuid.NewString()
			for j := range v.Hosts {
				baseConf[i].Hosts[j].UniqueId = uuid.NewString()
			}
		}
		conf = Config{
			BaseConf: baseConf,
		}
	})
	return &conf
}

func GetUniqueIDByPath(path string) string {
	for _, v := range GetConfig("").BaseConf {
		if v.Path == path {
			return v.UniqueId
		}
	}
	return ""
}

func GetUniqueIdsOfHostsByPath(path string) []string {
	result := []string{}
	for _, v := range GetConfig("").BaseConf {
		if v.Path == path {
			for _, h := range v.Hosts {
				result = append(result, h.UniqueId)
			}
			break
		}
	}
	return result
}
