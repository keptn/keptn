package actions

import (
	"os"

	configutils "github.com/keptn/go-utils/pkg/configuration-service/utils"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	"sigs.k8s.io/yaml"
)

func getVirtualServiceUri(service string) string {

	return "helm/" + service + "-generated/templates/" + service + "-istio-virtualservice.yaml"
}

func containsVirtualServices(project, stage, service string) (bool, error) {

	handler := configutils.NewResourceHandler(os.Getenv(envConfigSvcURL))
	resources, err := handler.GetAllServiceResources(project, stage, service)
	if err != nil {
		return false, err
	}
	for _, resource := range resources {
		if resource.ResourceURI != nil && *resource.ResourceURI == getVirtualServiceUri(service) {
			return true, nil
		}
	}
	return false, nil
}

func getServices(project string, stage string) ([]string, error) {

	handler := configutils.NewServiceHandler(os.Getenv(envConfigSvcURL))
	services, err := handler.GetAllServices(project, stage)
	if err != nil {
		return nil, err
	}
	var serviceNames []string
	for _, svc := range services {
		serviceNames = append(serviceNames, svc.ServiceName)
	}
	return serviceNames, nil
}

type ProblemDetails struct {
	ClientIP string `json:"ClientIP"`
}

func getIP(problem *keptnevents.ProblemEventData) (string, error) {

	details := ProblemDetails{}
	err := yaml.Unmarshal(problem.ProblemDetails, &details)
	if err != nil {
		return "", err
	}

	return details.ClientIP, nil
}
