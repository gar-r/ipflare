package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Entries   EntrySet `yaml:"entries"`
	ApiToken  string   `yaml:"api_token"`
	Frequency int      `yaml:"frequency"`
}

type EntrySet map[string][]string

func Parse(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		Entries: make(EntrySet),
	}
	err = yaml.Unmarshal(bytes, cfg)
	return cfg, err
}

func (c *Config) String() string {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
