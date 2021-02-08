package api

import (
	"github.com/gin-gonic/gin"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/statistics-service/config"
	"github.com/keptn/keptn/statistics-service/controller"
	"github.com/keptn/keptn/statistics-service/db"
	"github.com/keptn/keptn/statistics-service/operations"
	"net/http"
)

// GetStatistics godoc
// @Summary Get statistics
// @Description get statistics about Keptn installation
// @Tags Statistics
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   from     query    string     false        "From"
// @Param   to     query    string     false        "To"
// @Success 200 {object} operations.Statistics	"ok"
// @Failure 400 {object} operations.Error "Invalid payload"
// @Failure 500 {object} operations.Error "Internal error"
// @Router /statistics [get]
func GetStatistics(c *gin.Context) {
	logger := keptn.NewLogger("", "", "statistics-service")
	params := &operations.GetStatisticsParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		c.JSON(http.StatusBadRequest, operations.Error{
			ErrorCode: 400,
			Message:   "Invalid request format",
		})
		return
	}

	if !validateQueryTimestamps(params) {
		c.JSON(http.StatusBadRequest, operations.Error{
			ErrorCode: 400,
			Message:   "Invalid time frame: 'from' timestamp must not be greater than 'to' timestamp",
		})
		return
	}

	sb := controller.GetStatisticsBucketInstance()

	payload, err := getStatistics(params, sb)

	if err != nil && err == db.ErrNoStatisticsFound {
		c.JSON(http.StatusNotFound, operations.Error{
			Message:   "no statistics found for selected time frame",
			ErrorCode: 404,
		})
		return
	} else if err != nil {
		logger.Error("could not retrieve statistics: " + err.Error())
		c.JSON(http.StatusInternalServerError, operations.Error{
			Message:   "Internal server error",
			ErrorCode: 500,
		})
		return
	}
	payload.From = params.From
	payload.To = params.To

	c.JSON(http.StatusOK, payload)
}

func getStatistics(params *operations.GetStatisticsParams, sb controller.StatisticsInterface) (operations.GetStatisticsResponse, error) {
	var mergedStatistics = operations.Statistics{}

	cutoffTime := sb.GetCutoffTime()

	// check time
	if params.From.After(cutoffTime) {
		// case 1: time frame within "in-memory" interval (e.g. last 30 minutes)
		// -> return in-memory object
		mergedStatistics = sb.GetStatistics()

	} else {
		var statistics []operations.Statistics
		var err error
		if params.From.Before(cutoffTime) && params.To.Before(cutoffTime) {
			// case 2: time frame outside of "in-memory" interval
			// -> return results from database
			statistics, err = sb.GetRepo().GetStatistics(params.From, params.To)
			if err != nil && err == db.ErrNoStatisticsFound {
				return operations.GetStatisticsResponse{}, err
			}
		} else if params.From.Before(cutoffTime) && params.To.After(cutoffTime) {
			// case 3: time frame includes "in-memory" interval
			// -> get results from database and from in-memory and merge them
			statistics, err = sb.GetRepo().GetStatistics(params.From, params.To)
			if statistics == nil {
				statistics = []operations.Statistics{}
			}
			statistics = append(statistics, sb.GetStatistics())
		}

		mergedStatistics = operations.Statistics{
			From: params.From,
			To:   params.To,
		}
		mergedStatistics = operations.MergeStatistics(mergedStatistics, statistics)
	}
	return convertToGetStatisticsResponse(mergedStatistics)
}

func convertToGetStatisticsResponse(mergedStatistics operations.Statistics) (operations.GetStatisticsResponse, error) {
	result := operations.GetStatisticsResponse{
		From:     mergedStatistics.From,
		To:       mergedStatistics.To,
		Projects: []operations.GetStatisticsResponseProject{},
	}

	for projectName, project := range mergedStatistics.Projects {
		newProject := operations.GetStatisticsResponseProject{
			Name:     projectName,
			Services: []operations.GetStatisticsResponseService{},
		}

		for serviceName, service := range project.Services {
			newService := operations.GetStatisticsResponseService{
				Name:                   serviceName,
				Events:                 []operations.GetStatisticsResponseEvent{},
				KeptnServiceExecutions: []operations.GetStatisticsResponseKeptnService{},
			}

			for eventType, eventTypeCount := range service.Events {
				newService.Events = append(newService.Events, operations.GetStatisticsResponseEvent{
					Type:  eventType,
					Count: eventTypeCount,
				})
			}

			env := config.GetConfig()
			if env.NextGenEvents {
				newService.ExecutedSequencesPerType = []operations.GetStatisticsResponseEvent{}
				for eventType, count := range service.ExecutedSequencesPerType {
					newService.ExecutedSequencesPerType = append(newService.ExecutedSequencesPerType, operations.GetStatisticsResponseEvent{
						Type:  eventType,
						Count: count,
					})
				}
			}

			for keptnServiceName, keptnService := range service.KeptnServiceExecutions {
				newKeptnService := operations.GetStatisticsResponseKeptnService{
					Name:       keptnServiceName,
					Executions: []operations.GetStatisticsResponseEvent{},
				}

				for eventType, eventTypeCount := range keptnService.Executions {
					newKeptnService.Executions = append(newKeptnService.Executions, operations.GetStatisticsResponseEvent{
						Type:  eventType,
						Count: eventTypeCount,
					})
				}
				newService.KeptnServiceExecutions = append(newService.KeptnServiceExecutions, newKeptnService)
			}
			newProject.Services = append(newProject.Services, newService)
		}
		result.Projects = append(result.Projects, newProject)
	}

	return result, nil
}

func validateQueryTimestamps(params *operations.GetStatisticsParams) bool {
	if params.To.Before(params.From) {
		return false
	}
	return true
}
