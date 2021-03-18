package backend

import (
	"github.com/keptn/keptn/secret-service/internal/common"
	"github.com/keptn/keptn/secret-service/internal/model"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const BackendTypeK8s = "kubernetes"

type K8sSecretBackend struct {
	KubeAPI                kubernetes.Interface
	KeptnNamespaceProvider common.StringSupplier
	Scopes                 model.Scopes
}

func NewK8sSecretBackend(kubeAPI kubernetes.Interface) *K8sSecretBackend {
	return &K8sSecretBackend{
		KubeAPI:                kubeAPI,
		KeptnNamespaceProvider: common.EnvBasedStringSupplier("POD_NAMESPACE", DefaultNamespace),
	}
}

func (k K8sSecretBackend) CreateSecret(secret model.Secret) error {

	namespace := k.KeptnNamespaceProvider()
	kubeSecret := k.createK8sSecretObj(secret, namespace)
	_, err := k.KubeAPI.CoreV1().Secrets(namespace).Create(kubeSecret)
	if err != nil {
		return err
	}

	roles := k.createK8sRoleObj(secret, namespace)
	for _, role := range roles {
		if _, err = k.KubeAPI.RbacV1().Roles(namespace).Create(&role); err != nil {
			return err
		}
	}

	return nil
}

func (k K8sSecretBackend) UpdateSecret(secret model.Secret) error {
	panic("implement me")
}

func (k K8sSecretBackend) DeleteSecret(secret model.Secret) error {
	panic("implement me")
}

func (k K8sSecretBackend) createK8sRoleObj(secret model.Secret, namespace string) []rbacv1.Role {

	var k8sRolesToCreate []rbacv1.Role

	if scope, ok := k.Scopes.Scopes[secret.Scope]; ok {
		for capName, cap := range scope.Capabilities {
			capPermissions := cap.Permissions
			role := rbacv1.Role{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Role",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      capName,
					Namespace: namespace,
				},
				Rules: []rbacv1.PolicyRule{
					{
						Verbs:         capPermissions,
						Resources:     []string{"secrets"},
						ResourceNames: []string{secret.Name},
					},
				},
			}
			k8sRolesToCreate = append(k8sRolesToCreate, role)
		}
	}
	return k8sRolesToCreate
}

func (k K8sSecretBackend) createK8sSecretObj(secret model.Secret, namespace string) *corev1.Secret {
	return &corev1.Secret{
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

func init() {
	Register(BackendTypeK8s, func() SecretBackend {
		kubeAPI, _ := createKubeAPI()
		return NewK8sSecretBackend(kubeAPI)
	})
}
