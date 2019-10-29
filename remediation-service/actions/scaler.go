package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/keptn/keptn/remediation-service/pkg/utils"

	"github.com/ghodss/yaml"

	"k8s.io/helm/pkg/proto/hapi/chart"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnmodels "github.com/keptn/go-utils/pkg/models"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
)

type Scaler struct {
}

// NewScaler creates a new Scaler
func NewScaler() *Scaler {
	return &Scaler{}
}

func (s Scaler) GetAction() string {
	return "scaling"
}

func (s Scaler) ExecuteAction(problem *keptnevents.ProblemEventData, shkeptncontext string,
	action *keptnmodels.RemediationAction) error {

	replicaIncrement, err := strconv.Atoi(action.Value)
	if err != nil {
		return fmt.Errorf("could not parse Value of action: %v", err)
	}

	helmChartName := problem.Service + "-generated"
	// Read chart
	ch, err := keptnutils.GetChart(problem.Project, problem.Service, problem.Stage, helmChartName, os.Getenv(envConfigSvcURL))
	if err != nil {
		return fmt.Errorf("cannot get chart %s for service %s in stage %s of project %s: %v", helmChartName,
			problem.Service, problem.Stage, problem.Project, err)
	}

	changedTemplates, err := s.increaseReplicaCount(ch, replicaIncrement)
	if err != nil {
		return fmt.Errorf("failed to increase replica count: %v", err)
	}

	changedFiles := make(map[string]string)
	for _, template := range changedTemplates {
		changedFiles[template.Name] = string(template.Data)
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
	return nil
}

// increases the replica count in the deployments by the provided replicaIncrement
func (s Scaler) increaseReplicaCount(ch *chart.Chart, replicaIncrement int) ([]*chart.Template, error) {

	changedTemplates := make([]*chart.Template, 0, 0)

	for _, template := range ch.Templates {
		dec := kyaml.NewYAMLToJSONDecoder(bytes.NewReader(template.Data))
		newContent := make([]byte, 0, 0)
		containsDepl := false
		for {
			var document interface{}
			err := dec.Decode(&document)
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}

			doc, err := json.Marshal(document)
			if err != nil {
				return nil, err
			}

			var depl appsv1.Deployment
			if err := json.Unmarshal(doc, &depl); err == nil && keptnutils.IsDeployment(&depl) {
				// Deployment found
				containsDepl = true
				depl.Spec.Replicas = getPtr(*depl.Spec.Replicas + int32(replicaIncrement))
				newContent, err = appendAsYaml(newContent, depl)
				if err != nil {
					return nil, err
				}
			} else {
				newContent, err = appendAsYaml(newContent, document)
				if err != nil {
					return nil, err
				}
			}
		}
		if containsDepl {
			template.Data = newContent
			changedTemplates = append(changedTemplates, template)
		}
	}

	return changedTemplates, nil
}

func getPtr(x int32) *int32 {
	return &x
}

func appendAsYaml(content []byte, element interface{}) ([]byte, error) {

	jsonData, err := json.Marshal(element)
	if err != nil {
		return nil, err
	}
	yamlData, err := yaml.JSONToYAML(jsonData)
	if err != nil {
		return nil, err
	}
	content = append(content, []byte("---\n")...)
	return append(content, yamlData...), nil
}
