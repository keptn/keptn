package secretstore

import (
	"context"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/internal/common"

	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// SecretStore godoc
//go:generate moq -pkg fake -out ./fake/secretstore_mock.go . SecretStore
type SecretStore interface {
	// CreateSecret godoc
	CreateSecret(name string, content map[string][]byte) error
	// DeleteSecret godoc
	DeleteSecret(name string) error
	// GetSecret godoc
	GetSecret(name string) (map[string][]byte, error)
	// UpdateSecret godoc
	UpdateSecret(name string, content map[string][]byte) error
}

// K8sSecretStore godoc
type K8sSecretStore struct {
	client kubernetes.Interface
}

func New(kubeAPI kubernetes.Interface) *K8sSecretStore {
	return &K8sSecretStore{client: kubeAPI}
}

// CreateSecret godoc
func (k *K8sSecretStore) CreateSecret(name string, content map[string][]byte) error {
	namespace := common.GetKeptnNamespace()
	secret := k.createSecretObj(name, namespace, content)
	_, err := k.client.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// DeleteSecret godoc
func (k *K8sSecretStore) DeleteSecret(name string) error {
	namespace := common.GetKeptnNamespace()
	return k.client.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// GetSecret godoc
func (k *K8sSecretStore) GetSecret(name string) (map[string][]byte, error) {
	namespace := common.GetKeptnNamespace()
	get, err := k.client.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get git-credentials secret on kubernetes cluster")
	}
	return get.Data, nil
}

// UpdateSecret godoc
func (k *K8sSecretStore) UpdateSecret(name string, content map[string][]byte) error {
	namespace := common.GetKeptnNamespace()
	secret := k.createSecretObj(name, namespace, content)

	_, err := k.client.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		if err.(*k8serrors.StatusError).ErrStatus.Reason == "NotFound" {
			_, err := k.client.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (k *K8sSecretStore) createSecretObj(name string, namespace string, content map[string][]byte) *v1.Secret {
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
	return secret
}
