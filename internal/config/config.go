package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Environment struct {
	// Added `yaml:"..."` tags to ensure correct saving format
	Name          string `mapstructure:"name"           yaml:"name"`
	ContextMatch  string `mapstructure:"context_match"  yaml:"context_match"`
	BaseURL       string `mapstructure:"base_url"       yaml:"base_url"`
	PrometheusUID string `mapstructure:"prometheus_uid" yaml:"prometheus_uid"`
	Username      string `mapstructure:"username"       yaml:"username"`
	Password      string `mapstructure:"password"       yaml:"password"`
}

type Config struct {
	DefaultDashboard     string        `mapstructure:"default_dashboard"      yaml:"default_dashboard"`
	DefaultPrometheusUID string        `mapstructure:"default_prometheus_uid" yaml:"default_prometheus_uid"`
	Environments         []Environment `mapstructure:"environments"           yaml:"environments"`
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
