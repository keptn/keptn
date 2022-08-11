package handler

import (
	"github.com/keptn/keptn/shipyard-controller/internal/common"
	"github.com/keptn/keptn/shipyard-controller/internal/db"
	"time"

	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
)

const userFriendlyTimeFormat = "2006-01-02T15:04:05"

const (
	evaluationErrInvalidTimeframe = iota
	evaluationErrSendEventFailed
	evaluationErrServiceNotAvailable
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/evaluationmanager.go . IEvaluationManager
type IEvaluationManager interface {
	CreateEvaluation(project, stage, service string, params *models.CreateEvaluationParams) (*models.CreateEvaluationResponse, *models.Error)
}

type EvaluationManager struct {
	eventSender   keptn.EventSender
	projectMVRepo db.ProjectMVRepo
}

func NewEvaluationManager(eventSender keptn.EventSender, projectMVRepo db.ProjectMVRepo) (*EvaluationManager, error) {
	return &EvaluationManager{
		eventSender:   eventSender,
		projectMVRepo: projectMVRepo,
	}, nil
}

func (em *EvaluationManager) CreateEvaluation(project, stage, service string, params *models.CreateEvaluationParams) (*models.CreateEvaluationResponse, *models.Error) {
	_, err := em.projectMVRepo.GetService(project, stage, service)
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

	eventContext := &models.CreateEvaluationResponse{KeptnContext: keptnContext}

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
	if params.GitCommitID != "" {
		ce.SetExtension("gitcommitid", params.GitCommitID)
	}
	if err := ce.Context.SetSource("https://github.com/keptn/keptn/api"); err != nil {
		return nil, &models.Error{
			Code:    evaluationErrSendEventFailed,
			Message: common.Stringp(err.Error()),
		}
	}
	if err := em.eventSender.SendEvent(ce); err != nil {
		return nil, &models.Error{
			Code:    evaluationErrSendEventFailed,
			Message: common.Stringp(err.Error()),
		}
	}

	return eventContext, nil
}
