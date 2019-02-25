package main

import (
	"testing"

	"github.com/keptn/keptn/cli/utils"
	"github.com/knative/pkg/cloudevents"
)

type Example struct {
	Sequence int    `json:"id"`
	Message  string `json:"message"`
}

type Empty struct {
}

func TestKnativeCloudEvents(t *testing.T) {
	builder := cloudevents.Builder{
		Source:    "https://github.com/keptn/keptn/cli#cloudevents-example",
		EventType: "cloudevent.example",
	}

	// data := Example{
	// 	Message:  "hello, world!",
	// 	Sequence: 1,
	// }

	data := Empty{}

	err := utils.Send("http://control-andreas.keptn.35.239.5.164.xip.io/onboard", "***REMOVED***", builder, data)
	if err != nil {
		t.Errorf("An error occured %v", err)
	}
}
