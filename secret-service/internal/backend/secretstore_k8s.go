package backend

import (
	"github.com/keptn/keptn/secret-service/internal/common"
	"github.com/keptn/keptn/secret-service/internal/model"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const DefaultNamespace = "keptn"

type K8sSecretStore struct {
	KubeAPI                kubernetes.Interface
	KeptnNamespaceProvider common.StringSupplier
}

func NewK8sSecretStore(kubeAPI kubernetes.Interface) *K8sSecretStore {
	return &K8sSecretStore{
		KubeAPI:                kubeAPI,
		KeptnNamespaceProvider: common.EnvBasedStringSupplier("POD_NAMESPACE", DefaultNamespace),
	}
}

func (k K8sSecretStore) CreateSecret(secret model.Secret) error {

	namespace := k.KeptnNamespaceProvider()
	kubeSecret := k.createKubernetesSecret(secret, namespace)
	_, err := k.KubeAPI.CoreV1().Secrets(namespace).Create(kubeSecret)
	if err != nil {
		return err
	}
	return nil
}

func (k K8sSecretStore) UpdateSecret(secret model.Secret) error {
	panic("implement me")
}

func (k K8sSecretStore) DeleteSecret(secret model.Secret) error {
	panic("implement me")
}

func (k K8sSecretStore) createKubernetesSecret(secret model.Secret, namespace string) *v1.Secret {
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
