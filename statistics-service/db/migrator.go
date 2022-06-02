package db

import (
	"context"
	"github.com/globalsign/mgo/bson"
	"github.com/keptn/keptn/statistics-service/operations"
	"github.com/mitchellh/copystructure"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

// Migrator goes through each and every document and escapes keys which are known to dots (".")
// https://github.com/keptn/keptn/issues/6250
type Migrator struct {
	dbConnection    MongoDBConnection
	statsCollection *mongo.Collection
	batchSize       int
	interval        time.Duration
	migratedCount   uint
	nextBatchNum    int
}

// NewMigrator creates a new instance of a Migrator
func NewMigrator(batchSize int, interval time.Duration) *Migrator {
	return &Migrator{
		dbConnection: MongoDBConnection{},
		batchSize:    batchSize,
		interval:     interval,
		nextBatchNum: 1,
	}
}

func (m *Migrator) Run(ctx context.Context) (uint, error) {
	err := m.getCollection()
	if err != nil {
		return m.migratedCount, err
	}

	for {
		select {
		case <-ctx.Done():
			return m.migratedCount, nil
		case <-time.After(m.interval):
			done, err := m.migrateBatch(ctx)
			if done || err != nil {
				return m.migratedCount, err
			}
			log.Infof("Migrated documents: %d", m.migratedCount)
		}
	}
}

func (m *Migrator) migrateBatch(ctx context.Context) (bool, error) {
	skips := m.batchSize * (m.nextBatchNum - 1)
	m.nextBatchNum++
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cur, err := m.statsCollection.Find(ctx, bson.M{}, &options.FindOptions{
		Limit: crateInt64P(m.batchSize),
		Skip:  crateInt64P(skips),
	})
	defer func() {
		if cur == nil {
			return
		}
		cur.Close(ctx)
	}()
	if err != nil {
		return false, err
	}

	currBatchSize := 0
	for cur.Next(ctx) {
		stats := &operations.Statistics{}
		err := cur.Decode(stats)
		if err != nil {
			return false, err
		}
		currBatchSize++

		// if there is no project data, there is nothing to do
		if len(stats.Projects) == 0 {
			continue
		}
		// if there are no dots in keys, there is nothing to do
		if noDotsInKeys(stats) {
			continue
		}

		hexID := &HexID{}
		err = cur.Decode(hexID)
		if err != nil {
			return false, err
		}

		decodedKeys, err := decodeKeys([]operations.Statistics{*stats})
		if err != nil {
			return false, err
		}
		statsWithEncodedKeys, err := encodeKeys(&decodedKeys[0])
		if err != nil {
			return false, err
		}
		_, err = m.statsCollection.ReplaceOne(ctx, bson.M{"_id": hexID.ID}, statsWithEncodedKeys)
		if err != nil {
			return false, err
		}
		m.migratedCount++
	}
	if currBatchSize < m.batchSize {
		return true, nil
	}
	return false, nil
}
func noDotsInKeys(statistics *operations.Statistics) bool {
	for projKey, proj := range statistics.Projects {
		if strings.Contains(projKey, ".") {
			return false
		}
		for svcKey, svc := range proj.Services {
			if strings.Contains(svcKey, ".") {
				return false
			}
			for eventType := range svc.ExecutedSequencesPerType {
				if strings.Contains(eventType, ".") {
					return false
				}
			}
			for eventType := range svc.Events {
				if strings.Contains(eventType, ".") {
					return false
				}
			}
			for keptnServiceExecutionsKey, keptnService := range svc.KeptnServiceExecutions {
				if strings.Contains(keptnServiceExecutionsKey, ".") {
					return false
				}
				for eventType := range keptnService.Executions {
					if strings.Contains(eventType, ".") {
						return false
					}
				}
			}
		}
	}
	return true
}

func encodeKeys(statistics *operations.Statistics) (*operations.Statistics, error) {
	return transform(statistics, encodeKey)
}

func decodeKeys(statistics []operations.Statistics) ([]operations.Statistics, error) {
	newStatistics := []operations.Statistics{}
	for _, stat := range statistics {
		s, err := transform(&stat, decodeKey)
		if err != nil {
			return nil, err
		}
		newStatistics = append(newStatistics, *s)
	}
	return newStatistics, nil
}

func transform(statistics *operations.Statistics, transformFn func(string) string) (*operations.Statistics, error) {
	copiedStatistics, err := copystructure.Copy(statistics)
	if err != nil {
		return nil, err
	}
	newProjects := make(map[string]*operations.Project)
	for projKey, proj := range copiedStatistics.(*operations.Statistics).Projects {
		newServices := make(map[string]*operations.Service)
		for serviceKey, service := range proj.Services {
			newExecutedSequencesPerType := make(map[string]int)
			for eventType, numExecutedSequencesPerType := range service.ExecutedSequencesPerType {
				newExecutedSequencesPerType[transformFn(eventType)] = numExecutedSequencesPerType
			}
			service.ExecutedSequencesPerType = newExecutedSequencesPerType

			newEvents := make(map[string]int)
			for eventType, event := range service.Events {
				newEvents[transformFn(eventType)] = event
			}
			service.Events = newEvents

			newKeptnServiceExecutions := make(map[string]*operations.KeptnService)
			for keptnServiceExecutionKey, keptnService := range service.KeptnServiceExecutions {
				newServiceExecutions := make(map[string]int)
				for eventType2, numExecutions := range keptnService.Executions {
					newServiceExecutions[transformFn(eventType2)] = numExecutions
				}
				keptnService.Executions = newServiceExecutions
				newKeptnServiceExecutions[transformFn(keptnServiceExecutionKey)] = keptnService
			}
			service.KeptnServiceExecutions = newKeptnServiceExecutions
			newServices[transformFn(serviceKey)] = service
		}
		proj.Services = newServices
		newProjects[transformFn(projKey)] = proj
	}
	copiedStatistics.(*operations.Statistics).Projects = newProjects
	return copiedStatistics.(*operations.Statistics), nil
}

func encodeKey(key string) string {
	encodedKey := strings.ReplaceAll(strings.ReplaceAll(key, "~", "~t"), ".", "~p")
	return encodedKey
}
func decodeKey(key string) string {
	decodedKey := strings.ReplaceAll(strings.ReplaceAll(key, "~p", "."), "~t", "~")
	return decodedKey
}

func (m *Migrator) getCollection() error {
	err := m.dbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}

	if m.statsCollection == nil {
		m.statsCollection = m.dbConnection.Client.Database(databaseName).Collection(keptnStatsCollection)
	}
	return nil
}

func crateInt64P(x int) *int64 {
	i := int64(x)
	return &i
}
