package controller

import (
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

func TestHandleActionTriggeredEvent(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedBaseHandler := NewMockedHandler(createKeptn(), "")
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
	ce := cloudevents.NewEvent()
	_ = ce.SetData(cloudevents.ApplicationJSON, actionTriggeredEventData)

	expectedActionStartedEvent := cloudevents.NewEvent()
	expectedActionStartedEvent.SetType("sh.keptn.event.action.started")
	expectedActionStartedEvent.SetSource("helm-service")
	expectedActionStartedEvent.SetDataContentType(cloudevents.ApplicationJSON)
	expectedActionStartedEvent.SetExtension("triggeredid", "")
	expectedActionStartedEvent.SetExtension("shkeptncontext", "")
	expectedActionStartedEvent.SetData(cloudevents.ApplicationJSON, keptnv2.ActionStartedEventData{
		EventData: keptnv2.EventData{
			Project: "sockshop",
			Stage:   "production",
			Service: "carts",
			Status:  keptnv2.StatusSucceeded,
		},
	})

	expectedActionFinishedEvent := cloudevents.NewEvent()
	expectedActionFinishedEvent.SetType("sh.keptn.event.action.finished")
	expectedActionFinishedEvent.SetSource("helm-service")
	expectedActionFinishedEvent.SetDataContentType(cloudevents.ApplicationJSON)
	expectedActionFinishedEvent.SetExtension("triggeredid", "")
	expectedActionFinishedEvent.SetExtension("shkeptncontext", "")
	expectedActionFinishedEvent.SetData(cloudevents.ApplicationJSON, keptnv2.ActionFinishedEventData{
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
	})

	mockedConfigurationChanger.EXPECT().UpdateChart(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, "123-456", nil)

	instance.HandleEvent(ce, nilCloser)
	assert.Equal(t, expectedActionStartedEvent, mockedBaseHandler.sentCloudEvents[0])
	assert.Equal(t, expectedActionFinishedEvent, mockedBaseHandler.sentCloudEvents[1])

}

func TestHandleEvent_InvalidData(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedBaseHandler := NewMockedHandler(createKeptn(), "")
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

	expectedActionStartedEvent := cloudevents.NewEvent()
	expectedActionStartedEvent.SetType("sh.keptn.event.action.started")
	expectedActionStartedEvent.SetSource("helm-service")
	expectedActionStartedEvent.SetDataContentType(cloudevents.ApplicationJSON)
	expectedActionStartedEvent.SetExtension("triggeredid", "")
	expectedActionStartedEvent.SetExtension("shkeptncontext", "")
	expectedActionStartedEvent.SetData(cloudevents.ApplicationJSON, keptnv2.ActionStartedEventData{
		EventData: keptnv2.EventData{
			Project: "sockshop",
			Stage:   "production",
			Service: "carts",
			Status:  keptnv2.StatusSucceeded,
		},
	})

	expectedActionFinishedEvent := cloudevents.NewEvent()
	expectedActionFinishedEvent.SetType("sh.keptn.event.action.finished")
	expectedActionFinishedEvent.SetSource("helm-service")
	expectedActionFinishedEvent.SetDataContentType(cloudevents.ApplicationJSON)
	expectedActionFinishedEvent.SetExtension("triggeredid", "")
	expectedActionFinishedEvent.SetExtension("shkeptncontext", "")
	expectedActionFinishedEvent.SetData(cloudevents.ApplicationJSON, keptnv2.ActionFinishedEventData{
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
	})

	ce := cloudevents.NewEvent()
	_ = ce.SetData(cloudevents.ApplicationJSON, actionTriggeredEventData)

	instance.HandleEvent(ce, nilCloser)

	assert.Equal(t, expectedActionStartedEvent, mockedBaseHandler.sentCloudEvents[0])
	assert.Equal(t, expectedActionFinishedEvent, mockedBaseHandler.sentCloudEvents[1])

}

func TestHandleUnparsableEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedBaseHandler := NewMockedHandler(createKeptn(), "")
	mockedConfigurationChanger := mocks.NewMockIConfigurationChanger(ctrl)

	instance := ActionTriggeredHandler{
		Handler:       mockedBaseHandler,
		configChanger: mockedConfigurationChanger,
	}

	expectedFinishEventData := keptnv2.ActionFinishedEventData{
		EventData: keptnv2.EventData{
			Status:  "errored",
			Result:  "fail",
			Message: "failed to unmarshal data: [json] found bytes \"\"WEIRD_JSON_CONTENT\"\", but failed to unmarshal: json: cannot unmarshal string into Go value of type v0_2_0.ActionTriggeredEventData",
		},
		Action: keptnv2.ActionData{
			GitCommit: "",
		},
	}

	instance.HandleEvent(createUnparsableEvent(), nilCloser)
	assert.Equal(t, 1, len(mockedBaseHandler.handledErrorEvents))
	assert.Equal(t, expectedFinishEventData, mockedBaseHandler.handledErrorEvents[0])
}

func TestHandleEvent_SendStartEventFails(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Handle started event will fail
	opt := func(o *MockedHandlerOptions) {
		o.SendEventBehavior = func(eventType string) bool {
			if eventType == "sh.keptn.event.action.started" {
				return false
			}
			return true
		}
	}

	mockedBaseHandler := NewMockedHandler(createKeptn(), "", opt) //NewMockHandler(ctrl)
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
			Message: "Failed at sending event of type sh.keptn.event.action.started",
		},
		Action: keptnv2.ActionData{},
	}
	ce := cloudevents.NewEvent()
	_ = ce.SetData(cloudevents.ApplicationJSON, actionTriggeredEventData)

	instance.HandleEvent(ce, nilCloser)
	assert.Equal(t, 1, len(mockedBaseHandler.handledErrorEvents))
	assert.Equal(t, expectedFinishEvent, mockedBaseHandler.handledErrorEvents[0])
}

func TestHandleEvent_SendFinishEventFails(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Handle finished event will fail
	opt := func(o *MockedHandlerOptions) {
		o.SendEventBehavior = func(eventType string) bool {
			if eventType == "sh.keptn.event.action.finished" {
				return false
			}
			return true
		}
	}

	mockedBaseHandler := NewMockedHandler(createKeptn(), "", opt) //NewMockHandler(ctrl)
	mockedConfigurationChanger := mocks.NewMockIConfigurationChanger(ctrl)
	mockedConfigurationChanger.EXPECT().UpdateChart(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, "123-456", nil)

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
			Message: "Failed at sending event of type sh.keptn.event.action.finished",
		},
		Action: keptnv2.ActionData{},
	}
	ce := cloudevents.NewEvent()
	_ = ce.SetData(cloudevents.ApplicationJSON, actionTriggeredEventData)

	instance.HandleEvent(ce, nilCloser)
	assert.Equal(t, 1, len(mockedBaseHandler.handledErrorEvents))
	assert.Equal(t, expectedFinishEvent, mockedBaseHandler.handledErrorEvents[0])
}

func TestHandleEvent_WithMissingAction(t *testing.T) {

	ctrl := gomock.NewController(t)
	ctrl.Finish()
	mockedBaseHandler := NewMockedHandler(createKeptn(), "")
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

	ce := cloudevents.NewEvent()
	_ = ce.SetData(cloudevents.ApplicationJSON, actionTriggeredEventData)

	instance.HandleEvent(ce, nilCloser)

}
