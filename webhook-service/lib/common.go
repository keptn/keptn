package lib

import (
	"os"
	"strings"
)

const WebhookConfigMap = "keptn-webhook-config"

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
