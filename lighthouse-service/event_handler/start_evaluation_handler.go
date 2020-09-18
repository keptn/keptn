package event_handler

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnutils "github.com/keptn/go-utils/pkg/lib"
)

type StartEvaluationHandler struct {
	Event             cloudevents.Event
	KeptnHandler      *keptnutils.Keptn
	SLIProviderConfig SLIProviderConfig
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

	indicators := []string{}
	var filters = []*keptnevents.SLIFilter{}
	// get SLO file
	objectives, err := getSLOs(e.Project, e.Stage, e.Service)
	if err == nil && objectives != nil {
		eh.KeptnHandler.Logger.Info("SLO file found")
		for _, objective := range objectives.Objectives {
			indicators = append(indicators, objective.SLI)
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
	} else if err != nil && err != ErrSLOFileNotFound {
		var message string
		if err == ErrServiceNotFound {
			message = "error retrieving SLO file: service " + e.Service + " not found"
			eh.KeptnHandler.Logger.Error(message)
		} else if err == ErrStageNotFound {
			message = "error retrieving SLO file: stage " + e.Stage + " not found"
			eh.KeptnHandler.Logger.Error(message)
		} else if err == ErrProjectNotFound {
			message = "error retrieving SLO file: project " + e.Project + " not found"
			eh.KeptnHandler.Logger.Error(message)
		}
		evaluationDetails := keptnevents.EvaluationDetails{
			IndicatorResults: nil,
			TimeStart:        e.Start,
			TimeEnd:          e.End,
			Result:           message,
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
	} else if err != nil && err == ErrSLOFileNotFound {
		eh.KeptnHandler.Logger.Info("no SLO file found")
	}

	// get the SLI provider that has been configured for the project (e.g. 'dynatrace' or 'prometheus')
	var sliProvider string
	sliProvider, err = eh.SLIProviderConfig.GetSLIProvider(e.Project)
	if err != nil {
		sliProvider, err = eh.SLIProviderConfig.GetDefaultSLIProvider()
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
