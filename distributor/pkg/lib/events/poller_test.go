package events

import (
	"context"
	"encoding/json"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Test_PollAndForwardEvents(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {}))

	type args struct {
		envConfig config.EnvConfig
	}
	tests := []struct {
		name              string
		args              args
		serverHandlerFunc http.HandlerFunc
		eventSender       EventSender
		uniformWatch      IUniformWatch
	}{
		{
			name: "poll multiple topics",
			args: args{envConfig: config.EnvConfig{
				KeptnAPIEndpoint:    ts.URL,
				PubSubRecipient:     "http://127.0.0.1",
				PubSubTopic:         "sh.keptn.event.task.triggered,sh.keptn.event.task2.triggered",
				HTTPPollingInterval: "1",
			}},

			serverHandlerFunc: func(w http.ResponseWriter, request *http.Request) {
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
							},
						},
						NextPageKey: "",
						PageSize:    1,
						TotalCount:  1,
					}
				}

				marshal, _ := json.Marshal(events)
				w.Write(marshal)
			},
			eventSender: &keptnv2.TestSender{},
			uniformWatch: NewTestUniformWatch([]keptnmodels.TopicSubscription{{
				ID:    "id1",
				Topic: "sh.keptn.event.task.triggered",
			},
				{
					ID:    "id2",
					Topic: "sh.keptn.event.task2.triggered",
				}}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts.Config.Handler = tt.serverHandlerFunc

			poller := NewPoller(tt.args.envConfig, tt.eventSender, &http.Client{}, tt.uniformWatch)

			ctx, cancel := context.WithCancel(context.Background())
			executionContext := NewExecutionContext(ctx, 1)
			go poller.Start(executionContext)

			assert.Eventually(t, func() bool {
				return len(tt.eventSender.(*keptnv2.TestSender).SentEvents) == 2
			}, time.Second*time.Duration(5), time.Second)

			cancel()
			executionContext.Wg.Wait()
		})
	}
}
