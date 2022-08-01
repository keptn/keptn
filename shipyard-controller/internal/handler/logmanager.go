package handler

import (
	"github.com/keptn/keptn/shipyard-controller/internal/db"
	"github.com/keptn/keptn/shipyard-controller/models"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/logmanager.go . ILogManager
type ILogManager interface {
	CreateLogEntries(entry models.CreateLogsRequest) error
	GetLogEntries(filter models.GetLogParams) (*models.GetLogsResponse, error)
	DeleteLogEntries(params models.DeleteLogParams) error
}

type LogManager struct {
	logRepo db.LogRepo
}

func NewLogManager(logRepo db.LogRepo) *LogManager {
	return &LogManager{logRepo: logRepo}
}

func (lm *LogManager) CreateLogEntries(entries models.CreateLogsRequest) error {
	return lm.logRepo.CreateLogEntries(entries.Logs)
}

func (lm *LogManager) GetLogEntries(filter models.GetLogParams) (*models.GetLogsResponse, error) {
	return lm.logRepo.GetLogEntries(filter)
}

func (lm *LogManager) DeleteLogEntries(params models.DeleteLogParams) error {
	return lm.logRepo.DeleteLogEntries(params)
}
