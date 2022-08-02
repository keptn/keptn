package controller

import (
	"errors"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/internal/common"
	"github.com/keptn/keptn/shipyard-controller/internal/controller/fake"
	"github.com/keptn/keptn/shipyard-controller/internal/db"
	db_mock "github.com/keptn/keptn/shipyard-controller/internal/db/mock"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func Test_GetAllTriggeredEvents(t *testing.T) {
	triggered := GetTestTriggeredEvent()
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
		want    []apimodels.KeptnContextExtendedCE
		wantErr bool
	}{
		{
			name: "Get triggered events for all projects",
			fields: fields{
				projectRepo: &db_mock.ProjectMVRepoMock{GetProjectsFunc: func() ([]*apimodels.ExpandedProject, error) {
					return []*apimodels.ExpandedProject{{
						ProjectName: "sockshop",
					}, {
						ProjectName: "rockshop",
					}}, nil
				}},
				triggeredEventRepo: &db_mock.EventRepoMock{
					GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]apimodels.KeptnContextExtendedCE, error) {
						return []apimodels.KeptnContextExtendedCE{triggered}, nil
					},
					InsertEventFunc: nil,
					DeleteEventFunc: nil,
				},
			},
			args: args{},
			want: []apimodels.KeptnContextExtendedCE{
				triggered,
				triggered,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := &ShipyardController{
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
	triggered := GetTestTriggeredEvent()
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
		want    []apimodels.KeptnContextExtendedCE
		wantErr bool
	}{
		{
			name: "Get triggered events for project",
			fields: fields{
				projectRepo: &db_mock.ProjectMVRepoMock{GetProjectFunc: func(projectName string) (*apimodels.ExpandedProject, error) {
					return &apimodels.ExpandedProject{ProjectName: projectName}, nil
				}},
				triggeredEventRepo: &db_mock.EventRepoMock{
					GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]apimodels.KeptnContextExtendedCE, error) {
						return []apimodels.KeptnContextExtendedCE{triggered}, nil
					},
					InsertEventFunc: nil,
					DeleteEventFunc: nil,
				},
			},
			args: args{},
			want: []apimodels.KeptnContextExtendedCE{
				triggered,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := &ShipyardController{
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

func TestHandleTaskEvent(t *testing.T) {
	type fields struct {
		projectMvRepo         db.ProjectMVRepo
		eventRepo             db.EventRepo
		sequenceExecutionRepo *db_mock.SequenceExecutionRepoMock
		taskFinishedHook      *fake.ISequenceTaskFinishedHookMock
	}
	type args struct {
		event apimodels.KeptnContextExtendedCE
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		wantHookCalled bool
	}{
		{
			name: "received finished event with no matching sequence execution",
			fields: fields{
				projectMvRepo: nil,
				eventRepo: &db_mock.EventRepoMock{
					GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]apimodels.KeptnContextExtendedCE, error) {
						if status[0] == common.TriggeredEvent {
							return nil, nil
						} else if status[0] == common.StartedEvent {
							return nil, nil
						}
						return nil, errors.New("received unexpected request")
					},
					InsertEventFunc: func(project string, event apimodels.KeptnContextExtendedCE, status common.EventStatus) error {
						return nil
					},
					DeleteEventFunc: func(project string, eventID string, status common.EventStatus) error {
						return nil
					},
					GetStartedEventsForTriggeredIDFunc: func(eventScope models.EventScope) ([]apimodels.KeptnContextExtendedCE, error) {
						return nil, nil
					},
				},
				sequenceExecutionRepo: &db_mock.SequenceExecutionRepoMock{
					GetFunc: func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
						return nil, nil
					},
				},
				taskFinishedHook: &fake.ISequenceTaskFinishedHookMock{OnSequenceTaskFinishedFunc: func(event apimodels.KeptnContextExtendedCE) {}},
			},
			args: args{
				event: GetTestFinishedEventWithUnmatchedSource(),
			},
			wantErr:        true,
			wantHookCalled: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := &ShipyardController{
				projectMvRepo:         tt.fields.projectMvRepo,
				eventRepo:             tt.fields.eventRepo,
				sequenceExecutionRepo: tt.fields.sequenceExecutionRepo,
			}

			em.AddSequenceTaskFinishedHook(tt.fields.taskFinishedHook)
			if err := em.handleTaskEvent(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("handleTaskFinished() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantHookCalled {
				require.Len(t, tt.fields.taskFinishedHook.OnSequenceTaskFinishedCalls(), 1)
			} else {
				require.Empty(t, tt.fields.taskFinishedHook.OnSequenceTaskFinishedCalls())
			}
		})
	}
}
