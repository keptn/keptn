package lib

import (
	"context"
	"strings"

	logger "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DenyListProvider interface {
	Get() []string
}

type denyListProvider struct {
	getDeniedURLs GetDeniedURLsFunc
	kubeClient    kubernetes.Interface
}

type GetDeniedURLsFunc func(env map[string]string) []string

func NewDenyListProvider(kubeClient kubernetes.Interface) DenyListProvider {
	return denyListProvider{
		getDeniedURLs: GetDeniedURLs,
		kubeClient:    kubeClient,
	}
}

func (d denyListProvider) Get() []string {
	denyList := d.getDeniedURLs(GetEnv())

	configMap, err := d.kubeClient.CoreV1().ConfigMaps(GetNamespaceFromEnvVar()).Get(context.TODO(), WebhookConfigMap, metav1.GetOptions{})
	if err != nil {
		logger.Errorf("Could not get ConfigMap %s content: %v", WebhookConfigMap, err)
		return denyList
	}

	denyListString := configMap.Data["denyList"]
	denyListConfig := strings.Fields(denyListString)
	return append(denyList, denyListConfig...)
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
