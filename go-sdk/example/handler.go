package main

import (
	"bytes"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"html/template"
)

type GreetingsHandler struct {
}

type GreetingTriggeredData struct {
	v0_2_0.EventData
	Text string `json:"text"`
}

type GreetingFinishedData struct {
	v0_2_0.EventData
	GreetMessage string `json:"greetMessage"`
}

func NewGreetingsHandler() *GreetingsHandler {
	return &GreetingsHandler{}
}

func (g *GreetingsHandler) Execute(k sdk.IKeptn, event sdk.KeptnEvent) (interface{}, *sdk.Error) {
	greetingsTriggeredData := &GreetingTriggeredData{}
	if err := v0_2_0.Decode(event.Data, greetingsTriggeredData); err != nil {
		return nil, &sdk.Error{Err: err, StatusType: v0_2_0.StatusErrored, ResultType: v0_2_0.ResultFailed, Message: "Could not decode input event data"}
	}
	name := struct{ Name string }{"Keptn"}

	tmpl, err := template.New("").Parse(greetingsTriggeredData.Text)
	if err != nil {
		return nil, &sdk.Error{Err: err, StatusType: v0_2_0.StatusErrored, ResultType: v0_2_0.ResultFailed, Message: "Could not parse greeting message"}
	}

	var greetMessage bytes.Buffer
	if err = tmpl.Execute(&greetMessage, name); err != nil {
		return nil, &sdk.Error{Err: err, StatusType: v0_2_0.StatusErrored, ResultType: v0_2_0.ResultFailed, Message: "Could not parse process greeting message"}
	}
	finishedEventData := GreetingFinishedData{
		EventData:    greetingsTriggeredData.EventData,
		GreetMessage: greetMessage.String(),
	}
	return finishedEventData, nil
}
