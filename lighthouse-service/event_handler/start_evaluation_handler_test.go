package event_handler

import (
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

func TestHasTestPassedTrue(t *testing.T) {
	contentType := "application/json"
	source, _ := url.Parse("lighthouse-service")
	shkeptncontext := "0000-1111-2222-3333"

	// test TestsFinishedEvent
	testsFinishedEventData := keptnevents.TestsFinishedEventData{
		Result:             "pass",
		Project:            "sockshop",
		Service:            "carts",
		Stage:              "staging",
		TestStrategy:       "functional",
		DeploymentStrategy: "direct",
	}

	testsFinishedEvent := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptnevents.TestsFinishedEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: testsFinishedEventData,
	}

	logger := keptnutils.NewLogger(shkeptncontext, testsFinishedEvent.Context.GetID(), "lighthouse-service")
	eh := StartEvaluationHandler{Logger: logger, Event: testsFinishedEvent}
	assert.EqualValues(t, eh.hasTestPassed(), true)

	// test StartEvaluationEvent
	startEvaluationEventData := keptnevents.StartEvaluationEventData{
		Project:            "sockshop",
		Service:            "carts",
		Stage:              "staging",
		TestStrategy:       "functional",
		DeploymentStrategy: "direct",
	}

	startEvaluationEvent := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptnevents.StartEvaluationEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: startEvaluationEventData,
	}

	eh = StartEvaluationHandler{Logger: logger, Event: startEvaluationEvent}
	assert.EqualValues(t, eh.hasTestPassed(), true)
}

func TestHasTestPassedFalse(t *testing.T) {
	contentType := "application/json"
	source, _ := url.Parse("lighthouse-service")
	shkeptncontext := "0000-1111-2222-3333"

	testFinishedEvent := keptnevents.TestsFinishedEventData{
		Result:             "fail",
		Project:            "sockshop",
		Service:            "carts",
		Stage:              "staging",
		TestStrategy:       "functional",
		DeploymentStrategy: "direct",
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptnevents.EvaluationDoneEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: testFinishedEvent,
	}

	logger := keptnutils.NewLogger(shkeptncontext, event.Context.GetID(), "lighthouse-service")
	eh := StartEvaluationHandler{Logger: logger, Event: event}
	assert.EqualValues(t, eh.hasTestPassed(), false)
}