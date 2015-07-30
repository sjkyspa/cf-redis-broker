package metrics

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	BrokerEndpoint string `yaml:"broker_endpoint"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("config: Failed to open file: %s", err)
	}

	var config Config

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("config: Invalid file: %s", err)
	}

	if config.BrokerEndpoint == "" {
		return nil, fmt.Errorf("config: Missing BrokerEndpoint property")
	}
	return &config, nil
}
