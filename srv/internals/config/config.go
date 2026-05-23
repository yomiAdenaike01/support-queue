package config

import (
	"fmt"
	"os"
)

type Config struct{}

func New() *Config {
	return &Config{}
}

func (c *Config) GetEnv(key string) string {
	return os.Getenv(key)
}

func (c *Config) GetEnvOrDefault(key string, fallback string) string {
	value := c.GetEnv(key)
	if value == "" {
		return fallback
	}
	return value
}

func (c *Config) GetEnvOrFail(key string) string {
	value := c.GetEnv(key)
	if value == "" {
		panic(fmt.Sprintf("Missing required config value: %s", key))
	}
	return value
}
