package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/models"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const SequenceExecutionCollectionNameSuffix = "SequenceExecution"

type MongoDBSequenceExecutionRepo struct {
	DbConnection *MongoDBConnection
}

func NewMongoDBSequenceExecutionRepo(dbConnection *MongoDBConnection) *MongoDBSequenceExecutionRepo {
	return &MongoDBSequenceExecutionRepo{DbConnection: dbConnection}
}

func (mdbrepo *MongoDBSequenceExecutionRepo) Get(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext(filter.Scope.Project)
	if err != nil {
		return nil, err
	}
	defer cancel()

	searchOptions := mdbrepo.getSearchOptions(filter)

	cur, err := collection.Find(ctx, searchOptions)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	result := []models.SequenceExecution{}

	for cur.Next(ctx) {
		var outInterface interface{}
		err := cur.Decode(&outInterface)
		sequenceExecution, err := transformBSONToSequenceExecution(outInterface)
		if err != nil {
			log.Errorf("Could not decode sequenceExecution: %v", err)
			continue
		}
		result = append(result, *sequenceExecution)
	}

	return result, nil
}

func (mdbrepo *MongoDBSequenceExecutionRepo) GetByTriggeredID(project, triggeredID string) (*models.SequenceExecution, error) {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext(project)
	if err != nil {
		return nil, err
	}
	defer cancel()

	searchOptions := bson.M{}
	searchOptions = appendFilterAs(searchOptions, triggeredID, "scope.triggeredId")

	item := collection.FindOne(ctx, searchOptions)
	if item.Err() != nil && item.Err() != mongo.ErrNoDocuments {
		return nil, err
	} else if item.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}

	var outInterface interface{}
	if err := item.Decode(&outInterface); err != nil {
		return nil, err
	}
	sequenceExecution, err := transformBSONToSequenceExecution(outInterface)
	if err != nil {
		return nil, err
	}

	return sequenceExecution, nil
}

func transformBSONToSequenceExecution(outInterface interface{}) (*models.SequenceExecution, error) {
	outInterface, err := flattenRecursively(outInterface)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(outInterface)

	sequenceExecution := &models.SequenceExecution{}
	if err := json.Unmarshal(data, sequenceExecution); err != nil {
		return nil, err
	}
	//sequenceExecution.ID = outInterface.(map[string]interface{})["_id"].(string)
	return sequenceExecution, nil
}

func (mdbrepo *MongoDBSequenceExecutionRepo) Upsert(item models.SequenceExecution, upsertOptions *models.SequenceExecutionUpsertOptions) error {
	if item.Scope.Project == "" {
		return errors.New("project must be set")
	}
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext(item.Scope.Project)
	if err != nil {
		return err
	}
	defer cancel()

	if upsertOptions != nil && upsertOptions.CheckUniqueTriggeredID {
		existingSequence, err := mdbrepo.GetByTriggeredID(item.Scope.Project, item.Scope.TriggeredID)
		if err != nil {
			return fmt.Errorf("could not check for existing sequence with same triggeredID: %w", err)
		}
		if existingSequence != nil {
			return ErrSequenceWithTriggeredIDAlreadyExists
		}
	}
	opts := options.Update().SetUpsert(true)

	filter := bson.D{{"_id", item.ID}}
	update := bson.D{{"$set", item}}

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func (mdbrepo *MongoDBSequenceExecutionRepo) AppendTaskEvent(taskSequence models.SequenceExecution, event models.TaskEvent) (*models.SequenceExecution, error) {
	if taskSequence.Scope.Project == "" {
		return nil, errors.New("project must be set")
	}
	if taskSequence.ID == "" {
		return nil, errors.New("id of sequenceExecution must be set")
	}
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext(taskSequence.Scope.Project)
	if err != nil {
		return nil, err
	}
	defer cancel()

	// return the resulting document after the update
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	filter := bson.D{{"_id", taskSequence.ID}}

	update := bson.M{"$push": bson.M{"status.currentTask.events": event}}

	res := collection.FindOneAndUpdate(ctx, filter, update, opts)
	if res.Err() != nil {
		return nil, err
	}

	outInterface := map[string]interface{}{}
	err = res.Decode(outInterface)
	if err != nil {
		return nil, err
	}
	sequenceExecution, err := transformBSONToSequenceExecution(outInterface)
	if err != nil {
		return nil, err
	}
	return sequenceExecution, nil
}

func (mdbrepo *MongoDBSequenceExecutionRepo) UpdateStatus(taskSequence models.SequenceExecution, state string) (*models.SequenceExecution, error) {
	if taskSequence.Scope.Project == "" {
		return nil, errors.New("project must be set")
	}
	if taskSequence.ID == "" {
		return nil, errors.New("id of sequenceExecution must be set")
	}
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext(taskSequence.Scope.Project)
	if err != nil {
		return nil, err
	}
	defer cancel()

	// return the resulting document after the update
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	filter := bson.D{{"_id", taskSequence.ID}}

	update := bson.M{"$set": bson.M{"status.state": state}}

	res := collection.FindOneAndUpdate(ctx, filter, update, opts)
	if res.Err() != nil {
		return nil, err
	}

	outInterface := map[string]interface{}{}
	err = res.Decode(outInterface)
	if err != nil {
		return nil, err
	}
	sequenceExecution, err := transformBSONToSequenceExecution(outInterface)
	if err != nil {
		return nil, err
	}
	return sequenceExecution, nil
}

func (mdbrepo *MongoDBSequenceExecutionRepo) Clear(projectName string) error {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext(projectName)
	if err != nil {
		return err
	}
	defer cancel()

	_, err = collection.DeleteMany(ctx, bson.D{})
	return err
}

func (mdbrepo *MongoDBSequenceExecutionRepo) getCollectionAndContext(project string) (*mongo.Collection, context.Context, context.CancelFunc, error) {
	collectionName := fmt.Sprintf("%s-%s", project, SequenceExecutionCollectionNameSuffix)
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return nil, nil, nil, err
	}
	collection := mdbrepo.DbConnection.Client.Database(getDatabaseName()).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return collection, ctx, cancel, nil
}

func (mdbrepo *MongoDBSequenceExecutionRepo) getSearchOptions(filter models.SequenceExecutionFilter) bson.M {
	searchOptions := bson.M{}

	searchOptions = appendFilterAs(searchOptions, filter.Name, "sequence.name")
	searchOptions = appendFilterAs(searchOptions, filter.Scope.TriggeredID, "scope.triggeredId")
	searchOptions = appendFilterAs(searchOptions, filter.Scope.KeptnContext, "scope.keptnContext")
	searchOptions = appendFilterAs(searchOptions, filter.Scope.Project, "scope.project")
	searchOptions = appendFilterAs(searchOptions, filter.Scope.Stage, "scope.stage")
	searchOptions = appendFilterAs(searchOptions, filter.Scope.Service, "scope.service")
	searchOptions = appendFilterAs(searchOptions, filter.CurrentTriggeredID, "status.currentTask.triggeredID")

	if filter.Status != nil && len(filter.Status) > 0 {
		matchStates := []bson.M{}
		for _, status := range filter.Status {
			match := bson.M{"status.state": status}

			matchStates = append(matchStates, match)
		}
		searchOptions["$or"] = matchStates
	}

	return searchOptions
}

func appendFilterAs(filter bson.M, value, key string) bson.M {
	if value != "" {
		filter[key] = value
	}
	return filter
}
