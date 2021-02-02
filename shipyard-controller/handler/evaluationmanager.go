package handler

import (
	"fmt"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "2006-01-02T15:04:05"

const (
	evaluationErrInvalidTimeframe = iota
	evaluationErrSendEventFailed
	evaluationErrServiceNotAvailable
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/evaluationmanager.go . IEvaluationManager
type IEvaluationManager interface {
	CreateEvaluation(project, stage, service string, params *operations.CreateEvaluationParams) (*operations.CreateEvaluationResponse, *models.Error)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/serviceapi.go . IServiceAPI
type IServiceAPI interface {
	GetService(project, stage, service string) (*keptnapimodels.Service, error)
}

type EvaluationManager struct {
	EventSender keptn.EventSender
	ServiceAPI  IServiceAPI
	Logger      keptn.LoggerInterface
}

func NewEvaluationManager(eventSender keptn.EventSender, serviceAPI IServiceAPI, logger keptn.LoggerInterface) (*EvaluationManager, error) {
	return &EvaluationManager{
		EventSender: eventSender,
		ServiceAPI:  serviceAPI,
		Logger:      logger,
	}, nil

}

func (em *EvaluationManager) CreateEvaluation(project, stage, service string, params *operations.CreateEvaluationParams) (*operations.CreateEvaluationResponse, *models.Error) {
	// TODO: check service availability via materialized view after https://github.com/keptn/keptn/issues/2999 has been merged
	//_, err := em.ServiceAPI.GetService(project, stage, service)
	//if err != nil {
	//	return nil, &models.Error{
	//		Code:    evaluationErrServiceNotAvailable,
	//		Message: swag.String(err.Error()),
	//	}
	//}

	keptnContext := uuid.New().String()
	extensions := make(map[string]interface{})
	extensions["shkeptncontext"] = keptnContext

	start, end, err := getStartEndTime(params.Start, params.End, params.Timeframe)
	if err != nil {
		return nil, &models.Error{
			Code:    evaluationErrInvalidTimeframe,
			Message: swag.String(err.Error()),
		}
	}

	eventContext := &operations.CreateEvaluationResponse{KeptnContext: keptnContext}

	evaluationTriggeredEvent := keptnv2.EvaluationTriggeredEventData{
		EventData: keptnv2.EventData{
			Project: project,
			Service: service,
			Stage:   stage,
			Labels:  params.Labels,
		},
		Evaluation: struct {
			Start string `json:"start"`
			End   string `json:"end"`
		}{
			Start: start.String(),
			End:   end.String(),
		},
	}

	ce := common.CreateEventWithPayload(keptnContext, "", keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName), evaluationTriggeredEvent)
	if err := ce.Context.SetSource("https://github.com/keptn/keptn/api"); err != nil {
		return nil, &models.Error{
			Code:    evaluationErrSendEventFailed,
			Message: stringp(err.Error()),
		}
	}
	if err := em.EventSender.SendEvent(ce); err != nil {
		return nil, &models.Error{
			Code:    evaluationErrSendEventFailed,
			Message: stringp(err.Error()),
		}
	}

	return eventContext, nil
}

func getStartEndTime(startDatePoint, endDatePoint, timeframe string) (*time.Time, *time.Time, error) {
	// set default values for start and end time
	var err error

	minutes := 5 // default timeframe

	// input validation
	if startDatePoint != "" && endDatePoint == "" {
		// if a start date is set, but no end date is set, we require the timeframe to be set
		if timeframe == "" {
			errMsg := "Please provide a timeframe, e.g., --timeframe=5m, or an end date using --end=..."

			return nil, nil, fmt.Errorf(errMsg)
		}
	}
	if endDatePoint != "" && timeframe != "" {
		// can not use end date and timeframe at the same time
		errMsg := "You can not use --end together with --timeframe"

		return nil, nil, fmt.Errorf(errMsg)
	}
	if endDatePoint != "" && startDatePoint == "" {
		errMsg := "start date is required when using an end date"

		return nil, nil, fmt.Errorf(errMsg)
	}

	// parse timeframe
	if timeframe != "" {
		errMsg := "The time frame format is invalid. Use the format [duration]m, e.g.: 5m"

		i := strings.Index(timeframe, "m")

		if i > -1 {
			minutesStr := timeframe[:i]
			minutes, err = strconv.Atoi(minutesStr)
			if err != nil {
				return nil, nil, fmt.Errorf(errMsg)
			}
		} else {
			return nil, nil, fmt.Errorf(errMsg)
		}
	}

	// initialize default values for end and start time
	end := time.Now().UTC()
	start := time.Now().UTC().Add(-time.Duration(minutes) * time.Minute)

	// Parse start date
	if startDatePoint != "" {
		start, err = time.Parse(dateLayout, startDatePoint)

		if err != nil {
			return nil, nil, err
		}
	}

	// Parse end date
	if endDatePoint != "" {
		end, err = time.Parse(dateLayout, endDatePoint)

		if err != nil {
			return nil, nil, err
		}
	}

	// last but not least: if a start date and a timeframe is provided, we set the end date to start date + timeframe
	if startDatePoint != "" && endDatePoint == "" && timeframe != "" {
		minutesOffset := time.Minute * time.Duration(minutes)
		end = start.Add(minutesOffset)
	}

	// ensure end date is greater than start date
	diff := end.Sub(start).Minutes()

	if diff < 1 {
		errMsg := "end date must be at least 1 minute after start date"

		return nil, nil, fmt.Errorf(errMsg)
	}

	return &start, &end, nil
}
