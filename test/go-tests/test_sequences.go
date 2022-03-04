package go_tests

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	models2 "github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

type TestCaseOpt func(testCase *TestCase)

func WithShipyard(shipyard string) TestCaseOpt {
	return func(testCase *TestCase) {
		testCase.Shipyard = shipyard
	}
}

func WithService(service string) TestCaseOpt {
	return func(testCase *TestCase) {
		testCase.Service = service
	}
}

func NewTestCase(t *testing.T, projectName string, opts ...TestCaseOpt) *TestCase {
	tc := &TestCase{t: t, Project: projectName}

	for _, opt := range opts {
		opt(tc)
	}

	return tc
}

type SequenceValidation interface {
	Verify(tc *TestCase)
}

type Checkpoint interface {
	Do(tc *TestCase)
}

type TestCase struct {
	Project  string
	Service  string
	Shipyard string

	checkpoints []Checkpoint

	t               *testing.T
	keptnContext    string
	lastTriggeredID string
}

func (tc *TestCase) AddCheckpoint(cp Checkpoint) *TestCase {
	tc.checkpoints = append(tc.checkpoints, cp)
	return tc
}

func (tc *TestCase) Run() {
	shipyardFilePath, err := CreateTmpShipyardFile(tc.Shipyard)
	defer func() {
		err := os.Remove(shipyardFilePath)
		if err != nil {
			tc.t.Logf("Could not delete file: %s: %v", shipyardFilePath, err)
		}
	}()
	require.Nil(tc.t, err)

	tc.Project, err = CreateProject(tc.Project, shipyardFilePath)
	require.Nil(tc.t, err)

	if tc.Service != "" {
		output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", tc.Service, tc.Project))

		require.Nil(tc.t, err)
		require.Contains(tc.t, output, "created successfully")
	}

	for _, cp := range tc.checkpoints {
		cp.Do(tc)
	}
}

type TriggerSequenceCheckpoint struct {
	SequenceName string
	Stage        string
	Parameters   map[string]interface{}
}

func (ts TriggerSequenceCheckpoint) Do(tc *TestCase) {
	source := "golang-test"
	eventType := keptnv2.GetTriggeredEventType(ts.Stage + "." + ts.SequenceName)
	eventData := &keptnv2.EventData{}
	eventData.SetProject(tc.Project)
	eventData.SetService(tc.Service)
	eventData.SetStage(ts.Stage)

	resp, err := ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
		Contenttype:        "application/json",
		Data:               eventData,
		ID:                 uuid.NewString(),
		Shkeptnspecversion: KeptnSpecVersion,
		Source:             &source,
		Specversion:        "1.0",
		Type:               &eventType,
	}, 3)
	require.Nil(tc.t, err)

	eventContext := &models.EventContext{}
	err = resp.ToJSON(eventContext)
	require.Nil(tc.t, err)

	require.NotEmpty(tc.t, *eventContext.KeptnContext)
	tc.keptnContext = *eventContext.KeptnContext
}

type EventResponse struct {
	ResponseEvent models.KeptnContextExtendedCE
}

func (e EventResponse) Do(tc *TestCase) {
	if keptnv2.IsFinishedEventType(*e.ResponseEvent.Type) || keptnv2.IsStartedEventType(*e.ResponseEvent.Type) {
		if e.ResponseEvent.Triggeredid == "" {
			e.ResponseEvent.Triggeredid = tc.lastTriggeredID
		}
	}
	testSource := "go-tests"
	e.ResponseEvent.Shkeptncontext = tc.keptnContext
	e.ResponseEvent.Specversion = "1.0"
	e.ResponseEvent.ID = uuid.NewString()
	e.ResponseEvent.Contenttype = "application/json"
	e.ResponseEvent.Shkeptnspecversion = KeptnSpecVersion
	e.ResponseEvent.Source = &testSource
	resp, err := ApiPOSTRequest("/v1/event", e.ResponseEvent, 3)
	require.Nil(tc.t, err)
	require.Equal(tc.t, http.StatusOK, resp.Response().StatusCode)
}

type ExpectedEvent struct {
	Type             string
	Stage            string
	VerificationFunc func(t *testing.T, event models.KeptnContextExtendedCE)
}

func (e ExpectedEvent) Do(tc *TestCase) {
	var event *models.KeptnContextExtendedCE
	var err error
	require.Eventually(tc.t, func() bool {
		event, err = GetLatestEventOfType(tc.keptnContext, tc.Project, e.Stage, e.Type)
		if err != nil {
			return false
		}
		if event == nil {
			return false
		}
		return true
	}, 1*time.Minute, 5*time.Second)

	if keptnv2.IsTriggeredEventType(*event.Type) {
		tc.lastTriggeredID = event.ID
	}
	if e.VerificationFunc != nil {
		e.VerificationFunc(tc.t, *event)
	}
}

type ExpectedState struct {
	VerificationFunc func(t *testing.T, state models2.SequenceState)
}

func (e ExpectedState) Do(tc *TestCase) {
	var state models2.SequenceState
	require.Eventually(tc.t, func() bool {
		states, _, err := GetState(tc.Project)
		if err != nil {
			return false
		}
		if states == nil || len(states.States) == 0 {
			return false
		}

		for _, s := range states.States {
			if s.Shkeptncontext == tc.keptnContext {
				state = s
				return true
			}
		}

		return false
	}, 20*time.Second, 2*time.Second)

	if e.VerificationFunc != nil {
		e.VerificationFunc(tc.t, state)
	}
}

func Test_Sequences(t *testing.T) {

	tc := NewTestCase(t, "my-new-project", WithShipyard(sequenceStateShipyard), WithService("my-service"))

	tc.AddCheckpoint(&TriggerSequenceCheckpoint{
		SequenceName: "delivery",
		Stage:        "dev",
	})

	tc.AddCheckpoint(&ExpectedEvent{
		Type:  keptnv2.GetTriggeredEventType("delivery"),
		Stage: "dev",
		VerificationFunc: func(t *testing.T, event models.KeptnContextExtendedCE) {
			require.NotNil(t, event.GitCommitID)
		},
	})

	tc.AddCheckpoint(&ExpectedState{
		VerificationFunc: func(t *testing.T, state models2.SequenceState) {
			require.Equal(t, models2.SequenceStartedState, state.State)
		},
	})

	startedEvent := keptnv2.GetStartedEventType("delivery")
	tc.AddCheckpoint(&EventResponse{
		ResponseEvent: models.KeptnContextExtendedCE{
			Type: &startedEvent,
			Data: keptnv2.EventData{
				Project: "keptn-my-new-project",
				Stage:   "dev",
				Service: "my-service",
			},
		},
	})

	tc.Run()
}
