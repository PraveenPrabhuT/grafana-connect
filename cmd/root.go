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
	flagInteractiveNs  bool // -i
	flagInteractiveCtx bool // -I
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
			fmt.Println("   Run 'grafana-connect init' to generate a default config.")
			os.Exit(1)
		}

		var targetEnv *config.Environment
		var targetNamespace string

		// --- LOGIC BRANCHING ---

		if flagInteractiveCtx {
			// === MODE: -I (Choose Env -> Choose NS) ===

			// A. Select Environment from YAML
			// (We reuse the logic from list.go here, or import it.
			// For brevity, let's assume we implement a SelectEnvironment in ui package or inline it)
			// ... Fuzzy Find Environment Logic ...
			// Let's assume ui.SelectEnvironment returns the chosen config.Environment
			env, err := ui.SelectEnvironment(cfg.Environments) // *You need to add this to ui package*
			if err != nil {
				return // User aborted
			}
			targetEnv = env

			// B. Resolve Regex to actual Kube Context Name
			ctxName, err := kube.FindContextByRegex(env.ContextMatch)
			if err != nil {
				fmt.Printf("❌ Environment '%s' matches regex '%s', but no corresponding context found in local kubeconfig.\n", env.Name, env.ContextMatch)
				os.Exit(1)
			}

			// C. Fetch Namespaces for that context
			fmt.Printf("📡 Fetching namespaces from cluster [%s]...\n", ctxName)
			namespaces, err := kube.GetNamespaces(ctxName)
			if err != nil {
				fmt.Printf("❌ Failed to list namespaces: %v\n", err)
				os.Exit(1)
			}

			// D. Select Namespace
			ns, err := ui.SelectString("Select Namespace", namespaces)
			if err != nil {
				return
			}
			targetNamespace = ns

		} else if flagInteractiveNs {
			// === MODE: -i (Current Context -> Choose NS) ===

			// A. Get Current State
			state, err := kube.GetCurrentState()
			if err != nil {
				fmt.Printf("❌ Could not detect K8s state: %v\n", err)
				os.Exit(1)
			}

			// B. Find Matching Env (Automatic)
			env, err := kube.FindMatchingEnv(state.Context, cfg)
			if err != nil {
				fmt.Printf("⚠️  No config mapping found for context: %s\n", state.Context)
				os.Exit(1)
			}
			targetEnv = env

			// C. Fetch Namespaces
			fmt.Printf("📡 Fetching namespaces from cluster [%s]...\n", state.Context)
			namespaces, err := kube.GetNamespaces(state.Context)
			if err != nil {
				fmt.Printf("❌ Failed to list namespaces: %v\n", err)
				os.Exit(1)
			}

			// D. Select Namespace
			ns, err := ui.SelectString("Select Namespace", namespaces)
			if err != nil {
				return
			}
			targetNamespace = ns

		} else {
			// === MODE: Default (Auto-Detect Everything) ===
			state, err := kube.GetCurrentState()
			if err != nil {
				fmt.Printf("❌ Could not detect K8s state: %v\n", err)
				os.Exit(1)
			}

			env, err := kube.FindMatchingEnv(state.Context, cfg)
			if err != nil {
				fmt.Printf("⚠️  No mapping found for context: %s\n", state.Context)
				fmt.Println("   Use '-i' or '-I' to select an environment manually.")
				os.Exit(1)
			}

			targetEnv = env
			targetNamespace = state.Namespace
		}

		// --- EXECUTION ---
		// We pass the resolved Env and Namespace to the launcher
		launcher.Open(*targetEnv, targetNamespace)
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&flagInteractiveNs, "interactive-ns", "i", false, "Interactively select namespace from current cluster")
	rootCmd.Flags().BoolVarP(&flagInteractiveCtx, "interactive-full", "I", false, "Interactively select Environment AND Namespace")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
