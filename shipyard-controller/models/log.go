package models

import keptnmodels "github.com/keptn/go-utils/pkg/api/models"

type CreateLogsRequest struct {
	Logs []LogEntry `form:"logs" json"logs"`
}

type GetLogParams struct {
	NextPageKey int64 `form:"nextPageKey" json:"nextPageKey"`
	PageSize    int64 `form:"pageSize" json:"pageSize"`

	IntegrationID string `form:"integrationID" json:"integrationID"`
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
	Logs []LogEntry `json:"logs"`
}

type LogEntry keptnmodels.LogEntry
