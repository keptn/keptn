package handler_test

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSequenceMigrator_MigrateSequences(t *testing.T) {
	timestampForAllEvents := timeutils.GetKeptnTimeStamp(time.Now().UTC())
	type fields struct {
		eventRepo        *db_mock.EventRepoMock
		taskSequenceRepo *db_mock.SequenceStateRepoMock
		projectRepo      *db_mock.ProjectRepoMock
	}
	tests := []struct {
		name              string
		fields            fields
		wantSequenceState models.SequenceState
	}{
		{
			name: "recreate task sequence state - complete sequence",
			fields: fields{
				eventRepo: &db_mock.EventRepoMock{
					GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
						return []models.Event{
							{
								Data: keptnv2.EventData{
									Project: "my-project",
									Stage:   "staging",
									Service: "my-service",
									Result:  keptnv2.ResultFailed,
									Status:  keptnv2.StatusSucceeded,
								},
								ID:             fmt.Sprintf("my-root-event-id-1"),
								Shkeptncontext: fmt.Sprintf("my-keptn-context-1"),
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetFinishedEventType("staging.delivery")),
							},
							{
								Data: keptnv2.EvaluationFinishedEventData{
									EventData: keptnv2.EventData{
										Project: "my-project",
										Stage:   "staging",
										Service: "my-service",
										Result:  keptnv2.ResultFailed,
										Status:  keptnv2.StatusSucceeded,
									},
									Evaluation: keptnv2.EvaluationDetails{
										Score: 0,
									},
								},
								ID:             "staging-task-2-finished-id",
								Triggeredid:    "staging-task-2-triggered-id",
								Shkeptncontext: "my-keptn-context-1",
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName)),
							},
							{
								Data: keptnv2.EventData{
									Project: "my-project",
									Stage:   "staging",
									Service: "my-service",
								},
								ID:             "staging-task-2-started-id",
								Triggeredid:    "staging-task-2-triggered-id",
								Shkeptncontext: "my-keptn-context-1",
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetStartedEventType(keptnv2.EvaluationTaskName)),
							},
							{
								Data: keptnv2.DeploymentTriggeredEventData{
									EventData: keptnv2.EventData{
										Project: "my-project",
										Stage:   "staging",
										Service: "my-service",
									},
								},
								ID:             "staging-task-2-triggered-id",
								Shkeptncontext: "my-keptn-context-1",
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName)),
							},
							{
								Data: keptnv2.EventData{
									Project: "my-project",
									Stage:   "staging",
									Service: "my-service",
									Result:  keptnv2.ResultPass,
									Status:  keptnv2.StatusSucceeded,
								},
								ID:             "staging-task-1-started-id",
								Triggeredid:    "staging-task-1-triggered-id",
								Shkeptncontext: "my-keptn-context-1",
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName)),
							},
							{
								Data: keptnv2.EventData{
									Project: "my-project",
									Stage:   "staging",
									Service: "my-service",
								},
								ID:             "staging-task-1-started-id",
								Triggeredid:    "staging-task-1-triggered-id",
								Shkeptncontext: "my-keptn-context-1",
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetStartedEventType(keptnv2.DeploymentTaskName)),
							},
							{
								Data: keptnv2.DeploymentTriggeredEventData{
									EventData: keptnv2.EventData{
										Project: "my-project",
										Stage:   "staging",
										Service: "my-service",
									},
									ConfigurationChange: keptnv2.ConfigurationChange{Values: map[string]interface{}{
										"image": "my-image",
									}},
								},
								ID:             "staging-task-1-triggered-id",
								Shkeptncontext: "my-keptn-context-1",
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName)),
							},
							{
								Data: keptnv2.EventData{
									Project: "my-project",
									Stage:   "staging",
									Service: "my-service",
								},
								ID:             fmt.Sprintf("my-staging-triggered-id"),
								Shkeptncontext: fmt.Sprintf("my-keptn-context-1"),
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetTriggeredEventType("staging.delivery")),
							},
							// dev stage finished
							{
								Data: keptnv2.EventData{
									Project: "my-project",
									Stage:   "dev",
									Service: "my-service",
									Result:  keptnv2.ResultPass,
									Status:  keptnv2.StatusSucceeded,
								},
								ID:             fmt.Sprintf("my-root-event-id-1"),
								Shkeptncontext: fmt.Sprintf("my-keptn-context-1"),
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetFinishedEventType("dev.delivery")),
							},
							{
								Data: keptnv2.EvaluationFinishedEventData{
									EventData: keptnv2.EventData{
										Project: "my-project",
										Stage:   "dev",
										Service: "my-service",
										Result:  keptnv2.ResultPass,
										Status:  keptnv2.StatusSucceeded,
									},
									Evaluation: keptnv2.EvaluationDetails{
										Score: 100,
									},
								},
								ID:             "dev-task-2-finished-id",
								Triggeredid:    "dev-task-2-triggered-id",
								Shkeptncontext: "my-keptn-context-1",
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName)),
							},
							{
								Data: keptnv2.EventData{
									Project: "my-project",
									Stage:   "dev",
									Service: "my-service",
								},
								ID:             "dev-task-2-started-id",
								Triggeredid:    "dev-task-2-triggered-id",
								Shkeptncontext: "my-keptn-context-1",
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetStartedEventType(keptnv2.EvaluationTaskName)),
							},
							{
								Data: keptnv2.DeploymentTriggeredEventData{
									EventData: keptnv2.EventData{
										Project: "my-project",
										Stage:   "dev",
										Service: "my-service",
									},
									ConfigurationChange: keptnv2.ConfigurationChange{Values: map[string]interface{}{
										"image": "my-image",
									}},
								},
								ID:             "dev-task-2-triggered-id",
								Shkeptncontext: "my-keptn-context-1",
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName)),
							},
							{
								Data: keptnv2.EventData{
									Project: "my-project",
									Stage:   "dev",
									Service: "my-service",
									Result:  keptnv2.ResultPass,
									Status:  keptnv2.StatusSucceeded,
								},
								ID:             "dev-task-1-started-id",
								Triggeredid:    "dev-task-1-triggered-id",
								Shkeptncontext: "my-keptn-context-1",
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName)),
							},
							{
								Data: keptnv2.EventData{
									Project: "my-project",
									Stage:   "dev",
									Service: "my-service",
								},
								ID:             "dev-task-1-started-id",
								Triggeredid:    "dev-task-1-triggered-id",
								Shkeptncontext: "my-keptn-context-1",
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetStartedEventType(keptnv2.DeploymentTaskName)),
							},
							{
								Data: keptnv2.DeploymentTriggeredEventData{
									EventData: keptnv2.EventData{
										Project: "my-project",
										Stage:   "dev",
										Service: "my-service",
									},
									ConfigurationChange: keptnv2.ConfigurationChange{Values: map[string]interface{}{
										"image": "my-image",
									}},
								},
								ID:             "dev-task-1-triggered-id",
								Shkeptncontext: "my-keptn-context-1",
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName)),
							},
							{
								Data: keptnv2.EventData{
									Project: "my-project",
									Service: "my-service",
									Stage:   "dev",
								},
								ID:             fmt.Sprintf("my-root-event-id-1"),
								Shkeptncontext: fmt.Sprintf("my-keptn-context-1"),
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetTriggeredEventType("dev.delivery")),
							},
						}, nil
					},
					GetRootEventsFunc: func(params models.GetRootEventParams) (*models.GetEventsResult, error) {
						return &models.GetEventsResult{Events: []models.Event{
							{
								Data: keptnv2.EventData{
									Project: "my-project",
									Stage:   "dev",
									Service: "my-service",
								},
								ID:             fmt.Sprintf("my-root-event-id-1"),
								Shkeptncontext: fmt.Sprintf("my-keptn-context-1"),
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetTriggeredEventType("dev.delivery")),
							},
						}}, nil
					},
				},
				taskSequenceRepo: &db_mock.SequenceStateRepoMock{
					CreateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{}, nil
					},
				},
				projectRepo: &db_mock.ProjectRepoMock{
					GetProjectsFunc: func() ([]*models.ExpandedProject, error) {
						return []*models.ExpandedProject{
							{
								ProjectName: "my-project",
							},
						}, nil
					},
				},
			},
			wantSequenceState: models.SequenceState{
				Name:           "delivery",
				Service:        "my-service",
				Project:        "my-project",
				Time:           timestampForAllEvents,
				Shkeptncontext: "my-keptn-context-1",
				State:          models.SequenceFinished,
				Stages: []models.SequenceStateStage{
					{
						Name:  "dev",
						Image: "my-image",
						LatestEvaluation: &models.SequenceStateEvaluation{
							Result: string(keptnv2.ResultPass),
							Score:  100,
						},
						LatestEvent: &models.SequenceStateEvent{
							Type: keptnv2.GetFinishedEventType("dev.delivery"),
							ID:   "my-root-event-id-1",
							Time: timestampForAllEvents,
						},
						LatestFailedEvent: nil,
					},
					{
						Name:  "staging",
						Image: "my-image",
						LatestEvaluation: &models.SequenceStateEvaluation{
							Result: string(keptnv2.ResultFailed),
							Score:  0,
						},
						LatestEvent: &models.SequenceStateEvent{
							Type: keptnv2.GetFinishedEventType("staging.delivery"),
							ID:   "my-root-event-id-1",
							Time: timestampForAllEvents,
						},
						LatestFailedEvent: &models.SequenceStateEvent{
							Type: keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName),
							ID:   "staging-task-2-finished-id",
							Time: timestampForAllEvents,
						},
					},
				},
			},
		},
		{
			name: "recreate task sequence state - unfinished sequence",
			fields: fields{
				eventRepo: &db_mock.EventRepoMock{
					GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
						return []models.Event{
							{
								Data: keptnv2.EventData{
									Project: "my-project",
									Stage:   "dev",
									Service: "my-service",
									Result:  keptnv2.ResultPass,
									Status:  keptnv2.StatusSucceeded,
								},
								ID:             "dev-task-1-finished-id",
								Triggeredid:    "dev-task-1-triggered-id",
								Shkeptncontext: "my-keptn-context-1",
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName)),
							},
							{
								Data: keptnv2.EventData{
									Project: "my-project",
									Stage:   "dev",
									Service: "my-service",
								},
								ID:             "dev-task-1-started-id",
								Triggeredid:    "dev-task-1-triggered-id",
								Shkeptncontext: "my-keptn-context-1",
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetStartedEventType(keptnv2.DeploymentTaskName)),
							},
							{
								Data: keptnv2.DeploymentTriggeredEventData{
									EventData: keptnv2.EventData{
										Project: "my-project",
										Stage:   "dev",
										Service: "my-service",
									},
									ConfigurationChange: keptnv2.ConfigurationChange{Values: map[string]interface{}{
										"image": "my-image",
									}},
								},
								ID:             "dev-task-1-triggered-id",
								Shkeptncontext: "my-keptn-context-1",
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName)),
							},
							{
								Data: keptnv2.EventData{
									Project: "my-project",
									Service: "my-service",
									Stage:   "dev",
								},
								ID:             fmt.Sprintf("my-root-event-id-1"),
								Shkeptncontext: fmt.Sprintf("my-keptn-context-1"),
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetTriggeredEventType("dev.delivery")),
							},
						}, nil
					},
					GetRootEventsFunc: func(params models.GetRootEventParams) (*models.GetEventsResult, error) {
						return &models.GetEventsResult{Events: []models.Event{
							{
								Data: keptnv2.EventData{
									Project: "my-project",
									Stage:   "dev",
									Service: "my-service",
								},
								ID:             fmt.Sprintf("my-root-event-id-1"),
								Shkeptncontext: fmt.Sprintf("my-keptn-context-1"),
								Time:           timestampForAllEvents,
								Type:           common.Stringp(keptnv2.GetTriggeredEventType("dev.delivery")),
							},
						}}, nil
					},
				},
				taskSequenceRepo: &db_mock.SequenceStateRepoMock{
					CreateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{}, nil
					},
				},
				projectRepo: &db_mock.ProjectRepoMock{
					GetProjectsFunc: func() ([]*models.ExpandedProject, error) {
						return []*models.ExpandedProject{
							{
								ProjectName: "my-project",
							},
						}, nil
					},
				},
			},
			wantSequenceState: models.SequenceState{
				Name:           "delivery",
				Service:        "my-service",
				Project:        "my-project",
				Time:           timestampForAllEvents,
				Shkeptncontext: "my-keptn-context-1",
				State:          models.SequenceTriggeredState,
				Stages: []models.SequenceStateStage{
					{
						Name:             "dev",
						Image:            "my-image",
						LatestEvaluation: nil,
						LatestEvent: &models.SequenceStateEvent{
							Type: keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName),
							ID:   "dev-task-1-finished-id",
							Time: timestampForAllEvents,
						},
						LatestFailedEvent: nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := handler.NewSequenceMigrator(tt.fields.eventRepo, tt.fields.taskSequenceRepo, tt.fields.projectRepo)

			sm.MigrateSequences()

			require.Len(t, tt.fields.taskSequenceRepo.CreateSequenceStateCalls(), 1)
			gotState := tt.fields.taskSequenceRepo.CreateSequenceStateCalls()[0].State
			require.Equal(t, tt.wantSequenceState, gotState)
		})
	}
}
