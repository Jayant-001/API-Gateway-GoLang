package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	Routes []RouteConfig `yaml:"routes"`
	RateLimiter RateLimitConfig `yaml:"rate_limiter"`
	Auth AuthConfig `yaml:"auth"`
}

type ServerConfig struct {
	PORT string `yaml:"port"`
	ReadTimeout time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type RouteConfig struct {
	Name string `yaml:"name"`
	PathPrefix string `yaml:"path_prefix"`
	UpstreamURL string `yaml:"upstream_url"`
	AuthRequired bool `yaml:"auth_required"`
}

type RateLimitConfig struct {
	Enabled bool `yaml:"enabled"`
	RequestsPerSecond int `yaml:"requests_per_second"`
	Burst int `yaml:"burst"`
}

type AuthConfig struct {
	APIKEY string `yaml:"api_key"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}