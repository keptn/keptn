package db

import (
	"context"
	"fmt"
	"github.com/keptn/keptn/statistics-service/operations"
	logger "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"testing"
	"time"
)

var mongoDbVersion = "4.4.9"

func setupLocalMongoDB() func() {
	mongoServer, err := memongo.Start(mongoDbVersion)
	randomDbName := memongo.RandomDatabase()

	os.Setenv("MONGODB_DATABASE", randomDbName)
	os.Setenv("MONGODB_EXTERNAL_CONNECTION_STRING", fmt.Sprintf("%s/%s", mongoServer.URI(), randomDbName))

	var mongoDBClient *mongo.Client
	mongoDBClient, err = mongo.NewClient(options.Client().ApplyURI(mongoServer.URI()))
	if err != nil {
		logger.Fatalf("Mongo Client setup failed: %s", err)
	}
	err = mongoDBClient.Connect(context.TODO())
	if err != nil {
		log.Fatalf("Mongo Server setup failed: %s", err)
	}

	return func() { mongoServer.Stop() }
}

func TestStatisticsMongoDBRepo_Store_And_Get_Statistics(t *testing.T) {
	defer setupLocalMongoDB()()
	type args struct {
		statistics operations.Statistics
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "repo - store statistics",
			args: args{
				statistics: operations.Statistics{
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
										"my.keptn.event.type": 1,
									},
								},
							},
						},
					},
				},
			},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StatisticsMongoDBRepo{}
			if err := s.StoreStatistics(tt.args.statistics); (err != nil) != tt.wantErr {
				t.Errorf("StoreStatistics() error = %v, wantErr %v", err, tt.wantErr)
			}
			fetchedStats, err := s.GetStatistics(time.Time{}, time.Now().Add(time.Second*10))
			if err != nil {
				t.Errorf("GetStatistics() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, len([]operations.Statistics{tt.args.statistics}), len(fetchedStats))
			assert.Equal(t, tt.args.statistics.Projects, fetchedStats[0].Projects)
		})
	}
}

func TestStatisticsMongoDBRepo_MigrateKeys(t *testing.T) {
	defer setupLocalMongoDB()()
	type args struct {
		statistics operations.Statistics
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "repo - migrate keys",
			args: args{
				statistics: operations.Statistics{
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
										"my.keptn.event.type": 1,
									},
								},
							},
						},
					},
				},
			},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StatisticsMongoDBRepo{}

			insertStat(t, s, tt.args.statistics)
			numDocsMigrated, err := s.MigrateKeys()
			require.Nil(t, err)
			assert.Equal(t, uint(1), numDocsMigrated)

			fetchedStats, err := s.GetStatistics(time.Time{}, time.Now().Add(time.Second*10))
			if err != nil {
				t.Errorf("GetStatistics() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, len([]operations.Statistics{tt.args.statistics}), len(fetchedStats))
			assert.Equal(t, tt.args.statistics.Projects, fetchedStats[0].Projects)
		})
	}
}

func TestStatisticsMongoDBRepo_ReMigrateKeys(t *testing.T) {
	defer setupLocalMongoDB()()
	type args struct {
		statistics operations.Statistics
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "repo - migrate already migrated keys",
			args: args{
				statistics: operations.Statistics{
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
										"my.keptn.event.type": 1,
									},
								},
							},
						},
					},
				},
			},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StatisticsMongoDBRepo{}

			err := s.StoreStatistics(tt.args.statistics)
			require.Nil(t, err)

			numDocsMigrated, err := s.MigrateKeys()
			require.Nil(t, err)
			assert.Equal(t, uint(1), numDocsMigrated)

			fetchedStats, err := s.GetStatistics(time.Time{}, time.Now().Add(time.Second*10))
			if err != nil {
				t.Errorf("GetStatistics() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, len([]operations.Statistics{tt.args.statistics}), len(fetchedStats))
			assert.Equal(t, tt.args.statistics.Projects, fetchedStats[0].Projects)
		})
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
