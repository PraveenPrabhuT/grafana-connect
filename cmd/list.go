package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/PraveenPrabhuT/grafana-connect/internal/config"
	"github.com/PraveenPrabhuT/grafana-connect/internal/launcher"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Fuzzy search for an environment",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("❌ Error loading config: %v\n", err)
			os.Exit(1)
		}

		if len(cfg.Environments) == 0 {
			fmt.Println("⚠️  No environments defined in config.yaml")
			return
		}

		// The Fuzzy Finder Logic
		idx, err := fuzzyfinder.Find(
			cfg.Environments,
			func(i int) string {
				return cfg.Environments[i].Name
			},
			fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
				if i == -1 {
					return ""
				}
				env := cfg.Environments[i]

				// Build a nice preview string
				return fmt.Sprintf(
					"Environment: %s\n"+
						"-----------------------------------\n"+
						"🔗 URL:      %s\n"+
						"🆔 PromUID:  %s\n"+
						"👤 User:     %s\n"+
						"🔍 Matcher:  %s\n",
					strings.ToUpper(env.Name),
					env.BaseURL,
					env.PrometheusUID,
					env.Username,
					env.ContextMatch,
				)
			}),
		)

		if err != nil {
			// User aborted (Ctrl+C or Esc)
			return
		}

		// Launch with 'default' namespace since we are in manual mode
		selectedEnv := cfg.Environments[idx]
		launcher.Open(selectedEnv, "default")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
