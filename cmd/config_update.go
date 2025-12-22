package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/PraveenPrabhuT/grafana-connect/internal/config"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Create or update configuration interactively",
	Long:  "Starts a wizard to add new environments or update existing ones based on the Grafana Base URL.",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Setup Paths
		home, _ := os.UserHomeDir()
		configDir := filepath.Join(home, ".config", "grafana-connect")
		configPath := filepath.Join(configDir, "config.yaml")

		// 2. Load Existing or Create New
		var cfg config.Config
		if _, err := os.Stat(configPath); err == nil {
			fmt.Printf("📂 Loading config from: %s\n", configPath)
			data, err := os.ReadFile(configPath)
			if err == nil {
				_ = yaml.Unmarshal(data, &cfg)
			}
		}

		fmt.Println("🧙 Starting setup wizard...")
		// Removed Global Defaults Section

		for {
			prompt := promptui.Prompt{
				Label:     "Add/Update an Environment?",
				IsConfirm: true,
			}
			if _, err := prompt.Run(); err != nil {
				break // Exit loop
			}

			fmt.Println("\n--- 🌍 Environment Details ---")

			// Primary Key: URL
			pURL := promptui.Prompt{Label: "Grafana Base URL"}
			rawURL, _ := pURL.Run()
			baseURL := strings.TrimSuffix(rawURL, "/")

			// Check for existing
			var existingEnv *config.Environment
			idx := -1
			for i, e := range cfg.Environments {
				if e.BaseURL == baseURL {
					existingEnv = &cfg.Environments[i]
					idx = i
					fmt.Println("ℹ️  Updating existing environment entry.")
					break
				}
			}

			// Pre-fill defaults
			defName, defAlias, defCtx, defUID, defUser, defDash := "", "", "", "", "", ""
			if existingEnv != nil {
				defName = existingEnv.Name
				defAlias = existingEnv.Alias
				defCtx = existingEnv.ContextMatch
				defUID = existingEnv.PrometheusUID
				defUser = existingEnv.Username
				defDash = existingEnv.Dashboard
			}

			// --- THE PROMPTS ---
			pName := promptui.Prompt{Label: "Name (e.g. ackoprod)", Default: defName}
			name, _ := pName.Run()

			// NEW: Alias
			pAlias := promptui.Prompt{Label: "Alias (Shortcode e.g. prod)", Default: defAlias}
			alias, _ := pAlias.Run()

			pCtx := promptui.Prompt{Label: "Context Regex", Default: defCtx}
			if defCtx == "" {
				pCtx.Default = ".*" + name + ".*"
			}
			ctxMatch, _ := pCtx.Run()

			// NEW: Dashboard Path per env
			pDash := promptui.Prompt{
				Label:   "Dashboard Path (Slug)",
				Default: defDash,
			}
			if defDash == "" {
				pDash.Default = "k8s-pod-resources-clean/kubernetes-pod-resource-dashboard-v3"
			}
			dashboard, _ := pDash.Run()

			pUID := promptui.Prompt{Label: "Prometheus UID", Default: defUID}
			uid, _ := pUID.Run()

			pUser := promptui.Prompt{Label: "Username", Default: defUser}
			user, _ := pUser.Run()

			pPass := promptui.Prompt{Label: "Password (leave empty to keep)", Mask: '*'}
			pass, _ := pPass.Run()

			if existingEnv != nil && pass == "" {
				pass = existingEnv.Password
			}

			newEnv := config.Environment{
				Name: name, Alias: alias, ContextMatch: ctxMatch, BaseURL: baseURL,
				Dashboard: dashboard, PrometheusUID: uid, Username: user, Password: pass,
			}

			if idx >= 0 {
				cfg.Environments[idx] = newEnv
			} else {
				cfg.Environments = append(cfg.Environments, newEnv)
			}
			fmt.Println("✅ Saved environment.")
		}

		// 5. Write to disk
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Printf("❌ FS Error: %v\n", err)
			return
		}
		data, _ := yaml.Marshal(&cfg)
		if err := os.WriteFile(configPath, data, 0600); err != nil {
			fmt.Printf("❌ Save Error: %v\n", err)
			return
		}

		fmt.Printf("\n🎉 Config saved to: %s\n", configPath)
	},
}

func init() {
	configCmd.AddCommand(configUpdateCmd)
}
