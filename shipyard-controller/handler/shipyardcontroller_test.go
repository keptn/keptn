package handler

import (
	"errors"
	"github.com/go-test/deep"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	fakehooks "github.com/keptn/keptn/shipyard-controller/handler/sequencehooks/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func Test_GetAllTriggeredEvents(t *testing.T) {
	type fields struct {
		projectRepo        db.ProjectMVRepo
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
				projectRepo: &db_mock.ProjectMVRepoMock{GetProjectsFunc: func() ([]*models.ExpandedProject, error) {
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
				projectMvRepo: tt.fields.projectRepo,
				eventRepo:     tt.fields.triggeredEventRepo,
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

func Test_GetTriggeredEventsOfProject(t *testing.T) {
	type fields struct {
		projectRepo        db.ProjectMVRepo
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
				projectRepo: &db_mock.ProjectMVRepoMock{GetProjectFunc: func(projectName string) (*models.ExpandedProject, error) {
					return &models.ExpandedProject{ProjectName: projectName}, nil
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
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := &shipyardController{
				projectMvRepo: tt.fields.projectRepo,
				eventRepo:     tt.fields.triggeredEventRepo,
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

func Test_HandleStartedEvents(t *testing.T) {
	type fields struct {
		projectMvRepo    db.ProjectMVRepo
		eventRepo        db.EventRepo
		taskSequenceRepo db.TaskSequenceRepo
		taskStartedHook  *fakehooks.ISequenceTaskStartedHookMock
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
		wantHookCalled         bool
	}{
		{
			name: "received started event with matching triggered event",
			fields: fields{
				projectMvRepo: nil,
				eventRepo: &db_mock.EventRepoMock{
					GetEventsWithRetryFunc: func(project string, filter common.EventFilter, status common.EventStatus, retries int) ([]models.Event, error) {
						if status == common.TriggeredEvent {
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
				taskSequenceRepo: &db_mock.TaskSequenceRepoMock{GetTaskExecutionsFunc: func(project string, filter models.TaskExecution) ([]models.TaskExecution, error) {
					return []models.TaskExecution{
						{},
					}, nil
				}},
				taskStartedHook: &fakehooks.ISequenceTaskStartedHookMock{OnSequenceTaskStartedFunc: func(event models.Event) {}},
			},
			args: args{
				event: fake.GetTestStartedEvent(),
			},
			wantErr:        false,
			wantHookCalled: true,
		},
		{
			name: "received started event with no matching triggered event",
			fields: fields{
				projectMvRepo: nil,
				eventRepo: &db_mock.EventRepoMock{
					GetEventsWithRetryFunc: func(project string, filter common.EventFilter, status common.EventStatus, retries int) ([]models.Event, error) {
						if status == common.TriggeredEvent {
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
				taskSequenceRepo: &db_mock.TaskSequenceRepoMock{GetTaskExecutionsFunc: func(project string, filter models.TaskExecution) ([]models.TaskExecution, error) {
					return []models.TaskExecution{
						{},
					}, nil
				}},
				taskStartedHook: &fakehooks.ISequenceTaskStartedHookMock{OnSequenceTaskStartedFunc: func(event models.Event) {}},
			},
			args: args{
				event: fake.GetTestStartedEventWithUnmatchedTriggeredID(),
			},
			wantErr:                true,
			wantErrNoMatchingEvent: true,
			wantHookCalled:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := &shipyardController{
				projectMvRepo:    tt.fields.projectMvRepo,
				eventRepo:        tt.fields.eventRepo,
				taskSequenceRepo: tt.fields.taskSequenceRepo,
			}
			em.AddSequenceTaskStartedHook(tt.fields.taskStartedHook)
			err := em.handleStartedEvent(tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleStartedEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErrNoMatchingEvent && (err != ErrNoMatchingEvent) {
				t.Errorf("handleStartedEvent() expected ErrNoMatchingEvent but got %v", err)
			}

			if tt.wantHookCalled {
				require.Len(t, tt.fields.taskStartedHook.OnSequenceTaskStartedCalls(), 1)
			} else {
				require.Empty(t, tt.fields.taskStartedHook.OnSequenceTaskStartedCalls())
			}
		})
	}
}

func TestHandleFinishedEvent(t *testing.T) {
	type fields struct {
		projectMvRepo    db.ProjectMVRepo
		eventRepo        db.EventRepo
		taskSequenceRepo db.TaskSequenceRepo
		taskFinishedHook *fakehooks.ISequenceTaskFinishedHookMock
	}
	type args struct {
		event models.Event
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		wantHookCalled bool
	}{
		{
			name: "received finished event with no matching triggered event",
			fields: fields{
				projectMvRepo: nil,
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
					GetStartedEventsForTriggeredIDFunc: func(eventScope models.EventScope) ([]models.Event, error) {
						return nil, nil
					},
				},
				taskSequenceRepo: &db_mock.TaskSequenceRepoMock{GetTaskExecutionsFunc: func(project string, filter models.TaskExecution) ([]models.TaskExecution, error) {
					return []models.TaskExecution{
						{},
					}, nil
				}},
				taskFinishedHook: &fakehooks.ISequenceTaskFinishedHookMock{OnSequenceTaskFinishedFunc: func(event models.Event) {}},
			},
			args: args{
				event: fake.GetTestFinishedEventWithUnmatchedSource(),
			},
			wantErr:        true,
			wantHookCalled: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := &shipyardController{
				projectMvRepo:    tt.fields.projectMvRepo,
				eventRepo:        tt.fields.eventRepo,
				taskSequenceRepo: tt.fields.taskSequenceRepo,
			}

			em.AddSequenceTaskFinishedHook(tt.fields.taskFinishedHook)
			if err := em.handleFinishedEvent(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("handleFinishedEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantHookCalled {
				require.Len(t, tt.fields.taskFinishedHook.OnSequenceTaskFinishedCalls(), 1)
			} else {
				require.Empty(t, tt.fields.taskFinishedHook.OnSequenceTaskFinishedCalls())
			}
		})
	}
}
