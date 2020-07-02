package event_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnutils "github.com/keptn/go-utils/pkg/lib"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StartEvaluationHandler struct {
	Event        cloudevents.Event
	KeptnHandler *keptnutils.Keptn
}

func (eh *StartEvaluationHandler) HandleEvent() error {

	var keptnContext string
	_ = eh.Event.ExtensionAs("shkeptncontext", &keptnContext)

	e := &keptnevents.StartEvaluationEventData{}
	err := eh.Event.DataAs(e)
	if err != nil {
		eh.KeptnHandler.Logger.Error("Could not parse event payload: " + err.Error())
		return err
	}

	if e.TestStrategy == "" {
		eh.KeptnHandler.Logger.Debug("No test has been executed, no evaluation conducted")
		evaluationDetails := keptnevents.EvaluationDetails{
			IndicatorResults: nil,
			TimeStart:        e.Start,
			TimeEnd:          e.End,
			Result:           fmt.Sprintf("no evaluation performed by lighthouse because no test has been executed"),
		}
		// send the evaluation-done-event
		evaluationResult := keptnevents.EvaluationDoneEventData{
			EvaluationDetails:  &evaluationDetails,
			Result:             eh.getTestExecutionResult(),
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
		// no SLO file found (assumption that this is an empty SLO file) -> no need to evaluate
		eh.KeptnHandler.Logger.Debug("No SLO file found, no evaluation conducted")
		evaluationDetails := keptnevents.EvaluationDetails{
			IndicatorResults: nil,
			TimeStart:        e.Start,
			TimeEnd:          e.End,
			Result:           fmt.Sprintf("no evaluation performed by lighthouse because no SLO found for service %s", e.Service),
		}

		evaluationResult := keptnevents.EvaluationDoneEventData{
			EvaluationDetails:  &evaluationDetails,
			Result:             eh.getTestExecutionResult(),
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
			// blue-green deployed services should be evaluated based on data of either the primary or canary deployment
			if e.TestStrategy == "real-user" {
				// remediation use case will be tested by real users, therefore the evaluation needs to take place on the on the primary deployment
				deployment = "primary"
			} else {
				// while load-tests are running on the canary deployment
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
		eh.KeptnHandler.Logger.Error("no SLI-provider configured for project " + e.Project + ", no evaluation conducted")
		evaluationDetails := keptnevents.EvaluationDetails{
			IndicatorResults: nil,
			TimeStart:        e.Start,
			TimeEnd:          e.End,
			Result:           fmt.Sprintf("no evaluation performed by lighthouse because no SLI-provider configured for project %s", e.Project),
		}

		evaluationResult := keptnevents.EvaluationDoneEventData{
			EvaluationDetails:  &evaluationDetails,
			Result:             "failed",
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
	// send a new event to trigger the SLI retrieval
	eh.KeptnHandler.Logger.Debug("SLI provider for project " + e.Project + " is: " + sliProvider)
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

	eh.KeptnHandler.Logger.Debug("Send event: " + keptnevents.EvaluationDoneEventType)
	return eh.KeptnHandler.SendCloudEvent(event)
}

func getSLIProvider(project string) (string, error) {
	kubeAPI, err := getKubeAPI()
	if err != nil {
		return "", err
	}

	configMap, err := kubeAPI.CoreV1().ConfigMaps("keptn").Get("lighthouse-config-"+project, v1.GetOptions{})

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

	eh.KeptnHandler.Logger.Debug("Send event: " + keptnevents.InternalGetSLIEventType)
	return eh.KeptnHandler.SendCloudEvent(event)
}

func (eh *StartEvaluationHandler) getTestExecutionResult() string {
	dataByte, err := eh.Event.DataBytes()
	if err != nil {
		eh.KeptnHandler.Logger.Error("Could not get event as byte array: " + err.Error())
	}

	e := &keptnevents.TestsFinishedEventData{}
	err = json.Unmarshal(dataByte, e)
	if err != nil {
		eh.KeptnHandler.Logger.Error("Could not unmarshal event payload: " + err.Error())
	}
	return e.Result
}
