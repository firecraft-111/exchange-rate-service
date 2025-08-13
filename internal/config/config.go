package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`

	ExchangeRate struct {
		APIKey string `yaml:"api_key"`
	} `yaml:"exchange_rate"`
}

var App *Config

func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	App = &Config{}
	if err := yaml.Unmarshal(data, App); err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}

	return nil
}
