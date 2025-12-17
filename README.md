# 🚀 Grafana Connect

**The Context-Aware Dashboard Launcher for Kubernetes Engineers.**

Stop manually searching for Grafana URLs, copy-pasting credentials, and filtering by namespace. `grafana-connect` reads your current Kubernetes context (`kubectl config current-context`), matches it to your environment, and launches the correct dashboard instantly—with the right filters and credentials pre-loaded.

## ✨ Features

* **🧠 Context Aware:** Automatically detects which cluster you are working on.
* **⚡ Instant Launch:** Opens the dashboard filtered to your current **namespace**.
* **📋 Clipboard Integration:** Silently copies the environment password to your clipboard.
* **🔍 Interactive Explorer:**
  * `-i`: Pick a namespace from the current cluster using a fuzzy finder.
  * `-I`: Switch context *and* namespace entirely from the CLI.
* **⚙️ Highly Configurable:** Supports global defaults and per-environment overrides via YAML.

---

## 📦 Installation

### Via Homebrew (Recommended)
```bash
brew tap PraveenPrabhuT/homebrew-tap
brew install grafana-connect
```

### Via Nix (Flakes)
If you are a Nix user, you can run or install directly from the flake:

```bash
# Run once without installing
nix run github:PraveenPrabhuT/grafana-connect

# Install into your profile
nix profile install github:PraveenPrabhuT/grafana-connect

# Enter a dev shell with dependencies
nix develop github:PraveenPrabhuT/grafana-connect
```

### Manual Install
Download the binary for your OS from the [Releases](https://github.com/PraveenPrabhuT/grafana-connect/releases) page and add it to your path.

---

## 🛠 Configuration

You can generate a config file interactively:

```bash
grafana-connect config update
```

This will walk you through setting up your environments and save them to `~/.config/grafana-connect/config.yaml`.

### Configuration Structure (`config.yaml`)

```yaml
# Global Default (Fallback)
default_dashboard: "k8s-pod-resources-clean/kubernetes-pod-resource-dashboard-v3"

environments:
  - name: "ackodev"
    # Regex to match your K8s context (e.g., gke_project_dev-cluster)
    context_match: ".*-dev-cluster.*" 
    base_url: "[https://grafana-ng.internal.ackodev.com](https://grafana-ng.internal.ackodev.com)"
    prometheus_uid: "ebe84db7-b320-41f1-932c-f3e6bdb79432"
    username: "ackodev-grafana-ro"
    password: "dev-password-here"

  - name: "ackoprod"
    context_match: "ackoprod-cluster-01"
    base_url: "[https://central-dashboard.acko.com](https://central-dashboard.acko.com)"
    prometheus_uid: "3KacdaAUglvtBZI3"
    username: "ackoprod-grafana-ro"
    password: "prod-password-here"
```

| Field | Description |
| :--- | :--- |
| `context_match` | A Regex string. If your `kubectl` context matches this, the environment is selected. |
| `base_url` | The root URL of your Grafana instance. |
| `prometheus_uid` | The internal UID of the Datasource. Found in the dashboard URL as `var-DS_PROMETHEUS`. |

---

## 🚀 Usage

### 1. Auto-Detect Mode (Default)
Simply run the command. It uses your **current** kubectl context and namespace.

```bash
grafana-connect
```
> *Output:* 🚀 Detected ackodev (namespace: payments). Opening Dashboard...

### 2. Namespace Picker (`-i`)
Stay on the current cluster, but pick a specific namespace.

```bash
grafana-connect -i
```

### 3. Full Explorer (`-I`)
Switch to a different environment entirely (even if your terminal is pointing to a different context).

```bash
grafana-connect -I
```

### 4. Configuration Management
```bash
# View current config (passwords masked)
grafana-connect config get

# Add or update environments
grafana-connect config update
```

---

## 🧑‍💻 Development

### Prerequisites
* Go 1.25+

### Build locally
```bash
git clone [https://github.com/PraveenPrabhuT/grafana-connect.git](https://github.com/PraveenPrabhuT/grafana-connect.git)
cd grafana-connect
go build -o grafana-connect
```

### Linux Requirements
If you are on Linux, you need `xclip` or `xsel` installed for clipboard support (macOS works out of the box):
```bash
sudo apt-get install xclip
# or
sudo yum install xclip
```