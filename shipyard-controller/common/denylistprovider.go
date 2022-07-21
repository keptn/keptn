package common

import (
	"k8s.io/client-go/kubernetes"
)

type DenyListProvider interface {
	Get() []string
}

type denyListProvider struct {
	getDeniedURLs GetDeniedURLsFunc
	kubeClient    kubernetes.Interface
}

type GetDeniedURLsFunc func() []string

func NewDenyListProvider(kubeClient kubernetes.Interface) DenyListProvider {
	return denyListProvider{
		getDeniedURLs: GetDeniedURLs,
		kubeClient:    kubeClient,
	}
}

func (d denyListProvider) Get() []string {
	denyList := d.getDeniedURLs()
	return denyList

	// configMap, err := d.kubeClient.CoreV1().ConfigMaps(GetNamespaceFromEnvVar()).Get(context.TODO(), WebhookConfigMap, metav1.GetOptions{})
	// if err != nil {
	// 	logger.Errorf("Unable to get ConfigMap %s content: %s", WebhookConfigMap, err.Error())
	// 	return denyList
	// }

	// denyListString := configMap.Data["denyList"]
	// denyListConfig := strings.Fields(denyListString)
	// return append(denyList, denyListConfig...)
}

func GetDeniedURLs() []string {
	// kubeAPIHostIP := env[KubernetesSvcHostEnvVar]
	// kubeAPIPort := env[KubernetesAPIPortEnvVar]

	// urls := make([]string, 0)
	// if kubeAPIHostIP != "" {
	// 	urls = append(urls, kubeAPIHostIP)
	// }
	// if kubeAPIPort != "" {
	// 	urls = append(urls, "kubernetes"+":"+kubeAPIPort)
	// 	urls = append(urls, "kubernetes.default"+":"+kubeAPIPort)
	// 	urls = append(urls, "kubernetes.default.svc"+":"+kubeAPIPort)
	// 	urls = append(urls, "kubernetes.default.svc.cluster.local"+":"+kubeAPIPort)
	// }
	// if kubeAPIHostIP != "" && kubeAPIPort != "" {
	// 	urls = append(urls, kubeAPIHostIP+":"+kubeAPIPort)
	// }
	return []string{}
}
