package cmd

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"

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

		// 2. Detect Kubernetes Context
		state, err := kube.GetCurrentState()
		if err != nil {
			fmt.Printf("❌ Could not detect K8s state: %v\n", err)
			os.Exit(1)
		}

		// 3. Find Matching Environment
		env, err := kube.FindMatchingEnv(state.Context, cfg)
		if err != nil {
			fmt.Printf("⚠️  No mapping found for context: %s\n", state.Context)
			fmt.Println("   Use '-i' to select an environment manually.")
			os.Exit(1)
		}

		// 4. Construct the URL
		// Logic: Use Env specific UID -> Fallback to Config Default
		promUID := env.PrometheusUID
		if promUID == "" {
			promUID = cfg.DefaultPrometheusUID
		}

		// URL Encoding parameters
		params := url.Values{}
		params.Add("orgId", "1")
		params.Add("refresh", "30s")
		params.Add("var-DS_PROMETHEUS", promUID)
		params.Add("var-namespace", state.Namespace)
		params.Add("var-deployment", "All")
		params.Add("var-pod", "All")
		params.Add("var-container", "All")

		finalURL := fmt.Sprintf("%s/d/%s?%s",
			env.BaseURL,
			cfg.DefaultDashboard,
			params.Encode(),
		)

		// 5. Handle Credentials (Silent Copy)
		if env.Password != "" {
			err := clipboard.Init()
			if err == nil {
				clipboard.Write(clipboard.FmtText, []byte(env.Password))
				fmt.Println("📋 Password copied to clipboard!")
			} else {
				fmt.Println("⚠️  Clipboard unavailable (missing xclip/xsel?).")
			}
		}

		// 6. Launch
		fmt.Printf("🚀 Detected %s (%s). Opening Dashboard...\n", env.Name, state.Namespace)
		if err := browser.OpenURL(finalURL); err != nil {
			fmt.Printf("❌ Failed to open browser: %v\n", err)
			fmt.Printf("   Link: %s\n", finalURL)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
