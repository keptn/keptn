package handler

import (
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/remediation-service/pkg/sdk"
	"github.com/keptn/keptn/remediation-service/pkg/sdk/fake"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"testing"
)

func Test_Receiving_GetActionTriggeredEvent(t *testing.T) {

	resourceHandlerMock := &fake.ResourceHandlerMock{}
	resourceHandlerMock.GetServiceResourceFunc = func(project string, stage string, service string, resourceURI string) (*models.Resource, error) {
		return newRemediationResource("test/remediation.yaml"), nil
	}

	fakeKeptn := fake.NewFakeKeptn("test-remediation-svc", resourceHandlerMock, sdk.WithHandler(NewGetActionEventHandler(), "sh.keptn.event.get-action.triggered"))
	fakeKeptn.Start()
	fakeKeptn.NewEvent(newGetActionTriggeredEvent("test/get-action.triggered.json"))

	require.Equal(t, 2, len(fakeKeptn.GetEventSender().SentEvents))
	event, _ := keptnv2.ToKeptnEvent(fakeKeptn.GetEventSender().SentEvents[1])

	// verify started event
	require.Equal(t, "sh.keptn.event.get-action.started", fakeKeptn.GetEventSender().SentEvents[0].Type())

	// verify finished event
	require.Equal(t, "sh.keptn.event.get-action.finished", fakeKeptn.GetEventSender().SentEvents[1].Type())
	getActionFinishedData := keptnv2.GetActionFinishedEventData{}
	event.DataAs(&getActionFinishedData)
	require.Equal(t, 1, getActionFinishedData.ActionIndex)
}

func newRemediationResource(filename string) *models.Resource {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return &models.Resource{
		Metadata:        nil,
		ResourceContent: string(content),
		ResourceURI:     nil,
	}
}

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
