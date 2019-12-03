package event_handler

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StartEvaluationHandler struct {
	Logger *keptnutils.Logger
	Event  cloudevents.Event
}

func (eh *StartEvaluationHandler) HandleEvent() error {
	var keptnContext string
	_ = eh.Event.ExtensionAs("shkeptncontext", &keptnContext)

	e := &keptnevents.StartEvaluationEventData{}
	err := eh.Event.DataAs(e)

	if err != nil {
		eh.Logger.Error("Could not parse event payload: " + err.Error())
		return err
	}

	// functional tests dont need to be evaluated
	if e.TestStrategy == "functional" || e.TestStrategy == "" {
		eh.Logger.Debug("Functional tests are need evaluated")
		evaluationDetails := keptnevents.EvaluationDetails{
			IndicatorResults: nil,
			TimeStart:        e.Start,
			TimeEnd:          e.End,
			Result:           fmt.Sprintf("no evaluation performed by lighthouse service (TestStrategy=%s)", e.TestStrategy),
		}
		// send the evaluation-done-event
		evaluationResult := keptnevents.EvaluationDoneEventData{
			EvaluationDetails:  &evaluationDetails,
			Result:             "pass",
			Project:            e.Project,
			Service:            e.Service,
			Stage:              e.Stage,
			TestStrategy:       e.TestStrategy,
			DeploymentStrategy: e.DeploymentStrategy,
			Labels:             e.Labels,
		}
		err = eh.sendEvaluationDoneEvent(keptnContext, &evaluationResult)
		return err
	}

	// get SLO file
	objectives, err := getSLOs(e.Project, e.Stage, e.Service)
	if err != nil {
		// No SLO file found -> no need to evaluate
		eh.Logger.Debug("No SLO file found, no evaluation conducted")
		evaluationDetails := keptnevents.EvaluationDetails{
			IndicatorResults: nil,
			TimeStart:        e.Start,
			TimeEnd:          e.End,
			Result:           fmt.Sprintf("no evaluation performed by lighthouse service (no slo.yaml found)"),
		}
		// send the evaluation-done-event
		evaluationResult := keptnevents.EvaluationDoneEventData{
			EvaluationDetails:  &evaluationDetails,
			Result:             "pass",
			Project:            e.Project,
			Service:            e.Service,
			Stage:              e.Stage,
			TestStrategy:       e.TestStrategy,
			DeploymentStrategy: e.DeploymentStrategy,
			Labels:             e.Labels,
		}

		err = eh.sendEvaluationDoneEvent(keptnContext, &evaluationResult)
		return err
	}

	indicators := []string{}
	for _, objective := range objectives.Objectives {
		indicators = append(indicators, objective.SLI)
	}

	var filters = []*keptnevents.SLIFilter{}

	deployment := ""
	if e.DeploymentStrategy != "" {
		if e.DeploymentStrategy == "blue_green_service" {
			if e.TestStrategy == "real-user" { // e.g., remediation use case -> wait-service
				deployment = "primary"
			} else {
				deployment = "canary"
			}
		} else {
			// assert deployment_strategy == 'direct'
			deployment = "direct"
		}
	}

	if objectives.Filter != nil {
		for key, value := range objectives.Filter {
			filter := &keptnevents.SLIFilter{
				Key:   key,
				Value: value,
			}
			filters = append(filters, filter)
		}
	}

	// get the SLI provider that has been configured for the project (e.g. 'dynatrace' or 'prometheus')
	sliProvider, err := getSLIProvider(e.Project)
	if err != nil {
		// ToDo: We need to provide feedback to the user that this failed becuase no sli provider was set for project
		eh.Logger.Error("Could not determine SLI provider for project " + e.Project)
		return err
	}
	// send a new event to trigger the SLI retrieval
	eh.Logger.Debug("SLI provider for project " + e.Project + " is: " + sliProvider)
	err = eh.sendInternalGetSLIEvent(keptnContext, e.Project, e.Stage, e.Service, sliProvider, indicators, e.Start, e.End, e.TestStrategy, e.DeploymentStrategy, filters, e.Labels, deployment)
	return nil
}

func (eh *StartEvaluationHandler) sendEvaluationDoneEvent(shkeptncontext string, data *keptnevents.EvaluationDoneEventData) error {

	source, _ := url.Parse("lighthouse-service")
	contentType := "application/json"

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptnevents.EvaluationDoneEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: data,
	}

	eh.Logger.Debug("Send event: " + keptnevents.EvaluationDoneEventType)
	return sendEvent(event)
}

func getSLIProvider(project string) (string, error) {
	kubeClient, err := keptnutils.GetKubeAPI(true)

	if err != nil {
		return "", err
	}

	configMap, err := kubeClient.ConfigMaps("keptn").Get("lighthouse-config-"+project, v1.GetOptions{})

	if err != nil {
		return "", errors.New("No SLI provider specified for project " + project)
	}

	sliProvider := configMap.Data["sli-provider"]

	return sliProvider, nil
}

func (eh *StartEvaluationHandler) sendInternalGetSLIEvent(shkeptncontext string, project string, stage string, service string, sliProvider string, indicators []string, start string, end string, teststrategy string, deploymentStrategy string, filters []*keptnevents.SLIFilter, labels map[string]string, deployment string) error {
	source, _ := url.Parse("lighthouse-service")
	contentType := "application/json"

	getSLIEvent := keptnevents.InternalGetSLIEventData{
		SLIProvider:        sliProvider,
		Project:            project,
		Service:            service,
		Stage:              stage,
		Start:              start,
		End:                end,
		Indicators:         indicators,
		CustomFilters:      filters,
		TestStrategy:       teststrategy,
		DeploymentStrategy: deploymentStrategy,
		Labels:             labels,
		Deployment:         deployment,
	}
	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptnevents.InternalGetSLIEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: getSLIEvent,
	}

	eh.Logger.Debug("Send event: " + keptnevents.InternalGetSLIEventType)
	return sendEvent(event)
}
