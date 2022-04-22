package lib

import (
	"os"
	"strings"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

type Config struct {
	GetKubeAPI KubeAPIConfigFunc
}

var config *Config
var configOnce sync.Once

func GetConfig() *Config {
	configOnce.Do(func() {
		config = &Config{GetKubeAPI: getInClusterKubeClient}
	})
	return config
}

type KubeAPIConfigFunc func() (kubernetes.Interface, error)

func getInClusterKubeClient() (kubernetes.Interface, error) {
	var config *rest.Config
	config, err := rest.InClusterConfig()

	if err != nil {
		return nil, err
	}

	kubeAPI, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return kubeAPI, nil
}
