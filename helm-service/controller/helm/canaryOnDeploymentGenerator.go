package helm

import (
	keptnevents "github.com/keptn/go-utils/pkg/events"
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
