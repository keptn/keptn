package handlers

import (
	"errors"
	"github.com/go-test/deep"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"reflect"
	"testing"
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
			em := &eventManager{
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

func stringp(s string) *string {
	return &s
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
			em := &eventManager{
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
			em := &eventManager{
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
			em := &eventManager{
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
			em := &eventManager{
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
