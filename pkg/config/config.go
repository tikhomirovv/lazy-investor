package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type InstConf struct {
	Isin string
}
type Config struct {
	LogLevel    string     `yaml:"logLevel"`
	Instruments []InstConf `yaml:"instruments"`
}

func NewConfig(configPath string) (*Config, error) {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}
	var cfg Config
	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &cfg, nil
}
