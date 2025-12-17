package cmd

import (
	"fmt"

	"github.com/PraveenPrabhuT/grafana-connect/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Display current configuration",
	Long:  "Prints the current configuration YAML to stdout. Passwords are masked for security.",
	Run: func(cmd *cobra.Command, args []string) {
		// Load the config struct (which handles finding the file)
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("❌ Could not load config: %v\n", err)
			fmt.Println("   Run 'grafana-connect config update' to generate one.")
			return
		}

		// Create a copy to mask passwords without modifying the actual file
		safeCfg := *cfg
		maskedEnvs := make([]config.Environment, len(cfg.Environments))
		copy(maskedEnvs, cfg.Environments)

		for i := range maskedEnvs {
			if maskedEnvs[i].Password != "" {
				maskedEnvs[i].Password = "*****"
			}
		}
		safeCfg.Environments = maskedEnvs

		// Marshal to YAML for display
		data, err := yaml.Marshal(&safeCfg)
		if err != nil {
			fmt.Printf("❌ Error formatting config: %v\n", err)
			return
		}

		fmt.Println("# Current Configuration (Passwords Masked)")
		fmt.Println(string(data))
	},
}

func init() {
	configCmd.AddCommand(configGetCmd)
}
