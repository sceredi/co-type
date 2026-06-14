// Package config is responsible for setting up the configuration for the broker application.
package config

import (
	"sync"
)

// Config holds the configuration settings for the broker application.
type Config struct{}

var (
	globalCfg *Config
	once      sync.Once
)

// Get returns the global configuration instance, loading it if it hasn't been loaded yet.
func Get() *Config {
	once.Do(load)
	return globalCfg
}

func load() {
	globalCfg = &Config{}
}
