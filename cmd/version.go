package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// These are set by the linker during build
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print build information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Grafana Connect\n")
		fmt.Printf("  Version: %s\n", version)
		fmt.Printf("  Commit:  %s\n", commit)
		fmt.Printf("  Date:    %s\n", date)
		fmt.Printf("  Runtime: %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
