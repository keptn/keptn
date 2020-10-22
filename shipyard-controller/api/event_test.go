package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-test/deep"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

type getEventsMock func(project string, filter db.EventFilter, status db.EventStatus) ([]models.Event, error)
type insertEventMock func(project string, event models.Event, status db.EventStatus) error
type deleteEventMock func(project string, eventID string, status db.EventStatus) error

type mockEventRepo struct {
	getEvents   getEventsMock
	insertEvent insertEventMock
	deleteEvent deleteEventMock
}

func (t mockEventRepo) GetEvents(project string, filter db.EventFilter, status db.EventStatus) ([]models.Event, error) {
	return t.getEvents(project, filter, status)
}

func (t mockEventRepo) InsertEvent(project string, event models.Event, status db.EventStatus) error {
	return t.insertEvent(project, event, status)
}

func (t mockEventRepo) DeleteEvent(project string, eventID string, status db.EventStatus) error {
	return t.deleteEvent(project, eventID, status)
}

type mockTaskSequenceRepo struct {
	getTaskSequence           func(project, triggeredID string) (*models.TaskSequenceEvent, error)
	createTaskSequenceMapping func(project string, taskSequenceEvent models.TaskSequenceEvent) error
	deleteTaskSequenceMapping func(keptnContext, project, stage, taskSequenceName string) error
}

// GetTaskSequence godoc
func (mts mockTaskSequenceRepo) GetTaskSequence(project, triggeredID string) (*models.TaskSequenceEvent, error) {
	return mts.getTaskSequence(project, triggeredID)
}

// CreateTaskSequenceMapping godoc
func (mts mockTaskSequenceRepo) CreateTaskSequenceMapping(project string, taskSequenceEvent models.TaskSequenceEvent) error {
	return mts.createTaskSequenceMapping(project, taskSequenceEvent)
}

// DeleteTaskSequenceMapping godoc
func (mts mockTaskSequenceRepo) DeleteTaskSequenceMapping(keptnContext, project, stage, taskSequenceName string) error {
	return mts.deleteTaskSequenceMapping(keptnContext, project, stage, taskSequenceName)
}

type getProjectsMock func() ([]string, error)

type projectRepoMock struct {
	getProjects getProjectsMock
}

func (p projectRepoMock) GetProjects() ([]string, error) {
	return p.getProjects()
}

func getTestTriggeredEvent() models.Event {
	return models.Event{
		Contenttype:    "application/json",
		Data:           eventScope{Project: "test-project", Stage: "dev", Service: "carts"},
		Extensions:     nil,
		ID:             "test-triggered-id",
		Shkeptncontext: "test-context",
		Source:         stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    "",
		Type:           stringp("sh.keptn.event.approval.triggered"),
	}
}

func getTestStartedEvent() models.Event {
	return models.Event{
		Contenttype:    "application/json",
		Data:           eventScope{Project: "test-project", Stage: "dev", Service: "carts"},
		Extensions:     nil,
		ID:             "test-started-id",
		Shkeptncontext: "test-context",
		Source:         stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    "test-triggered-id",
		Type:           stringp("sh.keptn.event.approval.started"),
	}
}

func getTestStartedEventWithUnmatchedTriggeredID() models.Event {
	return models.Event{
		Contenttype:    "application/json",
		Data:           eventScope{Project: "test-project", Stage: "dev", Service: "carts"},
		Extensions:     nil,
		ID:             "test-started-id",
		Shkeptncontext: "test-context",
		Source:         stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    "unmatched-test-triggered-id",
		Type:           stringp("sh.keptn.event.approval.started"),
	}
}

func getTestFinishedEventWithUnmatchedSource() models.Event {
	return models.Event{
		Contenttype:    "application/json",
		Data:           eventScope{Project: "test-project", Stage: "dev", Service: "carts"},
		Extensions:     nil,
		ID:             "test-finished-id",
		Shkeptncontext: "test-context",
		Source:         stringp("unmatched-test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    "test-triggered-id",
		Type:           stringp("sh.keptn.event.approval.finished"),
	}
}

func Test_eventManager_GetAllTriggeredEvents(t *testing.T) {
	type fields struct {
		projectRepo        db.ProjectRepo
		triggeredEventRepo db.EventRepo
	}
	type args struct {
		filter db.EventFilter
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
				projectRepo: &projectRepoMock{getProjects: func() ([]string, error) {
					return []string{"sockshop", "rockshop"}, nil
				}},
				triggeredEventRepo: &mockEventRepo{
					getEvents: func(project string, filter db.EventFilter, status db.EventStatus) ([]models.Event, error) {
						return []models.Event{getTestTriggeredEvent()}, nil
					},
					insertEvent: nil,
					deleteEvent: nil,
				},
			},
			args: args{},
			want: []models.Event{
				getTestTriggeredEvent(),
				getTestTriggeredEvent(),
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
			got, err := em.getAllTriggeredEvents(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAllTriggeredEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAllTriggeredEvents() got = %v, want %v", got, tt.want)
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
		filter  db.EventFilter
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
				triggeredEventRepo: &mockEventRepo{
					getEvents: func(project string, filter db.EventFilter, status db.EventStatus) ([]models.Event, error) {
						return []models.Event{getTestTriggeredEvent()}, nil
					},
					insertEvent: nil,
					deleteEvent: nil,
				},
			},
			args: args{},
			want: []models.Event{
				getTestTriggeredEvent(),
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
			got, err := em.getTriggeredEventsOfProject(tt.args.project, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTriggeredEventsOfProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTriggeredEventsOfProject() got = %v, want %v", got, tt.want)
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
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "received started event with matching triggered event",
			fields: fields{
				projectRepo: nil,
				eventRepo: &mockEventRepo{
					getEvents: func(project string, filter db.EventFilter, status db.EventStatus) ([]models.Event, error) {
						if status == db.TriggeredEvent {
							return []models.Event{getTestTriggeredEvent()}, nil
						}
						return nil, errors.New("received unexpected request")
					},
					insertEvent: func(project string, event models.Event, status db.EventStatus) error {
						if len(deep.Equal(event, getTestStartedEvent())) != 0 {
							t.Errorf("received unexpected event in insertEvent func. wanted %v but got %v", getTestStartedEvent(), event)
							return nil
						}
						return nil
					},
					deleteEvent: func(project string, eventID string, status db.EventStatus) error {
						return nil
					},
				},
				logger: nil,
			},
			args: args{
				event: getTestStartedEvent(),
			},
			wantErr: false,
		},
		{
			name: "received started event with no matching triggered event",
			fields: fields{
				projectRepo: nil,
				eventRepo: &mockEventRepo{
					getEvents: func(project string, filter db.EventFilter, status db.EventStatus) ([]models.Event, error) {
						if status == db.TriggeredEvent {
							return nil, nil
						}
						return nil, errors.New("received unexpected request")
					},
					insertEvent: func(project string, event models.Event, status db.EventStatus) error {
						t.Error("event should not be stored in this case")
						return nil
					},
					deleteEvent: func(project string, eventID string, status db.EventStatus) error {
						return nil
					},
				},
				logger: nil,
			},
			args: args{
				event: getTestStartedEventWithUnmatchedTriggeredID(),
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
			if err := em.handleStartedEvent(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("handleStartedEvent() error = %v, wantErr %v", err, tt.wantErr)
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
				eventRepo: &mockEventRepo{
					getEvents: func(project string, filter db.EventFilter, status db.EventStatus) ([]models.Event, error) {
						if status == db.TriggeredEvent {
							return nil, nil
						} else if status == db.StartedEvent {
							return nil, nil
						}
						return nil, errors.New("received unexpected request")
					},
					insertEvent: func(project string, event models.Event, status db.EventStatus) error {
						return nil
					},
					deleteEvent: func(project string, eventID string, status db.EventStatus) error {
						return nil
					},
				},
				logger: nil,
			},
			args: args{
				event: getTestFinishedEventWithUnmatchedSource(),
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
		filter  db.EventFilter
		status  db.EventStatus
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
				eventRepo: &mockEventRepo{
					getEvents: func(project string, filter db.EventFilter, status db.EventStatus) ([]models.Event, error) {
						return []models.Event{getTestTriggeredEvent()}, nil
					},
				},
				logger: keptncommon.NewLogger("", "", ""),
			},
			args: args{
				project: "test-project",
				filter:  db.EventFilter{},
				status:  db.TriggeredEvent,
			},
			want:    []models.Event{getTestTriggeredEvent()},
			wantErr: false,
		},
		{
			name: "get event after retry",
			fields: fields{
				projectRepo: nil,
				eventRepo: &mockEventRepo{
					getEvents: func(project string, filter db.EventFilter, status db.EventStatus) ([]models.Event, error) {
						if eventAvailable {
							return []models.Event{getTestTriggeredEvent()}, nil
						}
						eventAvailable = true
						return nil, db.ErrNoEventFound
					},
				},
				logger: keptncommon.NewLogger("", "", ""),
			},
			args: args{
				project: "test-project",
				filter:  db.EventFilter{},
				status:  db.TriggeredEvent,
			},
			want:    []models.Event{getTestTriggeredEvent()},
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
			Deployment: keptnv2.DeploymentData{
				DeploymentStrategy: "direct",
			},
		},
		Extensions:     nil,
		ID:             "artifact-delivery-triggered-id",
		Shkeptncontext: "test-context",
		Source:         stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    "",
		Type:           stringp("sh.keptn.event.dev.artifact-delivery.triggered"),
	}
}

func getStartedEvent(stage string, triggeredID string, eventType string, source string) models.Event {
	return models.Event{
		Contenttype:    "application/json",
		Data:           eventScope{Project: "test-project", Stage: stage, Service: "carts"},
		Extensions:     nil,
		ID:             eventType + "-started-id",
		Shkeptncontext: "test-context",
		Source:         stringp(source),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    triggeredID,
		Type:           stringp(keptnv2.GetStartedEventType(eventType)),
	}
}

func getDeploymentFinishedEvent(stage string, triggeredID string, source string) models.Event {
	return models.Event{
		Contenttype: "application/json",
		Data: keptnv2.DeploymentFinishedEventData{
			EventData: keptnv2.EventData{
				Project: "test-project",
				Stage:   stage,
				Service: "carts",
				Status:  keptnv2.StatusSucceeded,
				Result:  keptnv2.ResultPass,
			},
			Deployment: keptnv2.DeploymentData{
				DeploymentURIsLocal:  []string{"uri-1", "uri-2"},
				DeploymentURIsPublic: []string{"public-uri-1", "public-uri-2"},
				DeploymentNames:      []string{"deployment-1"},
				GitCommit:            "commit-1",
			},
		},
		Extensions:     nil,
		ID:             "deployment-finished-id",
		Shkeptncontext: "test-context",
		Source:         stringp(source),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    triggeredID,
		Type:           stringp("sh.keptn.event.deployment.finished"),
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
			Deployment: keptnv2.DeploymentData{
				DeploymentURIsLocal:  []string{"uri-1", "uri-2"},
				DeploymentURIsPublic: []string{"public-uri-1", "public-uri-2"},
				DeploymentNames:      []string{"deployment-1"},
				GitCommit:            "commit-1",
			},
		},
		Extensions:     nil,
		ID:             "deployment-finished-id",
		Shkeptncontext: "test-context",
		Source:         stringp(source),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    triggeredID,
		Type:           stringp("sh.keptn.event.deployment.finished"),
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
		Source:         stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    triggeredID,
		Type:           stringp("sh.keptn.event.test.finished"),
	}
}

func getEvaluationTaskFinishedEvent(stage string, triggeredID string) models.Event {
	return models.Event{
		Contenttype: "application/json",
		Data: keptnv2.EvaluationFinishedEventData{
			EventData: keptnv2.EventData{
				Project: "test-project",
				Stage:   stage,
				Service: "carts",
				Status:  keptnv2.StatusSucceeded,
				Result:  keptnv2.ResultPass,
			},
			Evaluation: keptnv2.EvaluationDetails{
				Result: "pass",
			},
		},
		Extensions:     nil,
		ID:             "evaluation-finished-id",
		Shkeptncontext: "test-context",
		Source:         stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    triggeredID,
		Type:           stringp("sh.keptn.event.evaluation.finished"),
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
		Source:         stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    triggeredID,
		Type:           stringp("sh.keptn.event.release.finished"),
	}
}

func shouldContainEvent(t *testing.T, events []models.Event, eventType string, stage string, eval func(t *testing.T, event models.Event) bool) bool {
	var foundEvent *models.Event
	for index, event := range events {
		scope, _ := getEventScope(event)
		if *event.Type == eventType {
			if stage == "" {
				foundEvent = &events[index]
				break
			} else if stage != "" && scope.Stage == stage {
				foundEvent = &events[index]
				break
			}
		}
	}

	if foundEvent == nil {
		t.Errorf("event list does not contain event of type " + eventType)
		return true
	}
	if eval != nil {
		return eval(t, *foundEvent)
	}
	return false
}

func shouldNotContainEvent(t *testing.T, events []models.Event, eventType string, stage string) bool {
	for _, event := range events {
		if *event.Type == eventType {
			scope, _ := getEventScope(event)
			if stage == "" {
				t.Errorf("event list does contain event of type " + eventType)
				return true
			} else if stage != "" && scope.Stage == stage {
				t.Errorf("event list does contain event of type " + eventType)
				return true
			}
		}
	}
	return false
}

// Scenario 1: Complete task sequence execution + triggering of next task sequence. Events are received in order
func Test_shipyardController_Scenario1(t *testing.T) {

	t.Logf("Executing Shipyard Controller Scenario 1 with shipyard file %s", testShipyardFile)
	sc := getTestShipyardController()

	mockCS := newMockConfigurationService()
	defer mockCS.Close()

	done := false

	_ = os.Setenv("CONFIGURATION_SERVICE", mockCS.URL)

	mockEV := newMockEventbroker(t,
		func(meb *mockEventBroker, event *models.Event) {
			meb.receivedEvents = append(meb.receivedEvents, *event)
		},
		func(meb *mockEventBroker) {

		})
	defer mockEV.server.Close()
	_ = os.Setenv("EVENTBROKER", mockEV.server.URL)

	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.handleIncomingEvent(getArtifactDeliveryTriggeredEvent())
	if err != nil {
		t.Errorf("STEP 1 failed: handleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	// check event broker -> should contain deployment.triggered event with properties: [deployment]
	if len(mockEV.receivedEvents) != 1 {
		t.Errorf("STEP 1 failed: expected %d events in eventbroker, but got %d", 1, len(mockEV.receivedEvents))
		return
	}
	done = shouldContainEvent(t,
		mockEV.receivedEvents,
		keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		"",
		func(t *testing.T, event models.Event) bool {
			deploymentEvent := &keptnv2.DeploymentTriggeredEventData{}

			marshal, _ := json.Marshal(event.Data)
			if err := json.Unmarshal(marshal, deploymentEvent); err != nil {
				t.Error("could not parse incoming deployment.triggered event: " + err.Error())
				return true
			}

			if deploymentEvent.Deployment.DeploymentStrategy != "direct" {
				t.Error(fmt.Sprintf("did not receive correct deployment strategy. Expected 'direct' but got '%s'", deploymentEvent.Deployment.DeploymentStrategy))
				return true
			}
			if deploymentEvent.ConfigurationChange.Values["image"] != "carts" {
				t.Error(fmt.Sprintf("did not receive correct image. Expected 'carts' but got '%s'", deploymentEvent.ConfigurationChange.Values["image"]))
				return true
			}
			return false
		},
	)
	if done {
		return
	}
	// check triggeredEvent Collection -> should contain deployment.triggered event
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", db.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		Stage:   stringp("dev"),
		Service: stringp("carts"),
		Source:  stringp("shipyard-controller"),
	}, db.TriggeredEvent)
	done = shouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
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
		getDeploymentFinishedEvent("dev", triggeredID, "test-source"),
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
			if len(testData.Deployment.DeploymentURIsPublic) != 2 {
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
				t.Errorf("Test property was not transmitted corretly: %v", testData.Test)
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
	// send evaluation.finished event
	triggeredID, done = sendAndVerifyFinishedEvent(t, sc, getEvaluationTaskFinishedEvent("dev", triggeredID), keptnv2.EvaluationTaskName, keptnv2.ReleaseTaskName, mockEV, "", nil)
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
	done = shouldContainEvent(t, mockEV.receivedEvents, keptnv2.GetFinishedEventType("dev.artifact-delivery"), "dev", nil)
	if done {
		return
	}

	done = shouldContainEvent(t, mockEV.receivedEvents, keptnv2.GetTriggeredEventType("hardening.artifact-delivery"), "hardening", nil)
	if done {
		return
	}

	done = shouldContainEvent(
		t,
		mockEV.receivedEvents,
		keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		"hardening",
		func(t *testing.T, event models.Event) bool {
			marshal, _ := json.Marshal(event.Data)
			testData := &keptnv2.DeploymentTriggeredEventData{}

			err := json.Unmarshal(marshal, testData)

			if err != nil {
				t.Errorf("Expected test.triggered data but could not convert: %v: %s", event.Data, err.Error())
				return true
			}

			if len(testData.Deployment.DeploymentURIsLocal) != 2 {
				t.Errorf("DeploymentURIsLocal property was not transmitted correctly")
				return true
			}
			if len(testData.Deployment.DeploymentURIsPublic) != 2 {
				t.Errorf("DeploymentURIsLocal property was not transmitted correctly")
				return true
			}
			return false
		},
	)
	if done {
		return
	}

	finishedEvents, _ := sc.eventRepo.GetEvents("test-project", db.EventFilter{
		Stage: stringp("dev"),
	}, db.FinishedEvent)

	shouldNotContainEvent(t, finishedEvents, keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName), "dev")
	shouldNotContainEvent(t, finishedEvents, keptnv2.GetFinishedEventType(keptnv2.TestTaskName), "dev")
	shouldNotContainEvent(t, finishedEvents, keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), "dev")
	shouldNotContainEvent(t, finishedEvents, keptnv2.GetFinishedEventType(keptnv2.ReleaseTaskName), "dev")

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
	done = sendAndVerifyPartialFinishedEvent(t, sc, getDeploymentFinishedEvent("hardening", triggeredID, "test-source-1"), keptnv2.DeploymentTaskName, keptnv2.ReleaseTaskName, mockEV, "")
	if done {
		return
	}

	// STEP 10.2
	// send deployment.finished event 1 with ID 1
	triggeredID, done = sendAndVerifyFinishedEvent(t, sc, getDeploymentFinishedEvent("hardening", triggeredID, "test-source-2"), keptnv2.DeploymentTaskName, keptnv2.TestTaskName, mockEV, "", nil)
	if done {
		return
	}
}

// Scenario 2: Partial task sequence execution + triggering of next task sequence. Events are received out of order
func Test_shipyardController_Scenario2(t *testing.T) {

	t.Logf("Executing Shipyard Controller Scenario 1 with shipyard file %s", testShipyardFile)
	sc := getTestShipyardController()

	mockCS := newMockConfigurationService()
	defer mockCS.Close()

	done := false

	_ = os.Setenv("CONFIGURATION_SERVICE", mockCS.URL)

	mockEV := newMockEventbroker(t,
		func(meb *mockEventBroker, event *models.Event) {
			meb.receivedEvents = append(meb.receivedEvents, *event)
		},
		func(meb *mockEventBroker) {

		})
	defer mockEV.server.Close()
	_ = os.Setenv("EVENTBROKER", mockEV.server.URL)

	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.handleIncomingEvent(getArtifactDeliveryTriggeredEvent())
	if err != nil {
		t.Errorf("STEP 1 failed: handleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	// check event broker -> should contain deployment.triggered event with properties: [deployment]
	if len(mockEV.receivedEvents) != 1 {
		t.Errorf("STEP 1 failed: expected %d events in eventbroker, but got %d", 1, len(mockEV.receivedEvents))
		return
	}
	done = shouldContainEvent(t, mockEV.receivedEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
	if done {
		return
	}
	// check triggeredEvent Collection -> should contain deployment.triggered event
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", db.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		Stage:   stringp("dev"),
		Service: stringp("carts"),
		Source:  stringp("shipyard-controller"),
	}, db.TriggeredEvent)
	done = shouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
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
		getDeploymentFinishedEvent("dev", triggeredID, "test-source"),
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
			if len(testData.Deployment.DeploymentURIsPublic) != 2 {
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

	mockCS := newMockConfigurationService()
	defer mockCS.Close()

	done := false

	_ = os.Setenv("CONFIGURATION_SERVICE", mockCS.URL)

	mockEV := newMockEventbroker(t,
		func(meb *mockEventBroker, event *models.Event) {
			meb.receivedEvents = append(meb.receivedEvents, *event)
		},
		func(meb *mockEventBroker) {

		})
	defer mockEV.server.Close()
	_ = os.Setenv("EVENTBROKER", mockEV.server.URL)

	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.handleIncomingEvent(getArtifactDeliveryTriggeredEvent())
	if err != nil {
		t.Errorf("STEP 1 failed: handleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	// check event broker -> should contain deployment.triggered event with properties: [deployment]
	if len(mockEV.receivedEvents) != 1 {
		t.Errorf("STEP 1 failed: expected %d events in eventbroker, but got %d", 1, len(mockEV.receivedEvents))
		return
	}
	done = shouldContainEvent(t, mockEV.receivedEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
	if done {
		return
	}
	// check triggeredEvent Collection -> should contain deployment.triggered event
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", db.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		Stage:   stringp("dev"),
		Service: stringp("carts"),
		Source:  stringp("shipyard-controller"),
	}, db.TriggeredEvent)
	done = shouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
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

func sendAndVerifyFinishedEvent(t *testing.T, sc *shipyardController, finishedEvent models.Event, eventType, nextEventType string, mockEV *mockEventBroker, nextStage string, verifyTriggeredEvent func(t *testing.T, e models.Event) bool) (string, bool) {
	err := sc.handleIncomingEvent(finishedEvent)
	if err != nil {
		t.Errorf("STEP failed: handleIncomingEvent(%s) returned %v", *finishedEvent.Type, err)
		return "", true
	}

	scope, _ := getEventScope(finishedEvent)
	if nextStage == "" {
		nextStage = scope.Stage
	}
	// check triggeredEvent collection -> should not contain <eventType>.triggered event anymore
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", db.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(eventType),
		Stage:   &scope.Stage,
		Service: stringp("carts"),
		Source:  stringp("shipyard-controller"),
	}, db.TriggeredEvent)
	done := shouldNotContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(eventType), scope.Stage)
	if done {
		return "", true
	}

	// check triggeredEvent collection -> should contain <nextEventType>.triggered event
	triggeredEvents, _ = sc.eventRepo.GetEvents("test-project", db.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(nextEventType),
		Stage:   &nextStage,
		Service: stringp("carts"),
		Source:  stringp("shipyard-controller"),
	}, db.TriggeredEvent)

	triggeredID := triggeredEvents[0].ID
	done = shouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(nextEventType), nextStage, nil)
	if done {
		return "", true
	}

	// check startedEvent collection -> should not contain <eventType>.started event anymore
	startedEvents, _ := sc.eventRepo.GetEvents("test-project", db.EventFilter{
		Type:        keptnv2.GetStartedEventType(eventType),
		Stage:       &scope.Stage,
		Service:     stringp("carts"),
		TriggeredID: stringp(finishedEvent.Triggeredid),
	}, db.StartedEvent)
	done = shouldNotContainEvent(t, startedEvents, keptnv2.GetStartedEventType(eventType), scope.Stage)
	if done {
		return "", true
	}

	// check event broker -> should contain <nextEventType>.triggered event with properties
	done = shouldContainEvent(t, mockEV.receivedEvents, keptnv2.GetTriggeredEventType(nextEventType), nextStage, verifyTriggeredEvent)
	if done {
		return "", true
	}
	return triggeredID, false
}

func sendFinishedEventAndVerifyTaskSequenceCompletion(t *testing.T, sc *shipyardController, finishedEvent models.Event, eventType, taskSequence string, mockEV *mockEventBroker, nextStage string, verifyTriggeredEvent func(t *testing.T, e models.Event) bool) bool {
	err := sc.handleIncomingEvent(finishedEvent)
	if err != nil {
		t.Errorf("STEP failed: handleIncomingEvent(%s) returned %v", *finishedEvent.Type, err)
		return true
	}

	scope, _ := getEventScope(finishedEvent)
	if nextStage == "" {
		nextStage = scope.Stage
	}
	// check triggeredEvent collection -> should not contain <eventType>.triggered event anymore
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", db.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(eventType),
		Stage:   &scope.Stage,
		Service: stringp("carts"),
		Source:  stringp("shipyard-controller"),
	}, db.TriggeredEvent)
	done := shouldNotContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(eventType), scope.Stage)
	if done {
		return true
	}

	// check startedEvent collection -> should not contain <eventType>.started event anymore
	startedEvents, _ := sc.eventRepo.GetEvents("test-project", db.EventFilter{
		Type:        keptnv2.GetStartedEventType(eventType),
		Stage:       &scope.Stage,
		Service:     stringp("carts"),
		TriggeredID: stringp(finishedEvent.Triggeredid),
	}, db.StartedEvent)
	done = shouldNotContainEvent(t, startedEvents, keptnv2.GetStartedEventType(eventType), scope.Stage)
	if done {
		return true
	}

	// check event broker -> should contain <nextEventType>.triggered event with properties
	done = shouldContainEvent(t, mockEV.receivedEvents, keptnv2.GetFinishedEventType(taskSequence), nextStage, verifyTriggeredEvent)
	if done {
		return true
	}
	return false
}

func sendAndVerifyPartialFinishedEvent(t *testing.T, sc *shipyardController, finishedEvent models.Event, eventType, nextEventType string, mockEV *mockEventBroker, nextStage string) bool {
	err := sc.handleIncomingEvent(finishedEvent)
	if err != nil {
		t.Errorf("STEP failed: handleIncomingEvent(%s) returned %v", *finishedEvent.Type, err)
		return true
	}

	scope, _ := getEventScope(finishedEvent)
	if nextStage == "" {
		nextStage = scope.Stage
	}
	// check triggeredEvent collection -> should still contain <eventType>.triggered event
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", db.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(eventType),
		Stage:   &scope.Stage,
		Service: stringp("carts"),
		Source:  stringp("shipyard-controller"),
	}, db.TriggeredEvent)
	done := shouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(eventType), scope.Stage, nil)
	if done {
		return true
	}

	// check triggeredEvent collection -> should not contain <nextEventType>.triggered event
	triggeredEvents, _ = sc.eventRepo.GetEvents("test-project", db.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(nextEventType),
		Stage:   &nextStage,
		Service: stringp("carts"),
		Source:  stringp("shipyard-controller"),
	}, db.TriggeredEvent)

	done = shouldNotContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(nextEventType), nextStage)
	if done {
		return true
	}

	// check startedEvent collection -> should still contain one <eventType>.started event
	startedEvents, _ := sc.eventRepo.GetEvents("test-project", db.EventFilter{
		Type:        keptnv2.GetStartedEventType(eventType),
		Stage:       &scope.Stage,
		Service:     stringp("carts"),
		TriggeredID: stringp(finishedEvent.Triggeredid),
	}, db.StartedEvent)
	if len(startedEvents) != 1 {
		t.Errorf("List of started events does not hold proper number of events. Expected 1 but got %d", len(startedEvents))
		return true
	}
	done = shouldContainEvent(t, startedEvents, keptnv2.GetStartedEventType(eventType), scope.Stage, nil)
	if done {
		return true
	}

	// check event broker -> should not contain <nextEventType>.triggered event with properties
	done = shouldNotContainEvent(t, mockEV.receivedEvents, keptnv2.GetTriggeredEventType(nextEventType), nextStage)
	if done {
		return true
	}
	return false
}

func sendAndVerifyStartedEvent(t *testing.T, sc *shipyardController, taskName string, triggeredID string, stage string, fromSource string) bool {
	err := sc.handleIncomingEvent(getStartedEvent(stage, triggeredID, taskName, fromSource))
	if err != nil {
		t.Errorf("STEP failed: handleIncomingEvent(%s.started) returned %v", taskName, err)
		return true
	}
	// check startedEvent collection -> should contain <taskName>.started event
	startedEvents, _ := sc.eventRepo.GetEvents("test-project", db.EventFilter{
		Type:        keptnv2.GetStartedEventType(taskName),
		Stage:       stringp(stage),
		Service:     stringp("carts"),
		TriggeredID: stringp(triggeredID),
	}, db.StartedEvent)
	return shouldContainEvent(t, startedEvents, keptnv2.GetStartedEventType(taskName), stage, nil)
}

const testShipyardResource = `{
      "resourceContent": "YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJkCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5CiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBzdHJhdGVneTogZGlyZWN0CiAgICAgIC0gbmFtZTogdGVzdAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBraW5kOiBmdW5jdGlvbmFsCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbiAKICAgICAgLSBuYW1lOiByZWxlYXNlIAoKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg==",
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

  - name: hardening
    sequences:
    - name: artifact-delivery
      triggers:
      - dev.artifact-delivery.finished
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
      triggers:
      - hardening.artifact-delivery.finished
      tasks:
      - name: deployment
        properties:
          strategy: blue_green
      - name: release
      
    - name: remediation
      tasks:
      - name: remediation
      - name: evaluation`

func newMockConfigurationService() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(testShipyardResource))
	}))
	return ts
}

type mockEventBroker struct {
	server           *httptest.Server
	receivedEvents   []models.Event
	test             *testing.T
	handleEventFunc  func(meb *mockEventBroker, event *models.Event)
	verificationFunc func(meb *mockEventBroker)
}

func (meb *mockEventBroker) handleEvent(event *models.Event) {
	meb.handleEventFunc(meb, event)
}

func newMockEventbroker(test *testing.T, handleEventFunc func(meb *mockEventBroker, event *models.Event), verificationFunc func(meb *mockEventBroker)) *mockEventBroker {
	meb := &mockEventBroker{
		server:           nil,
		receivedEvents:   []models.Event{},
		test:             test,
		handleEventFunc:  handleEventFunc,
		verificationFunc: verificationFunc,
	}

	meb.server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		b, _ := ioutil.ReadAll(request.Body)
		defer func() {
			_ = request.Body.Close()
		}()
		event := &models.Event{}

		_ = json.Unmarshal(b, event)
		meb.handleEventFunc(meb, event)

	}))

	return meb
}

func getTestShipyardController() *shipyardController {
	triggeredEventsCollection := []models.Event{}
	startedEventsCollection := []models.Event{}
	finishedEventsCollection := []models.Event{}
	taskSequenceCollection := []models.TaskSequenceEvent{}

	em := &shipyardController{
		projectRepo: nil,
		eventRepo: &mockEventRepo{
			getEvents: func(project string, filter db.EventFilter, status db.EventStatus) ([]models.Event, error) {
				if status == db.TriggeredEvent {
					if triggeredEventsCollection == nil || len(triggeredEventsCollection) == 0 {
						return nil, db.ErrNoEventFound
					}
					return filterEvents(triggeredEventsCollection, filter)
				} else if status == db.StartedEvent {
					if startedEventsCollection == nil || len(startedEventsCollection) == 0 {
						return nil, db.ErrNoEventFound
					}
					return filterEvents(startedEventsCollection, filter)
				} else if status == db.FinishedEvent {
					if finishedEventsCollection == nil || len(finishedEventsCollection) == 0 {
						return nil, db.ErrNoEventFound
					}
					return filterEvents(finishedEventsCollection, filter)
				}
				return nil, nil
			},
			insertEvent: func(project string, event models.Event, status db.EventStatus) error {
				if status == db.TriggeredEvent {
					triggeredEventsCollection = append(triggeredEventsCollection, event)
				} else if status == db.StartedEvent {
					startedEventsCollection = append(startedEventsCollection, event)
				} else if status == db.FinishedEvent {
					finishedEventsCollection = append(finishedEventsCollection, event)
				}
				return nil
			},
			deleteEvent: func(project string, eventID string, status db.EventStatus) error {
				if status == db.TriggeredEvent {
					for index, event := range triggeredEventsCollection {
						if event.ID == eventID {
							triggeredEventsCollection = append(triggeredEventsCollection[:index], triggeredEventsCollection[index+1:]...)
							return nil
						}
					}
				} else if status == db.StartedEvent {
					for index, event := range startedEventsCollection {
						if event.ID == eventID {
							startedEventsCollection = append(startedEventsCollection[:index], startedEventsCollection[index+1:]...)
							return nil
						}
					}
				} else if status == db.FinishedEvent {
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
		taskSequenceRepo: &mockTaskSequenceRepo{
			getTaskSequence: func(project, triggeredID string) (*models.TaskSequenceEvent, error) {
				for _, ts := range taskSequenceCollection {
					if ts.TriggeredEventID == triggeredID {
						return &ts, nil
					}
				}
				return nil, nil
			},
			createTaskSequenceMapping: func(project string, taskSequenceEvent models.TaskSequenceEvent) error {
				taskSequenceCollection = append(taskSequenceCollection, taskSequenceEvent)
				return nil
			},
			deleteTaskSequenceMapping: func(keptnContext, project, stage, taskSequenceName string) error {
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
		logger: keptncommon.NewLogger("", "", ""),
	}
	return em
}

func filterEvents(eventsCollection []models.Event, filter db.EventFilter) ([]models.Event, error) {
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
		// TODO: Add test cases.
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTaskSequenceInStage() got = %v, want %v", got, tt.want)
			}
		})
	}
}
