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

// SchedulerConf holds Stage 0 scheduler settings (interval in seconds, run on start).
type SchedulerConf struct {
	// IntervalSeconds is how often to run the report pipeline in seconds (e.g. 3600 = hourly).
	IntervalSeconds int `yaml:"intervalSeconds"`
	// RunOnStart runs the pipeline once immediately after start.
	RunOnStart bool `yaml:"runOnStart"`
}

// CandlesConf holds candle-fetch settings for Stage 0.
type CandlesConf struct {
	// LookbackDays is how many days of history to fetch (e.g. 30).
	LookbackDays int `yaml:"lookbackDays"`
	// Interval is candle interval; Stage 0 uses daily ("1d") only.
	Interval string `yaml:"interval"`
}

// TelegramConf holds Telegram report settings. Token and ChatID come from env (TELEGRAM_BOT_TOKEN, TELEGRAM_CHAT_ID).
type TelegramConf struct {
	// Enabled turns on sending reports to Telegram. If false or env not set, reports are only logged.
	Enabled bool `yaml:"enabled"`
	// HandleCommands enables receiving updates and handling /candles etc. When true, a goroutine runs long polling.
	HandleCommands bool `yaml:"handleCommands"`
	// AllowedChatID if set (numeric) restricts command handling to this chat only; empty = respond to any chat.
	AllowedChatID string `yaml:"allowedChatID"`
}

// Config holds app-level settings (log level, instrument list, Stage 0 sections).
type Config struct {
	LogLevel    string        `yaml:"logLevel"`
	Instruments []InstConf    `yaml:"instruments"`
	Scheduler   SchedulerConf `yaml:"scheduler"`
	Candles     CandlesConf   `yaml:"candles"`
	Telegram    TelegramConf  `yaml:"telegram"`
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
	setConfigDefaults(&cfg)
	return &cfg, nil
}

// setConfigDefaults applies default values when optional sections are missing or zero.
func setConfigDefaults(c *Config) {
	if c.Scheduler.IntervalSeconds <= 0 {
		c.Scheduler.IntervalSeconds = 3600
	}
	if c.Candles.LookbackDays <= 0 {
		c.Candles.LookbackDays = 30
	}
	if c.Candles.Interval == "" {
		c.Candles.Interval = "1d"
	}
}
