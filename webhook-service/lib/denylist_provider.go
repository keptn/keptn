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

type DenyListProviderStruct struct {
	GetDeniedURLs GetDeniedURLsFunc
	KubeClient    kubernetes.Interface
}

type GetDeniedURLsFunc func(env map[string]string) []string

func NewDenyListProvider(kubeClient kubernetes.Interface) DenyListProvider {
	return DenyListProviderStruct{
		GetDeniedURLs: GetDeniedURLs,
		KubeClient:    kubeClient,
	}
}

func (d DenyListProviderStruct) Get() []string {
	denyList := d.GetDeniedURLs(GetEnv())

	configMap, err := d.KubeClient.CoreV1().ConfigMaps(GetNamespaceFromEnvVar()).Get(context.TODO(), WebhookConfigMap, metav1.GetOptions{})
	if err != nil {
		logger.Errorf("Unable to get ConfigMap %s content: %s", WebhookConfigMap, err.Error())
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
