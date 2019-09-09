package helm

import (
	"os"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/pkg/serviceutils"
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
func (c *CanaryOnDeploymentGenerator) DeleteRelease(project string, stage string, service string, generated bool) error {
	useInClusterConfig := false
	if os.Getenv("ENVIRONMENT") == "production" {
		useInClusterConfig = true
	}

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return err
	}

	ch, err := keptnutils.GetChart(project, service, stage, GetChartName(service, generated), url.String())
	if err != nil {
		return err
	}
	for _, dpl := range keptnutils.GetDeployments(ch) {
		if err := keptnutils.ScaleDeployment(useInClusterConfig, dpl.Name, c.GetNamespace(project, stage, generated), 0); err != nil {
			return err
		}
	}
	return nil
}
