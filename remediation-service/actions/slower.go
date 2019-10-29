package actions

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/keptn/keptn/remediation-service/pkg/utils"

	configutils "github.com/keptn/go-utils/pkg/configuration-service/utils"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnmodels "github.com/keptn/go-utils/pkg/models"
	"github.com/keptn/keptn/remediation-service/pkg/apis/networking/istio/v1alpha3"
)

type Slower struct {
}

// NewScaler creates a new Scaler
func NewSlower() *Slower {
	return &Slower{}
}

func (s Slower) GetAction() string {
	return "slowdown"
}

func (s Slower) ExecuteAction(problem *keptnevents.ProblemEventData, shkeptncontext string,
	action *keptnmodels.RemediationAction) error {

	slowDown := strings.TrimSpace(action.Value)
	validFormat := s.validateActionFormat(slowDown)
	if !validFormat {
		return fmt.Errorf("could not parse action: %s", action.Value)
	}

	ip, err := getIP(problem)
	if err != nil {
		return fmt.Errorf("could not parse ip from ProblemDetails: %v", err)
	}

	var services []string
	if problem.Service != "" {
		services = append(services, problem.Service)
	} else {
		s, err := s.getServices(problem.Project, problem.Stage)
		if err != nil {
			return fmt.Errorf("could not get services: %v", err)
		}
		services = s
	}

	for _, service := range services {

		containsVS, err := s.containsVirtualServices(problem.Project, problem.Stage, service)
		if err != nil {
			return fmt.Errorf("failed to check if service %s in project %s and stage %s"+
				" contains a VirtualService: %v", service, problem.Project, problem.Stage, err)
		}

		if containsVS {
			handler := configutils.NewResourceHandler(os.Getenv(envConfigSvcURL))
			resource, err := handler.GetServiceResource(problem.Project, problem.Stage, service,
				s.getVirtualServiceUri(service))
			if err != nil {
				return fmt.Errorf("could not get virutal service resource: %v", err)
			}

			newVS, err := s.addDelay(resource.ResourceContent, ip, slowDown)
			if err != nil {
				return fmt.Errorf("failed to add delay: %v", err)
			}

			changedFiles := map[string]string{
				s.getVirtualServiceUri(service): newVS,
			}

			data := keptnevents.ConfigurationChangeEventData{
				Project:                   problem.Project,
				Service:                   problem.Service,
				Stage:                     problem.Stage,
				FileChangesGeneratedChart: changedFiles,
			}

			err = utils.CreateAndSendConfigurationChangedEvent(problem, shkeptncontext, data)
			if err != nil {
				return fmt.Errorf("failed to send configuration change event: %v", err)
			}
		}
	}
	return nil
}

type ProblemDetails struct {
	ClientIP string `json:"ClientIP"`
}

func (s Slower) getIP(problem *keptnevents.ProblemEventData) (string, error) {

	details := ProblemDetails{}
	err := yaml.Unmarshal(problem.ProblemDetails, &details)
	if err != nil {
		return "", err
	}

	return details.ClientIP, nil
}

func (s Slower) addDelay(vsContent string, ip string, slowDown string) (string, error) {

	vs := v1alpha3.VirtualService{}
	err := yaml.Unmarshal([]byte(vsContent), &vs)
	if err != nil {
		return "", err
	}
	fault := v1alpha3.HTTPFaultInjection{
		Delay: &v1alpha3.HTTPFaultInjection_Delay{
			FixedDelay: slowDown,
			Percent:    int32(100),
		},
	}
	match := v1alpha3.HTTPMatchRequest{
		Headers: map[string]*v1alpha3.StringMatch_Exact{
			"X-Forwarded-For": &v1alpha3.StringMatch_Exact{
				Exact: ip,
			},
		},
	}

	for _, httpRoute := range vs.Spec.Http {
		httpRoute.Fault = &fault
		httpRoute.Match = append(httpRoute.Match, &match)
	}

	data, err := yaml.Marshal(vs)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func (s Slower) getVirtualServiceUri(service string) string {

	return "helm/" + service + "-generated/templates/" + service + "-istio-virtualservice.yaml"
}

func (s Slower) containsVirtualServices(project, stage, service string) (bool, error) {

	handler := configutils.NewResourceHandler(os.Getenv(envConfigSvcURL))
	resources, err := handler.GetAllServiceResources(project, stage, service)
	if err != nil {
		return false, err
	}
	for _, resource := range resources {
		if resource.ResourceURI != nil && *resource.ResourceURI == s.getVirtualServiceUri(service) {
			return true, nil
		}
	}
	return false, nil
}

func (s Slower) getServices(project string, stage string) ([]string, error) {

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

func (s Slower) validateActionFormat(action string) bool {

	if !strings.HasSuffix(action, "s") {
		return false
	}
	_, err := strconv.Atoi(action[:len(action)-1])
	if err != nil {
		return false
	}
	return true
}
