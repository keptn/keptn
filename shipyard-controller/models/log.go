package models

import apimodels "github.com/keptn/go-utils/pkg/api/models"

type CreateLogsRequest struct {
	// logs
	Logs []apimodels.LogEntry `form:"logs" json:"logs"`
}

type CreateLogsResponse struct{}

type GetLogParams struct {
	LogFilter
	PaginationParams
}

type DeleteLogParams struct {
	LogFilter
}

type DeleteLogResponse struct{}

type LogFilter struct {
	IntegrationID string `form:"integrationId" json:"integrationId"`
	FromTime      string `form:"fromTime" json:"fromTime"`
	BeforeTime    string `form:"beforeTime" json:"beforeTime"`
}

type GetLogsResponse struct {
	PaginationResult

	// logs
	Logs []apimodels.LogEntry `json:"logs"`
}
