package controller

import (
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/golang/mock/gomock"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateActionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ce := cloudevents.NewEvent()
	keptn, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})

	instance := NewActionTriggeredHandler(keptn, mocks.NewMockIConfigurationChanger(ctrl), "")
	assert.NotNil(t, instance)
}

func TestHandleEvent(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedBaseHandler := NewMockHandler(ctrl)
	mockedConfigurationChanger := mocks.NewMockIConfigurationChanger(ctrl)

	instance := ActionTriggeredHandler{
		Handler:       mockedBaseHandler,
		configChanger: mockedConfigurationChanger,
	}

	actionTriggeredEventData := keptnv2.ActionTriggeredEventData{
		EventData: keptnv2.EventData{
			Project: "sockshop",
			Stage:   "production",
			Service: "carts",
		},
		Action: keptnv2.ActionInfo{
			Name:        "my-scaling-action",
			Action:      "scaling",
			Description: "this is a unit test",
			Value:       "1",
		},
		Problem: keptnv2.ProblemDetails{},
	}

	expectedActionFinishedEvent := keptnv2.ActionFinishedEventData{
		EventData: keptnv2.EventData{
			Project: "sockshop",
			Stage:   "production",
			Service: "carts",
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultPass,
			Message: "Successfully executed scaling action",
		},
		Action: keptnv2.ActionData{
			GitCommit: "123-456",
		},
	}

	mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(createKeptn())
	var capturedEvents []string
	var capturedEventsData []interface{}
	mockedBaseHandler.EXPECT().sendEvent(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(triggerID, ceType string, data interface{}) {
		capturedEvents = append(capturedEvents, ceType)
		capturedEventsData = append(capturedEventsData, data)
	}).AnyTimes()
	mockedConfigurationChanger.EXPECT().UpdateChart(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, "123-456", nil)
	mockedBaseHandler.EXPECT().upgradeChart(gomock.Any(), actionTriggeredEventData.EventData, gomock.Any()).Return(nil)

	ce := cloudevents.NewEvent()
	ce.SetData(cloudevents.ApplicationJSON, actionTriggeredEventData)

	instance.HandleEvent(ce, nilCloser)

	assert.Equal(t, 2, len(capturedEvents))
	assert.Equal(t, "sh.keptn.event.action.started", capturedEvents[0])
	assert.Equal(t, "sh.keptn.event.action.finished", capturedEvents[1])
	assert.Equal(t, expectedActionFinishedEvent, capturedEventsData[1])
}

func TestHandleEvent_InvalidData(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedBaseHandler := NewMockHandler(ctrl)
	mockedConfigurationChanger := mocks.NewMockIConfigurationChanger(ctrl)

	instance := ActionTriggeredHandler{
		Handler:       mockedBaseHandler,
		configChanger: mockedConfigurationChanger,
	}

	actionTriggeredEventData := keptnv2.ActionTriggeredEventData{
		EventData: keptnv2.EventData{
			Project: "sockshop",
			Stage:   "production",
			Service: "carts",
		},
		Action: keptnv2.ActionInfo{
			Name:        "my-scaling-action",
			Action:      "scaling",
			Description: "this is a unit test",
			Value:       "one", // <<-- ohoh
		},
		Problem: keptnv2.ProblemDetails{},
	}

	expectedActionFinishedEvent := keptnv2.ActionFinishedEventData{
		EventData: keptnv2.EventData{
			Project: "sockshop",
			Stage:   "production",
			Service: "carts",
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultFailed,
			Message: "could not parse action.value to int",
		},
		Action: keptnv2.ActionData{
			GitCommit: "",
		},
	}

	mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(createKeptn())
	var capturedEvents []string
	var capturedEventsData []interface{}
	mockedBaseHandler.EXPECT().sendEvent(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(triggerID, ceType string, data interface{}) {
		capturedEvents = append(capturedEvents, ceType)
		capturedEventsData = append(capturedEventsData, data)
	}).AnyTimes()

	ce := cloudevents.NewEvent()
	ce.SetData(cloudevents.ApplicationJSON, actionTriggeredEventData)

	instance.HandleEvent(ce, nilCloser)

	assert.Equal(t, 2, len(capturedEvents))
	assert.Equal(t, "sh.keptn.event.action.started", capturedEvents[0])
	assert.Equal(t, "sh.keptn.event.action.finished", capturedEvents[1])
	assert.Equal(t, expectedActionFinishedEvent, capturedEventsData[1])
}

func TestHandleUnparsableEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedBaseHandler := NewMockHandler(ctrl)
	mockedConfigurationChanger := mocks.NewMockIConfigurationChanger(ctrl)

	instance := ActionTriggeredHandler{
		Handler:       mockedBaseHandler,
		configChanger: mockedConfigurationChanger,
	}

	expectedFinishEvent := keptnv2.ActionFinishedEventData{
		EventData: keptnv2.EventData{
			Status:  "errored",
			Result:  "fail",
			Message: "failed to unmarshal data: [json] found bytes \"\"WEIRD_JSON_CONTENT\"\", but failed to unmarshal: json: cannot unmarshal string into Go value of type v0_2_0.ActionTriggeredEventData",
		},
		Action: keptnv2.ActionData{
			GitCommit: "",
		},
	}

	mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(createKeptn())
	mockedBaseHandler.EXPECT().handleError("EVENT_ID", gomock.Any(), "action", expectedFinishEvent)

	instance.HandleEvent(createUnparsableEvent(), nilCloser)
}

func TestHandleEvent_SendStartEventFails(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedBaseHandler := NewMockHandler(ctrl)
	mockedConfigurationChanger := mocks.NewMockIConfigurationChanger(ctrl)

	instance := ActionTriggeredHandler{
		Handler:       mockedBaseHandler,
		configChanger: mockedConfigurationChanger,
	}

	actionTriggeredEventData := keptnv2.ActionTriggeredEventData{
		EventData: keptnv2.EventData{
			Project: "sockshop",
			Stage:   "production",
			Service: "carts",
		},
		Action: keptnv2.ActionInfo{
			Name:        "my-scaling-action",
			Action:      "scaling",
			Description: "this is a unit test",
			Value:       "1",
		},
		Problem: keptnv2.ProblemDetails{},
	}

	expectedFinishEvent := keptnv2.ActionFinishedEventData{
		EventData: keptnv2.EventData{
			Project: "sockshop",
			Stage:   "production",
			Service: "carts",
			Status:  "errored",
			Result:  "fail",
			Message: "failed to send event",
		},
		Action: keptnv2.ActionData{
			GitCommit: "",
		},
	}

	mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(createKeptn())
	mockedBaseHandler.EXPECT().sendEvent(gomock.Any(), "sh.keptn.event.action.started", gomock.Any()).Return(errors.New("failed to send event"))
	mockedBaseHandler.EXPECT().handleError("", errors.New("failed to send event"), "action", expectedFinishEvent)

	ce := cloudevents.NewEvent()
	ce.SetData(cloudevents.ApplicationJSON, actionTriggeredEventData)

	instance.HandleEvent(ce, nilCloser)
}

func TestHandleEvent_SendFinishEventFails(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedBaseHandler := NewMockHandler(ctrl)
	mockedConfigurationChanger := mocks.NewMockIConfigurationChanger(ctrl)

	instance := ActionTriggeredHandler{
		Handler:       mockedBaseHandler,
		configChanger: mockedConfigurationChanger,
	}

	actionTriggeredEventData := keptnv2.ActionTriggeredEventData{
		EventData: keptnv2.EventData{
			Project: "sockshop",
			Stage:   "production",
			Service: "carts",
		},
		Action: keptnv2.ActionInfo{
			Name:        "my-scaling-action",
			Action:      "scaling",
			Description: "this is a unit test",
			Value:       "1",
		},
		Problem: keptnv2.ProblemDetails{},
	}

	expectedFinishEvent := keptnv2.ActionFinishedEventData{
		EventData: keptnv2.EventData{
			Project: "sockshop",
			Stage:   "production",
			Service: "carts",
			Status:  "errored",
			Result:  "fail",
			Message: "OHOH",
		},
		Action: keptnv2.ActionData{
			GitCommit: "",
		},
	}

	mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(createKeptn())
	mockedBaseHandler.EXPECT().sendEvent(gomock.Any(), "sh.keptn.event.action.started", gomock.Any()).Return(nil)
	mockedBaseHandler.EXPECT().sendEvent(gomock.Any(), "sh.keptn.event.action.finished", gomock.Any()).Return(errors.New("OHOH"))
	mockedConfigurationChanger.EXPECT().UpdateChart(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, "123-456", nil)
	mockedBaseHandler.EXPECT().upgradeChart(gomock.Any(), actionTriggeredEventData.EventData, gomock.Any()).Return(nil)
	mockedBaseHandler.EXPECT().handleError("", gomock.Any(), "action", expectedFinishEvent)

	ce := cloudevents.NewEvent()
	ce.SetData(cloudevents.ApplicationJSON, actionTriggeredEventData)

	instance.HandleEvent(ce, nilCloser)
}

func TestHandleEvent_WithMissingAction(t *testing.T) {

	ctrl := gomock.NewController(t)
	ctrl.Finish()
	mockedBaseHandler := NewMockHandler(ctrl)
	mockedConfigurationChanger := mocks.NewMockIConfigurationChanger(ctrl)

	instance := ActionTriggeredHandler{
		Handler:       mockedBaseHandler,
		configChanger: mockedConfigurationChanger,
	}

	actionTriggeredEventData := keptnv2.ActionTriggeredEventData{
		EventData: keptnv2.EventData{},
		Action: keptnv2.ActionInfo{
			Action: "", // <<-- ohoh
		},
		Problem: keptnv2.ProblemDetails{},
	}

	mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(createKeptn())

	ce := cloudevents.NewEvent()
	ce.SetData(cloudevents.ApplicationJSON, actionTriggeredEventData)

	instance.HandleEvent(ce, nilCloser)

}
