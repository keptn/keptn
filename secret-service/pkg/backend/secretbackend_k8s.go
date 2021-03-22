package backend

import (
	"errors"
	"flag"
	"fmt"
	"github.com/keptn/keptn/secret-service/pkg/common"
	"github.com/keptn/keptn/secret-service/pkg/model"
	"github.com/keptn/keptn/secret-service/pkg/repository"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" //TODO: delete
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

const SecretBackendTypeK8s = "kubernetes"

var ErrSecretAlreadyExists = errors.New("Secret already exists")
var ErrSecretNotFound = errors.New("Secret not found")

type K8sSecretBackend struct {
	KubeAPI                kubernetes.Interface
	KeptnNamespaceProvider common.StringSupplier
	ScopesRepository       repository.ScopesRepository
}

func NewK8sSecretBackend(kubeAPI kubernetes.Interface, scopesRepository repository.ScopesRepository) *K8sSecretBackend {
	return &K8sSecretBackend{
		KubeAPI:                kubeAPI,
		KeptnNamespaceProvider: common.EnvBasedStringSupplier("POD_NAMESPACE", DefaultNamespace),
		ScopesRepository:       scopesRepository,
	}
}

func (k K8sSecretBackend) CreateSecret(secret model.Secret) error {

	scopes, err := k.ScopesRepository.Read()
	if err != nil {
		return err
	}
	if _, ok := scopes.Scopes[secret.Scope]; !ok {
		return fmt.Errorf("Scope %s not available for creation of Secret %s", secret.Scope, secret.Name)
	}

	namespace := k.KeptnNamespaceProvider()
	kubeSecret := k.createK8sSecretObj(secret, namespace)
	_, err = k.KubeAPI.CoreV1().Secrets(namespace).Create(kubeSecret)
	if err != nil {
		if statusError, isStatus := err.(*kubeerrors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
			return ErrSecretAlreadyExists
		}
		return err
	}

	roles := k.createK8sRoleObj(secret, scopes, namespace)
	for i := range roles {
		if _, err = k.KubeAPI.RbacV1().Roles(namespace).Create(&roles[i]); err != nil {
			return err
		}
	}

	return nil
}

func (k K8sSecretBackend) UpdateSecret(secret model.Secret) error {
	panic("implement me")
}

func (k K8sSecretBackend) DeleteSecret(secret model.Secret) error {
	namespace := k.KeptnNamespaceProvider()
	secretName := secret.Name

	err := k.KubeAPI.CoreV1().Secrets(namespace).Delete(secretName, &metav1.DeleteOptions{})
	if err != nil {
		if statusError, isStatus := err.(*kubeerrors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonNotFound {
			return ErrSecretNotFound
		}
		return err
	}
	return nil
}

func (k K8sSecretBackend) createK8sRoleObj(secret model.Secret, scopes model.Scopes, namespace string) []rbacv1.Role {

	var k8sRolesToCreate []rbacv1.Role

	if scope, ok := scopes.Scopes[secret.Scope]; ok {
		for capabilityName, capability := range scope.Capabilities {
			capabilityPermissions := capability.Permissions
			role := rbacv1.Role{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Role",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      capabilityName,
					Namespace: namespace,
				},
				Rules: []rbacv1.PolicyRule{
					{
						Verbs:         capabilityPermissions,
						APIGroups:     []string{""}, // TODO: what to put here?
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
	Register(SecretBackendTypeK8s, func() SecretBackend {

		// TODO: DELETE TILL *END*
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}

		// create the clientset
		kubeAPI, err := kubernetes.NewForConfig(config)
		// *END*

		//TODO: enable next line
		//kubeAPI, _ := createKubeAPI()
		scopesRepository := repository.NewFileBasedScopesRepository()
		return NewK8sSecretBackend(kubeAPI, scopesRepository)
	})
}
