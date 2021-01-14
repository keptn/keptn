package db

import (
	"errors"
	"github.com/keptn-sandbox/statistics-service/statistics-service/operations"
	"time"
)

// NoStatisticsFoundError godoc
var NoStatisticsFoundError = errors.New("no statistics found")

// StatisticsRepo godoc
type StatisticsRepo interface {
	// GetStatistics godoc
	GetStatistics(from, to time.Time) ([]operations.Statistics, error)
	// StoreStatistics godoc
	StoreStatistics(statistics operations.Statistics) error
	// DeleteStatistics godoc
	DeleteStatistics(from, to time.Time) error
}
