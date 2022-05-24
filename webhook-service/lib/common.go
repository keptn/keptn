package lib

import (
	"os"
	"strings"
)

const (
	WebhookConfigMap        = "keptn-webhook-config"
	KubernetesSvcHostEnvVar = "KUBERNETES_SERVICE_HOST"
	KubernetesAPIPortEnvVar = "KUBERNETES_SERVICE_PORT"
)

type AdrDomainNameMapping map[string][]string

func GetNamespaceFromEnvVar() string {
	return os.Getenv("POD_NAMESPACE")
}

func GetEnv() map[string]string {
	envMap := make(map[string]string)
	for _, e := range os.Environ() {
		if i := strings.Index(e, "="); i >= 0 {
			envMap[e[:i]] = e[i+1:]
		}
	}
	return envMap
}

func CreateListOfDeniedURLs(env map[string]string) []string {
	kubeAPIHostIP := env[KubernetesSvcHostEnvVar]
	kubeAPIPort := env[KubernetesAPIPortEnvVar]

	urls := []string{
		// Block access to Kubernetes API
		"kubernetes",
		"kubernetes.default",
		"kubernetes.default.svc",
		"kubernetes.default.svc.cluster.local",
		"svc.cluster.local",
		"cluster.local",
		// Block access to localhost
		"localhost",
		"127.0.0.1",
		"::1",
	}

	// try to add new host and port if existing
	if kubeAPIHostIP != "" {
		urls = append(urls, kubeAPIHostIP)
	}
	if kubeAPIPort != "" {
		urls = append(urls, "kubernetes"+":"+kubeAPIPort)
		urls = append(urls, "kubernetes.default"+":"+kubeAPIPort)
		urls = append(urls, "kubernetes.default.svc"+":"+kubeAPIPort)
		urls = append(urls, "kubernetes.default.svc.cluster.local"+":"+kubeAPIPort)
	}
	if kubeAPIHostIP != "" && kubeAPIPort != "" {
		urls = append(urls, kubeAPIHostIP+":"+kubeAPIPort)
	}
	return urls
}
