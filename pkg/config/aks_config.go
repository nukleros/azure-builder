package config

import (
	"fmt"
)

type AksConfig struct {
	Name   *string `yaml:"name"`
	Region *string `yaml:"region"`
}

func (config *AksConfig) ValidateNotNull() error {
	if config.Name == nil {
		return fmt.Errorf("could not find name in aks config")
	}

	if config.Region == nil {
		return fmt.Errorf("could not find region in aks config")
	}

	return nil
}
