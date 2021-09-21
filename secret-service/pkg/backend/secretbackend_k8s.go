package backend

import (
	"context"
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
	_, err = k.KubeAPI.CoreV1().Secrets(namespace).Create(context.TODO(), k.createK8sSecretObj(secret, namespace), metav1.CreateOptions{})
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
		_, err := k.KubeAPI.RbacV1().Roles(namespace).Create(context.TODO(), &roles[i], metav1.CreateOptions{})
		if err != nil {
			if statusError, isStatus := err.(*k8serr.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
				log.Infof("Try to update role %s as it already exists", roles[i].Name)
				role, err := k.KubeAPI.RbacV1().Roles(namespace).Get(context.TODO(), roles[i].Name, metav1.GetOptions{})
				if err != nil {
					log.Errorf("Unable to get details of role %s", roles[i].Name)
					return err
				}
				role.Rules[0].ResourceNames = append(role.Rules[0].ResourceNames, secret.Name)
				if _, err := k.KubeAPI.RbacV1().Roles(namespace).Update(context.TODO(), role, metav1.UpdateOptions{}); err != nil {
					log.Errorf("Unable to update role %s", roles[i].Name)
					return err
				}
			} else {
				return err
			}
		}
	}

	roleBinding := k.createK8sRoleBindingObj(secret, roles, namespace)
	_, err = k.KubeAPI.RbacV1().RoleBindings(namespace).Create(context.TODO(), &roleBinding, metav1.CreateOptions{})
	if err != nil {
		if statusError, isStatus := err.(*k8serr.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
			//no op
		} else {
			log.Errorf("Unable to create role binding: %s", err.Error())
			return err
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

	// delete secret
	err = k.KubeAPI.CoreV1().Secrets(namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
	if err != nil {
		log.Errorf("Unable to delete secret %s with scope %s: %s", secret.Name, secret.Scope, err)
		if statusError, isStatus := err.(*k8serr.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonNotFound {
			return ErrSecretNotFound
		}
		return err
	}

	// get current secrets with scope
	secretsWithScope, err := k.KubeAPI.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "app.kubernetes.io/scope=" + secret.Scope})
	if err != nil {
		return err
	}
	// if there are no secret with that scope
	// delete associated roles and rolebinding
	if len(secretsWithScope.Items) == 0 {
		log.Infof("No more secret with scope: %s. Deleting associated roles and role bindings", secret.Scope)
		if err := k.KubeAPI.RbacV1().Roles(namespace).DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: "app.kubernetes.io/scope=" + secret.Scope}); err != nil {
			log.Warnf("Unable to delete roles: %s", err.Error())
		}
		if err := k.KubeAPI.RbacV1().RoleBindings(namespace).DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: "app.kubernetes.io/scope=" + secret.Scope}); err != nil {
			log.Warnf("Unable to delete role bindings: %s", err.Error())
		}
	} else {
		// update the role resources, otherwise
		roles := k.createK8sRoleObj(secret, scopes, namespace)
		for i := range roles {
			log.Infof("Updating role %s", roles[i].Name)
			role, err := k.KubeAPI.RbacV1().Roles(namespace).Get(context.TODO(), roles[i].Name, metav1.GetOptions{})
			if err != nil {
				log.Errorf("Unable to get details of role %s", roles[i].Name)
				return err
			}
			role.Rules[0].ResourceNames = remove(role.Rules[0].ResourceNames, secret.Name)
			if _, err := k.KubeAPI.RbacV1().Roles(namespace).Update(context.TODO(), role, metav1.UpdateOptions{}); err != nil {
				log.Errorf("Unable to update role %s", roles[i].Name)
				return err
			}
		}
	}

	return nil
}

func (k K8sSecretBackend) GetSecrets() ([]model.GetSecretResponseItem, error) {
	result := []model.GetSecretResponseItem{}

	namespace := k.KeptnNamespaceProvider()
	list, err := k.KubeAPI.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/managed-by=keptn-secret-service",
	})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve secrets: %s", err.Error())
	}

	for _, secretItem := range list.Items {
		keys := []string{}
		for key := range secretItem.StringData {
			if key != "" {
				keys = insert(keys, key)
			}
		}
		for key := range secretItem.Data {
			if key != "" {
				keys = insert(keys, key)
			}
		}
		result = append(result, model.GetSecretResponseItem{
			SecretMetadata: model.SecretMetadata{
				Name:  secretItem.Name,
				Scope: secretItem.Labels["app.kubernetes.io/scope"],
			},
			Keys: keys,
		})
	}

	return result, nil
}

func (k K8sSecretBackend) UpdateSecret(secret model.Secret) error {
	log.Infof("Updating secret: %s with scope %s", secret.Name, secret.Scope)
	namespace := k.KeptnNamespaceProvider()
	kubeSecret := k.createK8sSecretObj(secret, namespace)

	_, err := k.KubeAPI.CoreV1().Secrets(namespace).Update(context.TODO(), kubeSecret, metav1.UpdateOptions{})
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
					Labels: map[string]string{
						"app.kubernetes.io/managed-by": "keptn-secret-service",
						"app.kubernetes.io/scope":      secret.Scope,
					},
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

func (k K8sSecretBackend) createK8sRoleBindingObj(secret model.Secret, roles []rbacv1.Role, namespace string) rbacv1.RoleBinding {
	roleBindingName := secret.Scope + "-rolebinding"
	log.Infof("creating role binding %s for secret %s with role %s and service account %s in namespace %s", roleBindingName, secret.Name, roles[0].Name, secret.Scope, namespace)

	roleBinding := rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      roleBindingName,
			Namespace: namespace,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "keptn-secret-service", // add a 'managed-by' label so we can identify secrets managed by the secret-service
				"app.kubernetes.io/scope":      secret.Scope,
			},
		},
		Subjects: []rbacv1.Subject{rbacv1.Subject{
			Kind: "ServiceAccount",
			Name: secret.Scope,
		}},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     roles[0].Name,
		},
	}
	return roleBinding
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
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "keptn-secret-service", // add a 'managed-by' label so we can identify secrets managed by the secret-service
				"app.kubernetes.io/scope":      secret.Scope,
			},
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

func insert(s []string, i string) []string {
	for _, val := range s {
		if val == i {
			return s
		}
	}
	s = append(s, i)
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
