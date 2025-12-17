package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the base command when called without any subcommands
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long:  `View or modify the grafana-connect configuration file.`,
}

func init() {
	rootCmd.AddCommand(configCmd)
}
