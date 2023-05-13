package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host         string   `yaml:"host"`
	Spiders      []int    `yaml:"spiders"`
	Taunts       []string `yaml:"taunts"`
	Whitelist    []string `yaml:"whitelist"`
	AttackLength int      `yaml:"attack_length"`
}

func LoadConfig(file string) (*Config, error) {
	file, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}
	if !fileExists(file) {
		return nil, fmt.Errorf("file does not exist")
	}
	c := &Config{
		Host:      "",
		Spiders:   []int{},
		Taunts:    []string{},
		Whitelist: []string{},
	}
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, c)
	return c, err
}

//go:embed services.yaml
var servicesFile string

type Services map[string]struct {
	Ports  []int  `yaml:"ports"`
	Banner string `yaml:"banner"`
}

func LoadServices() (*Services, error) {
	s := &Services{}
	err := yaml.Unmarshal([]byte(servicesFile), s)
	return s, err
}
