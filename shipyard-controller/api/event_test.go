package api

import (
	"errors"
	"github.com/go-test/deep"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
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
		Data:           eventData{Project: "test-project"},
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
		Data:           eventData{Project: "test-project"},
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
		Data:           eventData{Project: "test-project"},
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
		Data:           eventData{Project: "test-project"},
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
		Data:           eventData{Project: "test-project"},
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

func Test_getEventProject(t *testing.T) {
	type args struct {
		event models.Event
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "get project name",
			args: args{
				event: models.Event{
					Contenttype:    "",
					Data:           eventData{Project: "sockshop"},
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
			want:    "sockshop",
			wantErr: false,
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
			want:    "",
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
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getEventProject(tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEventProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getEventProject() got = %v, want %v", got, tt.want)
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

	triggeredEventsCollection := []models.Event{}
	startedEventsCollection := []models.Event{}

	em := getTestEventManager(triggeredEventsCollection, startedEventsCollection)

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

	triggeredEventsCollection := []models.Event{}
	startedEventsCollection := []models.Event{}

	em := getTestEventManager(triggeredEventsCollection, startedEventsCollection)

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

func getTestEventManager(triggeredEventsCollection []models.Event, startedEventsCollection []models.Event) *shipyardController {
	em := &shipyardController{
		projectRepo: nil,
		eventRepo: &mockEventRepo{
			getEvents: func(project string, filter db.EventFilter, status db.EventStatus) ([]models.Event, error) {
				if status == db.TriggeredEvent {
					if triggeredEventsCollection == nil || len(triggeredEventsCollection) == 0 {
						return nil, db.ErrNoEventFound
					}
					return triggeredEventsCollection, nil
				} else if status == db.StartedEvent {
					if startedEventsCollection == nil || len(startedEventsCollection) == 0 {
						return nil, db.ErrNoEventFound
					}
					return startedEventsCollection, nil
				}
				return nil, nil
			},
			insertEvent: func(project string, event models.Event, status db.EventStatus) error {
				if status == db.TriggeredEvent {
					triggeredEventsCollection = append(triggeredEventsCollection, event)
				} else if status == db.StartedEvent {
					startedEventsCollection = append(startedEventsCollection, event)
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
				}
				return nil
			},
		},
		logger: nil,
	}
	return em
}
