package kube

import (
	"fmt"
	"regexp"

	"github.com/PraveenPrabhuT/grafana-connect/internal/config"
)

func FindMatchingEnv(currentContext string, cfg *config.Config) (*config.Environment, error) {
	for _, env := range cfg.Environments {
		// If context_match is empty, we skip (or treat as a fallback)
		if env.ContextMatch == "" {
			continue
		}

		// Check if the current context matches the regex defined in YAML
		match, err := regexp.MatchString(env.ContextMatch, currentContext)
		if err != nil {
			return nil, fmt.Errorf("invalid regex in config for %s: %w", env.Name, err)
		}

		if match {
			return &env, nil
		}
	}
	return nil, fmt.Errorf("no matching environment found for context: %s", currentContext)
}
