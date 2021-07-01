package handler_test

import (
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/ghodss/yaml"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_1_4"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/remediation-service/handler"
	"github.com/keptn/keptn/remediation-service/internal/sdk"
	"github.com/keptn/keptn/remediation-service/internal/sdk/fake"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"reflect"
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

func Test_Receiving_GetActionTriggeredEvent_RemediationFromServiceLevel(t *testing.T) {

	fakeKeptn := fake.NewFakeKeptn("test-remediation-svc", sdk.WithHandler("sh.keptn.event.get-action.triggered", handler.NewGetActionEventHandler()))
	fakeKeptn.Start()
	fakeKeptn.NewEvent(newGetActionTriggeredEvent("test/events/get-action.triggered-0.json"))
	fakeKeptn.NewEvent(newGetActionTriggeredEvent("test/events/get-action.triggered-1.json"))
	fakeKeptn.NewEvent(newGetActionTriggeredEvent("test/events/get-action.triggered-2.json"))

	require.Equal(t, 6, len(fakeKeptn.GetEventSender().SentEvents))

	require.Equal(t, keptnv2.GetStartedEventType("get-action"), fakeKeptn.GetEventSender().SentEvents[0].Type())
	require.Equal(t, keptnv2.GetFinishedEventType("get-action"), fakeKeptn.GetEventSender().SentEvents[1].Type())
	require.Equal(t, keptnv2.GetStartedEventType("get-action"), fakeKeptn.GetEventSender().SentEvents[2].Type())
	require.Equal(t, keptnv2.GetFinishedEventType("get-action"), fakeKeptn.GetEventSender().SentEvents[3].Type())
	require.Equal(t, keptnv2.GetStartedEventType("get-action"), fakeKeptn.GetEventSender().SentEvents[4].Type())
	require.Equal(t, keptnv2.GetFinishedEventType("get-action"), fakeKeptn.GetEventSender().SentEvents[5].Type())

	finishedEvent, _ := keptnv2.ToKeptnEvent(fakeKeptn.GetEventSender().SentEvents[1])
	getActionFinishedData := keptnv2.GetActionFinishedEventData{}
	finishedEvent.DataAs(&getActionFinishedData)
	require.Equal(t, 1, getActionFinishedData.ActionIndex)
	require.Equal(t, keptnv2.StatusSucceeded, getActionFinishedData.Status)
	require.Equal(t, keptnv2.ResultPass, getActionFinishedData.Result)

	finishedEvent, _ = keptnv2.ToKeptnEvent(fakeKeptn.GetEventSender().SentEvents[3])
	getActionFinishedData = keptnv2.GetActionFinishedEventData{}
	finishedEvent.DataAs(&getActionFinishedData)
	require.Equal(t, 2, getActionFinishedData.ActionIndex)
	require.Equal(t, keptnv2.StatusSucceeded, getActionFinishedData.Status)
	require.Equal(t, keptnv2.ResultPass, getActionFinishedData.Result)

	finishedEvent, _ = keptnv2.ToKeptnEvent(fakeKeptn.GetEventSender().SentEvents[5])
	getActionFinishedData = keptnv2.GetActionFinishedEventData{}
	finishedEvent.DataAs(&getActionFinishedData)
	require.Equal(t, keptnv2.StatusSucceeded, getActionFinishedData.Status)
	require.Equal(t, keptnv2.ResultFailed, getActionFinishedData.Result)

}

func newRemediation(fileName string) *v0_1_4.Remediation {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	remediation := &v0_1_4.Remediation{}
	yaml.Unmarshal(content, remediation)
	return remediation
}

func newProblemDetails(problemTitle, rootCause string) keptnv2.ProblemDetails {
	return keptnv2.ProblemDetails{
		ProblemTitle: problemTitle,
		RootCause:    rootCause,
	}
}

func TestGetNextAction(t *testing.T) {

	type args struct {
		remediation    *v0_1_4.Remediation
		problemDetails keptnv2.ProblemDetails
		actionIndex    int
	}
	tests := []struct {
		name    string
		args    args
		want    *keptnv2.ActionInfo
		wantErr bool
	}{
		{
			"determine action - by rootCause",
			args{
				newRemediation("test/remediation.yaml"),
				newProblemDetails("", "problemType1"),
				0,
			},
			&keptnv2.ActionInfo{
				Name:        "actionName1",
				Action:      "action1",
				Description: "actionDescription1",
				Value:       map[string]interface{}{"foo": "bar"},
			},
			false,
		},
		{
			"determine-action - by problemTitle",
			args{
				newRemediation("test/remediation.yaml"),
				newProblemDetails("problemType1", ""),
				0,
			},
			&keptnv2.ActionInfo{
				Name:        "actionName1",
				Action:      "action1",
				Description: "actionDescription1",
				Value:       map[string]interface{}{"foo": "bar"},
			},
			false,
		},
		{
			"determine action - by rootCause2",
			args{
				newRemediation("test/remediation.yaml"),
				newProblemDetails("problemType1", "problemType2"),
				0,
			},
			&keptnv2.ActionInfo{
				Name:        "actionName11",
				Action:      "action11",
				Description: "actionDescription11",
				Value:       map[string]interface{}{"foo": "bar"},
			},
			false,
		},
		{
			"determine action - not found - Default",
			args{
				newRemediation("test/remediation.yaml"),
				newProblemDetails("", ""),
				0,
			},
			&keptnv2.ActionInfo{
				Name:        "escalateDefaultName",
				Action:      "escalateDefaultAction",
				Description: "escalateDefaultDescription",
				Value:       map[string]interface{}{"foo": "bar"},
			},
			false,
		},
		{
			"determine action - not found - no Default",
			args{
				newRemediation("test/remediation-without-default.yaml"),
				newProblemDetails("", ""),
				0,
			},
			nil,
			true,
		},
		{
			"determine action - action index out of bound",
			args{
				newRemediation("test/remediation.yaml"),
				newProblemDetails("", ""),
				2,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handler.GetNextAction(tt.args.remediation, tt.args.problemDetails, tt.args.actionIndex)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNextAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNextAction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRemediationResource(t *testing.T) {
	type args struct {
		resource *models.Resource
	}
	tests := []struct {
		name    string
		args    args
		want    *v0_1_4.Remediation
		wantErr bool
	}{
		{"", args{
			resource: newResourceFromFile("test/remediation.yaml"),
		}, newRemediation("test/remediation.yaml"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handler.ParseRemediationResource(tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRemediationResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRemediationResource() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func newResourceFromFile(filename string) *models.Resource {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to locate resources requested by the service: %s", err.Error())
	}

	return &models.Resource{
		Metadata:        nil,
		ResourceContent: string(content),
		ResourceURI:     nil,
	}
}
