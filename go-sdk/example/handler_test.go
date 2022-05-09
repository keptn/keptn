package main

import (
	"encoding/json"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

func Test_Handler(t *testing.T) {
	fakeKeptn := sdk.NewFakeKeptn("test-greeting-svc")
	fakeKeptn.AddTaskHandler(greetingsTriggeredEventType, NewGreetingsHandler())
	go fakeKeptn.Start()
	<-fakeKeptn.TestEventSource.Started
	fakeKeptn.NewEvent(newNewGreetingTriggeredEvent("test-assets/events/greeting.triggered-0.json"))

	require.Eventuallyf(t, func() bool {
		return len(fakeKeptn.GetEventSource().SentEvents) == 2
	}, time.Second, 10*time.Millisecond, "error message %s", "formatted")

	require.Eventuallyf(t, func() bool {
		return (keptnv2.GetStartedEventType("greeting") == *fakeKeptn.GetEventSource().SentEvents[0].Type)
	}, time.Second, 10*time.Millisecond, "error message %s", "formatted")

	require.Eventuallyf(t, func() bool {
		return (keptnv2.GetFinishedEventType("greeting") == *fakeKeptn.GetEventSource().SentEvents[1].Type)
	}, time.Second, 10*time.Millisecond, "error message %s", "formatted")

	finishedEvent := fakeKeptn.GetEventSource().SentEvents[1]
	greetingFinishedData := GreetingFinishedData{}
	finishedEvent.DataAs(&greetingFinishedData)
	require.Equal(t, "Hi, my name is Keptn", greetingFinishedData.GreetMessage)
}

func newNewGreetingTriggeredEvent(filename string) models.KeptnContextExtendedCE {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	event := models.KeptnContextExtendedCE{}
	err = json.Unmarshal(content, &event)
	_ = err
	return event
}
