package handlers

import (
	"github.com/keptn/keptn/mongodb-datastore/models"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"testing"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/magiconair/properties/assert"
)

// TestFlattenRecursivelyNestedDocuments checks whether the flattening works with nested bson.D (documents)
func TestFlattenRecursivelyNestedDocuments(t *testing.T) {
	logger := keptnutils.NewLogger("", "", "mongodb-service")

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
	logger := keptnutils.NewLogger("", "", "mongodb-service")

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
			want: "events",
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
				"shkeptncontext": primitive.Regex{Pattern: "test-context", Options: ""},
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
