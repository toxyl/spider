package main

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed services.yaml
var servicesFile string

type Services map[string]struct {
	Ports  []int  `yaml:"ports"`
	Banner string `yaml:"banner"`
}

func LoadServices() error {
	s := &Services{}
	err := yaml.Unmarshal([]byte(servicesFile), s)
	if err != nil {
		return err
	}
	services = s
	return nil
}
