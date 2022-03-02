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

const sequenceExecutionCollectionNameSuffix = "sequence-execution"

// eventQueueSequenceStateCollectionName contains information on whether a task sequence is currently paused and thus outgoing events should be blocked
const eventQueueSequenceStateCollectionName = "shipyard-controller-event-queue-sequence-state"

var ErrProjectNameMustNotBeEmpty = errors.New("project name must not be empty")
var ErrSequenceIDMustNotBeEmpty = errors.New("sequence ID must not be empty")

type MongoDBSequenceExecutionRepo struct {
	DbConnection *MongoDBConnection
}

func NewMongoDBSequenceExecutionRepo(dbConnection *MongoDBConnection) *MongoDBSequenceExecutionRepo {
	return &MongoDBSequenceExecutionRepo{DbConnection: dbConnection}
}

func (mdbrepo *MongoDBSequenceExecutionRepo) Get(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
	collection, ctx, cancel, err := mdbrepo.getSequenceExecutionStateCollection(filter.Scope.Project)
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
	collection, ctx, cancel, err := mdbrepo.getSequenceExecutionStateCollection(project)
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

func (mdbrepo *MongoDBSequenceExecutionRepo) Upsert(item models.SequenceExecution, upsertOptions *models.SequenceExecutionUpsertOptions) error {
	if item.Scope.Project == "" {
		return ErrProjectNameMustNotBeEmpty
	}
	collection, ctx, cancel, err := mdbrepo.getSequenceExecutionStateCollection(item.Scope.Project)
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
		return nil, ErrProjectNameMustNotBeEmpty
	}
	if taskSequence.ID == "" {
		return nil, ErrSequenceIDMustNotBeEmpty
	}
	collection, ctx, cancel, err := mdbrepo.getSequenceExecutionStateCollection(taskSequence.Scope.Project)
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

func (mdbrepo *MongoDBSequenceExecutionRepo) UpdateStatus(taskSequence models.SequenceExecution) (*models.SequenceExecution, error) {
	if taskSequence.Scope.Project == "" {
		return nil, ErrProjectNameMustNotBeEmpty
	}
	if taskSequence.ID == "" {
		return nil, ErrSequenceIDMustNotBeEmpty
	}
	collection, ctx, cancel, err := mdbrepo.getSequenceExecutionStateCollection(taskSequence.Scope.Project)
	if err != nil {
		return nil, err
	}
	defer cancel()

	// return the resulting document after the update
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	filter := bson.D{{"_id", taskSequence.ID}}

	update := bson.M{"$set": bson.M{
		"status.state":            taskSequence.Status.State,
		"status.stateBeforePause": taskSequence.Status.StateBeforePause,
	}}

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
	collection, ctx, cancel, err := mdbrepo.getSequenceExecutionStateCollection(projectName)
	if err != nil {
		return err
	}
	defer cancel()

	_, err = collection.DeleteMany(ctx, bson.D{})
	return err
}

func (mdbrepo *MongoDBSequenceExecutionRepo) PauseContext(eventScope models.EventScope) error {
	return mdbrepo.updateGlobalSequenceContext(eventScope, models.SequencePaused)
}

func (mdbrepo *MongoDBSequenceExecutionRepo) ResumeContext(eventScope models.EventScope) error {
	return mdbrepo.updateGlobalSequenceContext(eventScope, models.SequenceStartedState)
}

func (mdbrepo *MongoDBSequenceExecutionRepo) IsContextPaused(eventScope models.EventScope) bool {
	collection, ctx, cancel, err := mdbrepo.getCollection(eventQueueSequenceStateCollectionName)
	if err != nil {
		log.Errorf("Could not get collection: %v", err)
		return false
	}
	defer cancel()

	searchOptions := bson.M{}
	if eventScope.KeptnContext != "" {
		searchOptions[keptnContextScope] = eventScope.KeptnContext
	}
	cur, err := collection.Find(ctx, searchOptions)
	if err != nil {
		log.Errorf("Could not retrieve sequence context: %v", err)
		return false
	} else if cur.RemainingBatchLength() == 0 {
		return false
	}

	stateItems := []models.EventQueueSequenceState{}

	defer func() {
		err := cur.Close(ctx)
		if err != nil {
			log.Errorf("could not close cursor: %v", err)
		}
	}()
	for cur.Next(ctx) {
		stateItem := models.EventQueueSequenceState{}
		err := cur.Decode(&stateItem)
		if err != nil {
			log.Errorf("Could not decode item: %v", err)
			continue
		}
		stateItems = append(stateItems, stateItem)
	}

	for _, state := range stateItems {
		if state.Scope.Stage == "" && state.State == models.SequencePaused {
			// if the overall state is set to 'paused', this means that all stages are paused
			return true
		} else if state.Scope.Stage == eventScope.Stage && state.State == models.SequencePaused {
			// if not the overall state is 'paused', but specifically for this stage, we return true as well
			return true
		}
	}

	return false
}

func (mdbrepo *MongoDBSequenceExecutionRepo) getSequenceExecutionStateCollection(project string) (*mongo.Collection, context.Context, context.CancelFunc, error) {
	collectionName := fmt.Sprintf("%s-%s", project, sequenceExecutionCollectionNameSuffix)
	return mdbrepo.getCollection(collectionName)
}

func (mdbrepo *MongoDBSequenceExecutionRepo) getCollection(collectionName string) (*mongo.Collection, context.Context, context.CancelFunc, error) {
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

func (mdbrepo *MongoDBSequenceExecutionRepo) updateGlobalSequenceContext(eventScope models.EventScope, status string) error {
	collection, ctx, cancel, err := mdbrepo.getCollection(eventQueueSequenceStateCollectionName)
	if err != nil {
		return err
	}
	defer cancel()

	state := models.EventQueueSequenceState{
		State: status,
		Scope: eventScope,
	}

	opts := options.Update().SetUpsert(true)

	var filter bson.D
	if eventScope.Stage == "" {
		filter = bson.D{
			{keptnContextScope, eventScope.KeptnContext},
		}
	} else {
		filter = bson.D{
			{keptnContextScope, eventScope.KeptnContext},
			{stageScope, eventScope.Stage},
		}
	}
	update := bson.D{{"$set", state}}

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	return err
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

func appendFilterAs(filter bson.M, value, key string) bson.M {
	if value != "" {
		filter[key] = value
	}
	return filter
}
