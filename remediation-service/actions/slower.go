package actions

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/keptn/keptn/remediation-service/pkg/utils"

	configutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib"
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

func (s Slower) ExecuteAction(problem *keptn.ProblemEventData, keptnHandler *keptn.Keptn,
	action *keptn.RemediationAction) error {
	return s.executor(problem, keptnHandler, action, s.addDelay)
}

func (s Slower) ResolveAction(problem *keptn.ProblemEventData, keptnHandler *keptn.Keptn,
	action *keptn.RemediationAction) error {
	return s.executor(problem, keptnHandler, action, s.removeDelay)
}

func (s Slower) executor(problem *keptn.ProblemEventData, keptnHandler *keptn.Keptn,
	action *keptn.RemediationAction, editVs func(vsContent, ip string, slowDown string) (string, error)) error {

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
				return fmt.Errorf("could not get VirtualService resource: %v", err)
			}

			newVS, err := editVs(resource.ResourceContent, ip, slowDown)
			if err != nil {
				return fmt.Errorf("failed to add delay: %v", err)
			}

			changedFiles := map[string]string{
				getVirtualServiceUriInChart(service): newVS,
			}

			data := keptn.ConfigurationChangeEventData{
				Project:                   problem.Project,
				Service:                   problem.Service,
				Stage:                     problem.Stage,
				FileChangesGeneratedChart: changedFiles,
			}

			err = utils.CreateAndSendConfigurationChangedEvent(problem, keptnHandler, data)
			if err != nil {
				return fmt.Errorf("failed to send configuration change event: %v", err)
			}
		}
	}
	return nil
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
			"X-Forwarded-For": {
				Exact: ip,
			},
		},
	}

	if len(vs.Spec.Http) > 0 {
		newRoute := new(v1alpha3.HTTPRoute)
		deepCopy(vs.Spec.Http[len(vs.Spec.Http)-1], newRoute)

		newRoute.Fault = &fault
		newRoute.Match = append(newRoute.Match, &match)

		vs.Spec.Http = append([]*v1alpha3.HTTPRoute{newRoute}, vs.Spec.Http...)

		data, err := yaml.Marshal(vs)
		if err != nil {
			return "", err
		}
		return string(data), err
	}
	return "", errors.New("failed to add fault because no route is available")
}

func (s Slower) removeDelay(vsContent string, ip string, slowDown string) (string, error) {

	vs := v1alpha3.VirtualService{}
	err := yaml.Unmarshal([]byte(vsContent), &vs)
	if err != nil {
		return "", err
	}

	length := len(vs.Spec.Http)
	for i := 0; i < length; i++ {
		if vs.Spec.Http[i].Fault != nil && vs.Spec.Http[i].Fault.Delay != nil && len(vs.Spec.Http[i].Match) > 0 {
			for _, match := range vs.Spec.Http[i].Match {
				if val, ok := match.Headers["X-Forwarded-For"]; ok && val.Exact == ip {
					// found; delete the route
					copy(vs.Spec.Http[i:], vs.Spec.Http[i+1:])
					vs.Spec.Http = vs.Spec.Http[:len(vs.Spec.Http)-1]
					i--
					length--
					break
				}
			}
		}
	}

	data, err := yaml.Marshal(vs)
	if err != nil {
		return "", err
	}
	return string(data), nil

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
