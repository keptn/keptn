package backend

import (
	"errors"
	"fmt"
	"github.com/keptn/keptn/secret-service/pkg/common"
	"github.com/keptn/keptn/secret-service/pkg/model"
	"github.com/keptn/keptn/secret-service/pkg/repository"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const SecretBackendTypeK8s = "kubernetes"

var ErrSecretAlreadyExists = errors.New("secret already exists")
var ErrSecretNotFound = errors.New("secret not found")

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

func (k K8sSecretBackend) checkScopeDefined(secret model.Secret) (model.Scopes, error) {
	scopes, err := k.ScopesRepository.Read()
	if err != nil {
		return model.Scopes{}, err
	}
	if _, ok := scopes.Scopes[secret.Scope]; !ok {
		log.Errorf("Unable to find scope %s for secret %s", secret.Scope, secret.Name)
		return model.Scopes{}, fmt.Errorf("scope %s not available for creation of Secret %s", secret.Scope, secret.Name)
	}
	return scopes, nil
}

func (k K8sSecretBackend) CreateSecret(secret model.Secret) error {
	log.Infof("Creating secret: %s with scope %s", secret.Name, secret.Scope)
	scopes, err := k.checkScopeDefined(secret)
	if err != nil {
		return err
	}
	namespace := k.KeptnNamespaceProvider()
	_, err = k.KubeAPI.CoreV1().Secrets(namespace).Create(k.createK8sSecretObj(secret, namespace))
	if err != nil {
		log.Errorf("Unable to create secret %s with scope %s: %s", secret.Name, secret.Scope, err)
		if statusError, isStatus := err.(*k8serr.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
			return ErrSecretAlreadyExists
		}
		return err
	}

	roles := k.createK8sRoleObj(secret, scopes, namespace)
	for i := range roles {
		log.Infof("Creating role %s", roles[i].Name)
		_, err := k.KubeAPI.RbacV1().Roles(namespace).Create(&roles[i])
		if err != nil {
			if statusError, isStatus := err.(*k8serr.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
				log.Infof("Try to update role %s as it already exists", roles[i].Name)
				role, err := k.KubeAPI.RbacV1().Roles(namespace).Get(roles[i].Name, metav1.GetOptions{})
				if err != nil {
					log.Errorf("Unable to get details of role %s", roles[i].Name)
					return err
				}
				role.Rules[0].ResourceNames = append(role.Rules[0].ResourceNames, secret.Name)
				if _, err := k.KubeAPI.RbacV1().Roles(namespace).Update(role); err != nil {
					log.Errorf("Unable to update role %s", roles[i].Name)
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}

func (k K8sSecretBackend) DeleteSecret(secret model.Secret) error {
	log.Infof("Deleting secret: %s with scope %s", secret.Name, secret.Scope)
	scopes, err := k.checkScopeDefined(secret)
	if err != nil {
		return err
	}
	namespace := k.KeptnNamespaceProvider()
	secretName := secret.Name

	err = k.KubeAPI.CoreV1().Secrets(namespace).Delete(secretName, &metav1.DeleteOptions{})
	if err != nil {
		log.Errorf("Unable to delete secret %s with scope %s: %s", secret.Name, secret.Scope, err)
		if statusError, isStatus := err.(*k8serr.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonNotFound {
			return ErrSecretNotFound
		}
		return err
	}

	roles := k.createK8sRoleObj(secret, scopes, namespace)
	for i := range roles {
		log.Infof("Updating role %s", roles[i].Name)
		role, err := k.KubeAPI.RbacV1().Roles(namespace).Get(roles[i].Name, metav1.GetOptions{})
		if err != nil {
			log.Errorf("Unable to get details of role %s", roles[i].Name)
			return err
		}
		role.Rules[0].ResourceNames = remove(role.Rules[0].ResourceNames, secret.Name)
		if _, err := k.KubeAPI.RbacV1().Roles(namespace).Update(role); err != nil {
			log.Errorf("Unable to update role %s", roles[i].Name)
			return err
		}
	}

	return nil
}

func (k K8sSecretBackend) UpdateSecret(secret model.Secret) error {
	log.Infof("Updating secret: %s with scope %s", secret.Name, secret.Scope)
	namespace := k.KeptnNamespaceProvider()
	kubeSecret := k.createK8sSecretObj(secret, namespace)

	_, err := k.KubeAPI.CoreV1().Secrets(namespace).Update(kubeSecret)
	if err != nil {
		log.Errorf("Unable to update secret %s: %s", secret.Name, err)
		if statusError, isStatus := err.(*k8serr.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonNotFound {
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
						APIGroups:     []string{""},
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

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
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
	log.Info("Registering Secret Backend type: kubernetes")
	Register(SecretBackendTypeK8s, func() SecretBackend {
		kubeAPI, err := createKubeAPI()
		if err != nil {
			log.Fatalf("Unable to create kubernetes client: %s", err)
		}
		scopesRepository := repository.NewFileBasedScopesRepository()
		return NewK8sSecretBackend(kubeAPI, scopesRepository)
	})
}
