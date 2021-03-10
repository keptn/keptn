package handler

import (
	"encoding/json"
	"errors"
	"github.com/go-test/deep"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
	time "time"
)

const testShipyardResourceWithInvalidVersion = `{
      "resourceContent": "YXBpVmVyc2lvbjogMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJk",
      "resourceURI": "shipyard.yaml"
    }`

const testShipyardFileWithInvalidVersion = `apiVersion: 0
kind: Shipyard
metadata:
  name: test-shipyard`

const testShipyardResource = `{
      "resourceContent": "YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJkCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5CiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBzdHJhdGVneTogZGlyZWN0CiAgICAgIC0gbmFtZTogdGVzdAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBraW5kOiBmdW5jdGlvbmFsCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbiAKICAgICAgLSBuYW1lOiByZWxlYXNlIAogICAgLSBuYW1lOiByb2xsYmFjawogICAgICB0YXNrczoKICAgICAgLSBuYW1lOiByb2xsYmFjawogICAgICB0cmlnZ2VyZWRPbjoKICAgICAgICAtIGV2ZW50OiBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgICAgIHNlbGVjdG9yOgogICAgICAgICAgICBtYXRjaDoKICAgICAgICAgICAgICByZXN1bHQ6IGZhaWwKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyZWRPbjoKICAgICAgICAtIGV2ZW50OiBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCgogIC0gbmFtZTogcHJvZHVjdGlvbgogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeSAKICAgICAgdHJpZ2dlcmVkT246CiAgICAgICAgLSBldmVudDogaGFyZGVuaW5nLmFydGlmYWN0LWRlbGl2ZXJ5LmZpbmlzaGVkCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOgogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW4KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgIAogICAgLSBuYW1lOiByZW1lZGlhdGlvbgogICAgICB0YXNrczoKICAgICAgLSBuYW1lOiByZW1lZGlhdGlvbgogICAgICAtIG5hbWU6IGV2YWx1YXRpb24=",
      "resourceURI": "shipyard.yaml"
    }`

const testShipyardFile = `apiVersion: spec.keptn.sh/0.2.0
kind: Shipyard
metadata:
  name: test-shipyard
spec:
  stages:
  - name: dev
    sequences:
    - name: artifact-delivery
      tasks:
      - name: deployment
        properties:  
          strategy: direct
      - name: test
        properties:
          kind: functional
      - name: evaluation 
      - name: release 
    - name: rollback
      tasks:
      - name: rollback
      triggeredOn:
        - event: dev.artifact-delivery.finished
          selector:
            match:
              result: fail
  - name: hardening
    sequences:
    - name: artifact-delivery
      triggeredOn:
        - event: dev.artifact-delivery.finished
      tasks:
      - name: deployment
        properties: 
          strategy: blue_green_service
      - name: test
        properties:  
          kind: performance
      - name: evaluation
      - name: release

  - name: production
    sequences:
    - name: artifact-delivery 
      triggeredOn:
        - event: hardening.artifact-delivery.finished
      tasks:
      - name: deployment
        properties:
          strategy: blue_green
      - name: release
      
    - name: remediation
      tasks:
      - name: remediation
      - name: evaluation`

func Test_eventManager_GetAllTriggeredEvents(t *testing.T) {
	type fields struct {
		projectRepo        db.ProjectRepo
		triggeredEventRepo db.EventRepo
	}
	type args struct {
		filter common.EventFilter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Event
		wantErr bool
	}{
		{
			name: "Get triggered events for all projects",
			fields: fields{
				projectRepo: &db_mock.ProjectRepoMock{GetProjectsFunc: func() ([]*models.ExpandedProject, error) {
					return []*models.ExpandedProject{{
						ProjectName: "sockshop",
					}, {
						ProjectName: "rockshop",
					}}, nil
				}},
				triggeredEventRepo: &db_mock.EventRepoMock{
					GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
						return []models.Event{fake.GetTestTriggeredEvent()}, nil
					},
					InsertEventFunc: nil,
					DeleteEventFunc: nil,
				},
			},
			args: args{},
			want: []models.Event{
				fake.GetTestTriggeredEvent(),
				fake.GetTestTriggeredEvent(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := &shipyardController{
				projectRepo: tt.fields.projectRepo,
				eventRepo:   tt.fields.triggeredEventRepo,
			}
			got, err := em.GetAllTriggeredEvents(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllTriggeredEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllTriggeredEvents() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_eventManager_GetTriggeredEventsOfProject(t *testing.T) {
	type fields struct {
		projectRepo        db.ProjectRepo
		triggeredEventRepo db.EventRepo
	}
	type args struct {
		project string
		filter  common.EventFilter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Event
		wantErr bool
	}{
		{
			name: "Get triggered events for project",
			fields: fields{
				projectRepo: nil,
				triggeredEventRepo: &db_mock.EventRepoMock{
					GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
						return []models.Event{fake.GetTestTriggeredEvent()}, nil
					},
					InsertEventFunc: nil,
					DeleteEventFunc: nil,
				},
			},
			args: args{},
			want: []models.Event{
				fake.GetTestTriggeredEvent(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := &shipyardController{
				projectRepo: tt.fields.projectRepo,
				eventRepo:   tt.fields.triggeredEventRepo,
			}
			got, err := em.GetTriggeredEventsOfProject(tt.args.project, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTriggeredEventsOfProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTriggeredEventsOfProject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEventScope(t *testing.T) {
	type args struct {
		event models.Event
	}
	tests := []struct {
		name    string
		args    args
		want    *keptnv2.EventData
		wantErr bool
	}{
		{
			name: "get event scope",
			args: args{
				event: models.Event{
					Contenttype:    "",
					Data:           keptnv2.EventData{Project: "sockshop", Stage: "dev", Service: "carts"},
					Extensions:     nil,
					ID:             "",
					Shkeptncontext: "",
					Source:         nil,
					Specversion:    "",
					Time:           "",
					Triggeredid:    "",
					Type:           nil,
				},
			},
			want:    &keptnv2.EventData{Project: "sockshop", Stage: "dev", Service: "carts"},
			wantErr: false,
		},
		{
			name: "only project available, stage and service missing",
			args: args{
				event: models.Event{
					Contenttype:    "",
					Data:           keptnv2.EventData{Project: "sockshop"},
					Extensions:     nil,
					ID:             "",
					Shkeptncontext: "",
					Source:         nil,
					Specversion:    "",
					Time:           "",
					Triggeredid:    "",
					Type:           nil,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty data",
			args: args{
				event: models.Event{
					Contenttype:    "",
					Data:           nil,
					Extensions:     nil,
					ID:             "",
					Shkeptncontext: "",
					Source:         nil,
					Specversion:    "",
					Time:           "",
					Triggeredid:    "",
					Type:           nil,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "nonsense data",
			args: args{
				event: models.Event{
					Contenttype:    "",
					Data:           "invalid",
					Extensions:     nil,
					ID:             "",
					Shkeptncontext: "",
					Source:         nil,
					Specversion:    "",
					Time:           "",
					Triggeredid:    "",
					Type:           nil,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getEventScope(tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEventScope() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEventScope() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_eventManager_handleStartedEvent(t *testing.T) {
	type fields struct {
		projectRepo db.ProjectRepo
		eventRepo   db.EventRepo
		logger      *keptncommon.Logger
	}
	type args struct {
		event models.Event
	}
	tests := []struct {
		name                   string
		fields                 fields
		args                   args
		wantErr                bool
		wantErrNoMatchingEvent bool
	}{
		{
			name: "received started event with matching triggered event",
			fields: fields{
				projectRepo: nil,
				eventRepo: &db_mock.EventRepoMock{
					GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
						if status[0] == common.TriggeredEvent {
							return []models.Event{fake.GetTestTriggeredEvent()}, nil
						}
						return nil, errors.New("received unexpected request")
					},
					InsertEventFunc: func(project string, event models.Event, status common.EventStatus) error {
						if len(deep.Equal(event, fake.GetTestStartedEvent())) != 0 {
							t.Errorf("received unexpected event in insertEvent func. wanted %v but got %v", fake.GetTestStartedEvent(), event)
							return nil
						}
						return nil
					},
					DeleteEventFunc: func(project string, eventID string, status common.EventStatus) error {
						return nil
					},
				},
				logger: nil,
			},
			args: args{
				event: fake.GetTestStartedEvent(),
			},
			wantErr: false,
		},
		{
			name: "received started event with no matching triggered event",
			fields: fields{
				projectRepo: nil,
				eventRepo: &db_mock.EventRepoMock{
					GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
						if status[0] == common.TriggeredEvent {
							return nil, nil
						}
						return nil, errors.New("received unexpected request")
					},
					InsertEventFunc: func(project string, event models.Event, status common.EventStatus) error {
						t.Error("event should not be stored in this case")
						return nil
					},
					DeleteEventFunc: func(project string, eventID string, status common.EventStatus) error {
						return nil
					},
				},
				logger: nil,
			},
			args: args{
				event: fake.GetTestStartedEventWithUnmatchedTriggeredID(),
			},
			wantErr:                true,
			wantErrNoMatchingEvent: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := &shipyardController{
				projectRepo: tt.fields.projectRepo,
				eventRepo:   tt.fields.eventRepo,
				logger:      tt.fields.logger,
			}
			err := em.handleStartedEvent(tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleStartedEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErrNoMatchingEvent && (err != errNoMatchingEvent) {
				t.Errorf("handleStartedEvent() expected errNoMatchingEvent but got %v", err)
			}
		})
	}
}

func Test_eventManager_handleFinishedEvent(t *testing.T) {
	type fields struct {
		projectRepo db.ProjectRepo
		eventRepo   db.EventRepo
		logger      *keptncommon.Logger
	}
	type args struct {
		event models.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "received started event with no matching triggered event",
			fields: fields{
				projectRepo: nil,
				eventRepo: &db_mock.EventRepoMock{
					GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
						if status[0] == common.TriggeredEvent {
							return nil, nil
						} else if status[0] == common.StartedEvent {
							return nil, nil
						}
						return nil, errors.New("received unexpected request")
					},
					InsertEventFunc: func(project string, event models.Event, status common.EventStatus) error {
						return nil
					},
					DeleteEventFunc: func(project string, eventID string, status common.EventStatus) error {
						return nil
					},
				},
				logger: nil,
			},
			args: args{
				event: fake.GetTestFinishedEventWithUnmatchedSource(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := &shipyardController{
				projectRepo: tt.fields.projectRepo,
				eventRepo:   tt.fields.eventRepo,
				logger:      tt.fields.logger,
			}
			if err := em.handleFinishedEvent(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("handleFinishedEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_eventManager_getEvents(t *testing.T) {
	eventAvailable := false
	type fields struct {
		projectRepo db.ProjectRepo
		eventRepo   db.EventRepo
		logger      *keptncommon.Logger
	}
	type args struct {
		project string
		filter  common.EventFilter
		status  common.EventStatus
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Event
		wantErr bool
	}{
		{
			name: "get event",
			fields: fields{
				projectRepo: nil,
				eventRepo: &db_mock.EventRepoMock{
					GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
						return []models.Event{fake.GetTestTriggeredEvent()}, nil
					},
				},
				logger: keptncommon.NewLogger("", "", ""),
			},
			args: args{
				project: "test-project",
				filter:  common.EventFilter{},
				status:  common.TriggeredEvent,
			},
			want:    []models.Event{fake.GetTestTriggeredEvent()},
			wantErr: false,
		},
		{
			name: "get event after retry",
			fields: fields{
				projectRepo: nil,
				eventRepo: &db_mock.EventRepoMock{
					GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
						if eventAvailable {
							return []models.Event{fake.GetTestTriggeredEvent()}, nil
						}
						eventAvailable = true
						return nil, db.ErrNoEventFound
					},
				},
				logger: keptncommon.NewLogger("", "", ""),
			},
			args: args{
				project: "test-project",
				filter:  common.EventFilter{},
				status:  common.TriggeredEvent,
			},
			want:    []models.Event{fake.GetTestTriggeredEvent()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := &shipyardController{
				projectRepo: tt.fields.projectRepo,
				eventRepo:   tt.fields.eventRepo,
				logger:      tt.fields.logger,
			}
			got, err := em.getEvents(tt.args.project, tt.args.filter, tt.args.status, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEvents() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// Scenario 1: Complete task sequence execution + triggering of next task sequence. Events are received in order
func Test_shipyardController_Scenario1(t *testing.T) {

	t.Logf("Executing Shipyard Controller Scenario 1 with shipyard file %s", testShipyardFile)
	sc := getTestShipyardController()

	mockCS := fake.NewConfigurationService(testShipyardResource)
	defer mockCS.Close()

	done := false

	_ = os.Setenv("CONFIGURATION_SERVICE", mockCS.URL)

	mockEV := fake.NewEventBroker(t,
		func(meb *fake.EventBroker, event *models.Event) {
			meb.ReceivedEvents = append(meb.ReceivedEvents, *event)
		},
		func(meb *fake.EventBroker) {

		})
	defer mockEV.Server.Close()
	_ = os.Setenv("EVENTBROKER", mockEV.Server.URL)

	// STEP 1
	// send dev.artifact-delivery.triggered event
	sequenceTriggeredEvent := getArtifactDeliveryTriggeredEvent()
	err := sc.HandleIncomingEvent(sequenceTriggeredEvent)
	if err != nil {
		t.Errorf("STEP 1 failed: HandleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	// check event broker -> should contain deployment.triggered event with properties: [deployment]
	if len(mockEV.ReceivedEvents) != 1 {
		t.Errorf("STEP 1 failed: expected %d events in eventbroker, but got %d", 1, len(mockEV.ReceivedEvents))
		return
	}
	done = fake.ShouldContainEvent(t,
		mockEV.ReceivedEvents,
		keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		"",
		func(t *testing.T, event models.Event) bool {
			deploymentEvent := &keptnv2.DeploymentTriggeredEventData{}

			marshal, _ := json.Marshal(event.Data)
			if err := json.Unmarshal(marshal, deploymentEvent); err != nil {
				t.Error("could not parse incoming deployment.triggered event: " + err.Error())
				return true
			}

			if len(deploymentEvent.Deployment.DeploymentURIsPublic) > 1 {
				t.Errorf("ohoh")
				return true
			}

			if deploymentEvent.Deployment.DeploymentStrategy != "direct" {
				t.Errorf("did not receive correct deployment strategy. Expected 'direct' but got '%s'", deploymentEvent.Deployment.DeploymentStrategy)
				return true
			}
			if deploymentEvent.ConfigurationChange.Values["image"] != "carts" {
				t.Errorf("did not receive correct image. Expected 'carts' but got '%s'", deploymentEvent.ConfigurationChange.Values["image"])
				return true
			}
			return false
		},
	)
	if done {
		return
	}

	// check triggeredEvent Collection -> should contain deployment.triggered event
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		Stage:   common.Stringp("dev"),
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)
	done = fake.ShouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
	if done {
		return
	}
	triggeredID := triggeredEvents[0].ID

	// STEP 2
	// send deployment.started event
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.DeploymentTaskName, triggeredID, "dev", "test-source")
	if done {
		return
	}

	// STEP 3
	// send deployment.finished event
	triggeredID, done = sendAndVerifyFinishedEvent(
		t,
		sc,
		getDeploymentFinishedEvent("dev", triggeredID, "test-source", keptnv2.ResultPass),
		keptnv2.DeploymentTaskName,
		keptnv2.TestTaskName,
		mockEV,
		"",
		func(t *testing.T, e models.Event) bool {
			marshal, _ := json.Marshal(e.Data)
			testData := &keptnv2.TestTriggeredEventData{}

			err := json.Unmarshal(marshal, testData)

			if err != nil {
				t.Errorf("Expected test.triggered data but could not convert: %v: %s", e.Data, err.Error())
				return true
			}

			if len(testData.Deployment.DeploymentURIsLocal) != 2 {
				t.Errorf("DeploymentURIsLocal property was not transmitted correctly")
				return true
			}
			if len(testData.Deployment.DeploymentURIsPublic) != 3 {
				t.Errorf("DeploymentURIsLocal property was not transmitted correctly")
				return true
			}

			// also check if the payload of the .triggered event that started the sequence is present

			deploymentEvent := &keptnv2.DeploymentTriggeredEventData{}

			err = json.Unmarshal(marshal, deploymentEvent)
			if err != nil {
				t.Errorf("Expected deployment.triggered data but could not convert: %v: %s", e.Data, err.Error())
				return true
			}

			if deploymentEvent.ConfigurationChange.Values["image"] != "carts" {
				t.Errorf("did not receive correct image. Expected 'carts' but got '%s'", deploymentEvent.ConfigurationChange.Values["image"])
				return true
			}

			// check if the message property of the previous .finished event has been reset correctly
			if deploymentEvent.Message != "" {
				t.Errorf("expected message property to be empty, but got: '%s'", deploymentEvent.Message)
				return true
			}
			return false
		})
	if done {
		return
	}

	// STEP 4
	// send test.started event
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.TestTaskName, triggeredID, "dev", "test-source")
	if done {
		return
	}

	// STEP 5
	// send test.finished event
	triggeredID, done = sendAndVerifyFinishedEvent(
		t,
		sc,
		getTestTaskFinishedEvent("dev", triggeredID),
		keptnv2.TestTaskName,
		keptnv2.EvaluationTaskName,
		mockEV,
		"",
		func(t *testing.T, e models.Event) bool {
			marshal, _ := json.Marshal(e.Data)
			testData := &keptnv2.EvaluationTriggeredEventData{}

			err := json.Unmarshal(marshal, testData)

			if err != nil {
				t.Errorf("Expected test.triggered data but could not convert: %v: %s", e.Data, err.Error())
				return true
			}

			if len(testData.Deployment.DeploymentNames) != 1 {
				t.Errorf("Deployment property was not transmitted correctly: %v", testData.Deployment)
				return true
			}
			if testData.Test.Start != "start" || testData.Test.End != "end" {
				t.Errorf("Test property was not transmitted correctly: %v", testData.Test)
			}
			return false
		})
	if done {
		return
	}
	// STEP 6
	// send evaluation.started event
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.EvaluationTaskName, triggeredID, "dev", "test-source")
	if done {
		return
	}

	// STEP 7
	// send evaluation.finished event -> result = warning should not abort the task sequence
	triggeredID, done = sendAndVerifyFinishedEvent(t, sc, getEvaluationTaskFinishedEvent("dev", triggeredID, keptnv2.ResultWarning, keptnv2.StatusSucceeded), keptnv2.EvaluationTaskName, keptnv2.ReleaseTaskName, mockEV, "", nil)
	if done {
		return
	}

	// STEP 8
	// send release.started event
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.ReleaseTaskName, triggeredID, "dev", "test-source")
	if done {
		return
	}

	// STEP 9
	// send release.finished event
	triggeredID, done = sendAndVerifyFinishedEvent(t, sc, getReleaseTaskFinishedEvent("dev", triggeredID), keptnv2.ReleaseTaskName, keptnv2.DeploymentTaskName, mockEV, "hardening", nil)
	if done {
		return
	}

	// check if dev.artifact-delivery.finished has been sent
	done = fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetFinishedEventType("dev.artifact-delivery"), "dev", func(t *testing.T, event models.Event) bool {
		if event.Triggeredid != sequenceTriggeredEvent.ID {
			t.Error("expected triggeredid dev.artifact-delivery.finished to match ID of dev.artifact-delivery.triggered event")
			return true
		}
		return false
	})
	if done {
		return
	}

	done = fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetTriggeredEventType("hardening.artifact-delivery"), "hardening", func(t *testing.T, event models.Event) bool {
		marshal, _ := json.Marshal(event.Data)
		triggeredEvent := map[string]interface{}{}

		err := json.Unmarshal(marshal, &triggeredEvent)

		if err != nil {
			t.Errorf("Expected hardening.artifact-delivery.triggered data but could not convert: %v: %s", event.Data, err.Error())
			return true
		}

		if triggeredEvent["configurationChange"] == nil {
			t.Error("expected 'configurationChange' property to be present")
			return true
		}
		return false
	})
	if done {
		return
	}

	done = fake.ShouldContainEvent(
		t,
		mockEV.ReceivedEvents,
		keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		"hardening",
		func(t *testing.T, event models.Event) bool {
			marshal, _ := json.Marshal(event.Data)
			deploymentEvent := &keptnv2.DeploymentTriggeredEventData{}

			err := json.Unmarshal(marshal, deploymentEvent)

			if err != nil {
				t.Errorf("Expected test.triggered data but could not convert: %v: %s", event.Data, err.Error())
				return true
			}

			if deploymentEvent.ConfigurationChange.Values["image"] != "carts" {
				t.Errorf("did not receive correct image. Expected 'carts' but got '%s'", deploymentEvent.ConfigurationChange.Values["image"])
				return true
			}
			return false
		},
	)
	if done {
		return
	}

	finishedEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Stage: common.Stringp("dev"),
	}, common.FinishedEvent)

	fake.ShouldNotContainEvent(t, finishedEvents, keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName), "dev")
	fake.ShouldNotContainEvent(t, finishedEvents, keptnv2.GetFinishedEventType(keptnv2.TestTaskName), "dev")
	fake.ShouldNotContainEvent(t, finishedEvents, keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), "dev")
	fake.ShouldNotContainEvent(t, finishedEvents, keptnv2.GetFinishedEventType(keptnv2.ReleaseTaskName), "dev")

	// STEP 9.1
	// send deployment.started event 1 with ID 1
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.DeploymentTaskName, triggeredID, "hardening", "test-source-1")
	if done {
		return
	}

	// STEP 9.2
	// send deployment.started event 2 with ID 2
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.DeploymentTaskName, triggeredID, "hardening", "test-source-2")
	if done {
		return
	}

	// STEP 10.1
	// send deployment.finished event 1 with ID 1
	done = sendAndVerifyPartialFinishedEvent(t, sc, getDeploymentFinishedEvent("hardening", triggeredID, "test-source-1", keptnv2.ResultPass), keptnv2.DeploymentTaskName, keptnv2.ReleaseTaskName, mockEV, "")
	if done {
		return
	}

	// STEP 10.2
	// send deployment.finished event 1 with ID 1
	triggeredID, done = sendAndVerifyFinishedEvent(t, sc, getDeploymentFinishedEvent("hardening", triggeredID, "test-source-2", keptnv2.ResultPass), keptnv2.DeploymentTaskName, keptnv2.TestTaskName, mockEV, "", nil)
	if done {
		return
	}

	eventsDBMock := sc.eventsDbOperations.(*db_mock.EventsDbOperationsMock)
	// make sure that the UpdateEventOfServiceCalls has been called
	assert.NotEqual(t, 0, len(eventsDBMock.UpdateEventOfServiceCalls()))
	assert.NotEqual(t, 0, len(eventsDBMock.UpdateShipyardCalls()))
	assert.Equal(t, "test-project", eventsDBMock.UpdateShipyardCalls()[0].ProjectName)
	assert.NotEqual(t, "", eventsDBMock.UpdateShipyardCalls()[0].ShipyardContent)
}

// Scenario 2: Partial task sequence execution + triggering of next task sequence. Events are received out of order
func Test_shipyardController_Scenario2(t *testing.T) {

	t.Logf("Executing Shipyard Controller Scenario 1 with shipyard file %s", testShipyardFile)
	sc := getTestShipyardController()

	mockCS := fake.NewConfigurationService(testShipyardResource)
	defer mockCS.Close()

	done := false

	_ = os.Setenv("CONFIGURATION_SERVICE", mockCS.URL)

	mockEV := fake.NewEventBroker(t,
		func(meb *fake.EventBroker, event *models.Event) {
			meb.ReceivedEvents = append(meb.ReceivedEvents, *event)
		},
		func(meb *fake.EventBroker) {

		})
	defer mockEV.Server.Close()
	_ = os.Setenv("EVENTBROKER", mockEV.Server.URL)

	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.HandleIncomingEvent(getArtifactDeliveryTriggeredEvent())
	if err != nil {
		t.Errorf("STEP 1 failed: HandleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	// check event broker -> should contain deployment.triggered event with properties: [deployment]
	if len(mockEV.ReceivedEvents) != 1 {
		t.Errorf("STEP 1 failed: expected %d events in eventbroker, but got %d", 1, len(mockEV.ReceivedEvents))
		return
	}
	done = fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
	if done {
		return
	}
	// check triggeredEvent Collection -> should contain deployment.triggered event
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		Stage:   common.Stringp("dev"),
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)
	done = fake.ShouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
	if done {
		return
	}
	triggeredID := triggeredEvents[0].ID

	// STEP 2
	// send deployment.started event
	go func() {
		time.After(2 * time.Second)
		_ = sendAndVerifyStartedEvent(t, sc, keptnv2.DeploymentTaskName, triggeredID, "dev", "test-source")
	}()

	// STEP 3
	// send deployment.finished event
	triggeredID, done = sendAndVerifyFinishedEvent(
		t,
		sc,
		getDeploymentFinishedEvent("dev", triggeredID, "test-source", keptnv2.ResultPass),
		keptnv2.DeploymentTaskName,
		keptnv2.TestTaskName,
		mockEV,
		"",
		func(t *testing.T, e models.Event) bool {
			marshal, _ := json.Marshal(e.Data)
			testData := &keptnv2.TestTriggeredEventData{}

			err := json.Unmarshal(marshal, testData)

			if err != nil {
				t.Errorf("Expected test.triggered data but could not convert: %v: %s", e.Data, err.Error())
				return true
			}

			if len(testData.Deployment.DeploymentURIsLocal) != 2 {
				t.Errorf("DeploymentURIsLocal property was not transmitted correctly")
				return true
			}
			if len(testData.Deployment.DeploymentURIsPublic) != 3 {
				t.Errorf("DeploymentURIsLocal property was not transmitted correctly")
				return true
			}
			return false
		})
	if done {
		return
	}
}

// Scenario 3: Received .finished event with status "errored" should abort task sequence and send .finished event with status "errored"
func Test_shipyardController_Scenario3(t *testing.T) {

	t.Logf("Executing Shipyard Controller Scenario 1 with shipyard file %s", testShipyardFile)
	sc := getTestShipyardController()

	mockCS := fake.NewConfigurationService(testShipyardResource)
	defer mockCS.Close()

	done := false

	_ = os.Setenv("CONFIGURATION_SERVICE", mockCS.URL)

	mockEV := fake.NewEventBroker(t,
		func(meb *fake.EventBroker, event *models.Event) {
			meb.ReceivedEvents = append(meb.ReceivedEvents, *event)
		},
		func(meb *fake.EventBroker) {

		})
	defer mockEV.Server.Close()
	_ = os.Setenv("EVENTBROKER", mockEV.Server.URL)

	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.HandleIncomingEvent(getArtifactDeliveryTriggeredEvent())
	if err != nil {
		t.Errorf("STEP 1 failed: HandleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	// check event broker -> should contain deployment.triggered event with properties: [deployment]
	if len(mockEV.ReceivedEvents) != 1 {
		t.Errorf("STEP 1 failed: expected %d events in eventbroker, but got %d", 1, len(mockEV.ReceivedEvents))
		return
	}
	done = fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
	if done {
		return
	}
	// check triggeredEvent Collection -> should contain deployment.triggered event
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		Stage:   common.Stringp("dev"),
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)
	done = fake.ShouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
	if done {
		return
	}
	triggeredID := triggeredEvents[0].ID

	// STEP 2
	// send deployment.started event
	go func() {
		time.After(2 * time.Second)
		_ = sendAndVerifyStartedEvent(t, sc, keptnv2.DeploymentTaskName, triggeredID, "dev", "test-source")
	}()

	// STEP 3
	// send deployment.finished event
	done = sendFinishedEventAndVerifyTaskSequenceCompletion(
		t,
		sc,
		getErroredDeploymentFinishedEvent("dev", triggeredID, "test-source"),
		keptnv2.DeploymentTaskName,
		"dev.artifact-delivery",
		mockEV,
		"",
		func(t *testing.T, e models.Event) bool {
			marshal, _ := json.Marshal(e.Data)
			eventData := &keptnv2.EventData{}

			err := json.Unmarshal(marshal, eventData)

			if err != nil {
				t.Errorf("could not convert event data: %v: %s", e.Data, err.Error())
				return true
			}

			if eventData.Status != keptnv2.StatusErrored {
				t.Errorf("Expected Status %s, but got %s", keptnv2.StatusErrored, eventData.Status)
				return true
			}
			if eventData.Result != keptnv2.ResultFailed {
				t.Errorf("Expected Result %s, but got %s", keptnv2.ResultFailed, eventData.Result)
				return true
			}
			return false
		})
	if done {
		return
	}
}

// Scenario 4: Received .finished event with result "fail" - stop task sequence
func Test_shipyardController_Scenario4(t *testing.T) {

	t.Logf("Executing Shipyard Controller Scenario 1 with shipyard file %s", testShipyardFile)
	sc := getTestShipyardController()

	mockCS := fake.NewConfigurationService(testShipyardResource)
	defer mockCS.Close()

	done := false

	_ = os.Setenv("CONFIGURATION_SERVICE", mockCS.URL)

	mockEV := fake.NewEventBroker(t,
		func(meb *fake.EventBroker, event *models.Event) {
			meb.ReceivedEvents = append(meb.ReceivedEvents, *event)
		},
		func(meb *fake.EventBroker) {

		})
	defer mockEV.Server.Close()
	_ = os.Setenv("EVENTBROKER", mockEV.Server.URL)

	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.HandleIncomingEvent(getArtifactDeliveryTriggeredEvent())
	if err != nil {
		t.Errorf("STEP 1 failed: HandleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	// check event broker -> should contain deployment.triggered event with properties: [deployment]
	if len(mockEV.ReceivedEvents) != 1 {
		t.Errorf("STEP 1 failed: expected %d events in eventbroker, but got %d", 1, len(mockEV.ReceivedEvents))
		return
	}
	done = fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
	if done {
		return
	}
	// check triggeredEvent Collection -> should contain deployment.triggered event
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		Stage:   common.Stringp("dev"),
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)
	done = fake.ShouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
	if done {
		return
	}
	triggeredID := triggeredEvents[0].ID

	// STEP 2
	// send deployment.started event
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.DeploymentTaskName, triggeredID, "dev", "test-source")
	if done {
		return
	}

	// STEP 3
	// send deployment.finished event
	triggeredID, done = sendAndVerifyFinishedEvent(
		t,
		sc,
		getDeploymentFinishedEvent("dev", triggeredID, "test-source", keptnv2.ResultPass),
		keptnv2.DeploymentTaskName,
		keptnv2.TestTaskName,
		mockEV,
		"",
		func(t *testing.T, e models.Event) bool {
			marshal, _ := json.Marshal(e.Data)
			testData := &keptnv2.TestTriggeredEventData{}

			err := json.Unmarshal(marshal, testData)

			if err != nil {
				t.Errorf("Expected test.triggered data but could not convert: %v: %s", e.Data, err.Error())
				return true
			}

			if len(testData.Deployment.DeploymentURIsLocal) != 2 {
				t.Errorf("DeploymentURIsLocal property was not transmitted correctly")
				return true
			}
			if len(testData.Deployment.DeploymentURIsPublic) != 3 {
				t.Errorf("DeploymentURIsLocal property was not transmitted correctly")
				return true
			}
			return false
		})
	if done {
		return
	}

	// STEP 4
	// send test.started event
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.TestTaskName, triggeredID, "dev", "test-source")
	if done {
		return
	}

	// STEP 5
	// send test.finished event
	triggeredID, done = sendAndVerifyFinishedEvent(
		t,
		sc,
		getTestTaskFinishedEvent("dev", triggeredID),
		keptnv2.TestTaskName,
		keptnv2.EvaluationTaskName,
		mockEV,
		"",
		func(t *testing.T, e models.Event) bool {
			marshal, _ := json.Marshal(e.Data)
			testData := &keptnv2.EvaluationTriggeredEventData{}

			err := json.Unmarshal(marshal, testData)

			if err != nil {
				t.Errorf("Expected test.triggered data but could not convert: %v: %s", e.Data, err.Error())
				return true
			}

			if len(testData.Deployment.DeploymentNames) != 1 {
				t.Errorf("Deployment property was not transmitted correctly: %v", testData.Deployment)
				return true
			}
			if testData.Test.Start != "start" || testData.Test.End != "end" {
				t.Errorf("Test property was not transmitted correctly: %v", testData.Test)
			}
			return false
		})
	if done {
		return
	}
	// STEP 6
	// send evaluation.started event
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.EvaluationTaskName, triggeredID, "dev", "test-source")
	if done {
		return
	}

	// STEP 7
	// send evaluation.finished event with result=fail

	done = sendFinishedEventAndVerifyTaskSequenceCompletion(
		t,
		sc,
		getEvaluationTaskFinishedEvent("dev", triggeredID, keptnv2.ResultFailed, keptnv2.StatusSucceeded),
		keptnv2.EvaluationTaskName,
		"dev.artifact-delivery",
		mockEV,
		"",
		func(t *testing.T, e models.Event) bool {
			marshal, _ := json.Marshal(e.Data)
			eventData := &keptnv2.EventData{}

			err := json.Unmarshal(marshal, eventData)

			if err != nil {
				t.Errorf("could not convert event data: %v: %s", e.Data, err.Error())
				return true
			}

			if eventData.Status != keptnv2.StatusSucceeded {
				t.Errorf("Expected Status %s, but got %s", keptnv2.StatusSucceeded, eventData.Status)
				return true
			}
			if eventData.Result != keptnv2.ResultFailed {
				t.Errorf("Expected Result %s, but got %s", keptnv2.ResultFailed, eventData.Result)
				return true
			}
			return false
		})
	if done {
		return
	}

}

// Scenario 4: Received .finished event with result "fail" - stop task sequence and trigger next sequence based on result filter
func Test_shipyardController_TriggerOnFail(t *testing.T) {

	t.Logf("Executing Shipyard Controller with shipyard file %s", testShipyardFile)
	sc := getTestShipyardController()

	mockCS := fake.NewConfigurationService(testShipyardResource)
	defer mockCS.Close()

	done := false

	_ = os.Setenv("CONFIGURATION_SERVICE", mockCS.URL)

	mockEV := fake.NewEventBroker(t,
		func(meb *fake.EventBroker, event *models.Event) {
			meb.ReceivedEvents = append(meb.ReceivedEvents, *event)
		},
		func(meb *fake.EventBroker) {

		})
	defer mockEV.Server.Close()
	_ = os.Setenv("EVENTBROKER", mockEV.Server.URL)

	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.HandleIncomingEvent(getArtifactDeliveryTriggeredEvent())
	if err != nil {
		t.Errorf("STEP 1 failed: HandleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	// check event broker -> should contain deployment.triggered event with properties: [deployment]
	if len(mockEV.ReceivedEvents) != 1 {
		t.Errorf("STEP 1 failed: expected %d events in eventbroker, but got %d", 1, len(mockEV.ReceivedEvents))
		return
	}
	done = fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
	if done {
		return
	}
	// check triggeredEvent Collection -> should contain deployment.triggered event
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		Stage:   common.Stringp("dev"),
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)
	done = fake.ShouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
	if done {
		return
	}
	triggeredID := triggeredEvents[0].ID

	// STEP 2
	// send deployment.started event
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.DeploymentTaskName, triggeredID, "dev", "test-source")
	if done {
		return
	}

	// STEP 3
	// send deployment.finished event
	done = sendFinishedEventAndVerifyTaskSequenceCompletion(
		t,
		sc,
		getDeploymentFinishedEvent("dev", triggeredID, "test-source", keptnv2.ResultFailed),
		keptnv2.DeploymentTaskName,
		"dev.artifact-delivery",
		mockEV,
		"",
		func(t *testing.T, e models.Event) bool {
			marshal, _ := json.Marshal(e.Data)
			eventData := &keptnv2.EventData{}

			err := json.Unmarshal(marshal, eventData)

			if err != nil {
				t.Errorf("could not convert event data: %v: %s", e.Data, err.Error())
				return true
			}

			if eventData.Status != keptnv2.StatusSucceeded {
				t.Errorf("Expected Status %s, but got %s", keptnv2.StatusSucceeded, eventData.Status)
				return true
			}
			if eventData.Result != keptnv2.ResultFailed {
				t.Errorf("Expected Result %s, but got %s", keptnv2.ResultFailed, eventData.Result)
				return true
			}
			return false
		})
	if done {
		return
	}

	done = fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetTriggeredEventType("dev.rollback"), "dev", func(t *testing.T, event models.Event) bool {
		marshal, _ := json.Marshal(event.Data)
		triggeredEvent := map[string]interface{}{}

		err := json.Unmarshal(marshal, &triggeredEvent)

		if err != nil {
			t.Errorf("Expected dev.rollback.triggered data but could not convert: %v: %s", event.Data, err.Error())
			return true
		}
		return false
	})
	if done {
		return
	}

	done = fake.ShouldContainEvent(
		t,
		mockEV.ReceivedEvents,
		keptnv2.GetTriggeredEventType(keptnv2.RollbackTaskName),
		"dev",
		func(t *testing.T, event models.Event) bool {
			marshal, _ := json.Marshal(event.Data)
			deploymentEvent := &keptnv2.EventData{}

			err := json.Unmarshal(marshal, deploymentEvent)

			if err != nil {
				t.Errorf("Expected rollback.triggered data but could not convert: %v: %s", event.Data, err.Error())
				return true
			}
			return false
		},
	)
	if done {
		return
	}

	// hardening.artifact-delivery should not be triggered
	fake.ShouldNotContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetTriggeredEventType("hardening.artifact-delivery"), "hardening")
	// hardening.deployment should not be triggered
	fake.ShouldNotContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "hardening")

}

// Scenario 5: Received .triggered event for project with invalid shipyard version -> send .finished event with result = fail
func Test_shipyardController_Scenario5(t *testing.T) {

	t.Logf("Executing Shipyard Controller Scenario 5 with shipyard file %s", testShipyardFileWithInvalidVersion)
	sc := getTestShipyardController()

	mockCS := fake.NewConfigurationService(testShipyardResourceWithInvalidVersion)
	defer mockCS.Close()

	_ = os.Setenv("CONFIGURATION_SERVICE", mockCS.URL)

	mockEV := fake.NewEventBroker(t,
		func(meb *fake.EventBroker, event *models.Event) {
			meb.ReceivedEvents = append(meb.ReceivedEvents, *event)
		},
		func(meb *fake.EventBroker) {

		})
	defer mockEV.Server.Close()
	_ = os.Setenv("EVENTBROKER", mockEV.Server.URL)

	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.HandleIncomingEvent(getArtifactDeliveryTriggeredEvent())
	if err != nil {
		t.Errorf("STEP 1 failed: HandleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	// check event broker -> should contain deployment.triggered event with properties: [deployment]
	if len(mockEV.ReceivedEvents) != 1 {
		t.Errorf("STEP 1 failed: expected %d events in eventbroker, but got %d", 1, len(mockEV.ReceivedEvents))
		return
	}
	fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetFinishedEventType("dev.artifact-delivery"), "", nil)
}

// Updating shipyard content fails -> event handling should still happen
func Test_shipyardController_UpdateShipyardContentFails(t *testing.T) {

	t.Logf("Executing Shipyard Controller with shipyard file %s", testShipyardFileWithInvalidVersion)
	sc := getTestShipyardController()

	eventsOperations := sc.eventsDbOperations.(*db_mock.EventsDbOperationsMock)

	eventsOperations.UpdateShipyardFunc = func(projectName string, shipyardContent string) error {
		return errors.New("updating shipyard failed")
	}

	mockCS := fake.NewConfigurationService(testShipyardResourceWithInvalidVersion)
	defer mockCS.Close()

	_ = os.Setenv("CONFIGURATION_SERVICE", mockCS.URL)

	mockEV := fake.NewEventBroker(t,
		func(meb *fake.EventBroker, event *models.Event) {
			meb.ReceivedEvents = append(meb.ReceivedEvents, *event)
		},
		func(meb *fake.EventBroker) {

		})
	defer mockEV.Server.Close()
	_ = os.Setenv("EVENTBROKER", mockEV.Server.URL)

	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.HandleIncomingEvent(getArtifactDeliveryTriggeredEvent())
	if err != nil {
		t.Errorf("STEP 1 failed: HandleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	// check event broker -> should contain deployment.triggered event with properties: [deployment]
	if len(mockEV.ReceivedEvents) != 1 {
		t.Errorf("STEP 1 failed: expected %d events in eventbroker, but got %d", 1, len(mockEV.ReceivedEvents))
		return
	}
	fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetFinishedEventType("dev.artifact-delivery"), "", nil)
}

// Updating event of service fails -> event handling should still happen
func Test_shipyardController_UpdateEventOfServiceFailsFails(t *testing.T) {

	t.Logf("Executing Shipyard Controller with shipyard file %s", testShipyardFileWithInvalidVersion)
	sc := getTestShipyardController()

	eventsOperations := sc.eventsDbOperations.(*db_mock.EventsDbOperationsMock)

	eventsOperations.UpdateEventOfServiceFunc = func(event interface{}, eventType string, keptnContext string, eventID string, triggeredID string) error {
		return errors.New("updating event of service failed")
	}

	mockCS := fake.NewConfigurationService(testShipyardResourceWithInvalidVersion)
	defer mockCS.Close()

	_ = os.Setenv("CONFIGURATION_SERVICE", mockCS.URL)

	mockEV := fake.NewEventBroker(t,
		func(meb *fake.EventBroker, event *models.Event) {
			meb.ReceivedEvents = append(meb.ReceivedEvents, *event)
		},
		func(meb *fake.EventBroker) {

		})
	defer mockEV.Server.Close()
	_ = os.Setenv("EVENTBROKER", mockEV.Server.URL)

	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.HandleIncomingEvent(getArtifactDeliveryTriggeredEvent())
	if err != nil {
		t.Errorf("STEP 1 failed: HandleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	// check event broker -> should contain deployment.triggered event with properties: [deployment]
	if len(mockEV.ReceivedEvents) != 1 {
		t.Errorf("STEP 1 failed: expected %d events in eventbroker, but got %d", 1, len(mockEV.ReceivedEvents))
		return
	}
	fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetFinishedEventType("dev.artifact-delivery"), "", nil)
}

// Scenario 5: Received .triggered event for project with invalid shipyard version -> send .finished event with result = fail
func Test_shipyardController_UpdateServiceShouldNotBeCalledForEmptyService(t *testing.T) {

	t.Logf("Executing Shipyard Controller with shipyard file %s", testShipyardFileWithInvalidVersion)
	sc := getTestShipyardController()

	mockCS := fake.NewConfigurationService(testShipyardResourceWithInvalidVersion)
	defer mockCS.Close()

	_ = os.Setenv("CONFIGURATION_SERVICE", mockCS.URL)

	mockEV := fake.NewEventBroker(t,
		func(meb *fake.EventBroker, event *models.Event) {
			meb.ReceivedEvents = append(meb.ReceivedEvents, *event)
		},
		func(meb *fake.EventBroker) {

		})
	defer mockEV.Server.Close()
	_ = os.Setenv("EVENTBROKER", mockEV.Server.URL)

	event := getArtifactDeliveryTriggeredEvent()

	event.Data = keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "",
	}
	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.HandleIncomingEvent(event)

	assert.NotNil(t, err)

	eventsDBMock := sc.eventsDbOperations.(*db_mock.EventsDbOperationsMock)

	assert.Equal(t, 0, len(eventsDBMock.UpdateEventOfServiceCalls()))
}

func Test_shipyardController_getTaskSequenceInStage(t *testing.T) {
	type fields struct {
		projectRepo      db.ProjectRepo
		eventRepo        db.EventRepo
		taskSequenceRepo db.TaskSequenceRepo
		logger           *keptncommon.Logger
	}
	type args struct {
		stageName        string
		taskSequenceName string
		shipyard         *keptnv2.Shipyard
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *keptnv2.Sequence
		wantErr bool
	}{
		{
			name: "get built-in evaluation task sequence",
			fields: fields{
				projectRepo:      nil,
				eventRepo:        nil,
				taskSequenceRepo: nil,
				logger:           keptncommon.NewLogger("", "", ""),
			},
			args: args{
				stageName:        "dev",
				taskSequenceName: "evaluation",
				shipyard: &keptnv2.Shipyard{
					ApiVersion: "0.2.0",
					Kind:       "shipyard",
					Metadata:   keptnv2.Metadata{},
					Spec: keptnv2.ShipyardSpec{
						Stages: []keptnv2.Stage{
							{
								Name:      "dev",
								Sequences: []keptnv2.Sequence{},
							},
						},
					},
				},
			},
			want: &keptnv2.Sequence{
				Name:        "evaluation",
				TriggeredOn: nil,
				Tasks: []keptnv2.Task{
					{
						Name:       "evaluation",
						Properties: nil,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "get user-defined evaluation task sequence",
			fields: fields{
				projectRepo:      nil,
				eventRepo:        nil,
				taskSequenceRepo: nil,
				logger:           keptncommon.NewLogger("", "", ""),
			},
			args: args{
				stageName:        "dev",
				taskSequenceName: "evaluation",
				shipyard: &keptnv2.Shipyard{
					ApiVersion: "0.2.0",
					Kind:       "shipyard",
					Metadata:   keptnv2.Metadata{},
					Spec: keptnv2.ShipyardSpec{
						Stages: []keptnv2.Stage{
							{
								Name: "dev",
								Sequences: []keptnv2.Sequence{
									{
										Name:        "evaluation",
										TriggeredOn: nil,
										Tasks: []keptnv2.Task{
											{
												Name:       "evaluation",
												Properties: nil,
											},
											{
												Name:       "notify",
												Properties: nil,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: &keptnv2.Sequence{
				Name:        "evaluation",
				TriggeredOn: nil,
				Tasks: []keptnv2.Task{
					{
						Name:       "evaluation",
						Properties: nil,
					},
					{
						Name:       "notify",
						Properties: nil,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &shipyardController{
				projectRepo:      tt.fields.projectRepo,
				eventRepo:        tt.fields.eventRepo,
				taskSequenceRepo: tt.fields.taskSequenceRepo,
				logger:           tt.fields.logger,
			}
			got, err := sc.getTaskSequenceInStage(tt.args.stageName, tt.args.taskSequenceName, tt.args.shipyard)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTaskSequenceInStage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("getTaskSequenceInStage() got = %v, want %v", got, tt.want)
				for _, d := range diff {
					t.Log(d)
				}
			}
		})
	}
}

func Test_getTaskSequencesByTrigger(t *testing.T) {
	type args struct {
		eventScope            *keptnv2.EventData
		completedTaskSequence string
		shipyard              *keptnv2.Shipyard
	}
	tests := []struct {
		name string
		args args
		want []NextTaskSequence
	}{
		{
			name: "default behavior - get sequence triggered by result=pass,warning",
			args: args{
				eventScope: &keptnv2.EventData{
					Result: keptnv2.ResultPass,
					Stage:  "dev",
				},
				completedTaskSequence: "artifact-delivery",
				shipyard: &keptnv2.Shipyard{
					ApiVersion: shipyardVersion,
					Kind:       "shipyard",
					Metadata:   keptnv2.Metadata{},
					Spec: keptnv2.ShipyardSpec{
						Stages: []keptnv2.Stage{
							{
								Name: "dev",
								Sequences: []keptnv2.Sequence{
									{
										Name:        "artifact-delivery",
										TriggeredOn: nil,
										Tasks:       nil,
									},
								},
							},
							{
								Name: "hardening",
								Sequences: []keptnv2.Sequence{
									{
										Name: "artifact-delivery",
										TriggeredOn: []keptnv2.Trigger{
											{
												Event:    "dev.artifact-delivery.finished",
												Selector: keptnv2.Selector{},
											},
										},
										Tasks: nil,
									},
									{
										Name:        "artifact-delivery-2",
										TriggeredOn: nil,
										Tasks:       nil,
									},
								},
							},
						},
					},
				},
			},
			want: []NextTaskSequence{
				{
					Sequence: keptnv2.Sequence{
						Name: "artifact-delivery",
						TriggeredOn: []keptnv2.Trigger{
							{
								Event:    "dev.artifact-delivery.finished",
								Selector: keptnv2.Selector{},
							},
						},
						Tasks: nil,
					},
					StageName: "hardening",
				},
			},
		},
		{
			name: "get sequence triggered by result=fail",
			args: args{
				eventScope: &keptnv2.EventData{
					Result: keptnv2.ResultFailed,
					Stage:  "dev",
				},
				completedTaskSequence: "artifact-delivery",
				shipyard: &keptnv2.Shipyard{
					ApiVersion: shipyardVersion,
					Kind:       "shipyard",
					Metadata:   keptnv2.Metadata{},
					Spec: keptnv2.ShipyardSpec{
						Stages: []keptnv2.Stage{
							{
								Name: "dev",
								Sequences: []keptnv2.Sequence{
									{
										Name:        "artifact-delivery",
										TriggeredOn: nil,
										Tasks:       nil,
									},
								},
							},
							{
								Name: "hardening",
								Sequences: []keptnv2.Sequence{
									{
										Name: "artifact-delivery",
										TriggeredOn: []keptnv2.Trigger{
											{
												Event:    "dev.artifact-delivery.finished",
												Selector: keptnv2.Selector{},
											},
										},
										Tasks: nil,
									},
									{
										Name: "artifact-delivery-2",
										TriggeredOn: []keptnv2.Trigger{
											{
												Event: "dev.artifact-delivery.finished",
												Selector: keptnv2.Selector{
													Match: map[string]string{
														"result": string(keptnv2.ResultFailed),
													},
												},
											},
										},
										Tasks: nil,
									},
								},
							},
							{
								Name: "production",
								Sequences: []keptnv2.Sequence{
									{
										Name: "artifact-delivery",
										TriggeredOn: []keptnv2.Trigger{
											{
												Event:    "dev.artifact-delivery.finished",
												Selector: keptnv2.Selector{},
											},
										},
										Tasks: nil,
									},
									{
										Name: "artifact-delivery-2",
										TriggeredOn: []keptnv2.Trigger{
											{
												Event: "dev.artifact-delivery.finished",
												Selector: keptnv2.Selector{
													Match: map[string]string{
														"result": string(keptnv2.ResultFailed),
													},
												},
											},
										},
										Tasks: nil,
									},
								},
							},
						},
					},
				},
			},
			want: []NextTaskSequence{
				{
					Sequence: keptnv2.Sequence{
						Name: "artifact-delivery-2",
						TriggeredOn: []keptnv2.Trigger{
							{
								Event: "dev.artifact-delivery.finished",
								Selector: keptnv2.Selector{
									Match: map[string]string{
										"result": string(keptnv2.ResultFailed),
									},
								},
							},
						},
						Tasks: nil,
					},
					StageName: "hardening",
				},
				{
					Sequence: keptnv2.Sequence{
						Name: "artifact-delivery-2",
						TriggeredOn: []keptnv2.Trigger{
							{
								Event: "dev.artifact-delivery.finished",
								Selector: keptnv2.Selector{
									Match: map[string]string{
										"result": string(keptnv2.ResultFailed),
									},
								},
							},
						},
						Tasks: nil,
					},
					StageName: "production",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTaskSequencesByTrigger(tt.args.eventScope, tt.args.completedTaskSequence, tt.args.shipyard); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTaskSequencesByTrigger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getArtifactDeliveryTriggeredEvent() models.Event {
	return models.Event{
		Contenttype: "application/json",
		Data: keptnv2.DeploymentTriggeredEventData{
			EventData: keptnv2.EventData{
				Project: "test-project",
				Stage:   "dev",
				Service: "carts",
			},
			ConfigurationChange: struct {
				Values map[string]interface{} `json:"values"`
			}{
				Values: map[string]interface{}{
					"image": "carts",
				},
			},
			Deployment: keptnv2.DeploymentTriggeredData{
				DeploymentURIsPublic: []string{"uri"},
				DeploymentStrategy:   "direct",
			},
		},
		Extensions:     nil,
		ID:             "artifact-delivery-triggered-id",
		Shkeptncontext: "test-context",
		Source:         common.Stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    "",
		Type:           common.Stringp("sh.keptn.event.dev.artifact-delivery.triggered"),
	}
}

func getStartedEvent(stage string, triggeredID string, eventType string, source string) models.Event {
	return models.Event{
		Contenttype:    "application/json",
		Data:           fake.EventScope{Project: "test-project", Stage: stage, Service: "carts"},
		Extensions:     nil,
		ID:             eventType + "-started-id",
		Shkeptncontext: "test-context",
		Source:         common.Stringp(source),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    triggeredID,
		Type:           common.Stringp(keptnv2.GetStartedEventType(eventType)),
	}
}

func getDeploymentFinishedEvent(stage string, triggeredID string, source string, result keptnv2.ResultType) models.Event {
	return models.Event{
		Contenttype: "application/json",
		Data: keptnv2.DeploymentFinishedEventData{
			EventData: keptnv2.EventData{
				Project: "test-project",
				Stage:   stage,
				Service: "carts",
				Status:  keptnv2.StatusSucceeded,
				Result:  result,
				Message: "i am a message",
			},
			Deployment: keptnv2.DeploymentFinishedData{
				DeploymentURIsLocal:  []string{"uri-1", "uri-2"},
				DeploymentURIsPublic: []string{"public-uri-1", "public-uri-2"},
				DeploymentNames:      []string{"deployment-1"},
				GitCommit:            "commit-1",
			},
		},
		Extensions:     nil,
		ID:             "deployment-finished-id",
		Shkeptncontext: "test-context",
		Source:         common.Stringp(source),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    triggeredID,
		Type:           common.Stringp("sh.keptn.event.deployment.finished"),
	}
}

func getErroredDeploymentFinishedEvent(stage string, triggeredID string, source string) models.Event {
	return models.Event{
		Contenttype: "application/json",
		Data: keptnv2.DeploymentFinishedEventData{
			EventData: keptnv2.EventData{
				Project: "test-project",
				Stage:   stage,
				Service: "carts",
				Status:  keptnv2.StatusErrored,
				Result:  keptnv2.ResultFailed,
			},
			Deployment: keptnv2.DeploymentFinishedData{
				DeploymentURIsLocal:  []string{"uri-1", "uri-2"},
				DeploymentURIsPublic: []string{"public-uri-1", "public-uri-2"},
				DeploymentNames:      []string{"deployment-1"},
				GitCommit:            "commit-1",
			},
		},
		Extensions:     nil,
		ID:             "deployment-finished-id",
		Shkeptncontext: "test-context",
		Source:         common.Stringp(source),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    triggeredID,
		Type:           common.Stringp("sh.keptn.event.deployment.finished"),
	}
}

func getTestTaskFinishedEvent(stage string, triggeredID string) models.Event {
	return models.Event{
		Contenttype: "application/json",
		Data: keptnv2.TestFinishedEventData{
			EventData: keptnv2.EventData{
				Project: "test-project",
				Stage:   stage,
				Service: "carts",
				Status:  keptnv2.StatusSucceeded,
				Result:  keptnv2.ResultPass,
			},
			Test: struct {
				Start     string `json:"start"`
				End       string `json:"end"`
				GitCommit string `json:"gitCommit"`
			}{
				Start:     "start",
				End:       "end",
				GitCommit: "commit-id",
			},
		},
		Extensions:     nil,
		ID:             "test-finished-id",
		Shkeptncontext: "test-context",
		Source:         common.Stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    triggeredID,
		Type:           common.Stringp("sh.keptn.event.test.finished"),
	}
}

func getEvaluationTaskFinishedEvent(stage string, triggeredID string, result keptnv2.ResultType, status keptnv2.StatusType) models.Event {
	return models.Event{
		Contenttype: "application/json",
		Data: keptnv2.EvaluationFinishedEventData{
			EventData: keptnv2.EventData{
				Project: "test-project",
				Stage:   stage,
				Service: "carts",
				Status:  keptnv2.StatusSucceeded,
				Result:  result,
			},
			Evaluation: keptnv2.EvaluationDetails{
				Result: string(result),
			},
		},
		Extensions:     nil,
		ID:             "evaluation-finished-id",
		Shkeptncontext: "test-context",
		Source:         common.Stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    triggeredID,
		Type:           common.Stringp("sh.keptn.event.evaluation.finished"),
	}
}

func getReleaseTaskFinishedEvent(stage string, triggeredID string) models.Event {
	return models.Event{
		Contenttype: "application/json",
		Data: keptnv2.ReleaseFinishedEventData{
			EventData: keptnv2.EventData{
				Project: "test-project",
				Stage:   stage,
				Service: "carts",
				Status:  keptnv2.StatusSucceeded,
				Result:  keptnv2.ResultPass,
			},
		},
		Extensions:     nil,
		ID:             "release-finished-id",
		Shkeptncontext: "test-context",
		Source:         common.Stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    triggeredID,
		Type:           common.Stringp("sh.keptn.event.release.finished"),
	}
}

func sendAndVerifyFinishedEvent(t *testing.T, sc *shipyardController, finishedEvent models.Event, eventType, nextEventType string, mockEV *fake.EventBroker, nextStage string, verifyTriggeredEvent func(t *testing.T, e models.Event) bool) (string, bool) {
	err := sc.HandleIncomingEvent(finishedEvent)
	if err != nil {
		t.Errorf("STEP failed: HandleIncomingEvent(%s) returned %v", *finishedEvent.Type, err)
		return "", true
	}

	scope, _ := getEventScope(finishedEvent)
	if nextStage == "" {
		nextStage = scope.Stage
	}
	// check triggeredEvent collection -> should not contain <eventType>.triggered event anymore
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(eventType),
		Stage:   &scope.Stage,
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)
	done := fake.ShouldNotContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(eventType), scope.Stage)
	if done {
		return "", true
	}

	// check triggeredEvent collection -> should contain <nextEventType>.triggered event
	triggeredEvents, _ = sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(nextEventType),
		Stage:   &nextStage,
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)

	triggeredID := triggeredEvents[0].ID
	done = fake.ShouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(nextEventType), nextStage, nil)
	if done {
		return "", true
	}

	// check startedEvent collection -> should not contain <eventType>.started event anymore
	startedEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:        keptnv2.GetStartedEventType(eventType),
		Stage:       &scope.Stage,
		Service:     common.Stringp("carts"),
		TriggeredID: common.Stringp(finishedEvent.Triggeredid),
	}, common.StartedEvent)
	done = fake.ShouldNotContainEvent(t, startedEvents, keptnv2.GetStartedEventType(eventType), scope.Stage)
	if done {
		return "", true
	}

	// check event broker -> should contain <nextEventType>.triggered event with properties
	done = fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetTriggeredEventType(nextEventType), nextStage, verifyTriggeredEvent)
	if done {
		return "", true
	}
	return triggeredID, false
}

func sendFinishedEventAndVerifyTaskSequenceCompletion(t *testing.T, sc *shipyardController, finishedEvent models.Event, eventType, taskSequence string, mockEV *fake.EventBroker, nextStage string, verifyTriggeredEvent func(t *testing.T, e models.Event) bool) bool {
	err := sc.HandleIncomingEvent(finishedEvent)
	if err != nil {
		t.Errorf("STEP failed: HandleIncomingEvent(%s) returned %v", *finishedEvent.Type, err)
		return true
	}

	scope, _ := getEventScope(finishedEvent)
	if nextStage == "" {
		nextStage = scope.Stage
	}
	// check triggeredEvent collection -> should not contain <eventType>.triggered event anymore
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(eventType),
		Stage:   &scope.Stage,
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)
	done := fake.ShouldNotContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(eventType), scope.Stage)
	if done {
		return true
	}

	// check startedEvent collection -> should not contain <eventType>.started event anymore
	startedEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:        keptnv2.GetStartedEventType(eventType),
		Stage:       &scope.Stage,
		Service:     common.Stringp("carts"),
		TriggeredID: common.Stringp(finishedEvent.Triggeredid),
	}, common.StartedEvent)
	done = fake.ShouldNotContainEvent(t, startedEvents, keptnv2.GetStartedEventType(eventType), scope.Stage)
	if done {
		return true
	}

	// check event broker -> should contain <nextEventType>.triggered event with properties
	done = fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetFinishedEventType(taskSequence), nextStage, verifyTriggeredEvent)
	if done {
		return true
	}
	return false
}

func sendAndVerifyPartialFinishedEvent(t *testing.T, sc *shipyardController, finishedEvent models.Event, eventType, nextEventType string, mockEV *fake.EventBroker, nextStage string) bool {
	err := sc.HandleIncomingEvent(finishedEvent)
	if err != nil {
		t.Errorf("STEP failed: HandleIncomingEvent(%s) returned %v", *finishedEvent.Type, err)
		return true
	}

	scope, _ := getEventScope(finishedEvent)
	if nextStage == "" {
		nextStage = scope.Stage
	}
	// check triggeredEvent collection -> should still contain <eventType>.triggered event
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(eventType),
		Stage:   &scope.Stage,
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)
	done := fake.ShouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(eventType), scope.Stage, nil)
	if done {
		return true
	}

	// check triggeredEvent collection -> should not contain <nextEventType>.triggered event
	triggeredEvents, _ = sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(nextEventType),
		Stage:   &nextStage,
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)

	done = fake.ShouldNotContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(nextEventType), nextStage)
	if done {
		return true
	}

	// check startedEvent collection -> should still contain one <eventType>.started event
	startedEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:        keptnv2.GetStartedEventType(eventType),
		Stage:       &scope.Stage,
		Service:     common.Stringp("carts"),
		TriggeredID: common.Stringp(finishedEvent.Triggeredid),
	}, common.StartedEvent)
	if len(startedEvents) != 1 {
		t.Errorf("List of started events does not hold proper number of events. Expected 1 but got %d", len(startedEvents))
		return true
	}
	done = fake.ShouldContainEvent(t, startedEvents, keptnv2.GetStartedEventType(eventType), scope.Stage, nil)
	if done {
		return true
	}

	// check event broker -> should not contain <nextEventType>.triggered event with properties
	done = fake.ShouldNotContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetTriggeredEventType(nextEventType), nextStage)
	if done {
		return true
	}
	return false
}

func sendAndVerifyStartedEvent(t *testing.T, sc *shipyardController, taskName string, triggeredID string, stage string, fromSource string) bool {
	err := sc.HandleIncomingEvent(getStartedEvent(stage, triggeredID, taskName, fromSource))
	if err != nil {
		t.Errorf("STEP failed: HandleIncomingEvent(%s.started) returned %v", taskName, err)
		return true
	}
	// check startedEvent collection -> should contain <taskName>.started event
	startedEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:        keptnv2.GetStartedEventType(taskName),
		Stage:       common.Stringp(stage),
		Service:     common.Stringp("carts"),
		TriggeredID: common.Stringp(triggeredID),
	}, common.StartedEvent)
	return fake.ShouldContainEvent(t, startedEvents, keptnv2.GetStartedEventType(taskName), stage, nil)
}

func getTestShipyardController() *shipyardController {
	triggeredEventsCollection := []models.Event{}
	startedEventsCollection := []models.Event{}
	finishedEventsCollection := []models.Event{}
	taskSequenceCollection := []models.TaskSequenceEvent{}

	em := &shipyardController{
		projectRepo: nil,
		eventRepo: &db_mock.EventRepoMock{
			GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
				if status[0] == common.TriggeredEvent {
					if triggeredEventsCollection == nil || len(triggeredEventsCollection) == 0 {
						return nil, db.ErrNoEventFound
					}
					return filterEvents(triggeredEventsCollection, filter)
				} else if status[0] == common.StartedEvent {
					if startedEventsCollection == nil || len(startedEventsCollection) == 0 {
						return nil, db.ErrNoEventFound
					}
					return filterEvents(startedEventsCollection, filter)
				} else if status[0] == common.FinishedEvent {
					if finishedEventsCollection == nil || len(finishedEventsCollection) == 0 {
						return nil, db.ErrNoEventFound
					}
					return filterEvents(finishedEventsCollection, filter)
				}
				return nil, nil
			},
			InsertEventFunc: func(project string, event models.Event, status common.EventStatus) error {
				if status == common.TriggeredEvent {
					triggeredEventsCollection = append(triggeredEventsCollection, event)
				} else if status == common.StartedEvent {
					startedEventsCollection = append(startedEventsCollection, event)
				} else if status == common.FinishedEvent {
					finishedEventsCollection = append(finishedEventsCollection, event)
				}
				return nil
			},
			DeleteEventFunc: func(project string, eventID string, status common.EventStatus) error {
				if status == common.TriggeredEvent {
					for index, event := range triggeredEventsCollection {
						if event.ID == eventID {
							triggeredEventsCollection = append(triggeredEventsCollection[:index], triggeredEventsCollection[index+1:]...)
							return nil
						}
					}
				} else if status == common.StartedEvent {
					for index, event := range startedEventsCollection {
						if event.ID == eventID {
							startedEventsCollection = append(startedEventsCollection[:index], startedEventsCollection[index+1:]...)
							return nil
						}
					}
				} else if status == common.FinishedEvent {
					for index, event := range finishedEventsCollection {
						if event.ID == eventID {
							finishedEventsCollection = append(finishedEventsCollection[:index], finishedEventsCollection[index+1:]...)
							return nil
						}
					}
				}
				return nil
			},
		},
		taskSequenceRepo: &db_mock.TaskSequenceRepoMock{
			GetTaskSequenceFunc: func(project, triggeredID string) (*models.TaskSequenceEvent, error) {
				for _, ts := range taskSequenceCollection {
					if ts.TriggeredEventID == triggeredID {
						return &ts, nil
					}
				}
				return nil, nil
			},
			CreateTaskSequenceMappingFunc: func(project string, taskSequenceEvent models.TaskSequenceEvent) error {
				taskSequenceCollection = append(taskSequenceCollection, taskSequenceEvent)
				return nil
			},
			DeleteTaskSequenceMappingFunc: func(keptnContext, project, stage, taskSequenceName string) error {
				newTaskSequenceCollection := []models.TaskSequenceEvent{}

				for index, ts := range taskSequenceCollection {
					if ts.KeptnContext == keptnContext && ts.Stage == stage && ts.TaskSequenceName == taskSequenceName {
						continue
					}
					newTaskSequenceCollection = append(newTaskSequenceCollection, taskSequenceCollection[index])
				}
				taskSequenceCollection = newTaskSequenceCollection
				return nil
			},
		},
		eventsDbOperations: &db_mock.EventsDbOperationsMock{
			UpdateEventOfServiceFunc: func(event interface{}, eventType string, keptnContext string, eventID string, triggeredID string) error {
				return nil
			},
			UpdateShipyardFunc: func(projectName string, shipyardContent string) error {
				return nil
			},
		},
	}
	return em
}

func filterEvents(eventsCollection []models.Event, filter common.EventFilter) ([]models.Event, error) {
	result := []models.Event{}

	for _, event := range eventsCollection {
		scope, _ := getEventScope(event)
		if filter.Type != "" && *event.Type != filter.Type {
			continue
		}
		if filter.Stage != nil && *filter.Stage != scope.Stage {
			continue
		}
		if filter.Service != nil && *filter.Service != scope.Service {
			continue
		}
		if filter.TriggeredID != nil && *filter.TriggeredID != event.Triggeredid {
			continue
		}
		if filter.KeptnContext != nil && *filter.KeptnContext != event.Shkeptncontext {
			continue
		}
		result = append(result, event)
	}
	return result, nil
}
