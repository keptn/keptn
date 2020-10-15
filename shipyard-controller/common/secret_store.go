package common

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
)

type SecretStore interface {
	CreateSecret(name string, content map[string][]byte) error
	DeleteSecret(name string) error
	GetSecret(name string) (map[string][]byte, error)
}

type K8sSecretStore struct {
	client *kubernetes.Clientset
}

// NewK8sSecretStore
func NewK8sSecretStore() (*K8sSecretStore, error) {
	client, err := GetKubeAPI()
	if err != nil {
		return nil, err
	}
	return &K8sSecretStore{client: client}, nil
}

func (k *K8sSecretStore) CreateSecret(name string, content map[string][]byte) error {
	namespace := os.Getenv("POD_NAMESPACE")
	secret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: content,
		Type: "Opaque",
	}
	_, err := k.client.CoreV1().Secrets(namespace).Create(secret)
	if err != nil {
		return err
	}
	return nil
}

func (k *K8sSecretStore) DeleteSecret(name string) error {
	namespace := os.Getenv("POD_NAMESPACE")
	return k.client.CoreV1().Secrets(namespace).Delete(name, &metav1.DeleteOptions{})
}

func (K8sSecretStore) GetSecret(name string) (map[string][]byte, error) {
	panic("implement me")
}

func GetKubeAPI() (*kubernetes.Clientset, error) {
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
