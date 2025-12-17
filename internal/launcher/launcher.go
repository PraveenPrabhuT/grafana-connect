package launcher

import (
	"fmt"
	"net/url"

	"github.com/PraveenPrabhuT/grafana-connect/internal/config"
	"github.com/pkg/browser"
	"golang.design/x/clipboard"
)

// Open prepares the environment and launches the browser
func Open(env config.Environment, globalCfg *config.Config, namespace string) {
	// 1. Resolve Prometheus UID (Env override > Global default)
	promUID := env.PrometheusUID
	if promUID == "" {
		promUID = globalCfg.DefaultPrometheusUID
	}

	// 2. Build URL
	params := url.Values{}
	params.Add("orgId", "1")
	params.Add("refresh", "30s")
	params.Add("var-DS_PROMETHEUS", promUID)
	params.Add("var-namespace", namespace)
	params.Add("var-deployment", "All")
	params.Add("var-pod", "All")

	finalURL := fmt.Sprintf("%s/d/%s?%s",
		env.BaseURL,
		globalCfg.DefaultDashboard,
		params.Encode(),
	)

	// 3. Handle Clipboard
	if env.Password != "" {
		// Init returns an error if the system clipboard is missing (e.g. headless linux)
		if err := clipboard.Init(); err == nil {
			clipboard.Write(clipboard.FmtText, []byte(env.Password))
			fmt.Println("📋 Password copied to clipboard!")
		} else {
			// Fail silently or log debug
		}
	}

	// 4. Launch
	fmt.Printf("🚀 Opening %s [%s]...\n", env.Name, namespace)
	if err := browser.OpenURL(finalURL); err != nil {
		fmt.Printf("❌ Failed to open browser: %v\n", err)
		fmt.Printf("   Link: %s\n", finalURL)
	}
}
