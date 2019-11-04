package helm

import (
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	corev1 "k8s.io/api/core/v1"
)

// CanaryOnDeploymentGenerator implements functions for doing a canary on the deployment level
type CanaryOnDeploymentGenerator struct {
}

// NewCanaryOnDeploymentGenerator creates a new CanaryOnDeploymentGenerator
func NewCanaryOnDeploymentGenerator() *CanaryOnDeploymentGenerator {
	return &CanaryOnDeploymentGenerator{}
}

// GetCanaryService returns a service which can be used for canary releases
func (*CanaryOnDeploymentGenerator) GetCanaryService(originalSvc corev1.Service, project string, stageName string) (canaryService *corev1.Service) {
	canaryService = originalSvc.DeepCopy()
	canaryService.Name = canaryService.Name + "-canary"
	return
}

// IsK8sResourceDuplicated shows whether a resource is duplicated or not
func (*CanaryOnDeploymentGenerator) IsK8sResourceDuplicated() bool {
	return false
}

// GetNamespace returns the namespace for a specific project and stage
func (*CanaryOnDeploymentGenerator) GetNamespace(project string, stage string, generated bool) string {
	return project + "-" + stage
}

// DeleteRelease deletes the release by scaling the deployments down to zero
func (c *CanaryOnDeploymentGenerator) DeleteCanaryRelease(project string, stage string, service string) error {
	releaseName := GetReleaseName(project, stage, service, false)
	if _, err := keptnutils.ExecuteCommand("helm", []string{"delete", releaseName, "--purge"}); err != nil {
		return err
	}
	return nil
}
