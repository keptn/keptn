package fake

import (
	"context"

	"github.com/keptn/go-utils/pkg/api/models"
)

type LogInterfaceMock struct {
	LogFn func(logs []models.LogEntry)
}

func (l *LogInterfaceMock) Log(logs []models.LogEntry) {
	if l.LogFn != nil {
		l.LogFn(logs)
	}
}

func (l *LogInterfaceMock) Flush() error {
	panic("implement me")
}

func (l *LogInterfaceMock) GetLogs(params models.GetLogsParams) (*models.GetLogsResponse, error) {
	panic("implement me")
}

func (l *LogInterfaceMock) DeleteLogs(filter models.LogFilter) error {
	panic("implement me")
}

func (l *LogInterfaceMock) Start(ctx context.Context) {
	panic("implement me")
}
