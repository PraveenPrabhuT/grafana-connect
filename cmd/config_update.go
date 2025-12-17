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
		} else {
			fmt.Println("✨ No config found. Creating a new one.")
		}

		// 3. Global Settings Wizard
		fmt.Println("\n--- ⚙️  Global Defaults ---")

		pDash := promptui.Prompt{
			Label:   "Default Dashboard Slug",
			Default: cfg.DefaultDashboard,
		}
		if cfg.DefaultDashboard == "" {
			pDash.Default = "k8s-pod-resources-clean/kubernetes-pod-resource-dashboard-v3"
		}
		cfg.DefaultDashboard, _ = pDash.Run()

		pProm := promptui.Prompt{
			Label:   "Default Prometheus UID",
			Default: cfg.DefaultPrometheusUID,
		}
		cfg.DefaultPrometheusUID, _ = pProm.Run()

		// 4. Environments Wizard loop
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
			defName, defCtx, defUID, defUser := "", "", "", ""
			if existingEnv != nil {
				defName = existingEnv.Name
				defCtx = existingEnv.ContextMatch
				defUID = existingEnv.PrometheusUID
				defUser = existingEnv.Username
			}

			pName := promptui.Prompt{Label: "Name (e.g., prod)", Default: defName}
			name, _ := pName.Run()

			pCtx := promptui.Prompt{Label: "Context Regex", Default: defCtx}
			if defCtx == "" {
				pCtx.Default = ".*" + name + ".*"
			}
			ctxMatch, _ := pCtx.Run()

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
				Name: name, ContextMatch: ctxMatch, BaseURL: baseURL,
				PrometheusUID: uid, Username: user, Password: pass,
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
