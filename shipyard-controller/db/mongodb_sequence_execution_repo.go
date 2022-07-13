package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/db/models/sequence_execution"
	v02 "github.com/keptn/keptn/shipyard-controller/db/models/sequence_execution/v1"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const sequenceExecutionCollectionNameSuffix = "sequence-execution"

// eventQueueSequenceStateCollectionName contains information on whether a task sequence is currently paused and thus outgoing events should be blocked
const eventQueueSequenceStateCollectionName = "shipyard-controller-event-queue-sequence-state"

var ErrProjectNameMustNotBeEmpty = errors.New("project name must not be empty")
var ErrSequenceIDMustNotBeEmpty = errors.New("sequence ID must not be empty")

type SequenceExecutionRepoOpt func(repo *MongoDBSequenceExecutionRepo)

func WithSequenceExecutionModelTransformer(transformer sequence_execution.ModelTransformer) SequenceExecutionRepoOpt {
	return func(repo *MongoDBSequenceExecutionRepo) {
		repo.ModelTransformer = transformer
	}
}

type MongoDBSequenceExecutionRepo struct {
	DbConnection *MongoDBConnection
	// ModelTransformer allows transforming sequence execution objects before storing and after retrieving to and from the database.
	// This is needed if the items for that collection should be stored in a different format while retaining its structure used outside the package
	ModelTransformer sequence_execution.ModelTransformer
}

func NewMongoDBSequenceExecutionRepo(dbConnection *MongoDBConnection, opts ...SequenceExecutionRepoOpt) *MongoDBSequenceExecutionRepo {
	repo := &MongoDBSequenceExecutionRepo{
		DbConnection: dbConnection,
		// use the v1 ModelTransformer by default
		ModelTransformer: &v02.ModelTransformer{},
	}

	for _, opt := range opts {
		opt(repo)
	}
	return repo
}

// Get returns all matching sequence executions, based on the given filter
func (mdbrepo *MongoDBSequenceExecutionRepo) Get(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
	collection, ctx, cancel, err := mdbrepo.getSequenceExecutionStateCollection(filter.Scope.Project)
	if err != nil {
		return nil, err
	}
	defer cancel()

	searchOptions := mdbrepo.getSearchOptions(filter)

	cur, err := collection.Find(ctx, searchOptions)
	defer closeCursor(ctx, cur)

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	result := []models.SequenceExecution{}

	for cur.Next(ctx) {
		var outInterface interface{}
		err := cur.Decode(&outInterface)
		sequenceExecution, err := mdbrepo.transformBSONToSequenceExecution(outInterface)
		if err != nil {
			log.Errorf("Could not decode sequenceExecution: %v", err)
			continue
		}
		result = append(result, *sequenceExecution)
	}

	return result, nil
}

// GetPaginated returns all matching sequence executions within the specified range, based on the given filter
func (mdbrepo *MongoDBSequenceExecutionRepo) GetPaginated(filter models.SequenceExecutionFilter, paginationParams models.PaginationParams) ([]models.SequenceExecution, *models.PaginationResult, error) {
	collection, ctx, cancel, err := mdbrepo.getSequenceExecutionStateCollection(filter.Scope.Project)
	if err != nil {
		return nil, nil, err
	}
	defer cancel()

	searchOptions := mdbrepo.getSearchOptions(filter)

	totalCount, err := collection.CountDocuments(ctx, searchOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("error counting elements in sequence execution collection: %w", err)
	}

	sortOptions := options.Find().SetSort(bson.D{{Key: "time", Value: -1}}).SetSkip(paginationParams.NextPageKey)

	if paginationParams.PageSize > 0 {
		sortOptions = sortOptions.SetLimit(paginationParams.PageSize)
	}

	cur, err := collection.Find(ctx, searchOptions, sortOptions)
	defer closeCursor(ctx, cur)

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, nil, err
	}

	result := []models.SequenceExecution{}
	paginationResult := &models.PaginationResult{
		TotalCount: totalCount,
	}

	for cur.Next(ctx) {
		var outInterface interface{}
		err := cur.Decode(&outInterface)
		sequenceExecution, err := mdbrepo.transformBSONToSequenceExecution(outInterface)
		if err != nil {
			log.Errorf("Could not decode sequenceExecution: %v", err)
			continue
		}
		result = append(result, *sequenceExecution)
	}

	paginationResult.PageSize = int64(len(result))

	if paginationParams.PageSize > 0 && paginationParams.PageSize+paginationParams.NextPageKey < totalCount {
		paginationResult.NextPageKey = paginationParams.PageSize + paginationParams.NextPageKey
	}
	return result, paginationResult, nil
}

// GetByTriggeredID searches for a sequence execution with the given triggeredID.
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
	sequenceExecution, err := mdbrepo.transformBSONToSequenceExecution(outInterface)
	if err != nil {
		return nil, err
	}

	return sequenceExecution, nil
}

// Upsert inserts or updates a sequence execution into the sequence execution collection.
// By setting the CheckUniqueTriggeredID of the upsertOptions to true, this function will return a ErrSequenceWithTriggeredIDAlreadyExists,
// if a sequence with the same triggeredID already exists (can be useful to avoid storing duplicate sequences).
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

	var internalDBItem interface{}
	if mdbrepo.ModelTransformer != nil {
		internalDBItem = mdbrepo.ModelTransformer.TransformToDBModel(item)
	} else {
		internalDBItem = item
	}

	filter := bson.D{{"_id", item.ID}}

	if upsertOptions != nil && upsertOptions.Replace {
		_, err = collection.ReplaceOne(ctx, filter, internalDBItem)
	} else {
		update := bson.D{{"$set", internalDBItem}}
		_, err = collection.UpdateOne(ctx, filter, update, opts)
	}
	if err != nil {
		return err
	}
	return nil
}

// AppendTaskEvent adds an event that is relevant to the execution of the current task.
// This function needs to be thread safe since it can  potentially be invoked by multiple threads at the same time.
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

	// by using the $push operator in the FindOneAndUpdate function, we ensure that we follow an append-only approach to this property,
	// since this is the one property that can potentially be updated by multiple threads handling .finished/.started events for the same task
	var eventItem interface{}
	if mdbrepo.ModelTransformer != nil {
		// if we have a transformer that defines how the event should be stored internally, use it to transform the item
		eventItem = mdbrepo.ModelTransformer.TransformEventToDBModel(event)
	} else {
		eventItem = event
	}
	update := bson.M{"$push": bson.M{"status.currentTask.events": eventItem}}

	res := collection.FindOneAndUpdate(ctx, filter, update, opts)
	if res.Err() != nil {
		return nil, err
	}

	outInterface := map[string]interface{}{}
	err = res.Decode(outInterface)
	if err != nil {
		return nil, err
	}
	sequenceExecution, err := mdbrepo.transformBSONToSequenceExecution(outInterface)
	if err != nil {
		return nil, err
	}
	return sequenceExecution, nil
}

// UpdateStatus is used to update the overall state of the sequence, e.g. when it was paused via the API.
// This will not update a complete sequence execution, but just the attributes representing the overall state of the sequence
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
	sequenceExecution, err := mdbrepo.transformBSONToSequenceExecution(outInterface)
	if err != nil {
		return nil, err
	}
	return sequenceExecution, nil
}

// Clear deletes the sequence execution collection of the given project
func (mdbrepo *MongoDBSequenceExecutionRepo) Clear(projectName string) error {
	collection, ctx, cancel, err := mdbrepo.getSequenceExecutionStateCollection(projectName)
	if err != nil {
		return err
	}
	defer cancel()

	_, err = collection.DeleteMany(ctx, bson.D{})
	return err
}

// PauseContext pauses all sequence executions for the given Keptn Context
func (mdbrepo *MongoDBSequenceExecutionRepo) PauseContext(eventScope models.EventScope) error {
	return mdbrepo.updateGlobalSequenceContext(eventScope, apimodels.SequencePaused)
}

// ResumeContext resumes all sequence executions for the given Keptn Context
func (mdbrepo *MongoDBSequenceExecutionRepo) ResumeContext(eventScope models.EventScope) error {
	return mdbrepo.updateGlobalSequenceContext(eventScope, apimodels.SequenceStartedState)
}

// IsContextPaused checks whether a sequence that belongs to the given Keptn Context is currently paused
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
	defer closeCursor(ctx, cur)

	if err != nil {
		log.Errorf("Could not retrieve sequence context: %v", err)
		return false
	} else if cur.RemainingBatchLength() == 0 {
		return false
	}

	stateItems := []models.EventQueueSequenceState{}

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
		if state.Scope.Stage == "" && state.State == apimodels.SequencePaused {
			// if the overall state is set to 'paused', this means that all stages are paused
			return true
		} else if state.Scope.Stage == eventScope.Stage && state.State == apimodels.SequencePaused {
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
	if !filter.TriggeredAt.IsZero() {
		searchOptions["triggeredAt"] = bson.M{
			"$lt": filter.TriggeredAt,
		}
	}

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

func (mdbrepo *MongoDBSequenceExecutionRepo) transformBSONToSequenceExecution(outInterface interface{}) (*models.SequenceExecution, error) {
	outInterface, err := flattenRecursively(outInterface)
	if err != nil {
		return nil, err
	}

	// if we have a model transformer, use that one to transform the item into a model.SequenceExecution object
	if mdbrepo.ModelTransformer != nil {
		return mdbrepo.ModelTransformer.TransformToSequenceExecution(outInterface)
	}

	// if we have no model transformer, directly decode to models.SequenceExecution
	data, err := json.Marshal(outInterface)
	if err != nil {
		return nil, err
	}

	sequenceExecution := &models.SequenceExecution{}
	if err := json.Unmarshal(data, sequenceExecution); err != nil {
		return nil, err
	}

	return sequenceExecution, nil
}

func appendFilterAs(filter bson.M, value, key string) bson.M {
	if value != "" {
		filter[key] = value
	}
	return filter
}
