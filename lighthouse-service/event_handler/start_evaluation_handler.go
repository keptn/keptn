package event_handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	logger "github.com/sirupsen/logrus"
	"net/url"
	"sync"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

type StartEvaluationHandler struct {
	Event             cloudevents.Event
	KeptnHandler      *keptnv2.Keptn
	SLIProviderConfig SLIProviderConfig
	SLOFileRetriever  SLOFileRetriever `deep:"-"`
}

func (eh *StartEvaluationHandler) HandleEvent(ctx context.Context) error {
	var keptnContext string
	_ = eh.Event.ExtensionAs("shkeptncontext", &keptnContext)

	e := &keptnv2.EvaluationTriggeredEventData{}
	err := eh.Event.DataAs(e)
	if err != nil {
		logger.Error("Could not parse event payload: " + err.Error())
		return err
	}

	startedEvent := keptnv2.EvaluationStartedEventData{
		EventData: e.EventData,
	}
	startedEvent.EventData.Status = keptnv2.StatusSucceeded

	// send evaluation.started event
	err = sendEvent(keptnContext, eh.Event.ID(), keptnv2.GetStartedEventType(keptnv2.EvaluationTaskName), eh.KeptnHandler, startedEvent)
	if err != nil {
		logger.Error("Could not send evaluation.started event: " + err.Error())
		return err
	}

	// try to parse timestamps
	evaluationStartTimestamp, evaluationEndTimestamp, err := getEvaluationTimestamps(e)
	if err != nil {
		return eh.sendEvaluationFinishedWithErrorEvent(evaluationStartTimestamp, evaluationEndTimestamp, e, err.Error())
	}
	ctx.Value(GracefulShutdownKey).(*sync.WaitGroup).Add(1)
	go eh.sendGetSliCloudEvent(ctx, keptnContext, e, evaluationStartTimestamp, evaluationEndTimestamp)

	return nil
}

// fetch SLO and send the internal get-sli event
func (eh *StartEvaluationHandler) sendGetSliCloudEvent(ctx context.Context, keptnContext string, e *keptnv2.EvaluationTriggeredEventData, evaluationStartTimestamp string, evaluationEndTimestamp string) error {
	defer func() {
		ctx.Value(GracefulShutdownKey).(*sync.WaitGroup).Done()
		eh.KeptnHandler.Logger.Info("Terminating Start-evaluation handler")
	}()

	indicators := []string{}
	var filters = []*keptnv2.SLIFilter{}

	// collect objectives from SLO file
	objectives, err := eh.SLOFileRetriever.GetSLOs(e.Project, e.Stage, e.Service)
	if err == nil && objectives != nil {
		logger.Info("SLO file found")
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
		} else if err == ErrStageNotFound {
			message = "error retrieving SLO file: stage " + e.Stage + " not found"
		} else if err == ErrProjectNotFound {
			message = "error retrieving SLO file: project " + e.Project + " not found"
		} else {
			message = fmt.Sprintf("error retrieving SLO file: %s", err.Error())
		}
		logger.Error(message)
		return eh.sendEvaluationFinishedWithErrorEvent(evaluationStartTimestamp, evaluationEndTimestamp, e, message)
	} else if err != nil && err == ErrSLOFileNotFound {
		logger.Error("no SLO file found")
	}

	// get the SLI provider that has been configured for the project (e.g. 'dynatrace' or 'prometheus') from the respective configmap
	var sliProvider string
	sliProvider, err = eh.SLIProviderConfig.GetSLIProvider(e.Project)
	if err != nil {
		// no provider found - fallback to default SLI provider
		sliProvider, err = eh.SLIProviderConfig.GetDefaultSLIProvider()
		if err != nil {
			// no default SLI provider configured
			logger.Error("no SLI-provider configured for project " + e.Project + ", no evaluation conducted")
			evaluationDetails := keptnv2.EvaluationDetails{
				IndicatorResults: nil,
				TimeStart:        evaluationStartTimestamp,
				TimeEnd:          evaluationEndTimestamp,
				Result:           string(keptnv2.ResultPass),
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
	logger.Debug("SLI provider for project " + e.Project + " is: " + sliProvider)
	err = eh.sendInternalGetSLIEvent(keptnContext, e, sliProvider, indicators, evaluationStartTimestamp, evaluationEndTimestamp, filters)
	return nil
}

func (eh *StartEvaluationHandler) sendEvaluationFinishedWithErrorEvent(start, end string, e *keptnv2.EvaluationTriggeredEventData, message string) error {
	evaluationDetails := keptnv2.EvaluationDetails{
		IndicatorResults: nil,
		TimeStart:        start,
		TimeEnd:          end,
		Result:           string(keptnv2.ResultFailed),
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

	_, err := eh.KeptnHandler.SendTaskFinishedEvent(&evaluationFinishedData, "lighthouse-service")
	return err
}

func getEvaluationTimestamps(e *keptnv2.EvaluationTriggeredEventData) (string, string, error) {
	if (e.Evaluation.Start != "" && e.Evaluation.End != "") || (e.Evaluation.Timeframe != "") {
		params := timeutils.GetStartEndTimeParams{
			StartDate: e.Evaluation.Start,
			EndDate:   e.Evaluation.End,
			Timeframe: e.Evaluation.Timeframe,
		}
		start, end, err := timeutils.GetStartEndTime(params)
		if err != nil {
			return "", "", err
		}
		return timeutils.GetKeptnTimeStamp(*start), timeutils.GetKeptnTimeStamp(*end), nil
	} else if e.Test.Start != "" && e.Test.End != "" {
		return e.Test.Start, e.Test.End, nil
	}
	return "", "", errors.New("evaluation.triggered event does not contain evaluation timeframe")
}

func (eh *StartEvaluationHandler) sendInternalGetSLIEvent(shkeptncontext string, e *keptnv2.EvaluationTriggeredEventData, sliProvider string, indicators []string, start string, end string, filters []*keptnv2.SLIFilter) error {
	source, _ := url.Parse("lighthouse-service")

	getSLITriggeredEventData := keptnv2.GetSLITriggeredEventData{
		EventData: keptnv2.EventData{
			Project: e.Project,
			Stage:   e.Stage,
			Service: e.Service,
			Labels:  e.Labels,
		},
		GetSLI: keptnv2.GetSLI{
			SLIProvider:   sliProvider,
			Start:         start,
			End:           end,
			Indicators:    indicators,
			CustomFilters: filters,
		},
	}

	if e.Deployment.DeploymentNames != nil && len(e.Deployment.DeploymentNames) > 0 {
		getSLITriggeredEventData.Deployment = e.Deployment.DeploymentNames[0]
	}

	event := cloudevents.NewEvent()
	event.SetType(keptnv2.GetTriggeredEventType(keptnv2.GetSLITaskName))
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", shkeptncontext)
	event.SetData(cloudevents.ApplicationJSON, getSLITriggeredEventData)

	logger.Debug("Send event: " + keptnv2.GetTriggeredEventType(keptnv2.GetSLITaskName))
	return eh.KeptnHandler.SendCloudEvent(event)
}
