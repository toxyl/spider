package main

import (
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
	MetricNexus  struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		Key  string `yaml:"key"`
	} `yaml:"metric_nexus"`
}

func LoadConfig(file string) error {
	file, err := filepath.Abs(file)
	if err != nil {
		return err
	}
	if !fileExists(file) {
		return fmt.Errorf("file does not exist")
	}
	c := &Config{
		Host:      "",
		Spiders:   []int{},
		Taunts:    []string{},
		Whitelist: []string{},
	}
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return err
	}
	config = c
	return nil
}
