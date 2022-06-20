package main

import (
	"encoding/json"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"io/ioutil"
	"log"
	"testing"
)

func Test_Handler(t *testing.T) {
	fakeKeptn := sdk.NewFakeKeptn("test-greeting-svc", nil)
	fakeKeptn.AddTaskHandler(greetingsTriggeredEventType, NewGreetingsHandler())
	fakeKeptn.NewEvent(newNewGreetingTriggeredEvent("test-assets/events/greeting.triggered-0.json"))
	fakeKeptn.AssertNumberOfEventSent(t, 2)
	fakeKeptn.AssertSentEventType(t, 0, keptnv2.GetStartedEventType("greeting"))
	fakeKeptn.AssertSentEventType(t, 1, keptnv2.GetFinishedEventType("greeting"))
	fakeKeptn.AssertSentEvent(t, 1, func(e models.KeptnContextExtendedCE) bool {
		greetingFinishedData := GreetingFinishedData{}
		e.DataAs(&greetingFinishedData)
		return "Hi, my name is Keptn" == greetingFinishedData.GreetMessage
	})
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
