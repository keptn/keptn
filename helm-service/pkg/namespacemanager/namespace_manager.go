package namespacemanager

import (
	"errors"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	keptn "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
)

type NamespaceManager struct {
	logger keptn.LoggerInterface
}

func NewNamespaceManager(logger keptn.LoggerInterface) *NamespaceManager {
	return &NamespaceManager{logger: logger}
}

// InitNamespaces initializes namespaces if they do not exist yet
func (p *NamespaceManager) InitNamespaces(project string, stages []string) error {

	for _, stage := range stages {

		namespace := project + "-" + stage
		exists, err := keptnutils.ExistsNamespace(true, namespace)
		if err != nil {
			return fmt.Errorf("error when checking availability of namespace: %v", err)
		}
		if exists {
			p.logger.Debug(fmt.Sprintf("Reuse existing namespace %s", namespace))
		} else {
			p.logger.Debug(fmt.Sprintf("Create new namespace %s", namespace))
			if err != keptnutils.CreateNamespace(true, namespace) {
				return fmt.Errorf("error when creating namespace %s: %v", namespace, err)
			}
		}
	}
	return nil
}

// InjectIstio injects Istio into the namespace used for the project and stage
func (p *NamespaceManager) InjectIstio(project string, stage string) error {
	kubeClient, err := keptnutils.GetKubeAPI(true)
	if err != nil {
		return fmt.Errorf("error when getting kube API: %v", err)
	}
	namespaceName := project + "-" + stage
	namespace, err := kubeClient.Namespaces().Get(namespaceName, v1.GetOptions{})
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

	namespace.ObjectMeta.Labels["istio-injection"] = "enabled"
	_, err = kubeClient.Namespaces().Update(namespace)
	return err
}
