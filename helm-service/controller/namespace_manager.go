package controller

import (
	"errors"
	"fmt"

	"github.com/keptn/keptn/helm-service/controller/helm"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/keptn/go-utils/pkg/configuration-service/models"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

type NamespaceManager struct {
	logger keptnutils.LoggerInterface
}

func NewNamespaceManager(logger keptnutils.LoggerInterface) *NamespaceManager {
	return &NamespaceManager{logger: logger}
}

// InitNamespaces initializes namespaces if they do not exist yet
func (p *NamespaceManager) InitNamespaces(project string, stages []*models.Stage) error {

	for _, shipyardStage := range stages {

		namespace := helm.GetUmbrellaNamespace(project, shipyardStage.StageName)
		exists, err := keptnutils.ExistsNamespace(true, namespace)
		if err != nil {
			return fmt.Errorf("error when checking availablity of namespace: %v", err)
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
	namespace, err := kubeClient.Namespaces().Get(helm.GetUmbrellaNamespace(project, stage), v1.GetOptions{})
	if err != nil {
		return err
	}
	if namespace == nil {
		return errors.New("error when getting namespace")
	}

	p.logger.Debug(fmt.Sprintf("Inject Istio to the %s namespace for blue-green deployments", helm.GetUmbrellaNamespace(project, stage)))

	if namespace.ObjectMeta.Labels == nil {
		namespace.ObjectMeta.Labels = make(map[string]string)
	}

	namespace.ObjectMeta.Labels["istio-injection"] = "enabled"
	_, err = kubeClient.Namespaces().Update(namespace)
	return err
}
