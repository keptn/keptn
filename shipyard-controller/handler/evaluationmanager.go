package handler

import (
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"time"
)

const userFriendlyTimeFormat = "2006-01-02T15:04:05"

const (
	evaluationErrInvalidTimeframe = iota
	evaluationErrSendEventFailed
	evaluationErrServiceNotAvailable
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/evaluationmanager.go . IEvaluationManager
type IEvaluationManager interface {
	CreateEvaluation(project, stage, service string, params *operations.CreateEvaluationParams) (*operations.CreateEvaluationResponse, *models.Error)
}

type EvaluationManager struct {
	EventSender keptn.EventSender
	ServiceAPI  db.ServicesDbOperations
}

func NewEvaluationManager(eventSender keptn.EventSender, serviceAPI db.ServicesDbOperations) (*EvaluationManager, error) {
	return &EvaluationManager{
		EventSender: eventSender,
		ServiceAPI:  serviceAPI,
	}, nil

}

func (em *EvaluationManager) CreateEvaluation(project, stage, service string, params *operations.CreateEvaluationParams) (*operations.CreateEvaluationResponse, *models.Error) {
	_, err := em.ServiceAPI.GetService(project, stage, service)
	if err != nil {
		return nil, &models.Error{
			Code:    evaluationErrServiceNotAvailable,
			Message: strutils.Stringp(err.Error()),
		}
	}

	keptnContext := uuid.New().String()
	extensions := make(map[string]interface{})
	extensions["shkeptncontext"] = keptnContext

	var start, end *time.Time
	start, end, err = timeutils.GetStartEndTime(timeutils.GetStartEndTimeParams{
		StartDate: params.Start,
		EndDate:   params.End,
		Timeframe: params.Timeframe,
	})
	if err != nil {
		// if we got an error, try again with other time format
		start, end, err = timeutils.GetStartEndTime(timeutils.GetStartEndTimeParams{
			StartDate:  params.Start,
			EndDate:    params.End,
			Timeframe:  params.Timeframe,
			TimeFormat: userFriendlyTimeFormat,
		})
		if err != nil {
			return nil, &models.Error{
				Code:    evaluationErrInvalidTimeframe,
				Message: strutils.Stringp(err.Error()),
			}
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
		Evaluation: keptnv2.Evaluation{
			Start: timeutils.GetKeptnTimeStamp(*start),
			End:   timeutils.GetKeptnTimeStamp(*end),
		},
	}

	ce := common.CreateEventWithPayload(keptnContext, "", keptnv2.GetTriggeredEventType(stage+"."+keptnv2.EvaluationTaskName), evaluationTriggeredEvent)
	if err := ce.Context.SetSource("https://github.com/keptn/keptn/api"); err != nil {
		return nil, &models.Error{
			Code:    evaluationErrSendEventFailed,
			Message: common.Stringp(err.Error()),
		}
	}
	if err := em.EventSender.SendEvent(ce); err != nil {
		return nil, &models.Error{
			Code:    evaluationErrSendEventFailed,
			Message: common.Stringp(err.Error()),
		}
	}

	return eventContext, nil
}
