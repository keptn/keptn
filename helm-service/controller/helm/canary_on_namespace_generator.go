package helm

import (
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	corev1 "k8s.io/api/core/v1"
)

// CanaryOnNamespaceGenerator implements functions for doing a canary on the namespace level
type CanaryOnNamespaceGenerator struct {
}

// NewCanaryOnNamespaceGenerator creates a new CanaryOnNamespaceGenerator
func NewCanaryOnNamespaceGenerator() *CanaryOnNamespaceGenerator {
	return &CanaryOnNamespaceGenerator{}
}

// GetCanaryService returns a service which can be used for canary releases
func (c *CanaryOnNamespaceGenerator) GetCanaryService(originalSvc corev1.Service, project string, stageName string) (canaryService *corev1.Service) {
	canaryService = &corev1.Service{}

	canaryService.Kind = "Service"
	canaryService.APIVersion = "v1"
	canaryService.Name = originalSvc.Name + "-canary"
	canaryService.Spec.Type = "ExternalName"
	canaryService.Spec.ExternalName = originalSvc.Name + "." + c.GetNamespace(project, stageName, false) + ".svc.cluster.local"
	canaryService.Spec.Ports = originalSvc.Spec.Ports
	return
}

// IsK8sResourceDuplicated shows whether a resource is duplicated or not
func (*CanaryOnNamespaceGenerator) IsK8sResourceDuplicated() bool {
	return true
}

// GetNamespace returns the namespace for a specific project and stage
func (*CanaryOnNamespaceGenerator) GetNamespace(project string, stage string, generated bool) string {
	suffix := ""
	if generated {
		suffix = "-generated"
	}
	return project + "-" + stage + suffix
}

// DeleteRelease deletes the release by deleting the helm chart
func (*CanaryOnNamespaceGenerator) DeleteCanaryRelease(project string, stage string, service string) error {
	releaseName := GetReleaseName(project, stage, service, false)
	if _, err := keptnutils.ExecuteCommand("helm", []string{"delete", releaseName, "--purge"}); err != nil {
		return err
	}
	return nil
}
