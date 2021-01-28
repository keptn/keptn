package controller

import (
	"fmt"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/statistics-service/config"
	"github.com/keptn/keptn/statistics-service/db"
	"github.com/keptn/keptn/statistics-service/operations"
	"strings"
	"sync"
	"time"
)

var statisticsBucketInstance *StatisticsBucket

// StatisticsBucket godoc
type StatisticsBucket struct {
	// StatisticsRepo interface for accessing the repository
	StatisticsRepo db.StatisticsRepo
	// Statistics in-memory statistics
	Statistics      operations.Statistics
	uniqueSequences map[string]bool
	logger          keptn.LoggerInterface
	lock            sync.Mutex
	cutoffTime      time.Time
	nextGenEvents   bool
}

// GetStatisticsBucketInstance godoc
func GetStatisticsBucketInstance() *StatisticsBucket {
	if statisticsBucketInstance == nil {
		env := config.GetConfig()
		statisticsBucketInstance = &StatisticsBucket{
			StatisticsRepo: &db.StatisticsMongoDBRepo{},
			logger:         keptn.NewLogger("", "", "statistics service"),
			nextGenEvents:  env.NextGenEvents,
		}

		statisticsBucketInstance.createNewBucket()
		go func() {
			bucketInterval := time.Duration(env.AggregationIntervalSeconds) * time.Second
			bucketTimer := time.NewTimer(bucketInterval)
			defer bucketTimer.Stop()
			for {
				bucketTimer.Reset(bucketInterval)
				<-bucketTimer.C
				statisticsBucketInstance.logger.Info(fmt.Sprintf("%d seconds have passed. Creating a new statistics bucket\n", env.AggregationIntervalSeconds))
				statisticsBucketInstance.storeCurrentBucket()
				statisticsBucketInstance.createNewBucket()
			}
		}()
	}
	return statisticsBucketInstance
}

// GetCutoffTime returns the cutoff time (=time window in which the service holds its data in memory)
func (sb *StatisticsBucket) GetCutoffTime() time.Time {
	return sb.cutoffTime
}

// GetStatistics godoc
func (sb *StatisticsBucket) GetStatistics() operations.Statistics {
	return sb.Statistics
}

// GetRepo godoc
func (sb *StatisticsBucket) GetRepo() db.StatisticsRepo {
	return sb.StatisticsRepo
}

// AddEvent godoc
func (sb *StatisticsBucket) AddEvent(event operations.Event) {
	sb.lock.Lock()
	defer sb.lock.Unlock()

	if event.Data.Project == "" || event.Data.Service == "" || event.Type == "" || event.Source == "" {
		return
	}
	sb.logger.Info("updating statistics for service " + event.Data.Service + " in project " + event.Data.Project)
	sb.uniqueSequences[event.Shkeptncontext] = true

	sb.Statistics.IncreaseEventTypeCount(event.Data.Project, event.Data.Service, event.Type, 1)

	if sb.nextGenEvents {
		// increase service execution count using .started events
		if strings.HasSuffix(event.Type, ".started") {
			sb.Statistics.IncreaseKeptnServiceExecutionCount(
				event.Data.Project,
				event.Data.Service,
				event.Source,
				strings.TrimSuffix(event.Type, ".started"), 1,
			)
		}
		if strings.HasSuffix(event.Type, ".finished") && event.Source == "shipyard-controller" {
			// when the shipyard controller sends a .finished event, this means that a task sequence has been completed
			sb.Statistics.IncreaseExecutedSequencesCount(event.Data.Project, event.Data.Service, 1)
			sb.Statistics.IncreaseExecutedSequenceCountForType(event.Data.Project, event.Data.Service, strings.TrimSuffix(event.Type, ".finished"), 1)
		}
	} else {
		// increase service execution count using 'source' property from event
		sb.Statistics.IncreaseKeptnServiceExecutionCount(
			event.Data.Project,
			event.Data.Service,
			event.Source,
			event.Type, 1,
		)
	}
}

func (sb *StatisticsBucket) storeCurrentBucket() {
	sb.lock.Lock()
	defer sb.lock.Unlock()
	sb.logger.Info(fmt.Sprintf("Storing statistics for time frame %s - %s\n\n", sb.Statistics.From.String(), sb.Statistics.To.String()))
	sb.Statistics.To = time.Now().Round(time.Second)
	if err := sb.StatisticsRepo.StoreStatistics(sb.Statistics); err != nil {
		sb.logger.Error(fmt.Sprintf("Could not store statistics: " + err.Error()))
	}
	sb.logger.Info(fmt.Sprintf("Statistics stored successfully"))
}

func (sb *StatisticsBucket) createNewBucket() {
	sb.lock.Lock()
	defer sb.lock.Unlock()
	sb.cutoffTime = time.Now().Round(time.Second)
	sb.uniqueSequences = map[string]bool{}
	sb.Statistics = operations.Statistics{
		From: time.Now().Round(time.Second),
	}
}
