package config

import (
	"github.com/spf13/viper"
)

type Environment struct {
	Name          string `mapstructure:"name"`
	ContextMatch  string `mapstructure:"context_match"`
	BaseURL       string `mapstructure:"base_url"`
	PrometheusUID string `mapstructure:"prometheus_uid"`
	Username      string `mapstructure:"username"`
	Password      string `mapstructure:"password"`
}

type Config struct {
	DefaultDashboard     string        `mapstructure:"default_dashboard"`
	DefaultPrometheusUID string        `mapstructure:"default_prometheus_uid"`
	Environments         []Environment `mapstructure:"environments"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/grafana-connect/")

	// Hardcoded Defaults
	viper.SetDefault("default_dashboard", "k8s-pod-resources-clean/kubernetes-pod-resource-dashboard-v3")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	err := viper.Unmarshal(&cfg)
	return &cfg, err
}
