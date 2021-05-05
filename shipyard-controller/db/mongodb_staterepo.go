package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const taskSequenceStateCollectionSuffix = "-taskSequenceStates"

var ErrNoStateFound = errors.New("no sequence state found")
var ErrStateAlreadyExists = errors.New("sequence state already exists")

type MongoDBStateRepo struct {
	DbConnection MongoDBConnection
}

func (mdbrepo *MongoDBStateRepo) CreateState(state models.SequenceState) error {

	if state.Project == "" {
		return errors.New("project must be set")
	}
	if state.Shkeptncontext == "" {
		return errors.New("shkeptncontext must be set")
	}
	if state.Name == "" {
		return errors.New("name must be set")
	}
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	databaseName := getDatabaseName()
	collection := mdbrepo.DbConnection.Client.Database(databaseName).Collection(state.Project + taskSequenceStateCollectionSuffix)

	existingSequence := collection.FindOne(ctx, bson.M{"shkeptncontext": state.Shkeptncontext})
	if existingSequence.Err() == nil || existingSequence.Err() != mongo.ErrNoDocuments {
		return ErrStateAlreadyExists
	}

	if _, err := collection.InsertOne(ctx, state); err != nil {
		return err
	}
	return nil
}

func (mdbrepo *MongoDBStateRepo) FindStates(filter models.StateFilter) (*models.SequenceStates, error) {
	if filter.Project == "" {
		return nil, errors.New("project must be set")
	}
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return nil, err
	}
	collection := mdbrepo.DbConnection.Client.Database(getDatabaseName()).Collection(filter.Project + taskSequenceStateCollectionSuffix)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	searchOptions := mdbrepo.getSearchOptions(filter)

	totalCount, err := collection.CountDocuments(ctx, searchOptions)
	if err != nil {
		return nil, fmt.Errorf("error counting elements in events collection: %v", err)
	}

	sortOptions := options.Find().SetSort(bson.D{{Key: "time", Value: -1}}).SetSkip(filter.NextPageKey)

	if filter.PageSize > 0 {
		sortOptions = sortOptions.SetLimit(filter.PageSize)
	}

	cur, err := collection.Find(ctx, searchOptions, sortOptions)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	result := &models.SequenceStates{
		States:      []models.SequenceState{},
		NextPageKey: 0,
		PageSize:    0,
		TotalCount:  totalCount,
	}
	states := []models.SequenceState{}

	if filter.PageSize > 0 && filter.PageSize+filter.NextPageKey < totalCount {
		result.NextPageKey = filter.PageSize + filter.NextPageKey
	}

	for cur.Next(ctx) {
		sequenceState := &models.SequenceState{}
		if err := cur.Decode(sequenceState); err != nil {
			// TODO log
		}
		states = append(states, *sequenceState)
	}
	result.States = states
	return result, nil
}

func (mdbrepo *MongoDBStateRepo) getSearchOptions(filter models.StateFilter) bson.M {
	searchOptions := bson.M{
		"project": filter.Project,
	}

	if filter.Shkeptncontext != "" {
		searchOptions["shkeptncontext"] = filter.Shkeptncontext
	}

	if filter.Name != "" {
		searchOptions["name"] = filter.Name
	}
	return searchOptions
}

func (mdbrepo *MongoDBStateRepo) UpdateState(state models.SequenceState) error {
	if state.Project == "" {
		return errors.New("project must be set")
	}
	if state.Shkeptncontext == "" {
		return errors.New("shkeptncontext must be set")
	}
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.DbConnection.Client.Database(getDatabaseName()).Collection(state.Project + taskSequenceStateCollectionSuffix)
	_, err = collection.ReplaceOne(ctx, bson.M{"shkeptncontext": state.Shkeptncontext}, state)
	if err != nil {
		return err
	}
	return nil
}

func (mdbrepo *MongoDBStateRepo) DeleteStates(filter models.StateFilter) error {
	if filter.Project == "" {
		return errors.New("project must be set")
	}
	if filter.Shkeptncontext == "" {
		return errors.New("shkeptncontext must be set")
	}
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}

	searchOptions := mdbrepo.getSearchOptions(filter)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.DbConnection.Client.Database(getDatabaseName()).Collection(filter.Project + taskSequenceStateCollectionSuffix)
	_, err = collection.DeleteMany(ctx, searchOptions)
	if err != nil {
		return err
	}
	return nil
}
