package config

import (
	"encoding/json"
	"github.com/pivotal-cf/brokerapi"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	ListeningAddr      string                              `json:"listeningAddr"`
	LogLevel           string                              `json:"logLevel"`
	Credentials        brokerapi.BrokerCredentials         `json:"brokerCredentials"`
	ServiceCatalog     []brokerapi.Service                 `json:"serviceCatalog"`
	ServiceInstances   map[string]brokerapi.ServiceDetails `json:"serviceInstances"`
	BindingCredentials map[string]interface{}              `json:"bindingCredentials"`
}

func LoadFromFile(path string) (*Config, error) {
	if len(path) == 0 {
		binDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return nil, err
		}
		path = filepath.Join(binDir, "cf-config-broker.json")
	}
	jsonConf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseJson(jsonConf)
}

func ParseJson(jsonConf []byte) (*Config, error) {
	config := &Config{}

	err := json.Unmarshal(jsonConf, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
