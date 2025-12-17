package kube

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeState struct {
	Context   string
	Namespace string
}

func GetCurrentState() (*KubeState, error) {
	// 1. Load the default kubeconfig (~/.kube/config)
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	// 2. Extract RawConfig to get the CurrentContext name
	rawConfig, err := kubeConfig.RawConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load kubeconfig: %w", err)
	}

	currentCtxName := rawConfig.CurrentContext
	if currentCtxName == "" {
		return nil, fmt.Errorf("no current-context set in kubeconfig")
	}

	// 3. Extract the Namespace from the current context
	// Note: .Namespace() handles the logic of falling back to "default" if not set
	ns, _, err := kubeConfig.Namespace()
	if err != nil {
		ns = "default" // Safe fallback
	}

	return &KubeState{
		Context:   currentCtxName,
		Namespace: ns,
	}, nil
}

func GetNamespaces(contextName string) ([]string, error) {
	// 1. Build Config for the specific context
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{CurrentContext: contextName}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to build config for context %s: %w", contextName, err)
	}

	// 2. Create Clientset
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	// 3. Call K8s API
	// Set a timeout to avoid hanging if the cluster is unreachable
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	nsList, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	// 4. Extract and Sort
	var namespaces []string
	for _, ns := range nsList.Items {
		namespaces = append(namespaces, ns.Name)
	}
	sort.Strings(namespaces)

	return namespaces, nil
}

// FindContextByRegex looks through ~/.kube/config and returns the first context matching the regex
func FindContextByRegex(regexStr string) (string, error) {
	config, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return "", err
	}

	r, err := regexp.Compile(regexStr)
	if err != nil {
		return "", err
	}

	for ctxName := range config.Contexts {
		if r.MatchString(ctxName) {
			return ctxName, nil
		}
	}
	return "", fmt.Errorf("no kubeconfig context found matching regex: %s", regexStr)
}
