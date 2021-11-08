package db

import (
	"context"
	"github.com/keptn/keptn/mongodb-datastore/models"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

const (
	eventsCollectionName              = "keptnUnmappedEvents" // TODO do we need this still?
	invalidatedEventsCollectionSuffix = "-invalidatedEvents"
)

type MongoDBEventRepo struct {
	DBConnection *MongoDBConnection
}

func NewMongoDBEventRepo(dbConnection *MongoDBConnection) *MongoDBEventRepo {
	return &MongoDBEventRepo{DBConnection: dbConnection}
}

func (mr *MongoDBEventRepo) InsertEvent(event models.KeptnContextExtendedCE) error {
	collection, ctx, cancel, err := mr.getCollectionAndContext()
	if err != nil {
		return err
	}
	defer cancel()
	panic("implement me")

}

func (mr *MongoDBEventRepo) DropProjectCollections(project string) error {
	panic("implement me")
}

func (mr *MongoDBEventRepo) GetEvents(params event.GetEventParams) (event.GetEventsOKBody, error) {
	panic("implement me")
}

func (mr *MongoDBEventRepo) GetEventsByType(params event.GetEventsByTypeParams) (*event.GetEventsByTypeOKBody, error) {
	panic("implement me")
}

func (mr *MongoDBEventRepo) getCollectionAndContext(collectionName string) (*mongo.Collection, context.Context, context.CancelFunc, error) {
	err := mr.DBConnection.EnsureDBConnection()
	if err != nil {
		return nil, nil, nil, err
	}
	collection := mr.DBConnection.Client.Database(getDatabaseName()).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return collection, ctx, cancel, nil
}

func getCollectionNameForEvent(event *models.KeptnContextExtendedCE) string {
	collectionName := eventsCollectionName
	// check if the data object contains the project name.
	// if yes, store the event in the collection for the project, otherwise in /events
	eventData, ok := event.Data.(map[string]interface{})
	if ok && eventData["project"] != nil {
		collectionNameStr, ok := eventData["project"].(string)
		if ok && collectionNameStr != "" {
			collectionName = collectionNameStr
		}
	}

	return collectionName
}

func getAggregationPipeline(params event.GetEventsByTypeParams, collectionName string, matchFields bson.M) mongo.Pipeline {
	// TODO: find better name for this function
	matchStage := bson.D{
		{Key: "$match", Value: matchFields},
	}

	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.M{
			"from": getInvalidatedCollectionName(collectionName),
			"let": bson.M{
				"event_id":          "$id",
				"event_triggeredid": "$triggeredid",
			},
			"pipeline": []bson.M{
				{
					"$match": bson.M{
						"$expr": bson.M{
							"$or": []bson.M{
								{
									// backwards-compatibility to 0.7.x -> triggeredid of .invalidated event refers to the id of the evaluation-done event
									"$eq": []string{"$triggeredid", "$$event_id"},
								},
								{
									// logic for 0.8: triggeredid of .invalidated event refers to the triggeredid of the evaluation.finished event (both are related to the same .triggered event)
									"$eq": []string{"$triggeredid", "$$event_triggeredid"},
								},
							},
						},
					},
				},
				{
					"$limit": 1,
				},
			},
			"as": "invalidated",
		}},
	}

	matchInvalidatedStage := bson.D{
		{Key: "$match", Value: bson.M{
			"invalidated": bson.M{
				"$size": 0,
			},
		}},
	}
	sortStage := bson.D{
		{Key: "$sort", Value: bson.M{
			"time": -1,
		}},
	}
	var aggregationPipeline mongo.Pipeline
	if params.Limit != nil && *params.Limit > 0 {
		limitStage := bson.D{
			{Key: "$limit", Value: *params.Limit},
		}
		aggregationPipeline = mongo.Pipeline{matchStage, lookupStage, matchInvalidatedStage, sortStage, limitStage}
	} else {
		aggregationPipeline = mongo.Pipeline{matchStage, lookupStage, matchInvalidatedStage, sortStage}
	}

	return aggregationPipeline
}

func parseFilter(filter string) bson.M {
	filterObject := bson.M{}
	keyValues := strings.Split(filter, " AND ")

	for _, keyValuePair := range keyValues {
		split := strings.Split(keyValuePair, ":")
		if len(split) == 2 {
			splitValue := strings.Split(split[1], ",")
			if len(splitValue) == 1 {
				filterObject[split[0]] = split[1]
			} else {
				filterObject[split[0]] = bson.M{
					"$in": splitValue,
				}
			}
		}
	}

	return filterObject
}

func validateFilter(searchOptions bson.M) bool {
	if (searchOptions["data.project"] == nil || searchOptions["data.project"] == "") && (searchOptions["shkeptncontext"] == nil || searchOptions["shkeptncontext"] == "") {
		return false
	}

	return true
}

func getInvalidatedCollectionName(collectionName string) string {
	invalidatedCollectionName := collectionName + invalidatedEventsCollectionSuffix
	return invalidatedCollectionName
}
