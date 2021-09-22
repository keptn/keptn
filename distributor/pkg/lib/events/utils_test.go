package events

import (
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestAddEvent(t *testing.T) {
	cache := NewCloudEventsCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")
	cache.Add("t2", "e3")

	assert.True(t, cache.Contains("t1", "e1"))
	assert.True(t, cache.Contains("t1", "e2"))
	assert.False(t, cache.Contains("t1", "e3"))
	assert.True(t, cache.Contains("t2", "e3"))
}

func TestAddEventTwice(t *testing.T) {
	cache := NewCloudEventsCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")
	cache.Add("t1", "e2")
	assert.Equal(t, 2, cache.Length("t1"))
	assert.Equal(t, 2, len(cache.Get("t1")))
}

func TestAddRemoveEvent(t *testing.T) {
	cache := NewCloudEventsCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")
	cache.Add("t1", "e3")

	assert.Equal(t, 3, cache.Length("t1"))

	cache.Remove("t1", "e1")
	assert.Equal(t, 2, cache.Length("t1"))
	assert.True(t, cache.Contains("t1", "e2"))
	assert.True(t, cache.Contains("t1", "e3"))

	cache.Remove("t1", "e3")
	assert.Equal(t, 1, cache.Length("t1"))
	assert.True(t, cache.Contains("t1", "e2"))
}

func TestKeep(t *testing.T) {
	cache := NewCloudEventsCache()
	cache.Add("t1", "e1")
	cache.Add("t1", "e2")
	cache.Add("t2", "e3")
	cache.Add("t2", "e4")
	cache.Add("t2", "e5")

	cache.Keep("t1", []*models.KeptnContextExtendedCE{ce("e2")})
	cache.Keep("t2", []*models.KeptnContextExtendedCE{ce("e3"), ce("e5")})

	assert.Equal(t, 1, cache.Length("t1"))
	assert.Equal(t, 2, cache.Length("t2"))
	assert.False(t, cache.Contains("t1", "e1"))
	assert.True(t, cache.Contains("t1", "e2"))
	assert.True(t, cache.Contains("t2", "e3"))
	assert.False(t, cache.Contains("t2", "e4"))
	assert.True(t, cache.Contains("t2", "e5"))
}

func Test_decodeCloudEvent(t *testing.T) {
	type args struct {
		data []byte
	}
	var tests = []struct {
		name    string
		args    args
		want    *cloudevents.Event
		wantErr bool
	}{
		{
			name: "Get V1.0 CloudEvent",
			args: args{
				data: []byte(`{
				"data": "",
				"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
				"source": "helm-service",
				"specversion": "1.0",
				"type": "sh.keptn.events.deployment-finished",
				"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
			}`),
			},
			want:    getExpectedCloudEvent(),
			wantErr: false,
		},
		{
			name: "Get V1.0 CloudEvent",
			args: args{
				data: []byte(""),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeCloudEvent(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeCloudEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !reflect.DeepEqual(got.Context.GetSpecVersion(), tt.want.Context.GetSpecVersion()) {
					t.Errorf("decodeCloudEvent() specVersion: got = %v, want1 %v", got.Context.GetSpecVersion(), tt.want.Context.GetSpecVersion())
				}
				if !reflect.DeepEqual(got.Context.GetType(), tt.want.Context.GetType()) {
					t.Errorf("decodeCloudEvent() type: got = %v, want1 %v", got.Context.GetType(), tt.want.Context.GetType())
				}
			}
		})
	}
}

func Test_matchesFilter(t *testing.T) {
	type args struct {
		e cloudevents.Event
	}
	tests := []struct {
		name         string
		args         args
		eventMatcher EventMatcher
		want         bool
	}{
		{
			name: "no filter - should match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			eventMatcher: EventMatcher{
				Project: "",
				Stage:   "",
				Service: "",
			},
			want: true,
		},
		{
			name: "project filter - should match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			eventMatcher: EventMatcher{
				Project: "my-project",
			},
			want: true,
		},
		{
			name: "project filter (comma-separated list) - should match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			eventMatcher: EventMatcher{
				Project: "my-project,my-project-2,my-project-3",
			},
			want: true,
		},
		{
			name: "project filter - should not match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			eventMatcher: EventMatcher{
				Project: "my-other-project",
			},
			want: false,
		},
		{
			name: "project filter (comma-separated list) - should not match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			eventMatcher: EventMatcher{
				Project: "my-other-project,my-second-project",
			},
			want: false,
		},
		{
			name: "stage filter - should match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			eventMatcher: EventMatcher{
				Stage: "my-stage",
			},
			want: true,
		},
		{
			name: "stage filter (comma-separated list) - should match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			eventMatcher: EventMatcher{
				Stage: "my-first-stage,my-stage",
			},
			want: true,
		},
		{
			name: "stage filter - should not match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			eventMatcher: EventMatcher{
				Stage: "my-other-stage",
			},
			want: false,
		},
		{
			name: "service filter - should match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			eventMatcher: EventMatcher{
				Project: "",
				Stage:   "",
				Service: "my-service",
			},
			want: true,
		},
		{
			name: "service filter (comma-separated list) - should match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			eventMatcher: EventMatcher{
				Service: "my-other-service,my-service",
			},
			want: true,
		},
		{
			name: "service filter - should not match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			eventMatcher: EventMatcher{
				Service: "my-other-service",
			},
			want: false,
		},
		{
			name: "combined filter - should match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			eventMatcher: EventMatcher{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			want: true,
		},
		{
			name: "combined filter - should not match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			eventMatcher: EventMatcher{
				Project: "my-other-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			want: false,
		},
		{
			name: "combined filter (comma-separated list) - should not match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			eventMatcher: EventMatcher{
				Project: "my-project,project-1",
				Stage:   "my-stage",
				Service: "my-service",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.eventMatcher.Matches(tt.args.e); got != tt.want {
				t.Errorf("matchesFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddAdditionalEventData(t *testing.T) {
	data := keptnv2.EvaluationFinishedEventData{
		EventData: keptnv2.EventData{Project: "project", Service: "service", Stage: "stage"},
	}

	event, _ := keptnv2.KeptnEvent("sh.keptn.event.dev.delivery.triggered", "source", data).Build()
	AddAdditionalEventData(&event, AdditionalEventData{"something": "additional"})

	dataAsMap := map[string]interface{}{}
	marshal, _ := json.Marshal(event.Data)
	json.Unmarshal(marshal, &dataAsMap)

	assert.Equal(t, map[string]interface{}{"something": "additional"}, dataAsMap["additionalData"])
}

func getCloudEventWithEventData(eventData keptnv2.EventData) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSource("helm-service")
	event.SetType("sh.keptn.events.deployment-finished")
	event.SetID("6de83495-4f83-481c-8dbe-fcceb2e0243b")
	event.SetExtension("shkeptncontext", "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb")
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetData(cloudevents.ApplicationJSON, eventData)
	return event
}

func ce(id string) *models.KeptnContextExtendedCE {
	return &models.KeptnContextExtendedCE{
		ID: id,
	}
}

func getExpectedCloudEvent() *cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSource("helm-service")
	event.SetType("sh.keptn.events.deployment-finished")
	event.SetID("6de83495-4f83-481c-8dbe-fcceb2e0243b")
	event.SetExtension("shkeptncontext", "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb")
	event.SetData(cloudevents.TextPlain, `""`)
	return &event
}
