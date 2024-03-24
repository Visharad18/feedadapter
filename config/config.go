package config

import (
	"os"
	"time"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	HTTPPort           string `yaml:"http_port"`
	AlphaVantageConfig `yaml:"alpha_vantage_config"`
	PolygonConfig      `yaml:"polygon_config"`
	Symbols            []string      `yaml:"symbols"`
	FetchInterval      time.Duration `yaml:"fetch_interval"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	f, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return nil, err
	}
	return &cfg, yaml.Unmarshal(f, &cfg)
}
