package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type BasicAuthUser struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	Port      int `yaml:"port"`
	BasicAuth struct {
		Enabled bool            `yaml:"enabled"`
		Users   []BasicAuthUser `yaml:"users"`
	} `yaml:"basic_auth"`
}

var AppConfig Config

func LoadConfig(configPath string) error {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &AppConfig)
	return err
}
