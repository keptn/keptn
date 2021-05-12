package helm

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"io"
	"strings"

	utils "github.com/keptn/go-utils/pkg/api/utils"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
)

const (
	serviceNameMaxLen = 53
)

// GetServices returns all services contained in the Helm manifest
func GetServices(helmManifest string) []*corev1.Service {

	services := []*corev1.Service{}
	dec := kyaml.NewYAMLToJSONDecoder(strings.NewReader(helmManifest))
	for {
		var svc corev1.Service
		err := dec.Decode(&svc)
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		if keptnutils.IsService(&svc) {
			services = append(services, &svc)
		}
	}

	return services
}

// GetDeployments returns all deployments contained in the Helm manifest
func GetDeployments(helmManifest string) []*appsv1.Deployment {

	deployments := []*appsv1.Deployment{}
	dec := kyaml.NewYAMLToJSONDecoder(strings.NewReader(helmManifest))
	for {
		var dpl appsv1.Deployment
		err := dec.Decode(&dpl)
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		if keptnutils.IsDeployment(&dpl) {
			deployments = append(deployments, &dpl)
		}
	}
	return deployments
}

// GetChartName returns the name of the chart, e.g., ${SERVICE}[-generated]
func GetChartName(service string, generated bool) string {
	suffix := ""
	if generated {
		suffix = "-generated"
	}
	return service + suffix
}

// GetReleaseName returns the name of the Helm release, e.g., ${PROJECT}-${STAGE}-${SERVICE}[-generated]
func GetReleaseName(project string, stage string, service string, generated bool) string {
	suffix := ""
	if generated {
		suffix = "-generated"
	}

	proposedReleaseName := project + "-" + stage + "-" + service + suffix

	// check if it fits
	if len(proposedReleaseName) < serviceNameMaxLen {
		return proposedReleaseName
	}

	// alternative 1: remove project name
	proposedReleaseName = stage + "-" + service + suffix

	// check if it fits
	if len(proposedReleaseName) < serviceNameMaxLen {
		return proposedReleaseName
	}

	// alternative 2: also remove the stage name
	proposedReleaseName = service + suffix

	// check if it fits
	if len(proposedReleaseName) < serviceNameMaxLen {
		return proposedReleaseName
	}

	// It still doesn't fit ... We should really not get here, but it's an edge case.
	// Our last chance is to remove characters from right side of service name
	// Note: this can lead to ambiguous release names

	proposedReleaseName = service[:serviceNameMaxLen-len(suffix)-1] + suffix

	return proposedReleaseName
}

// DoesChartExist checks if the GIT repo contains the specified chart
func DoesChartExist(event keptnv2.EventData, chartName string, configServiceURL string) (bool, error) {
	resourceHandler := utils.NewResourceHandler(configServiceURL)

	helmChartURI := "helm/" + chartName + ".tgz"
	_, err := resourceHandler.GetServiceResource(event.Project, event.Stage, event.Service, helmChartURI)
	if err == utils.ResourceNotFoundError {
		return false, nil
	} else if err == nil {
		return true, nil
	}
	return false, err
}
