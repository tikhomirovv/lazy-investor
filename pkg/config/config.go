// Package config provides application configuration loading (YAML).
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// InstConf describes a single instrument entry in config.
type InstConf struct {
	Isin string `yaml:"isin"`
}

// Config holds app-level settings (log level, instrument list).
type Config struct {
	LogLevel    string     `yaml:"logLevel"`
	Instruments []InstConf `yaml:"instruments"`
}

// NewConfig reads and parses the YAML config file at configPath.
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
