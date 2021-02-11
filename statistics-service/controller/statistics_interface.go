package controller

import (
	"time"

	"github.com/keptn/keptn/statistics-service/db"
	"github.com/keptn/keptn/statistics-service/operations"
)

// StatisticsInterface godoc
type StatisticsInterface interface {
	// GetCutoffTime godoc
	GetCutoffTime() time.Time
	// GetStatistics godoc
	GetStatistics() operations.Statistics
	// AddEvent godoc
	AddEvent(event operations.Event)
	// GetRepo godoc
	GetRepo() db.StatisticsRepo
}
