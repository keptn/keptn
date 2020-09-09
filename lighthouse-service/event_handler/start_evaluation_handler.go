package event_handler

import (
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/url"
)

type StartEvaluationHandler struct {
	Event        cloudevents.Event
	KeptnHandler *keptnv2.Keptn
}

func (eh *StartEvaluationHandler) HandleEvent() error {

	var keptnContext string
	_ = eh.Event.ExtensionAs("shkeptncontext", &keptnContext)

	e := &keptnv2.EvaluationTriggeredEventData{}
	err := eh.Event.DataAs(e)
	if err != nil {
		eh.KeptnHandler.Logger.Error("Could not parse event payload: " + err.Error())
		return err
	}

	startedEvent := keptnv2.EvaluationStartedEventData{
		EventData: e.EventData,
	}
	startedEvent.EventData.Status = keptnv2.StatusSucceeded

	err = sendEvent(keptnContext, eh.Event.ID(), keptnv2.GetStartedEventType(keptnv2.EvaluationTaskName), eh.KeptnHandler, startedEvent)
	if err != nil {
		eh.KeptnHandler.Logger.Error("Could not send evaluation.started event: " + err.Error())
		return err
	}

	// get SLO file
	objectives, err := getSLOs(e.Project, e.Stage, e.Service)
	if err != nil {
		// no SLO file found (assumption that this is an empty SLO file) -> no need to evaluate
		eh.KeptnHandler.Logger.Debug("No SLO file found, no evaluation conducted")
		evaluationDetails := keptnv2.EvaluationDetails{
			IndicatorResults: nil,
			TimeStart:        e.Test.Start,
			TimeEnd:          e.Test.End,
			Result:           fmt.Sprintf("no evaluation performed by lighthouse because no SLO found for service %s", e.Service),
		}

		evaluationFinishedData := keptnv2.EvaluationFinishedEventData{
			EventData: keptnv2.EventData{
				Project: e.Project,
				Stage:   e.Stage,
				Service: e.Service,
				Labels:  e.Labels,
				Status:  keptnv2.StatusSucceeded,
				Result:  keptnv2.ResultPass,
				Message: fmt.Sprintf("no evaluation performed by lighthouse because no SLO found for service %s", e.Service),
			},
			Evaluation: evaluationDetails,
		}

		err = sendEvent(keptnContext, eh.Event.ID(), keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), eh.KeptnHandler, &evaluationFinishedData)
		return err
	}

	indicators := []string{}
	for _, objective := range objectives.Objectives {
		indicators = append(indicators, objective.SLI)
	}

	var filters = []*keptnevents.SLIFilter{}

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
		evaluationDetails := keptnv2.EvaluationDetails{
			IndicatorResults: nil,
			TimeStart:        e.Test.Start,
			TimeEnd:          e.Test.End,
			Result:           fmt.Sprintf("no evaluation performed by lighthouse because no SLI-provider configured for project %s", e.Project),
		}

		evaluationFinishedData := keptnv2.EvaluationFinishedEventData{
			EventData: keptnv2.EventData{
				Project: e.Project,
				Stage:   e.Stage,
				Service: e.Service,
				Labels:  e.Labels,
				Status:  keptnv2.StatusSucceeded,
				Result:  keptnv2.ResultPass,
				Message: fmt.Sprintf("no evaluation performed by lighthouse because no SLI-provider configured for project %s", e.Project),
			},
			Evaluation: evaluationDetails,
		}

		return sendEvent(keptnContext, eh.Event.ID(), keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), eh.KeptnHandler, &evaluationFinishedData)
	}
	// send a new event to trigger the SLI retrieval
	eh.KeptnHandler.Logger.Debug("SLI provider for project " + e.Project + " is: " + sliProvider)
	err = eh.sendInternalGetSLIEvent(keptnContext, e.Project, e.Stage, e.Service, sliProvider, indicators, e.Test.Start, e.Test.End, filters, e.Labels)
	return nil
}

func getSLIProvider(project string) (string, error) {
	kubeAPI, err := getKubeAPI()
	if err != nil {
		return "", err
	}

	configMap, err := kubeAPI.CoreV1().ConfigMaps(namespace).Get("lighthouse-config-"+project, v1.GetOptions{})

	if err != nil {
		return "", errors.New("No SLI provider specified for project " + project)
	}

	sliProvider := configMap.Data["sli-provider"]

	return sliProvider, nil
}

func (eh *StartEvaluationHandler) sendInternalGetSLIEvent(shkeptncontext string, project string, stage string, service string, sliProvider string, indicators []string, start string, end string, filters []*keptnevents.SLIFilter, labels map[string]string) error {
	source, _ := url.Parse("lighthouse-service")

	getSLIEvent := keptnevents.InternalGetSLIEventData{
		SLIProvider:   sliProvider,
		Project:       project,
		Service:       service,
		Stage:         stage,
		Start:         start,
		End:           end,
		Indicators:    indicators,
		CustomFilters: filters,
		Labels:        labels,
	}

	event := cloudevents.NewEvent()
	event.SetType(keptnevents.InternalGetSLIEventType)
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", shkeptncontext)
	event.SetData(cloudevents.ApplicationJSON, getSLIEvent)

	eh.KeptnHandler.Logger.Debug("Send event: " + keptnevents.InternalGetSLIEventType)
	return eh.KeptnHandler.SendCloudEvent(event)
}

func (eh *StartEvaluationHandler) getTestExecutionResult() string {

	e := &keptnevents.TestsFinishedEventData{}
	err := eh.Event.DataAs(e)
	if err != nil {
		eh.KeptnHandler.Logger.Error("Could not unmarshal event payload: " + err.Error())
	}
	return e.Result
}
