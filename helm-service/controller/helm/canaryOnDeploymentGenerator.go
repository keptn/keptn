package helm

import (
	"os"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	corev1 "k8s.io/api/core/v1"
)

type CanaryOnDeploymentGenerator struct {
}

func NewCanaryOnDeploymentGenerator() *CanaryOnDeploymentGenerator {
	return &CanaryOnDeploymentGenerator{}
}

func (*CanaryOnDeploymentGenerator) GetCanaryService(originalSvc corev1.Service, event *keptnevents.ServiceCreateEventData, stageName string) (canaryService *corev1.Service) {

	canaryService = originalSvc.DeepCopy()
	canaryService.Name = canaryService.Name + "-canary"
	return
}

func (*CanaryOnDeploymentGenerator) IsK8sResourceDuplicated() bool {

	return false
}

// GetNamespace returns the namespace for a specific project and stage
func (*CanaryOnDeploymentGenerator) GetNamespace(project string, stage string, generated bool) string {
	return project + "-" + stage
}

func (*CanaryOnDeploymentGenerator) DeleteRelease(project string, service string, stage string, generated bool, configServiceURL string) error {

	useInClusterConfig := false
	if os.Getenv("env") == "production" {
		useInClusterConfig = true
	}

	ch, err := GetChart(project, service, stage, GetChartName(service, generated), configServiceURL)
	if err != nil {
		return err
	}
	for _, dpl := range GetDeployments(ch) {
		if err := keptnutils.ScaleDeployment(useInClusterConfig, dpl.Name, dpl.Namespace, 0); err != nil {
			return err
		}
	}
	return nil
}
