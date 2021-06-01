package handler

import (
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
)

type ILogManager interface {
	CreateLogEntries(entry models.CreateLogsRequest) error
	GetLogEntries(filter models.GetLogParams) (models.GetLogsResponse, error)
}

type LogManager struct {
	logRepo db.LogRepo
}

func NewLogManager(logRepo db.LogRepo) *LogManager {
	return &LogManager{logRepo: logRepo}
}

func (LogManager) CreateLogEntries(entry models.CreateLogsRequest) error {
	panic("implement me")
}

func (LogManager) GetLogEntries(filter models.GetLogParams) (models.GetLogsResponse, error) {
	panic("implement me")
}
