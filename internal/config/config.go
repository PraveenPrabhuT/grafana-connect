package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Environment struct {
	Name          string `mapstructure:"name"           yaml:"name"`
	Alias         string `mapstructure:"alias"          yaml:"alias"` // Added for next feature
	ContextMatch  string `mapstructure:"context_match"  yaml:"context_match"`
	BaseURL       string `mapstructure:"base_url"       yaml:"base_url"`
	Dashboard     string `mapstructure:"dashboard"      yaml:"dashboard"` // Moved here
	PrometheusUID string `mapstructure:"prometheus_uid" yaml:"prometheus_uid"`
	Username      string `mapstructure:"username"       yaml:"username"`
	Password      string `mapstructure:"password"       yaml:"password"`
}

type Config struct {
	// Global defaults are gone. Only the list remains.
	Environments []Environment `mapstructure:"environments" yaml:"environments"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	home, _ := os.UserHomeDir()
	viper.AddConfigPath(filepath.Join(home, ".config", "grafana-connect"))
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	err := viper.Unmarshal(&cfg)
	return &cfg, err
}

// Helper to find env by Alias
func (c *Config) FindByAlias(alias string) *Environment {
	for _, env := range c.Environments {
		if env.Alias == alias {
			return &env
		}
	}
	return nil
}
