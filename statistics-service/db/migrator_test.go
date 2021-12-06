package db

import (
	"context"
	"github.com/keptn/keptn/statistics-service/operations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_Migrate(t *testing.T) {
	defer setupLocalMongoDB(t)()
	migrator := NewMigrator(1, time.Millisecond*100)

	statistics := operations.Statistics{
		From: time.Now(),
		To:   time.Now().Add(time.Second),
		Projects: map[string]*operations.Project{
			"my-project": {
				Name: "my-project",
				Services: map[string]*operations.Service{
					"my-service": {
						Name: "my-service",
						Events: map[string]int{
							"my.keptn.event.type": 2,
						},
						KeptnServiceExecutions: map[string]*operations.KeptnService{
							"my-keptn-service": {
								Name: "my-keptn-service",
								Executions: map[string]int{
									"my.keptn.event.type": 1,
								},
							},
						},
						ExecutedSequencesPerType: map[string]int{
							"my~keptn.event.type": 1,
						},
					},
				},
			},
		},
	}

	repo := StatisticsMongoDBRepo{}

	statistics.From = time.Now().Add(time.Second * 1)
	statistics.From = time.Now().Add(time.Second * 2)
	insertStat(t, &repo, statistics)

	statistics.From = time.Now().Add(time.Second * 3)
	statistics.From = time.Now().Add(time.Second * 4)
	repo.StoreStatistics(statistics)

	statistics.From = time.Now().Add(time.Second * 5)
	statistics.From = time.Now().Add(time.Second * 6)
	insertStat(t, &repo, statistics)

	statistics.From = time.Now().Add(time.Second * 7)
	statistics.From = time.Now().Add(time.Second * 7)
	repo.StoreStatistics(statistics)

	migratedDocs, err := migrator.Run(context.TODO())
	require.Nil(t, err)
	assert.Equal(t, uint(2), migratedDocs)

	fetchedStats, err := repo.GetStatistics(time.Time{}, time.Now().Add(10*time.Hour))
	require.Nil(t, err)
	require.Equal(t, 4, len(fetchedStats))
	for _, f := range fetchedStats {
		assert.Equal(t, statistics.Projects, f.Projects)
	}
}

func insertStat(t *testing.T, s *StatisticsMongoDBRepo, statistics operations.Statistics) {
	err := s.getCollection()
	require.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	_, err = s.statsCollection.InsertOne(ctx, statistics)
	require.Nil(t, err)
}
