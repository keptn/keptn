package main

import (
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"testing"
)

func Test_Handler(t *testing.T) {
	fakeKeptn := sdk.NewFakeKeptn("test-greeting-svc")
	fakeKeptn.AddTaskHandler(greetingsTriggeredEventType, NewGreetingsHandler())
	fakeKeptn.Start()
	fakeKeptn.NewEvent(newNewGreetingTriggeredEvent("test-assets/events/greeting.triggered-0.json"))

	require.Equal(t, 2, len(fakeKeptn.GetEventSender().SentEvents))
	require.Equal(t, keptnv2.GetStartedEventType("greeting"), fakeKeptn.GetEventSender().SentEvents[0].Type())
	require.Equal(t, keptnv2.GetFinishedEventType("greeting"), fakeKeptn.GetEventSender().SentEvents[1].Type())

	finishedEvent, _ := keptnv2.ToKeptnEvent(fakeKeptn.GetEventSender().SentEvents[1])
	greetingFinishedData := GreetingFinishedData{}
	finishedEvent.DataAs(&greetingFinishedData)
	require.Equal(t, "Hi, my name is Keptn", greetingFinishedData.GreetMessage)
}

func newNewGreetingTriggeredEvent(filename string) cloudevents.Event {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	event := models.KeptnContextExtendedCE{}
	err = json.Unmarshal(content, &event)
	_ = err
	return keptnv2.ToCloudEvent(event)
}
