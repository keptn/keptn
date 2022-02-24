package events

import (
	"context"
	"encoding/json"
	"fmt"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/common/strutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	keptnfake "github.com/keptn/go-utils/pkg/lib/v0_2_0/fake"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_PollAndForwardEvents1(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		eventType := strings.TrimPrefix(request.URL.Path, "/controlPlane/v1/event/triggered/")
		var events keptnmodels.Events
		if eventType == "sh.keptn.event.task.triggered" {
			events = keptnmodels.Events{
				Events: []*keptnmodels.KeptnContextExtendedCE{
					{
						ID:          "id-1",
						Type:        strutils.Stringp("sh.keptn.event.task.triggered"),
						Source:      strutils.Stringp("source"),
						Specversion: "1.0",
						Data:        map[string]interface{}{"some": "property"},
					},
				},
				NextPageKey: "",
				PageSize:    1,
				TotalCount:  1,
			}
		}
		if eventType == "sh.keptn.event.task2.triggered" {
			events = keptnmodels.Events{
				Events: []*keptnmodels.KeptnContextExtendedCE{
					{
						ID:          "id-2",
						Type:        strutils.Stringp("sh.keptn.event.task2.triggered"),
						Source:      strutils.Stringp("source"),
						Specversion: "1.0",
						Data:        map[string]interface{}{"some": "property"},
					},
				},
				NextPageKey: "",
				PageSize:    1,
				TotalCount:  1,
			}
		}

		marshal, _ := json.Marshal(events)
		w.Write(marshal)
	}))

	envConfig := config.EnvConfig{
		KeptnAPIEndpoint:    server.URL,
		PubSubRecipient:     "http://127.0.0.1",
		PubSubTopic:         "sh.keptn.event.task.triggered,sh.keptn.event.task2.triggered",
		HTTPPollingInterval: "1",
	}
	eventSender := keptnfake.EventSender{}
	apiset, _ := keptnapi.New(server.URL)
	poller := NewPoller(envConfig, apiset.ShipyardControlV1(), &eventSender)

	ctx, cancel := context.WithCancel(context.Background())
	executionContext := NewExecutionContext(ctx, 1)
	go poller.Start(executionContext)

	poller.UpdateSubscriptions([]keptnmodels.EventSubscription{
		{
			ID:    "id1",
			Event: "sh.keptn.event.task.triggered",
		},
		{
			ID:    "id2",
			Event: "sh.keptn.event.task2.triggered",
		},
		{
			ID:    "id3",
			Event: "sh.keptn.event.task2.triggered",
		},
	})

	assert.Eventually(t, func() bool {
		if len(eventSender.SentEvents) != 3 {
			return false
		}
		firstSentEvent := eventSender.SentEvents[0]
		event1, _ := keptnv2.ToKeptnEvent(firstSentEvent)
		var event1TmpData map[string]interface{}
		event1.GetTemporaryData("distributor", &event1TmpData)
		subscriptionIDInFirstEvent := event1TmpData["subscriptionID"]

		secondSentEvent := eventSender.SentEvents[1]
		event, _ := keptnv2.ToKeptnEvent(secondSentEvent)
		var event2TmpData map[string]interface{}
		event.GetTemporaryData("distributor", &event2TmpData)
		subscriptionIDInSecondEvent := event2TmpData["subscriptionID"]

		thirdSentEvent := eventSender.SentEvents[2]
		event3, _ := keptnv2.ToKeptnEvent(thirdSentEvent)
		var event3TmpData map[string]interface{}
		event3.GetTemporaryData("distributor", &event3TmpData)
		subscriptionIDInThirdEvent := event3TmpData["subscriptionID"]

		return subscriptionIDInFirstEvent == "id1" && subscriptionIDInSecondEvent == "id2" && subscriptionIDInThirdEvent == "id3"
	}, time.Second*time.Duration(5), time.Second)
	cancel()
	executionContext.Wg.Wait()
}

func Test_PollAndForwardEvents2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		eventType := strings.TrimPrefix(request.URL.Path, "/controlPlane/v1/event/triggered/")
		var events keptnmodels.Events
		if eventType == "sh.keptn.event.task.triggered" {
			events = keptnmodels.Events{
				Events: []*keptnmodels.KeptnContextExtendedCE{
					{
						ID:          "id-1",
						Type:        strutils.Stringp("sh.keptn.event.task.triggered"),
						Source:      strutils.Stringp("source"),
						Specversion: "1.0",
						Data:        map[string]interface{}{"some": "property"},
					},
				},
				NextPageKey: "",
				PageSize:    1,
				TotalCount:  1,
			}
		}

		marshal, _ := json.Marshal(events)
		w.Write(marshal)
	}))

	envConfig := config.EnvConfig{
		KeptnAPIEndpoint:    server.URL,
		PubSubRecipient:     "http://127.0.0.1",
		PubSubTopic:         "sh.keptn.event.task.triggered,sh.keptn.event.task2.triggered",
		HTTPPollingInterval: "1",
	}
	eventSender := keptnfake.EventSender{}

	apiset, _ := keptnapi.New(server.URL)
	poller := NewPoller(envConfig, apiset.ShipyardControlV1(), &eventSender)

	ctx, cancel := context.WithCancel(context.Background())
	executionContext := NewExecutionContext(ctx, 1)
	go poller.Start(executionContext)

	numSubscriptions := 100
	subscriptions := []keptnmodels.EventSubscription{}
	for i := 0; i < numSubscriptions; i++ {
		subscriptions = append(subscriptions, keptnmodels.EventSubscription{
			ID:    fmt.Sprintf("id%d", i),
			Event: "sh.keptn.event.task.triggered",
		})
	}

	poller.UpdateSubscriptions(subscriptions)

	assert.Eventually(t, func() bool {
		if len(eventSender.SentEvents) != numSubscriptions {
			return false
		}
		return true
	}, time.Second*time.Duration(5), time.Second)
	cancel()
	executionContext.Wg.Wait()
}

func Test_getEventFilterForSubscription(t *testing.T) {
	type args struct {
		subscription keptnmodels.EventSubscription
	}
	tests := []struct {
		name string
		args args
		want keptnapi.EventFilter
	}{
		{
			name: "get default filter",
			args: args{
				subscription: keptnmodels.EventSubscription{
					Event: "my-event",
				},
			},
			want: keptnapi.EventFilter{
				EventType: "my-event",
			},
		},
		{
			name: "multiple projects - get default filter",
			args: args{
				subscription: keptnmodels.EventSubscription{
					Event: "my-event",
					Filter: keptnmodels.EventSubscriptionFilter{
						Projects: []string{"a", "b"},
					},
				},
			},
			want: keptnapi.EventFilter{
				EventType: "my-event",
			},
		},
		{
			name: "one project",
			args: args{
				subscription: keptnmodels.EventSubscription{
					Event: "my-event",
					Filter: keptnmodels.EventSubscriptionFilter{
						Projects: []string{"a"},
					},
				},
			},
			want: keptnapi.EventFilter{
				EventType: "my-event",
				Project:   "a",
			},
		},
		{
			name: "one project, one stage",
			args: args{
				subscription: keptnmodels.EventSubscription{
					Event: "my-event",
					Filter: keptnmodels.EventSubscriptionFilter{
						Projects: []string{"a"},
						Stages:   []string{"stage-a"},
					},
				},
			},
			want: keptnapi.EventFilter{
				EventType: "my-event",
				Project:   "a",
				Stage:     "stage-a",
			},
		},
		{
			name: "one project, one stage, one service",
			args: args{
				subscription: keptnmodels.EventSubscription{
					Event: "my-event",
					Filter: keptnmodels.EventSubscriptionFilter{
						Projects: []string{"a"},
						Stages:   []string{"stage-a"},
						Services: []string{"service-a"},
					},
				},
			},
			want: keptnapi.EventFilter{
				EventType: "my-event",
				Project:   "a",
				Stage:     "stage-a",
				Service:   "service-a",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEventFilterForSubscription(tt.args.subscription); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEventFilterForSubscription() = %v, want %v", got, tt.want)
			}
		})
	}
}
