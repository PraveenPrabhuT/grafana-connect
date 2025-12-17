package ui

import (
	"fmt"
	"github.com/PraveenPrabhuT/grafana-connect/internal/config"
	"github.com/ktr0731/go-fuzzyfinder"
	"strings"
)

func SelectEnvironment(envs []config.Environment) (*config.Environment, error) {
	idx, err := fuzzyfinder.Find(
		envs,
		func(i int) string {
			return envs[i].Name
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			env := envs[i]
			return fmt.Sprintf("Environment: %s\nURL: %s\nUser: %s",
				strings.ToUpper(env.Name), env.BaseURL, env.Username)
		}),
	)
	if err != nil {
		return nil, err
	}
	return &envs[idx], nil
}
