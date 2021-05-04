package handler_test

import (
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/remediation-service/handler"
	"github.com/keptn/keptn/remediation-service/internal/sdk"
	"github.com/keptn/keptn/remediation-service/internal/sdk/fake"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"testing"
)

func newGetActionTriggeredEvent(filename string) cloudevents.Event {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	event := models.KeptnContextExtendedCE{}
	err = json.Unmarshal(content, &event)
	_ = err
	return keptnv2.ToCloudEvent(event)
}

func Test_Receiving_GetActionTriggeredEvent_RemedöiationFromServiceLevel(t *testing.T) {

	fakeKeptn := fake.NewFakeKeptn("test-remediation-svc", sdk.WithHandler(handler.NewGetActionEventHandler(), "sh.keptn.event.get-action.triggered"))
	fakeKeptn.Start()
	fakeKeptn.NewEvent(newGetActionTriggeredEvent("test/events/get-action.triggered-0.json"))
	fakeKeptn.NewEvent(newGetActionTriggeredEvent("test/events/get-action.triggered-1.json"))

	require.Equal(t, 4, len(fakeKeptn.GetEventSender().SentEvents))

	finishedEvent, _ := keptnv2.ToKeptnEvent(fakeKeptn.GetEventSender().SentEvents[1])
	require.Equal(t, "sh.keptn.event.get-action.started", fakeKeptn.GetEventSender().SentEvents[0].Type())
	require.Equal(t, "sh.keptn.event.get-action.finished", fakeKeptn.GetEventSender().SentEvents[1].Type())
	getActionFinishedData := keptnv2.GetActionFinishedEventData{}
	finishedEvent.DataAs(&getActionFinishedData)
	require.Equal(t, 1, getActionFinishedData.ActionIndex)

	finishedEvent, _ = keptnv2.ToKeptnEvent(fakeKeptn.GetEventSender().SentEvents[3])
	require.Equal(t, "sh.keptn.event.get-action.started", fakeKeptn.GetEventSender().SentEvents[2].Type())
	require.Equal(t, "sh.keptn.event.get-action.finished", fakeKeptn.GetEventSender().SentEvents[3].Type())
	getActionFinishedData = keptnv2.GetActionFinishedEventData{}
	finishedEvent.DataAs(&getActionFinishedData)
	require.Equal(t, 2, getActionFinishedData.ActionIndex)
}
