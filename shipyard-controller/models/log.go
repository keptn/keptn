package models

import apimodels "github.com/keptn/go-utils/pkg/api/models"

type CreateLogsRequest struct {
	// logs
	Logs []apimodels.LogEntry `form:"logs" json:"logs"`
}

type CreateLogsResponse struct{}

type GetLogParams struct {
	LogFilter

	NextPageKey int64 `form:"nextPageKey" json:"nextPageKey"`
	PageSize    int64 `form:"pageSize" json:"pageSize"`
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
	// Pointer to next page
	NextPageKey int64 `json:"nextPageKey,omitempty"`

	// Size of returned page
	PageSize int64 `json:"pageSize,omitempty"`

	// Total number of logs
	TotalCount int64 `json:"totalCount,omitempty"`

	// logs
	Logs []apimodels.LogEntry `json:"logs"`
}
