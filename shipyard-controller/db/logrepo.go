package db

import "github.com/keptn/keptn/shipyard-controller/models"

type LogRepo interface {
	CreateLogEntries(entries []models.LogEntry) error
	GetLogEntries(filter models.GetLogParams) (*models.GetLogsResponse, error)
}
