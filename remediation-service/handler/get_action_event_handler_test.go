package handler_test

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_1_4"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/keptn/keptn/remediation-service/handler"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

func newGetActionTriggeredEvent(filename string) models.KeptnContextExtendedCE {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	event := models.KeptnContextExtendedCE{}
	json.Unmarshal(content, &event)
	return event
}

func Test_Receiving_GetActionTriggeredEvent_RemediationFromServiceLevel(t *testing.T) {
	fakeKeptn := sdk.NewFakeKeptn("test-remediation-svc")
	fakeKeptn.AddTaskHandler("sh.keptn.event.get-action.triggered", handler.NewGetActionEventHandler())
	fakeKeptn.NewEvent(newGetActionTriggeredEvent("test/events/get-action.triggered-0.json"))
	fakeKeptn.NewEvent(newGetActionTriggeredEvent("test/events/get-action.triggered-1.json"))
	fakeKeptn.NewEvent(newGetActionTriggeredEvent("test/events/get-action.triggered-2.json"))

	fakeKeptn.AssertNumberOfEventSent(t, 6)
	fakeKeptn.AssertSentEventType(t, 0, keptnv2.GetStartedEventType("get-action"))
	fakeKeptn.AssertSentEventType(t, 1, keptnv2.GetFinishedEventType("get-action"))
	fakeKeptn.AssertSentEventType(t, 2, keptnv2.GetStartedEventType("get-action"))
	fakeKeptn.AssertSentEventType(t, 3, keptnv2.GetFinishedEventType("get-action"))
	fakeKeptn.AssertSentEventType(t, 4, keptnv2.GetStartedEventType("get-action"))
	fakeKeptn.AssertSentEventType(t, 5, keptnv2.GetFinishedEventType("get-action"))

	fakeKeptn.AssertSentEventStatus(t, 1, keptnv2.StatusSucceeded)
	fakeKeptn.AssertSentEventResult(t, 1, keptnv2.ResultPass)

	fakeKeptn.AssertSentEvent(t, 1, func(ce models.KeptnContextExtendedCE) bool {
		getActionFinishedData := keptnv2.GetActionFinishedEventData{}
		ce.DataAs(&getActionFinishedData)
		return getActionFinishedData.GetAction.ActionIndex == 1
	})
	fakeKeptn.AssertSentEvent(t, 3, func(ce models.KeptnContextExtendedCE) bool {
		getActionFinishedData := keptnv2.GetActionFinishedEventData{}
		ce.DataAs(&getActionFinishedData)
		return getActionFinishedData.GetAction.ActionIndex == 2
	})
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
