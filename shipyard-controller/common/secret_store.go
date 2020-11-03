package common

import (
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// SecretStore godoc
type SecretStore interface {
	// CreateSecret godoc
	CreateSecret(name string, content map[string][]byte) error
	// DeleteSecret godoc
	DeleteSecret(name string) error
	// GetSecret godoc
	GetSecret(name string) (map[string][]byte, error)
}

// K8sSecretStore godoc
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

// CreateSecret godoc
func (k *K8sSecretStore) CreateSecret(name string, content map[string][]byte) error {
	namespace := GetKeptnNamespace()
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

// DeleteSecret godoc
func (k *K8sSecretStore) DeleteSecret(name string) error {
	namespace := GetKeptnNamespace()
	return k.client.CoreV1().Secrets(namespace).Delete(name, &metav1.DeleteOptions{})
}

// GetSecret godoc
func (k *K8sSecretStore) GetSecret(name string) (map[string][]byte, error) {
	namespace := GetKeptnNamespace()
	get, err := k.client.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return get.Data, nil
}

// GetKubeAPI godoc
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
