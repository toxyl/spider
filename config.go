package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host         string   `yaml:"host"`
	Ports        []int    `yaml:"ports"`
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
		Ports:     []int{},
		Taunts:    []string{},
		Whitelist: []string{},
	}
	b, err := os.ReadFile(file)
	yaml.Unmarshal(b, c)
	return c, err
}
