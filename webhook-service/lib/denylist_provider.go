package lib

import (
	"context"
	"strings"

	keptnkubeutils "github.com/keptn/kubernetes-utils/pkg"
	logger "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IDenyListProvider interface {
	GetDenyList() []string
}

type DenyListProvider struct {
}

func (d DenyListProvider) GetDenyList() []string {
	denyList := GetDeniedURLs(GetEnv())
	kubeAPI, err := keptnkubeutils.GetKubeAPI(true)
	if err != nil {
		logger.Errorf("Unable to read ConfigMap %s: cannot get kubeAPI: %s", WebhookConfigMap, err.Error())
		return denyList
	}

	configMap, err := kubeAPI.ConfigMaps(GetNamespaceFromEnvVar()).Get(context.TODO(), WebhookConfigMap, v1.GetOptions{})
	if err != nil {
		logger.Errorf("Unable to get ConfigMap %s content: %s", WebhookConfigMap, err.Error())
		return denyList
	}

	denyListString := configMap.Data["denyList"]
	denyList = strings.Fields(denyListString)
	return denyList
}

func GetDeniedURLs(env map[string]string) []string {
	kubeAPIHostIP := env[KubernetesSvcHostEnvVar]
	kubeAPIPort := env[KubernetesAPIPortEnvVar]

	urls := make([]string, 0)
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
