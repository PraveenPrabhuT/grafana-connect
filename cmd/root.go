package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/PraveenPrabhuT/grafana-connect/internal/config"
	"github.com/PraveenPrabhuT/grafana-connect/internal/kube"
	"github.com/PraveenPrabhuT/grafana-connect/internal/launcher"
	"github.com/PraveenPrabhuT/grafana-connect/internal/ui"
)

var (
	flagInteractiveNs  bool   // -i
	flagInteractiveCtx bool   // -I
	flagAlias          string // -e
	flagNamespace      string // -n
)

var rootCmd = &cobra.Command{
	Use:   "grafana-connect",
	Short: "Context-aware Grafana launcher",
	Long:  `Automatically detects your K8s context and opens the relevant Grafana dashboard with filters applied.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Load Configuration
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("❌ Error loading config: %v\n", err)
			os.Exit(1)
		}

		var targetEnv *config.Environment
		var targetNamespace string

		// --- LOGIC FLOW ---

		// 1. Check for Alias Flag (-e)
		if flagAlias != "" {
			targetEnv = cfg.FindByAlias(flagAlias)
			if targetEnv == nil {
				fmt.Printf("❌ No environment found with alias: '%s'\n", flagAlias)
				os.Exit(1)
			}
			// Default NS for alias mode is "default" unless overridden later
			targetNamespace = "default"
		}

		// 2. Check for Interactive Flags (-I / -i) ONLY if alias wasn't provided
		if targetEnv == nil {
			if flagInteractiveCtx {
				// -I: Full Selection
				env, err := ui.SelectEnvironment(cfg.Environments)
				if err != nil {
					return
				}
				targetEnv = env

				// Resolve context for NS fetching
				ctxName, err := kube.FindContextByRegex(env.ContextMatch)
				if err == nil {
					// Only try to fetch namespaces if we found a matching local context
					fmt.Printf("📡 Fetching namespaces from [%s]...\n", ctxName)
					nss, err := kube.GetNamespaces(ctxName)
					if err == nil {
						targetNamespace, _ = ui.SelectString("Select Namespace", nss)
					}
				}
				if targetNamespace == "" {
					targetNamespace = "default"
				}

			} else if flagInteractiveNs {
				// === MODE: -i (Current Context -> Choose NS) ===

				// A. Get Current State
				state, err := kube.GetCurrentState()
				if err != nil {
					fmt.Printf("❌ Could not detect K8s state: %v\n", err)
					os.Exit(1)
				}
				targetEnv, _ = kube.FindMatchingEnv(state.Context, cfg)

				fmt.Printf("📡 Fetching namespaces from [%s]...\n", state.Context)
				nss, err := kube.GetNamespaces(state.Context)
				if err == nil {
					targetNamespace, _ = ui.SelectString("Select Namespace", nss)
				}
			}
		}

		// 3. Fallback to Auto-Detect
		if targetEnv == nil && !flagInteractiveNs {
			state, err := kube.GetCurrentState()
			if err != nil {
				fmt.Printf("❌ Could not detect K8s state: %v\n", err)
				os.Exit(1)
			}
			targetEnv, err = kube.FindMatchingEnv(state.Context, cfg)
			if err != nil {
				fmt.Printf("⚠️  No mapping found for context: %s\n", state.Context)
				os.Exit(1)
			}
			targetNamespace = state.Namespace
		}

		// 4. Apply Namespace Override Flag (-n)
		// This applies to ANY mode above (Alias, Interactive, or Auto)
		if flagNamespace != "" {
			targetNamespace = flagNamespace
		}

		// Final Launch
		if targetEnv != nil {
			launcher.Open(*targetEnv, targetNamespace)
		}
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&flagInteractiveNs, "interactive-ns", "i", false, "Pick namespace interactively")
	rootCmd.Flags().BoolVarP(&flagInteractiveCtx, "interactive-full", "I", false, "Pick environment and namespace interactively")
	rootCmd.Flags().StringVarP(&flagAlias, "env", "e", "", "Select environment by alias (e.g. 'prod')")
	rootCmd.Flags().StringVarP(&flagNamespace, "namespace", "n", "", "Override namespace")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
