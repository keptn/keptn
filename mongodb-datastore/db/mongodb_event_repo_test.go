package db

import (
	"context"
	"fmt"
	keptnapi "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/mongodb-datastore/common"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"reflect"
	"testing"
	"time"
)

var mongoDbVersion = "4.4.9"

func TestMain(m *testing.M) {
	mongoServer, err := setupLocalMongoDB()
	if err != nil {
		log.Fatalf("Mongo Server setup failed: %s", err)
	}
	defer mongoServer.Stop()
	m.Run()
}

func setupLocalMongoDB() (*memongo.Server, error) {
	mongoServer, err := memongo.Start(mongoDbVersion)
	if err != nil {
		return nil, err
	}

	randomDbName := memongo.RandomDatabase()

	_ = os.Setenv("MONGODB_DATABASE", randomDbName)
	_ = os.Setenv("MONGODB_EXTERNAL_CONNECTION_STRING", fmt.Sprintf("%s/%s", mongoServer.URI(), randomDbName))

	var mongoClient *mongo.Client
	mongoClient, err = mongo.NewClient(options.Client().ApplyURI(mongoServer.URI()))
	if err != nil {
		return nil, err
	}
	err = mongoClient.Connect(context.TODO())
	if err != nil {
		return nil, err
	}

	return mongoServer, err
}

func TestMongoDBEventRepo_InsertAndRetrieve(t *testing.T) {
	repo := NewMongoDBEventRepo(GetMongoDBConnectionInstance())
	time := time.Time{}
	filter := "data.project:my-project"
	pageSize := int64(0)
	events, err := repo.GetEventsByType(
		event.GetEventsByTypeParams{
			EventType: "test",
			Filter:    filter,
			Limit:     &pageSize,
		},
	)
	require.Nil(t, err)
	require.Empty(t, events.Events)

	evaluationEventType := keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName)

	keptnContext := "my-context"
	testEvent := keptnapi.KeptnContextExtendedCE{
		Contenttype:        "application/cloudevents+json",
		Data:               map[string]interface{}{"project": "my-project", "service": "my-service", "stage": "my-stage"},
		ID:                 "my-evaluation-id",
		Source:             stringp("test-source"),
		Specversion:        "1.0",
		Time:               time,
		Type:               stringp(evaluationEventType),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: "0.2.3",
		Triggeredid:        "my-triggered-id",
	}

	err = repo.InsertEvent(testEvent)
	require.Nil(t, err)

	invalidatedEvent := keptnapi.KeptnContextExtendedCE{
		Contenttype:        "application/cloudevents+json",
		Data:               map[string]interface{}{"project": "my-project", "service": "my-service", "stage": "my-stage"},
		ID:                 "my-invalidated-id",
		Source:             stringp("test-source"),
		Specversion:        "1.0",
		Time:               time,
		Type:               stringp(keptnv2.GetInvalidatedEventType(keptnv2.EvaluationTaskName)),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: "0.2.3",
		Triggeredid:        "my-triggered-id",
	}

	err = repo.InsertEvent(invalidatedEvent)

	require.Nil(t, err)

	events, err = repo.GetEvents(
		event.GetEventsParams{
			KeptnContext: &keptnContext,
			PageSize:     &pageSize,
			Type:         &evaluationEventType,
		},
	)

	require.Nil(t, err)
	require.NotNil(t, events)
	require.Len(t, events.Events, 1)
	require.Equal(t, testEvent, events.Events[0])

	filter = "data.project:my-project"
	eventsByType, err := repo.GetEventsByType(
		event.GetEventsByTypeParams{
			EventType: evaluationEventType,
			Filter:    filter,
			Limit:     &pageSize,
		},
	)

	require.Nil(t, err)
	require.NotNil(t, eventsByType)
	require.Len(t, eventsByType.Events, 1)

	excludeInvalidated := true
	// now try to query again with excluded invalidated events
	eventsByType, err = repo.GetEventsByType(
		event.GetEventsByTypeParams{
			EventType:          evaluationEventType,
			ExcludeInvalidated: &excludeInvalidated,
			Filter:             filter,
			Limit:              &pageSize,
		},
	)

	require.Nil(t, err)
	require.NotNil(t, eventsByType)
	require.Empty(t, eventsByType.Events)
}

func TestMongoDBEventRepo_Retrieve_NoProjectOrKeptnContext(t *testing.T) {
	repo := NewMongoDBEventRepo(GetMongoDBConnectionInstance())

	pageSize := int64(0)
	filter := ""
	eventsByType, err := repo.GetEventsByType(
		event.GetEventsByTypeParams{
			EventType: keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName),
			Filter:    filter,
			Limit:     &pageSize,
		},
	)

	require.NotNil(t, err)

	require.ErrorIs(t, err, common.ErrInvalidEventFilter)
	require.Nil(t, eventsByType)
}

func TestMongoDBEventRepo_DropCollections(t *testing.T) {
	repo := NewMongoDBEventRepo(GetMongoDBConnectionInstance())

	evaluationEventType := keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName)

	keptnContext := "my-context"
	testEvent := keptnapi.KeptnContextExtendedCE{

		Contenttype:        "application/cloudevents+json",
		Data:               map[string]interface{}{"project": "my-project", "service": "my-service", "stage": "my-stage"},
		ID:                 "my-evaluation-id",
		Source:             stringp("test-source"),
		Specversion:        "1.0",
		Time:               time.Time{},
		Type:               stringp(evaluationEventType),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: "0.2.3",
		Triggeredid:        "my-triggered-id",
	}

	err := repo.InsertEvent(testEvent)
	require.Nil(t, err)

	invalidatedEvent := keptnapi.KeptnContextExtendedCE{

		Contenttype:        "application/cloudevents+json",
		Data:               map[string]interface{}{"project": "my-project", "service": "my-service", "stage": "my-stage"},
		ID:                 "my-invalidated-id",
		Source:             stringp("test-source"),
		Specversion:        "1.0",
		Time:               time.Time{},
		Type:               stringp(keptnv2.GetInvalidatedEventType(keptnv2.EvaluationTaskName)),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: "0.2.3",
		Triggeredid:        "my-evaluation-id",
	}

	err = repo.InsertEvent(invalidatedEvent)

	require.Nil(t, err)

	err = repo.DropProjectCollections(testEvent)

	require.Nil(t, err)

	pageSize := int64(0)
	events, err := repo.GetEvents(
		event.GetEventsParams{
			KeptnContext: &keptnContext,
			PageSize:     &pageSize,
			Type:         &evaluationEventType,
		},
	)

	require.Nil(t, err)
	require.Empty(t, events.Events)
}

// TestFlattenRecursivelyNestedDocuments checks whether the flattening works with nested bson.D (documents)
func TestFlattenRecursivelyNestedDocuments(t *testing.T) {
	grandchild := bson.D{{Key: "apple", Value: "red"}, {Key: "orange", Value: "orange"}}
	child := bson.D{{Key: "foo", Value: "bar"}, {Key: "grandchild", Value: grandchild}}
	parent := bson.D{{Key: "hello", Value: "world"}, {Key: "child", Value: child}}

	// checks:
	flattened, _ := flattenRecursively(parent)
	parentMap, _ := flattened.(map[string]interface{})
	assert.Equal(t, parentMap["hello"], "world", "flatting failed")

	childMap := parentMap["child"].(map[string]interface{})
	assert.Equal(t, childMap["foo"], "bar", "flatting failed")

	grandchildMap := childMap["grandchild"].(map[string]interface{})
	assert.Equal(t, grandchildMap["orange"], "orange", "flatting failed")
}

// TestFlattenRecursivelyNestedDocuments checks whether the flattening works with nested bson.D (documents)
// and bson.A (arrays)
func TestFlattenRecursivelyNestedDocumentsWithArray(t *testing.T) {
	grandchild := bson.D{{Key: "apple", Value: "red"}, {Key: "orange", Value: "orange"}}
	child := bson.A{grandchild, "foo", "bar"}
	parent := bson.D{{Key: "hello", Value: "world"}, {Key: "child", Value: child}}

	// checks:
	flattened, _ := flattenRecursively(parent)
	parentMap, _ := flattened.(map[string]interface{})
	assert.Equal(t, parentMap["hello"], "world", "flatting failed")

	childMap := parentMap["child"].(bson.A)
	assert.Equal(t, len(childMap), 3, "flatting failed")

	grandchildMap := childMap[0].(map[string]interface{})
	assert.Equal(t, grandchildMap["apple"], "red", "flatting failed")
}

func Test_getProjectOfEvent(t *testing.T) {
	type args struct {
		event keptnapi.KeptnContextExtendedCE
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Use project property in data object",
			args: args{
				event: keptnapi.KeptnContextExtendedCE{
					Contenttype: "",
					Data: map[string]interface{}{
						"project": "sockshop",
					},
					Extensions:     nil,
					ID:             "",
					Source:         nil,
					Specversion:    "",
					Time:           time.Time{},
					Type:           nil,
					Shkeptncontext: "",
				},
			},
			want: "sockshop",
		},
		{
			name: "Use generic events collection",
			args: args{
				event: keptnapi.KeptnContextExtendedCE{

					Contenttype:    "",
					Data:           nil,
					Extensions:     nil,
					ID:             "",
					Source:         nil,
					Specversion:    "",
					Time:           time.Time{},
					Type:           nil,
					Shkeptncontext: "",
				},
			},
			want: unmappedEventsCollectionName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getProjectOfEvent(tt.args.event); got != tt.want {
				t.Errorf("getProjectOfEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSearchOptions(t *testing.T) {
	type args struct {
		params event.GetEventsParams
	}
	tests := []struct {
		name string
		args args
		want bson.M
	}{
		{
			name: "get search options",
			args: args{
				params: event.GetEventsParams{
					HTTPRequest:  nil,
					FromTime:     stringp("1"),
					BeforeTime:   nil,
					KeptnContext: stringp("test-context"),
					NextPageKey:  nil,
					PageSize:     nil,
					Project:      stringp("sockshop"),
					Root:         nil,
					Service:      stringp("carts"),
					Source:       stringp("test-service"),
					Stage:        stringp("dev"),
					Type:         stringp("test-event"),
				},
			},
			want: bson.M{
				"data.project":   "sockshop",
				"data.stage":     "dev",
				"data.service":   "carts",
				"source":         "test-service",
				"type":           "test-event",
				"shkeptncontext": "test-context",
				"time": bson.M{
					"$gt": "1",
				},
			},
		},
		{
			name: "get search options for evaluation.finished events query",
			args: args{
				params: event.GetEventsParams{
					HTTPRequest:  nil,
					FromTime:     stringp("1"),
					BeforeTime:   nil,
					KeptnContext: stringp("test-context"),
					NextPageKey:  nil,
					PageSize:     nil,
					Project:      stringp("sockshop"),
					Root:         nil,
					Service:      stringp("carts"),
					Source:       stringp("test-service"),
					Stage:        stringp("dev"),
					Type:         stringp(keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName)),
				},
			},
			want: bson.M{
				"data.project":   "sockshop",
				"data.stage":     "dev",
				"data.service":   "carts",
				"source":         "test-service",
				"type":           keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName),
				"shkeptncontext": "test-context",
				"time": bson.M{
					"$gt": "1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSearchOptions(tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSearchOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func stringp(s string) *string {
	return &s
}

func Test_transformEventToInterface(t *testing.T) {
	time := time.Time{}
	type args struct {
		event *keptnapi.KeptnContextExtendedCE
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "transform event",
			args: args{
				event: &keptnapi.KeptnContextExtendedCE{

					Contenttype:    "application/json",
					Data:           "test-content",
					Extensions:     nil,
					ID:             "1",
					Source:         stringp("test-source"),
					Specversion:    "0.2",
					Time:           time,
					Type:           stringp("test-type"),
					Shkeptncontext: "123",
				},
			},
			want: map[string]interface{}{
				"contenttype":    "application/json",
				"data":           "test-content",
				"id":             "1",
				"shkeptncontext": "123",
				"source":         "test-source",
				"specversion":    "0.2",
				"time":           "0001-01-01T00:00:00Z",
				"type":           "test-type",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformEventToInterface(tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("transformEventToInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("transformEventToInterface() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseFilter(t *testing.T) {
	type args struct {
		filter string
	}
	tests := []struct {
		name string
		args args
		want bson.M
	}{
		{
			name: "get key values",
			args: args{
				filter: "data.project:sockshop AND shkeptncontext:test-context",
			},
			want: bson.M{
				"data.project":   "sockshop",
				"shkeptncontext": "test-context",
			},
		},
		{
			name: "get key values",
			args: args{
				filter: "data.project:sockshop AND data.result:pass,warn",
			},
			want: bson.M{
				"data.project": "sockshop",
				"data.result": bson.M{
					"$in": []string{"pass", "warn"},
				},
			},
		},
		{
			name: "empty input",
			args: args{
				filter: "",
			},
			want: bson.M{},
		},
		{
			name: "nonsense input",
			args: args{
				filter: "bla",
			},
			want: bson.M{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseFilter(tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateFilter(t *testing.T) {
	type args struct {
		searchOptions bson.M
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "data.project provided",
			args: args{
				searchOptions: bson.M{
					"data.project": "test",
				},
			},
			wantErr: false,
		},
		{
			name: "data.project empty string",
			args: args{
				searchOptions: bson.M{
					"data.project": "",
				},
			},
			wantErr: true,
		},
		{
			name: "shkeptncontext provided",
			args: args{
				searchOptions: bson.M{
					"shkeptncontext": "test",
				},
			},
			wantErr: false,
		},
		{
			name: "shkeptncontext empty string",
			args: args{
				searchOptions: bson.M{
					"shkeptncontext": "",
				},
			},
			wantErr: true,
		},
		{
			name: "empty",
			args: args{
				searchOptions: bson.M{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFilter(tt.args.searchOptions)

			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func Test_getAggregationPipeline(t *testing.T) {
	limit := int64(2)
	type args struct {
		params         event.GetEventsByTypeParams
		collectionName string
		matchFields    bson.M
	}
	tests := []struct {
		name string
		args args
		want mongo.Pipeline
	}{
		{
			name: "",
			args: args{
				params: event.GetEventsByTypeParams{
					Limit:     &limit,
					EventType: "my-type",
				},
				collectionName: "test-collection",
				matchFields: bson.M{
					"project": "test-project",
				},
			},
			want: mongo.Pipeline{
				bson.D{
					{Key: "$match", Value: bson.M{
						"project": "test-project",
					}},
				},
				bson.D{
					{Key: "$lookup", Value: bson.M{
						"from":         "test-collection-invalidatedEvents",
						"localField":   "triggeredid",
						"foreignField": "triggeredid",
						"as":           "invalidated",
					}},
				},
				bson.D{
					{Key: "$match", Value: bson.M{
						"invalidated": bson.M{
							"$size": 0,
						},
					}},
				},
				bson.D{
					{Key: "$sort",
						Value: bson.M{
							"time": -1,
						},
					},
				},
				bson.D{
					{
						Key: "$limit", Value: limit,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getInvalidatedEventQuery(tt.args.params, tt.args.collectionName, tt.args.matchFields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getInvalidatedEventQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
