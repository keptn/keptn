package handlers

import (
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"reflect"
	"testing"
)

type GetEventsMock func(project string, filter db.EventFilter) ([]models.Event, error)
type InsertEventMock func(project string, event models.Event) error
type DeleteEventMock func(project string, eventId string) error

type TriggeredEventMock struct {
	getEvents   GetEventsMock
	insertEvent InsertEventMock
	deleteEvent DeleteEventMock
}

func (t TriggeredEventMock) GetEvents(project string, filter db.EventFilter) ([]models.Event, error) {
	return t.getEvents(project, filter)
}

func (t TriggeredEventMock) InsertEvent(project string, event models.Event) error {
	return t.insertEvent(project, event)
}

func (t TriggeredEventMock) DeleteEvent(project string, eventId string) error {
	return t.deleteEvent(project, eventId)
}

type GetProjectsMock func() ([]string, error)

type ProjectRepoMock struct {
	getProjects GetProjectsMock
}

func (p ProjectRepoMock) GetProjects() ([]string, error) {
	return p.getProjects()
}

func getTestEvent() models.Event {
	return models.Event{
		Contenttype:    "application/json",
		Data:           nil,
		Extensions:     nil,
		ID:             "test-id",
		Shkeptncontext: "test-context",
		Source:         stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    "test-triggered-id",
		Type:           stringp("sh.keptn.event.approval.triggered"),
	}
}

func Test_eventManager_GetAllTriggeredEvents(t *testing.T) {
	type fields struct {
		projectRepo        db.ProjectRepo
		triggeredEventRepo db.TriggeredEventRepo
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
				projectRepo: &ProjectRepoMock{getProjects: func() ([]string, error) {
					return []string{"sockshop", "rockshop"}, nil
				}},
				triggeredEventRepo: &TriggeredEventMock{
					getEvents: func(project string, filter db.EventFilter) ([]models.Event, error) {
						return []models.Event{getTestEvent()}, nil
					},
					insertEvent: nil,
					deleteEvent: nil,
				},
			},
			args: args{},
			want: []models.Event{
				getTestEvent(),
				getTestEvent(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := &eventManager{
				projectRepo:        tt.fields.projectRepo,
				triggeredEventRepo: tt.fields.triggeredEventRepo,
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

func stringp(s string) *string {
	return &s
}

func Test_eventManager_GetTriggeredEventsOfProject(t *testing.T) {
	type fields struct {
		projectRepo        db.ProjectRepo
		triggeredEventRepo db.TriggeredEventRepo
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
				triggeredEventRepo: &TriggeredEventMock{
					getEvents: func(project string, filter db.EventFilter) ([]models.Event, error) {
						return []models.Event{getTestEvent()}, nil
					},
					insertEvent: nil,
					deleteEvent: nil,
				},
			},
			args: args{},
			want: []models.Event{
				getTestEvent(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := &eventManager{
				projectRepo:        tt.fields.projectRepo,
				triggeredEventRepo: tt.fields.triggeredEventRepo,
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

func Test_eventManager_InsertEvent(t *testing.T) {
	type fields struct {
		projectRepo        db.ProjectRepo
		triggeredEventRepo db.TriggeredEventRepo
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
				triggeredEventRepo: &TriggeredEventMock{
					insertEvent: func(project string, event models.Event) error {
						return nil
					},
				},
			},
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := &eventManager{
				projectRepo:        tt.fields.projectRepo,
				triggeredEventRepo: tt.fields.triggeredEventRepo,
			}
			if err := em.InsertEvent(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("InsertEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
