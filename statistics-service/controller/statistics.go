package controller

import (
	"fmt"
	"github.com/keptn-sandbox/statistics-service/statistics-service/config"
	"github.com/keptn-sandbox/statistics-service/statistics-service/db"
	"github.com/keptn-sandbox/statistics-service/statistics-service/operations"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"strings"
	"sync"
	"time"
)

var statisticsBucketInstance *statisticsBucket

type statisticsBucket struct {
	StatisticsRepo  db.StatisticsRepo
	Statistics      operations.Statistics
	uniqueSequences map[string]bool
	logger          keptn.LoggerInterface
	lock            sync.Mutex
	cutoffTime      time.Time
	nextGenEvents   bool
}

// GetStatisticsBucketInstance godoc
func GetStatisticsBucketInstance() *statisticsBucket {
	if statisticsBucketInstance == nil {
		env := config.GetConfig()
		statisticsBucketInstance = &statisticsBucket{
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

// GetCutoffTime
func (sb *statisticsBucket) GetCutoffTime() time.Time {
	return sb.cutoffTime
}

// GetStatistics godoc
func (sb *statisticsBucket) GetStatistics() operations.Statistics {
	return sb.Statistics
}

// GetRepo godoc
func (sb *statisticsBucket) GetRepo() db.StatisticsRepo {
	return sb.StatisticsRepo
}

// AddEvent godoc
func (sb *statisticsBucket) AddEvent(event operations.Event) {
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

func (sb *statisticsBucket) storeCurrentBucket() {
	sb.lock.Lock()
	defer sb.lock.Unlock()
	sb.logger.Info(fmt.Sprintf("Storing statistics for time frame %s - %s\n\n", sb.Statistics.From.String(), sb.Statistics.To.String()))
	sb.Statistics.To = time.Now().Round(time.Second)
	if err := sb.StatisticsRepo.StoreStatistics(sb.Statistics); err != nil {
		sb.logger.Error(fmt.Sprintf("Could not store statistics: " + err.Error()))
	}
	sb.logger.Info(fmt.Sprintf("Statistics stored successfully"))
}

func (sb *statisticsBucket) createNewBucket() {
	sb.lock.Lock()
	defer sb.lock.Unlock()
	sb.cutoffTime = time.Now().Round(time.Second)
	sb.uniqueSequences = map[string]bool{}
	sb.Statistics = operations.Statistics{
		From: time.Now().Round(time.Second),
	}
}
