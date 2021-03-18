package backend

import (
	"github.com/keptn/keptn/secret-service/internal/common"
	"github.com/keptn/keptn/secret-service/internal/model"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const BackendTypeK8s = "kubernetes"

type K8sSecretBackend struct {
	KubeAPI                kubernetes.Interface
	KeptnNamespaceProvider common.StringSupplier
}

func NewK8sSecretBackend(kubeAPI kubernetes.Interface) *K8sSecretBackend {
	return &K8sSecretBackend{
		KubeAPI:                kubeAPI,
		KeptnNamespaceProvider: common.EnvBasedStringSupplier("POD_NAMESPACE", DefaultNamespace),
	}
}

func (k K8sSecretBackend) CreateSecret(secret model.Secret) error {

	namespace := k.KeptnNamespaceProvider()
	kubeSecret := k.createKubernetesSecret(secret, namespace)
	_, err := k.KubeAPI.CoreV1().Secrets(namespace).Create(kubeSecret)
	if err != nil {
		return err
	}
	return nil
}

func (k K8sSecretBackend) UpdateSecret(secret model.Secret) error {
	panic("implement me")
}

func (k K8sSecretBackend) DeleteSecret(secret model.Secret) error {
	panic("implement me")
}

func (k K8sSecretBackend) createKubernetesSecret(secret model.Secret, namespace string) *v1.Secret {
	return &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secret.Name,
			Namespace: namespace,
		},
		StringData: secret.Data,
		Type:       "Opaque",
	}
}

func init() {
	Register(BackendTypeK8s, func() SecretBackend {
		kubeAPI, _ := createKubeAPI()
		return NewK8sSecretBackend(kubeAPI)
	})
}

func createKubeAPI() (*kubernetes.Clientset, error) {
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
