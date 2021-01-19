package event_handler

import (
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"net/url"
)

type StartEvaluationHandler struct {
	Event             cloudevents.Event
	KeptnHandler      *keptnv2.Keptn
	SLIProviderConfig SLIProviderConfig
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

	evaluationStartTimestamp, evaluationEndTimestamp, err := getEvaluationTimestamps(e)
	if err != nil {
		return eh.sendEvaluationFinishedWithErrorEvent(keptnContext, evaluationStartTimestamp, evaluationEndTimestamp, e, err.Error())
	}

	// get SLO file
	indicators := []string{}
	var filters = []*keptnv2.SLIFilter{}
	// get SLO file
	objectives, err := getSLOs(e.Project, e.Stage, e.Service)
	if err == nil && objectives != nil {
		eh.KeptnHandler.Logger.Info("SLO file found")
		for _, objective := range objectives.Objectives {
			indicators = append(indicators, objective.SLI)
		}

		if objectives.Filter != nil {
			for key, value := range objectives.Filter {
				filter := &keptnv2.SLIFilter{
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
		return eh.sendEvaluationFinishedWithErrorEvent(keptnContext, evaluationStartTimestamp, evaluationEndTimestamp, e, message)
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
			evaluationDetails := keptnv2.EvaluationDetails{
				IndicatorResults: nil,
				TimeStart:        evaluationStartTimestamp,
				TimeEnd:          evaluationEndTimestamp,
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
	}
	// send a new event to trigger the SLI retrieval
	eh.KeptnHandler.Logger.Debug("SLI provider for project " + e.Project + " is: " + sliProvider)
	err = eh.sendInternalGetSLIEvent(keptnContext, e.Project, e.Stage, e.Service, sliProvider, indicators, evaluationStartTimestamp, evaluationEndTimestamp, filters, e.Labels)
	return nil
}

func (eh *StartEvaluationHandler) sendEvaluationFinishedWithErrorEvent(keptnContext, start, end string, e *keptnv2.EvaluationTriggeredEventData, message string) error {
	evaluationDetails := keptnv2.EvaluationDetails{
		IndicatorResults: nil,
		TimeStart:        start,
		TimeEnd:          end,
		Result:           message,
	}

	evaluationFinishedData := keptnv2.EvaluationFinishedEventData{
		EventData: keptnv2.EventData{
			Project: e.Project,
			Stage:   e.Stage,
			Service: e.Service,
			Labels:  e.Labels,
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
			Message: message,
		},
		Evaluation: evaluationDetails,
	}

	return sendEvent(keptnContext, eh.Event.ID(), keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), eh.KeptnHandler, &evaluationFinishedData)
}

func getEvaluationTimestamps(e *keptnv2.EvaluationTriggeredEventData) (string, string, error) {
	if e.Evaluation.Start != "" && e.Evaluation.End != "" {
		return e.Evaluation.Start, e.Evaluation.End, nil
	} else if e.Test.Start != "" && e.Test.End != "" {
		return e.Test.Start, e.Test.End, nil
	}
	return "", "", errors.New("evaluation.triggered event does not contain evaluation timeframe")
}

func (eh *StartEvaluationHandler) sendInternalGetSLIEvent(shkeptncontext string, project string, stage string, service string, sliProvider string, indicators []string, start string, end string, filters []*keptnv2.SLIFilter, labels map[string]string) error {
	source, _ := url.Parse("lighthouse-service")

	getSLIEvent := keptnv2.GetSLITriggeredEventData{
		EventData: keptnv2.EventData{
			Project: project,
			Stage:   stage,
			Service: service,
			Labels:  labels,
		},
		GetSLI: keptnv2.GetSLI{
			SLIProvider:   sliProvider,
			Start:         start,
			End:           end,
			Indicators:    indicators,
			CustomFilters: filters,
		},
	}

	event := cloudevents.NewEvent()
	event.SetType(keptnv2.GetTriggeredEventType(keptnv2.GetSLITaskName))
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", shkeptncontext)
	event.SetData(cloudevents.ApplicationJSON, getSLIEvent)

	eh.KeptnHandler.Logger.Debug("Send event: " + keptnv2.GetTriggeredEventType(keptnv2.GetSLITaskName))
	return eh.KeptnHandler.SendCloudEvent(event)
}
