package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pivotal-cf/brokerapi"
)

type Config struct {
	ListeningAddr      string                                `json:"listeningAddr"`
	LogLevel           string                                `json:"logLevel"`
	Credentials        brokerapi.BrokerCredentials           `json:"brokerCredentials"`
	ServiceCatalog     []brokerapi.Service                   `json:"serviceCatalog"`
	ServiceInstances   map[string]brokerapi.ProvisionDetails `json:"provisionDetails"`
	BindingCredentials map[string]interface{}                `json:"bindingCredentials"`
}

// LoadFromFile loads config file from provided file path location
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
	return parseJSON(jsonConf)
}

func parseJSON(jsonConf []byte) (*Config, error) {
	config := &Config{}

	err := json.Unmarshal(jsonConf, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
