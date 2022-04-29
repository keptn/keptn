package lib

import (
	"context"
	"strings"

	keptnkubeutils "github.com/keptn/kubernetes-utils/pkg"
	logger "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type IDenyListProvider interface {
	GetDenyList() []string
}

type DenyListProvider struct {
	GetKubeAPI                  GetKubeAPIFunc
	GetDeniedURLs               GetDeniedURLsFunc
	GetWebhookDenyListConfigMap GetWebhookDenyListConfigMapFunc
}

type GetKubeAPIFunc func(useInClusterConfig bool) (v1.CoreV1Interface, error)
type GetDeniedURLsFunc func(env map[string]string) []string
type GetWebhookDenyListConfigMapFunc func(kubeAPI v1.CoreV1Interface) (*corev1.ConfigMap, error)

func NewDenyListProvider() DenyListProvider {
	return DenyListProvider{
		GetKubeAPI:                  keptnkubeutils.GetKubeAPI,
		GetDeniedURLs:               GetDeniedURLs,
		GetWebhookDenyListConfigMap: getWebhookDenyListConfigMap,
	}
}

func (d DenyListProvider) GetDenyList() []string {
	denyList := d.GetDeniedURLs(GetEnv())
	kubeAPI, err := d.GetKubeAPI(true)
	if err != nil {
		logger.Errorf("Unable to read ConfigMap %s: cannot get kubeAPI: %s", WebhookConfigMap, err.Error())
		return denyList
	}

	configMap, err := d.GetWebhookDenyListConfigMap(kubeAPI)
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

func getWebhookDenyListConfigMap(kubeAPI v1.CoreV1Interface) (*corev1.ConfigMap, error) {
	return kubeAPI.ConfigMaps(GetNamespaceFromEnvVar()).Get(context.TODO(), WebhookConfigMap, metav1.GetOptions{})
}
