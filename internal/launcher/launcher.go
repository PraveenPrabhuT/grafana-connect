package launcher

import (
	"fmt"
	"net/url"

	"github.com/PraveenPrabhuT/grafana-connect/internal/config"
	"github.com/atotto/clipboard"
	"github.com/pkg/browser"
)

// Open now only needs the Environment and the Namespace
func Open(env config.Environment, namespace string) {
	// 1. Build URL (Using env-specific fields)
	params := url.Values{}
	params.Add("orgId", "1")
	params.Add("refresh", "30s")
	params.Add("var-DS_PROMETHEUS", env.PrometheusUID)
	params.Add("var-namespace", namespace)
	params.Add("var-deployment", "All")
	params.Add("var-pod", "All")
	params.Add("var-container", "All")

	// Fallback for dashboard if user forgot to set it in the new config structure
	dashSlug := env.Dashboard
	if dashSlug == "" {
		// Reasonable safety net or panic? Let's warn.
		fmt.Println("⚠️  Warning: No dashboard path set for this environment.")
	}

	finalURL := fmt.Sprintf("%s/d/%s?%s",
		env.BaseURL,
		dashSlug,
		params.Encode(),
	)

	// 2. Handle Clipboard
	if env.Password != "" {
		// Init returns an error if the system clipboard is missing (e.g. headless linux)
		if err := clipboard.WriteAll(env.Password); err == nil {
			fmt.Println("📋 Password copied to clipboard!")
		} else {
			// On Linux, this might fail if xclip/xsel isn't installed. Warn the user.
			fmt.Printf("⚠️  Clipboard error: %v (Do you have xclip/xsel installed?)\n", err)
		}
	}

	// 3. Launch
	fmt.Printf("🚀 Opening %s [%s]...\n", env.Name, namespace)
	if err := browser.OpenURL(finalURL); err != nil {
		fmt.Printf("❌ Failed to open browser: %v\n", err)
		fmt.Printf("   Link: %s\n", finalURL)
	}
}
