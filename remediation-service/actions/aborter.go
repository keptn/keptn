package actions

import (
	"fmt"
	"os"

	"sigs.k8s.io/yaml"

	"github.com/keptn/keptn/remediation-service/pkg/utils"

	configutils "github.com/keptn/go-utils/pkg/configuration-service/utils"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnmodels "github.com/keptn/go-utils/pkg/models"
	"github.com/keptn/keptn/remediation-service/pkg/apis/networking/istio/v1alpha3"
)

type Aborter struct {
}

// NewAborter creates a new Aborter
func NewAborter() *Aborter {
	return &Aborter{}
}

func (a Aborter) GetAction() string {
	return "abort"
}

func (a Aborter) ExecuteAction(problem *keptnevents.ProblemEventData, shkeptncontext string,
	action *keptnmodels.RemediationAction) error {

	ip, err := getIP(problem)
	if err != nil {
		return fmt.Errorf("could not parse ip from ProblemDetails: %v", err)
	}

	var services []string
	if problem.Service != "" {
		services = append(services, problem.Service)
	} else {
		s, err := getServices(problem.Project, problem.Stage)
		if err != nil {
			return fmt.Errorf("could not get services: %v", err)
		}
		services = s
	}

	for _, service := range services {

		containsVS, err := containsVirtualServices(problem.Project, problem.Stage, service)
		if err != nil {
			return fmt.Errorf("failed to check if service %s in project %s and stage %s"+
				" contains a VirtualService: %v", service, problem.Project, problem.Stage, err)
		}

		if containsVS {
			handler := configutils.NewResourceHandler(os.Getenv(envConfigSvcURL))
			resource, err := handler.GetServiceResource(problem.Project, problem.Stage, service,
				getVirtualServiceUri(service))
			if err != nil {
				return fmt.Errorf("could not get virutal service resource: %v", err)
			}

			newVS, err := a.addAbort(resource.ResourceContent, ip)
			if err != nil {
				return fmt.Errorf("failed to add delay: %v", err)
			}

			changedFiles := map[string]string{
				getVirtualServiceUri(service): newVS,
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

func (a Aborter) addAbort(vsContent string, ip string) (string, error) {

	vs := v1alpha3.VirtualService{}
	err := yaml.Unmarshal([]byte(vsContent), &vs)
	if err != nil {
		return "", err
	}
	fault := v1alpha3.HTTPFaultInjection{
		Abort: &v1alpha3.HTTPFaultInjection_Abort{
			Percent:    int32(100),
			HttpStatus: int32(403),
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
