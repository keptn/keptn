package handlers

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/keptn/keptn/mongodb-datastore/models"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"
	"github.com/magiconair/properties/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

// TestFlattenRecursivelyNestedDocuments checks whether the flattening works with nested bson.D (documents)
func TestFlattenRecursivelyNestedDocuments(t *testing.T) {
	logger := keptncommon.NewLogger("", "", "mongodb-service")

	grandchild := bson.D{{"apple", "red"}, {"orange", "orange"}}
	child := bson.D{{"foo", "bar"}, {"grandchild", grandchild}}
	parent := bson.D{{"hello", "world"}, {"child", child}}

	// checks:
	flattened, _ := flattenRecursively(parent, logger)
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
	logger := keptncommon.NewLogger("", "", "mongodb-service")

	grandchild := bson.D{{"apple", "red"}, {"orange", "orange"}}
	child := bson.A{grandchild, "foo", "bar"}
	parent := bson.D{{"hello", "world"}, {"child", child}}

	// checks:
	flattened, _ := flattenRecursively(parent, logger)
	parentMap, _ := flattened.(map[string]interface{})
	assert.Equal(t, parentMap["hello"], "world", "flatting failed")

	childMap := parentMap["child"].(bson.A)
	assert.Equal(t, len(childMap), 3, "flatting failed")

	grandchildMap := childMap[0].(map[string]interface{})
	assert.Equal(t, grandchildMap["apple"], "red", "flatting failed")
}

func Test_getProjectOfEvent(t *testing.T) {
	type args struct {
		event *models.KeptnContextExtendedCE
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Use project property in data object",
			args: args{
				event: &models.KeptnContextExtendedCE{
					Event: models.Event{
						Contenttype: "",
						Data: map[string]interface{}{
							"project": "sockshop",
						},
						Extensions:  nil,
						ID:          "",
						Source:      "",
						Specversion: "",
						Time:        models.Time{},
						Type:        "",
					},
					Shkeptncontext: "",
				},
			},
			want: "sockshop",
		},
		{
			name: "Use generic events collection",
			args: args{
				event: &models.KeptnContextExtendedCE{
					Event: models.Event{
						Contenttype: "",
						Data:        nil,
						Extensions:  nil,
						ID:          "",
						Source:      "",
						Specversion: "",
						Time:        models.Time{},
						Type:        "",
					},
					Shkeptncontext: "",
				},
			},
			want: eventsCollectionName,
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
				"data.project": "sockshop",
				"data.stage":   "dev",
				"data.service": "carts",
				"source":       "test-service",
				"$or": []bson.M{
					{"type": keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName)},
					{"type": keptn07EvaluationDoneEventType},
				},
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
	type args struct {
		event *models.KeptnContextExtendedCE
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
				event: &models.KeptnContextExtendedCE{
					Event: models.Event{
						Contenttype: "application/json",
						Data:        "test-content",
						Extensions:  nil,
						ID:          "1",
						Source:      "test-source",
						Specversion: "0.2",
						Time:        models.Time{},
						Type:        "test-type",
					},
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
				"time":           "0001-01-01T00:00:00.000Z",
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
		name string
		args args
		want bool
	}{
		{
			name: "data.project provided",
			args: args{
				searchOptions: bson.M{
					"data.project": "test",
				},
			},
			want: true,
		},
		{
			name: "data.project empty string",
			args: args{
				searchOptions: bson.M{
					"data.project": "",
				},
			},
			want: false,
		},
		{
			name: "shkeptncontext provided",
			args: args{
				searchOptions: bson.M{
					"shkeptncontext": "test",
				},
			},
			want: true,
		},
		{
			name: "shkeptncontext empty string",
			args: args{
				searchOptions: bson.M{
					"shkeptncontext": "",
				},
			},
			want: false,
		},
		{
			name: "empty",
			args: args{
				searchOptions: bson.M{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateFilter(tt.args.searchOptions); got != tt.want {
				t.Errorf("validateFilter() = %v, want %v", got, tt.want)
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
					{"$match", bson.M{
						"project": "test-project",
					}},
				},
				bson.D{
					{"$lookup", bson.M{
						"from": "test-collection-invalidatedEvents",
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
												"$eq": []string{"$triggeredid", "$$event_id"},
											},
											{
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
				},
				bson.D{
					{"$match", bson.M{
						"invalidated": bson.M{
							"$size": 0,
						},
					}},
				},
				bson.D{
					{"$sort",
						bson.M{
							"time": -1,
						},
					},
				},
				bson.D{
					{
						"$limit", limit,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAggregationPipeline(tt.args.params, tt.args.collectionName, tt.args.matchFields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAggregationPipeline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getInvalidatedEventType(t *testing.T) {
	type args struct {
		params string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "arbitrary event type",
			args: args{
				params: "my-type",
			},
			want: "my-type.invalidated",
		},
		{
			name: "arbitrary event type",
			args: args{
				params: "my-type.foo.bar",
			},
			want: "my-type.foo.invalidated",
		},
		{
			name: "evaluation-done",
			args: args{
				params: keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName),
			},
			want: "sh.keptn.event.evaluation.invalidated",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getInvalidatedEventType(tt.args.params); got != tt.want {
				t.Errorf("getInvalidatedEventType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_transformEvaluationDonEvent(t *testing.T) {
	type args struct {
		keptnEvent *models.KeptnContextExtendedCE
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantEvent *models.KeptnContextExtendedCE
	}{
		{
			name: "transform evaluation-done event",
			args: args{
				keptnEvent: &models.KeptnContextExtendedCE{
					Event: models.Event{
						Contenttype: "",
						Data: map[string]interface{}{
							"result":  "pass",
							"project": "my-project",
							"stage":   "my-stage",
							"service": "my-service",
							"labels": map[string]interface{}{
								"foo": "bar",
							},
							"evaluationdetails": keptnv2.EvaluationDetails{
								Result: string(keptnv2.ResultPass),
								Score:  10,
							},
						},
						Extensions:  nil,
						ID:          "",
						Source:      "lighthouse-service",
						Specversion: "0.2",
						Time:        models.Time{},
						Type:        keptn07EvaluationDoneEventType,
					},
					Shkeptncontext: "my-context",
					Triggeredid:    "my-triggeredid",
				},
			},
			wantEvent: &models.KeptnContextExtendedCE{
				Event: models.Event{
					Contenttype: "",
					Data: &keptnv2.EvaluationFinishedEventData{
						EventData: keptnv2.EventData{
							Project: "my-project",
							Stage:   "my-stage",
							Service: "my-service",
							Labels: map[string]string{
								"foo": "bar",
							},
							Result: keptnv2.ResultPass,
						},
						Evaluation: keptnv2.EvaluationDetails{
							Result: string(keptnv2.ResultPass),
							Score:  10,
						},
					},
					Extensions:  nil,
					ID:          "",
					Source:      "lighthouse-service",
					Specversion: "1.0",
					Time:        models.Time{},
					Type:        models.Type(keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName)),
				},
				Shkeptncontext: "my-context",
				Triggeredid:    "my-triggeredid",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := transformEvaluationDonEvent(tt.args.keptnEvent); (err != nil) != tt.wantErr {
				t.Errorf("transformEvaluationDonEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := deep.Equal(tt.args.keptnEvent, tt.wantEvent); len(diff) > 0 {
				t.Errorf("transformEvaluationDonEvent() did not return expected event:")
				for _, d := range diff {
					t.Log(d)
				}
			}
		})
	}
}
