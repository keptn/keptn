package handlers

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"

	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptnutils "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/evaluation"
	"github.com/keptn/keptn/api/restapi/operations/event"
	"github.com/keptn/keptn/api/utils"
)

// TriggerEvaluationHandlerFunc triggers a new evaluation by sending a start-evaluation event
func TriggerEvaluationHandlerFunc(params evaluation.TriggerEvaluationParams, principal *models.Principal) middleware.Responder {
	serviceHandler := keptnapi.NewServiceHandler("http://configuration-service:8080")
	_, err := serviceHandler.GetService(params.ProjectName, params.StageName, params.ServiceName)
	if err != nil {
		return evaluation.NewTriggerEvaluationBadRequest().WithPayload(&models.Error{
			Code:    400,
			Message: swag.String(err.Error()),
		})
	}

	keptnContext := uuid.New().String()
	extensions := make(map[string]interface{})
	extensions["shkeptncontext"] = keptnContext

	logger := keptnutils.NewLogger(keptnContext, "", "api")
	logger.Info("API received a trigger-evaluation request")

	start, end, err := getStartEndTime(params.Evaluation.Start, params.Evaluation.End, params.Evaluation.Timeframe)
	if err != nil {
		return evaluation.NewTriggerEvaluationBadRequest().WithPayload(&models.Error{
			Code:    400,
			Message: swag.String(err.Error()),
		})
	}

	source, _ := url.Parse("https://github.com/keptn/keptn/api")

	eventContext := models.EventContext{KeptnContext: &keptnContext}

	startEvaluationEvent := keptnv2.EvaluationTriggeredEventData{
		EventData: keptnv2.EventData{
			Project: params.ProjectName,
			Service: params.ServiceName,
			Stage:   params.StageName,
			Labels:  params.Evaluation.Labels,
		},
		Evaluation: struct {
			Start string `json:"start"`
			End   string `json:"end"`
		}{
			Start: start.Format("2006-01-02T15:04:05.000Z"),
			End:   end.Format("2006-01-02T15:04:05.000Z"),
		},
	}

	err = utils.SendEvent(
		keptnContext,
		"",
		keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName),
		source.String(),
		startEvaluationEvent,
		logger,
	)
	if err != nil {
		return sendInternalErrorForPost(err, logger)
	}

	return event.NewPostEventOK().WithPayload(&eventContext)
}

func getStartEndTime(startDatePoint string, endDatePoint string, timeframe string) (*time.Time, *time.Time, error) {
	// set default values for start and end time
	dateLayout := "2006-01-02T15:04:05"
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
