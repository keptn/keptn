package namespacemanager

import (
	"context"
	"errors"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/keptn/go-utils/pkg/common/kubeutils"
	keptn "github.com/keptn/go-utils/pkg/lib/keptn"
)

// INamespaceManager defines operations for initializing and configuring namespaces
type INamespaceManager interface {
	CreateNamespaceIfNotExists(nsName string) error
	InjectIstio(project string, stage string) error
}

// NamespaceManager is an implementation of INamespaceManager
type NamespaceManager struct {
	logger keptn.LoggerInterface
}

// NewNamespaceManager creates a new instance of a NamespaceManager
func NewNamespaceManager(logger keptn.LoggerInterface) *NamespaceManager {
	return &NamespaceManager{logger: logger}
}

// InitNamespaces initializes namespaces if they do not exist yet
func (p *NamespaceManager) CreateNamespaceIfNotExists(nsName string) error {
	namespaceManager, err := kubeutils.NewNamespaceManager(true)
	if err != nil {
		return err
	}
	exists, err := namespaceManager.ExistsNamespace(context.TODO(), nsName)
	if err != nil {
		return fmt.Errorf("error when checking availability of namespace: %v", err)
	}
	if exists {
		p.logger.Debug(fmt.Sprintf("Reuse existing namespace %s", nsName))
	} else {
		p.logger.Debug(fmt.Sprintf("Create new namespace %s", nsName))
		if err != namespaceManager.CreateNamespace(context.TODO(), nsName) {
			return fmt.Errorf("error when creating namespace %s: %v", nsName, err)
		}
	}
	return nil
}

// InjectIstio injects Istio into the namespace used for the project and stage by adding the label istio-injection
func (p *NamespaceManager) InjectIstio(project string, stage string) error {
	kubeClient, err := kubeutils.GetClientSet(true)
	if err != nil {
		return fmt.Errorf("error when getting kube API: %v", err)
	}
	namespaceName := project + "-" + stage
	namespace, err := kubeClient.CoreV1().Namespaces().Get(context.TODO(), namespaceName, v1.GetOptions{})
	if err != nil {
		return err
	}
	if namespace == nil {
		return errors.New("error when getting namespace")
	}

	p.logger.Info(fmt.Sprintf("Inject Istio to the %s namespace for blue-green deployments", namespaceName))

	if namespace.ObjectMeta.Labels == nil {
		namespace.ObjectMeta.Labels = make(map[string]string)
	}

	// add the label istio-injection to the namespace
	namespace.ObjectMeta.Labels["istio-injection"] = "enabled"
	_, err = kubeClient.CoreV1().Namespaces().Update(context.TODO(), namespace, v1.UpdateOptions{})
	return err
}
