package api

import (
	"encoding/json"
	"errors"
	"github.com/go-test/deep"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
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

func getTestFinishedEvent() models.Event {
	return models.Event{
		Contenttype:    "application/json",
		Data:           eventScope{Project: "test-project", Stage: "dev", Service: "carts"},
		Extensions:     nil,
		ID:             "test-finished-id",
		Shkeptncontext: "test-context",
		Source:         stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    "test-triggered-id",
		Type:           stringp("sh.keptn.event.approval.finished"),
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

func Test_eventManager_HandleTriggeredEvent(t *testing.T) {
	type fields struct {
		projectRepo        db.ProjectRepo
		triggeredEventRepo db.EventRepo
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
			name: "insert event",
			fields: fields{
				triggeredEventRepo: &mockEventRepo{
					insertEvent: func(project string, event models.Event, status db.EventStatus) error {
						return nil
					},
				},
			},
			args: args{
				event: getTestTriggeredEvent(),
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
			if err := em.handleTriggeredEvent(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("handleTriggeredEvent() error = %v, wantErr %v", err, tt.wantErr)
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
		want    *eventScope
		wantErr bool
	}{
		{
			name: "get event scope",
			args: args{
				event: models.Event{
					Contenttype:    "",
					Data:           eventScope{Project: "sockshop", Stage: "dev", Service: "carts"},
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
			want:    &eventScope{Project: "sockshop", Stage: "dev", Service: "carts"},
			wantErr: false,
		},
		{
			name: "only project available, stage and service missing",
			args: args{
				event: models.Event{
					Contenttype:    "",
					Data:           eventScope{Project: "sockshop"},
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
		logger      *keptn.Logger
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
		logger      *keptn.Logger
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
			name: "received finished event with matching started and triggered event",
			fields: fields{
				projectRepo: nil,
				eventRepo: &mockEventRepo{
					getEvents: func(project string, filter db.EventFilter, status db.EventStatus) ([]models.Event, error) {
						if status == db.TriggeredEvent {
							return []models.Event{getTestTriggeredEvent()}, nil
						} else if status == db.StartedEvent {
							return []models.Event{getTestStartedEvent()}, nil
						}
						return nil, errors.New("received unexpected request")
					},
					insertEvent: func(project string, event models.Event, status db.EventStatus) error {
						t.Error("insertEvent() should not be called in this case")
						return nil
					},
					deleteEvent: func(project string, eventID string, status db.EventStatus) error {
						if status == db.TriggeredEvent {
							if eventID != getTestTriggeredEvent().ID {
								t.Errorf("received unexpected ID for deletion of triggered event. wanted %s but got %s", getTestTriggeredEvent().ID, eventID)
							}
							return nil
						} else if status == db.StartedEvent {
							if eventID != getTestStartedEvent().ID {
								t.Errorf("received unexpected ID for deletion of started event. wanted %s but got %s", getTestTriggeredEvent().ID, eventID)
							}
							return nil
						}
						return nil
					},
				},
				logger: nil,
			},
			args: args{
				event: getTestFinishedEvent(),
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
						} else if status == db.StartedEvent {
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
		logger      *keptn.Logger
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
				logger: keptn.NewLogger("", "", ""),
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
				logger: keptn.NewLogger("", "", ""),
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

// integration test Scenario 1: all events received in expected order
func Test_eventManager_Scenario1(t *testing.T) {

	mockCS := newMockConfigurationService()
	mockCS.Start()
	defer mockCS.Close()

	os.Setenv("CONFIGURATION_SERVICE", mockCS.URL)

	mockEV := newMockEventbroker(t,
		func(meb *mockEventBroker, event *models.Event) {

		},
		func(meb *mockEventBroker) {

		})
	mockEV.server.Start()
	defer mockEV.server.Close()
	os.Setenv("EVENTBROKER", mockEV.server.URL)

	em := getTestShipyardController()

	// STEP 1: send a triggered event -> should be persisted in collection
	triggeredEvent := getTestTriggeredEvent()
	wantEventsInTriggeredCollection := []models.Event{getTestTriggeredEvent()}
	err := em.handleIncomingEvent(triggeredEvent)

	if err != nil {
		t.Errorf("handleIncomingEvent(triggeredEvent) error = %v", err)
	}

	triggeredEvents, err := em.getEvents("", db.EventFilter{}, db.TriggeredEvent, 0)

	if err != nil {
		t.Errorf("GetEvents() error = %v", err)
	}
	if !reflect.DeepEqual(triggeredEvents, wantEventsInTriggeredCollection) {
		t.Errorf("STEP 1 failed: got triggeredEvents = %v, want %v", triggeredEvents, wantEventsInTriggeredCollection)
	}

	// STEP 2: send started event -> event should be persisted in collection
	startedEvent := getTestStartedEvent()
	wantStartedEventsCollection := []models.Event{getTestStartedEvent()}
	err = em.handleIncomingEvent(startedEvent)

	if err != nil {
		t.Errorf("handleIncomingEvent(startedEvent) error = %v", err)
	}

	startedEvents, err := em.getEvents("", db.EventFilter{}, db.StartedEvent, 0)

	if err != nil {
		t.Errorf("GetEvents() error = %v", err)
	}
	if !reflect.DeepEqual(startedEvents, wantStartedEventsCollection) {
		t.Errorf("STEP 2 failed: got startedEvents = %v, want %v", startedEvents, wantStartedEventsCollection)
	}

	// STEP 3: send finished event -> started and triggered event should be deleted from collections
	finishedEvent := getTestFinishedEvent()
	wantEventsInTriggeredCollection = []models.Event{}
	wantStartedEventsCollection = []models.Event{}
	err = em.handleIncomingEvent(finishedEvent)

	if err != nil {
		t.Errorf("handleIncomingEvent(finishedEvent) error = %v", err)
	}

	startedEvents, err = em.getEvents("", db.EventFilter{}, db.StartedEvent, 0)

	if err != nil {
		t.Errorf("GetEvents() error = %v", err)
	}
	if startedEvents != nil && len(startedEvents) > 0 {
		t.Errorf("STEP 3 failed: got startedEvents = %v, want %v", startedEvents, wantStartedEventsCollection)
	}

	triggeredEvents, err = em.getEvents("", db.EventFilter{}, db.TriggeredEvent, 0)

	if err != nil {
		t.Errorf("GetEvents() error = %v", err)
	}
	if triggeredEvents != nil && len(triggeredEvents) > 0 {
		t.Errorf("STEP 3 failed: got triggeredEvents = %v, want %v", triggeredEvents, wantEventsInTriggeredCollection)
	}
}

// integration test Scenario 2: receive triggered event after started event
func Test_eventManager_Scenario2(t *testing.T) {

	var err error
	var wantEventsInTriggeredCollection []models.Event
	var triggeredEvent models.Event
	var triggeredEvents []models.Event

	mockCS := newMockConfigurationService()
	mockCS.Start()
	defer mockCS.Close()

	os.Setenv("CONFIGURATION_SERVICE", mockCS.URL)

	mockEV := newMockEventbroker(t,
		func(meb *mockEventBroker, event *models.Event) {

		},
		func(meb *mockEventBroker) {

		})
	mockEV.server.Start()
	defer mockEV.server.Close()
	os.Setenv("EVENTBROKER", mockEV.server.URL)

	em := getTestShipyardController()

	go func() {
		<-time.After(2 * time.Second)
		// STEP 1: send a triggered event -> should be persisted in collection
		triggeredEvent = getTestTriggeredEvent()
		wantEventsInTriggeredCollection := []models.Event{getTestTriggeredEvent()}
		err := em.handleIncomingEvent(triggeredEvent)
		if err != nil {
			t.Errorf("handleIncomingEvent(triggeredEvent) error = %v", err)
		}

		triggeredEvents, err := em.getEvents("", db.EventFilter{}, db.TriggeredEvent, 0)

		if err != nil {
			t.Errorf("GetEvents() error = %v", err)
		}
		if !reflect.DeepEqual(triggeredEvents, wantEventsInTriggeredCollection) {
			t.Errorf("STEP 1 failed: got triggeredEvents = %v, want %v", triggeredEvents, wantEventsInTriggeredCollection)
		}
	}()

	// STEP 2: send started event -> event should be persisted in collection
	startedEvent := getTestStartedEvent()
	wantStartedEventsCollection := []models.Event{getTestStartedEvent()}
	err = em.handleIncomingEvent(startedEvent)

	if err != nil {
		t.Errorf("handleIncomingEvent(startedEvent) error = %v", err)
	}

	startedEvents, err := em.getEvents("", db.EventFilter{}, db.StartedEvent, 0)

	if err != nil {
		t.Errorf("GetEvents() error = %v", err)
	}
	if !reflect.DeepEqual(startedEvents, wantStartedEventsCollection) {
		t.Errorf("STEP 2 failed: got startedEvents = %v, want %v", startedEvents, wantStartedEventsCollection)
	}

	// STEP 3: send finished event -> started and triggered event should be deleted from collections
	finishedEvent := getTestFinishedEvent()
	wantEventsInTriggeredCollection = []models.Event{}
	wantStartedEventsCollection = []models.Event{}
	err = em.handleIncomingEvent(finishedEvent)

	if err != nil {
		t.Errorf("handleIncomingEvent(finishedEvent) error = %v", err)
	}

	startedEvents, err = em.getEvents("", db.EventFilter{}, db.StartedEvent, 0)

	if err != nil {
		t.Errorf("GetEvents() error = %v", err)
	}
	if startedEvents != nil && len(startedEvents) > 0 {
		t.Errorf("STEP 3 failed: got startedEvents = %v, want %v", startedEvents, wantStartedEventsCollection)
	}

	triggeredEvents, err = em.getEvents("", db.EventFilter{}, db.TriggeredEvent, 0)

	if err != nil {
		t.Errorf("GetEvents() error = %v", err)
	}
	if triggeredEvents != nil && len(triggeredEvents) > 0 {
		t.Errorf("STEP 3 failed: got triggeredEvents = %v, want %v", triggeredEvents, wantEventsInTriggeredCollection)
	}
}

func Test_shipyardController_Scenario1(t *testing.T) {

	_ = getTestShipyardController()

	mockCS := newMockConfigurationService()
	mockCS.Start()
	defer mockCS.Close()

	os.Setenv("CONFIGURATION_SERVICE", mockCS.URL)

	mockEV := newMockEventbroker(t,
		func(meb *mockEventBroker, event *models.Event) {
			meb.receivedEvents = append(meb.receivedEvents, event)
		},
		func(meb *mockEventBroker) {

		})
	mockEV.server.Start()
	defer mockEV.server.Close()
	os.Setenv("EVENTBROKER", mockEV.server.URL)

	// STEP 1
	// send dev.artifact-delivery.triggered event

	// check event broker -> should contain deployment.triggered event with properties: [deployment]

	// check triggeredEvent Collection -> should contain deployment.triggered event

	// STEP 2
	// send deployment.started event

	// check startedEvent collection -> should contain deployment.started event

	// STEP 3
	// send deployment.finished event

	// check triggeredEvent collection -> should not contain deployment.triggered event anymore

	// check startedEvent collection -> should not contain deployment.started event anymore

	// check event broker -> should contain test.triggered event with properties: [deployment, test]

	// STEP 4
	// send test.started event

	// check startedEvent collection -> should contain test.started

	// STEP 5
	// send test.finished event

	// check triggeredEvent collection -> should not contain test.triggered event anymore

	// check startedEvent collection -> should not contain test.started event anymore

	// check event broker -> should contain evaluation.triggered event with properties: [deployment, test, evaluation]

	// STEP 6
	// send evaluation.started event

	// check startedEvent collection -> should contain evaluation.started

	// STEP 7
	// send evaluation.finished event

	// check triggeredEvent collection -> should not contain evaluation.triggered event anymore

	// check startedEvent collection -> should not contain evaluation.started event anymore

	// check event broker -> should contain release.triggered event with properties: [deployment, test, evaluation, release]

	// STEP 8
	// send release.started event

	// check startedEvent collection -> should contain release.started

	// STEP 9
	// send release.finished event

	// check triggeredEvent collection -> should not contain release.triggered event anymore

	// check startedEvent collection -> should not contain release.started event anymore

	// check event broker -> should contain dev.artifact-delivery.finished event

	// check event broker -> should contain hardening.artifact-delivery.triggered event

	// check event broker -> should contain deployment.triggered event with properties: [deployment]

	// STEP 9.1
	// send deployment.started event 1 with ID 1

	// check startedEvent collection -> should contain deployment.started with ID 1

	// STEP 9.2
	// send deployment.started event 2 with ID 2

	// check startedEvent collection -> should contain deployment.started with ID 2

	// STEP 10.1
	// send deployment.finished event 1 with ID 1

	// check triggeredEvents collection -> should still contain deployment.triggered event

	// check startedEvents collection -> should contain deployment.started 2, but not deployment.started 1

	// check finishedEvents collection -> should contain deployment.finished 1

	// check event broker: should not contain test.triggered

	// STEP 10.2
	// send deployment.finished event 1 with ID 1

	// check triggeredEvents collection -> should not contain deployment.triggered event

	// check startedEvents collection -> should not deployment.started

	// check finishedEvents collection -> should not deployment.finished

	// check event broker: should contain test.triggered with properties: [deployment{1,2}, test]

	/*
		case keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName):
			triggered := keptnv2.DeploymentTriggeredEventData{}
			b, _ := json.Marshal(event.Data)
			_ = json.Unmarshal(b, triggered)

			started := keptnv2.DeploymentStartedEventData{
				EventData: triggered.EventData,
			}
			newType := keptnv2.GetStartedEventType(keptnv2.DeploymentTaskName)
			event.Type = &newType
			event.Data = started
			sc.handleIncomingEvent(*event)

			newType = keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName)
			event.Type = &newType
			event.Data = keptnv2.DeploymentFinishedEventData{
				EventData: triggered.EventData,
				Deployment: keptnv2.DeploymentData{
					DeploymentURIsLocal:  []string{"test-uri-local"},
					DeploymentURIsPublic: []string{"test-uri-public"},
					DeploymentNames:      []string{"name-1"},
					GitCommit:            "123",
				},
			}
			sc.handleIncomingEvent(*event)
			break
		case keptnv2.GetTriggeredEventType(keptnv2.TestTaskName):
			break
		case keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName):
			break
		default:
			break
		}

	*/
}

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
		w.Write([]byte(testShipyardFile))
	}))
	return ts
}

type mockEventBroker struct {
	server           *httptest.Server
	receivedEvents   []*models.Event
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
		receivedEvents:   []*models.Event{},
		test:             test,
		handleEventFunc:  handleEventFunc,
		verificationFunc: verificationFunc,
	}

	meb.server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var b []byte
		request.Body.Read(b)
		defer request.Body.Close()
		event := &models.Event{}

		_ = json.Unmarshal(b, event)
		go func() {
			meb.handleEventFunc(meb, event)
		}()
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
		logger: keptn.NewLogger("", "", ""),
	}
	return em
}

func filterEvents(eventsCollection []models.Event, filter db.EventFilter) ([]models.Event, error) {
	result := []models.Event{}

	for _, event := range eventsCollection {
		scope, _ := getEventScope(event)
		if *event.Type != filter.Type {
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
		result = append(result, event)
	}
	return result, nil
}
