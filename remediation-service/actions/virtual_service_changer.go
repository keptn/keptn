package actions

import (
	"encoding/json"
	"fmt"
	"os"

	configutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib"
	"sigs.k8s.io/yaml"
)

func getVirtualServiceUriInChart(service string) string {
	return "templates/" + service + "-istio-virtualservice.yaml"
}

func getVirtualServiceUri(service string) string {

	return "helm/" + service + "-generated/" + getVirtualServiceUriInChart(service)
}

type respError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func deepCopy(a, b interface{}) {
	byt, _ := json.Marshal(a)
	json.Unmarshal(byt, b)
}

func containsVirtualServices(project, stage, service string) (bool, error) {

	handler := configutils.NewResourceHandler(os.Getenv(envConfigSvcURL))
	_, err := handler.GetServiceResource(project, stage, service,
		getVirtualServiceUri(service))
	if err != nil {
		respError := respError{}
		json.Unmarshal([]byte(err.Error()), &respError)
		if respError.Code == 404 {
			return false, nil
		}
		return false, fmt.Errorf("could not get VirtualService resource: %v", err)
	}
	return true, nil
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

func getIP(problem *keptn.ProblemEventData) (string, error) {

	details := ProblemDetails{}
	err := yaml.Unmarshal(problem.ProblemDetails, &details)
	if err != nil {
		return "", err
	}

	return details.ClientIP, nil
}
