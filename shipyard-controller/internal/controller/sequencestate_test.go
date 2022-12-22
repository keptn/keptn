package controller_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/keptn/keptn/shipyard-controller/internal/common"
	"github.com/keptn/keptn/shipyard-controller/internal/controller"
	"github.com/keptn/keptn/shipyard-controller/internal/db"
	db_mock "github.com/keptn/keptn/shipyard-controller/internal/db/mock"

	"github.com/keptn/go-utils/pkg/api/models"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	scmodels "github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
)

type SequenceStateMVTestFields struct {
	SequenceStateRepo *db_mock.SequenceStateRepoMock
}

func TestSequenceStateMaterializedView_OnSequenceStarted(t *testing.T) {
	type args struct {
		event models.KeptnContextExtendedCE
	}
	tests := []struct {
		name                   string
		fields                 SequenceStateMVTestFields
		args                   args
		expectUpdateToBeCalled bool
	}{
		{
			name: "start sequence",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{
							States: []models.SequenceState{
								{
									Name:           "my-sequence",
									Service:        "my-service",
									Project:        "my-project",
									Shkeptncontext: "my-context",
									State:          "triggered",
									Stages:         nil,
								},
							},
						}, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			args: args{
				event: models.KeptnContextExtendedCE{
					Data: keptnv2.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
					},
					Shkeptncontext: "my-context",
					Type:           common.Stringp("my-type"),
				},
			},
			expectUpdateToBeCalled: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smv := controller.NewSequenceStateMaterializedView(tt.fields.SequenceStateRepo)
			smv.OnSequenceStarted(tt.args.event)

			if tt.expectUpdateToBeCalled {
				require.NotEmpty(t, tt.fields.SequenceStateRepo.UpdateSequenceStateCalls())
				require.Equal(t, models.SequenceStartedState, tt.fields.SequenceStateRepo.UpdateSequenceStateCalls()[0].State.State)
			} else {
				require.Empty(t, tt.fields.SequenceStateRepo.UpdateSequenceStateCalls())
			}
		})
	}
}

func TestSequenceStateMaterializedView_OnSequenceWaiting(t *testing.T) {
	type args struct {
		event models.KeptnContextExtendedCE
	}
	tests := []struct {
		name                   string
		fields                 SequenceStateMVTestFields
		args                   args
		expectUpdateToBeCalled bool
	}{
		{
			name: "start sequence",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{
							States: []models.SequenceState{
								{
									Name:           "my-sequence",
									Service:        "my-service",
									Project:        "my-project",
									Shkeptncontext: "my-context",
									State:          "triggered",
									Stages:         nil,
								},
							},
						}, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			args: args{
				event: models.KeptnContextExtendedCE{
					Data: keptnv2.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
					},
					Shkeptncontext: "my-context",
					Type:           common.Stringp("my-type"),
				},
			},
			expectUpdateToBeCalled: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smv := controller.NewSequenceStateMaterializedView(tt.fields.SequenceStateRepo)
			smv.OnSequenceWaiting(tt.args.event)

			if tt.expectUpdateToBeCalled {
				require.NotEmpty(t, tt.fields.SequenceStateRepo.UpdateSequenceStateCalls())
				require.Equal(t, models.SequenceWaitingState, tt.fields.SequenceStateRepo.UpdateSequenceStateCalls()[0].State.State)
			} else {
				require.Empty(t, tt.fields.SequenceStateRepo.UpdateSequenceStateCalls())
			}
		})
	}
}

func TestSequenceStateMaterializedView_OnSequenceTimeOud(t *testing.T) {
	type args struct {
		event models.KeptnContextExtendedCE
	}
	tests := []struct {
		name                   string
		fields                 SequenceStateMVTestFields
		args                   args
		expectUpdateToBeCalled bool
	}{
		{
			name: "sequence timed out",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{
							States: []models.SequenceState{
								{
									Name:           "my-sequence",
									Service:        "my-service",
									Project:        "my-project",
									Shkeptncontext: "my-context",
									State:          "triggered",
									Stages:         nil,
								},
							},
						}, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			args: args{
				event: models.KeptnContextExtendedCE{
					Data: keptnv2.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
					},
					Shkeptncontext: "my-context",
					Type:           common.Stringp("my-type"),
				},
			},
			expectUpdateToBeCalled: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smv := controller.NewSequenceStateMaterializedView(tt.fields.SequenceStateRepo)
			smv.OnSequenceTimeout(tt.args.event)

			if tt.expectUpdateToBeCalled {
				require.NotEmpty(t, tt.fields.SequenceStateRepo.UpdateSequenceStateCalls())
				require.Equal(t, models.TimedOut, tt.fields.SequenceStateRepo.UpdateSequenceStateCalls()[0].State.State)
			} else {
				require.Empty(t, tt.fields.SequenceStateRepo.UpdateSequenceStateCalls())
			}
		})
	}
}

func TestSequenceStateMaterializedView_OnSequenceFinished(t *testing.T) {

	type args struct {
		event models.KeptnContextExtendedCE
	}
	tests := []struct {
		name                   string
		fields                 SequenceStateMVTestFields
		args                   args
		expectUpdateToBeCalled bool
	}{
		{
			name: "finish sequence",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{
							States: []models.SequenceState{
								{
									Name:           "my-sequence",
									Service:        "my-service",
									Project:        "my-project",
									Shkeptncontext: "my-context",
									State:          "triggered",
									Stages: []models.SequenceStateStage{
										{
											Name:  "dev",
											State: "succeeded",
										},
										{
											Name:  "dev",
											State: "succeeded",
										},
									},
								},
							},
						}, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			args: args{
				event: models.KeptnContextExtendedCE{
					Data: keptnv2.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
					},
					Shkeptncontext: "my-context",
					Type:           common.Stringp("my-type"),
				},
			},
			expectUpdateToBeCalled: true,
		},
		{
			name: "try to finish sequence - not all stages finished yet",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{
							States: []models.SequenceState{
								{
									Name:           "my-sequence",
									Service:        "my-service",
									Project:        "my-project",
									Shkeptncontext: "my-context",
									State:          "triggered",
									Stages: []models.SequenceStateStage{
										{
											Name:  "dev",
											State: "succeeded",
										},
										{
											Name:  "dev",
											State: "triggered",
										},
									},
								},
							},
						}, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			args: args{
				event: models.KeptnContextExtendedCE{
					Data: keptnv2.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
					},
					Shkeptncontext: "my-context",
					Type:           common.Stringp("my-type"),
				},
			},
			expectUpdateToBeCalled: false,
		},
		{
			name: "invalid event scope - do not update",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{
							States: []models.SequenceState{
								{
									Name:           "my-sequence",
									Service:        "my-service",
									Project:        "my-project",
									Shkeptncontext: "my-context",
									State:          "triggered",
									Stages:         nil,
								},
							},
						}, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			args: args{
				event: models.KeptnContextExtendedCE{
					Data:           keptnv2.EventData{},
					Shkeptncontext: "my-context",
				},
			},
			expectUpdateToBeCalled: false,
		},
		{
			name: "cannot find sequence - do not update",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return nil, errors.New("oops")
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			args: args{
				event: models.KeptnContextExtendedCE{
					Data: keptnv2.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
					},
					Shkeptncontext: "my-context",
				},
			},
			expectUpdateToBeCalled: false,
		},
		{
			name: "cannot find sequence - do not update (2)",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{}, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			args: args{
				event: models.KeptnContextExtendedCE{
					Data: keptnv2.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
					},
					Shkeptncontext: "my-context",
					Type:           common.Stringp("my-type"),
				},
			},
			expectUpdateToBeCalled: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smv := controller.NewSequenceStateMaterializedView(tt.fields.SequenceStateRepo)
			smv.OnSequenceFinished(tt.args.event)

			if tt.expectUpdateToBeCalled {
				require.NotEmpty(t, tt.fields.SequenceStateRepo.UpdateSequenceStateCalls())
				require.Equal(t, "finished", tt.fields.SequenceStateRepo.UpdateSequenceStateCalls()[0].State.State)
			} else {
				require.Empty(t, tt.fields.SequenceStateRepo.UpdateSequenceStateCalls())
			}
		})
	}
}

func TestSequenceStateMaterializedView_OnSequenceTaskFinished(t *testing.T) {
	tests := []struct {
		name                         string
		fields                       SequenceStateMVTestFields
		eventId                      string
		eventData                    keptncommon.EventProperties
		keptnContext                 string
		eventType                    string
		eventSource                  string
		expectUpdateStateToBeCalled  bool
		expectEvaluationToBeUpdated  bool
		expectFailedEventToBeUpdated bool
	}{
		{
			name: "update evaluation",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{
							States: []models.SequenceState{
								{
									Name:           "my-sequence",
									Service:        "my-service",
									Project:        "my-project",
									Shkeptncontext: "my-context",
									State:          "triggered",
									Stages:         nil,
								},
							},
						}, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			expectUpdateStateToBeCalled: true,
			keptnContext:                "my-context",
			eventType:                   keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName),
			eventSource:                 controller.SequenceEvaluationService,
			eventId:                     "my-id",
			eventData: &keptnv2.EvaluationFinishedEventData{
				EventData: keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-state",
					Service: "my-service",
					Result:  keptnv2.ResultPass,
				},
				Evaluation: keptnv2.EvaluationDetails{
					Score: 100.0,
				},
			},
			expectEvaluationToBeUpdated: true,
		},
		{
			name: "update evaluation fails: not a lighthouse finished event",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{
							States: []models.SequenceState{
								{
									Name:           "my-sequence",
									Service:        "my-service",
									Project:        "my-project",
									Shkeptncontext: "my-context",
									State:          "triggered",
									Stages:         nil,
								},
							},
						}, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			expectUpdateStateToBeCalled: true,
			keptnContext:                "my-context",
			eventType:                   keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName),
			eventSource:                 "not-lighthouse",
			eventId:                     "my-id",
			eventData: &keptnv2.EvaluationFinishedEventData{
				EventData: keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-state",
					Service: "my-service",
					Result:  keptnv2.ResultPass,
				},
				Evaluation: keptnv2.EvaluationDetails{
					Score: 100.0,
				},
			},
			expectEvaluationToBeUpdated: false,
		},
		{
			name: "failed task",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{
							States: []models.SequenceState{
								{
									Name:           "my-sequence",
									Service:        "my-service",
									Project:        "my-project",
									Shkeptncontext: "my-context",
									State:          "triggered",
									Stages:         nil,
								},
							},
						}, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			expectUpdateStateToBeCalled: true,
			keptnContext:                "my-context",
			eventType:                   keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
			eventSource:                 controller.SequenceEvaluationService,
			eventId:                     "my-id",
			eventData: &keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-state",
				Service: "my-service",
				Result:  keptnv2.ResultFailed,
			},
			expectFailedEventToBeUpdated: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := models.KeptnContextExtendedCE{
				Data:           tt.eventData,
				ID:             tt.eventId,
				Shkeptncontext: tt.keptnContext,
				Type:           &tt.eventType,
				Source:         &tt.eventSource,
			}

			smv := controller.NewSequenceStateMaterializedView(tt.fields.SequenceStateRepo)

			smv.OnSequenceTaskEvent(event)

			if tt.expectUpdateStateToBeCalled {
				require.Equal(t, 1, len(tt.fields.SequenceStateRepo.UpdateSequenceStateCalls()))
				call := tt.fields.SequenceStateRepo.UpdateSequenceStateCalls()[0]
				require.Equal(t, tt.eventData.GetProject(), call.State.Project)
				require.Equal(t, tt.eventData.GetService(), call.State.Service)
				require.Equal(t, tt.keptnContext, call.State.Shkeptncontext)
				require.Equal(t, tt.eventType, call.State.Stages[0].LatestEvent.Type)
				require.Equal(t, tt.eventId, call.State.Stages[0].LatestEvent.ID)

				if tt.expectEvaluationToBeUpdated {
					require.NotEmpty(t, 1, call.State.Stages[0].LatestEvaluation)
					evaluationFinishedData := tt.eventData.(*keptnv2.EvaluationFinishedEventData)
					require.Equal(t, evaluationFinishedData.Evaluation.Score, call.State.Stages[0].LatestEvaluation.Score)
					require.Equal(t, string(evaluationFinishedData.Result), call.State.Stages[0].LatestEvaluation.Result)
				}
				if tt.expectFailedEventToBeUpdated {
					require.NotEmpty(t, call.State.Stages[0].LatestFailedEvent)
					require.Equal(t, tt.eventId, call.State.Stages[0].LatestFailedEvent.ID)
				}

			} else {
				require.Equal(t, 0, len(tt.fields.SequenceStateRepo.UpdateSequenceStateCalls()))
			}

		})
	}
}

func TestSequenceStateMaterializedView_MultipleScoresOnSequenceFinished(t *testing.T) {

	eventType := keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName)
	badSource := "not the right one"
	goodSource := controller.SequenceEvaluationService

	events := []models.KeptnContextExtendedCE{
		{Shkeptncontext: "my-context",
			Type:   &eventType,
			Source: &badSource,
			ID:     "my-id1",
			Data: &keptnv2.EvaluationFinishedEventData{
				EventData: keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-state",
					Service: "my-service",
					Result:  keptnv2.ResultPass,
				},
				Evaluation: keptnv2.EvaluationDetails{
					Score: 75.0,
				},
			},
		},
		{Shkeptncontext: "my-context",
			Type:   &eventType,
			Source: &goodSource,
			ID:     "my-id2",
			Data: &keptnv2.EvaluationFinishedEventData{
				EventData: keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-state",
					Service: "my-service",
					Result:  keptnv2.ResultPass,
				},
				Evaluation: keptnv2.EvaluationDetails{
					Score: 100.0,
				},
			},
		},
		{
			Shkeptncontext: "my-context",
			Type:           &eventType,
			Source:         &badSource,
			ID:             "my-id3",
			Data: &keptnv2.EvaluationFinishedEventData{
				EventData: keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-state",
					Service: "my-service",
					Result:  keptnv2.ResultPass,
				},
				Evaluation: keptnv2.EvaluationDetails{
					Score: 55.0,
				},
			},
		},
	}

	t.Run("multiple score test", func(t *testing.T) {

		SequenceStateRepo := &db_mock.SequenceStateRepoMock{
			FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
				return &models.SequenceStates{
					States: []models.SequenceState{
						{
							Name:           "my-sequence",
							Service:        "my-service",
							Project:        "my-project",
							Shkeptncontext: "my-context",
							State:          "triggered",
							Stages:         nil,
						},
					},
				}, nil
			},
			UpdateSequenceStateFunc: func(state models.SequenceState) error {
				return nil
			},
		}
		smv := controller.NewSequenceStateMaterializedView(SequenceStateRepo)

		for i, event := range events {

			smv.OnSequenceTaskEvent(event)
			call := SequenceStateRepo.UpdateSequenceStateCalls()[i]
			require.Equal(t, event.ID, call.State.Stages[0].LatestEvent.ID)
			if *event.Source == goodSource {
				goodScore := event.Data.(*keptnv2.EvaluationFinishedEventData).Evaluation.Score
				require.Equal(t, goodScore, call.State.Stages[0].LatestEvaluation.Score)
			} else {
				require.Nil(t, call.State.Stages[0].LatestEvaluation)
			}
		}
	})
}

func TestSequenceStateMaterializedView_OnSequenceTaskTriggered(t *testing.T) {
	tests := []struct {
		name                        string
		fields                      SequenceStateMVTestFields
		expectUpdateStateToBeCalled bool
		expectImageToBeUpdated      bool
		project                     string
		service                     string
		stage                       string
		eventId                     string
		keptnContext                string
		eventType                   string
	}{
		{
			name: "update sequence state - insert new stage",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{
							States: []models.SequenceState{
								{
									Name:           "my-sequence",
									Service:        "my-service",
									Project:        "my-project",
									Shkeptncontext: "my-context",
									State:          "triggered",
									Stages:         nil,
								},
							},
						}, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			expectImageToBeUpdated:      true,
			expectUpdateStateToBeCalled: true,
			project:                     "my-project",
			service:                     "my-service",
			stage:                       "my-stage",
			keptnContext:                "my-context",
			eventType:                   keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
			eventId:                     "my-id",
		},
		{
			name: "update sequence state with existing stage",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{
							States: []models.SequenceState{
								{
									Name:           "my-sequence",
									Service:        "my-service",
									Project:        "my-project",
									Shkeptncontext: "my-context",
									State:          "triggered",
									Stages: []models.SequenceStateStage{
										{
											Name: "my-stage",
											LatestEvent: &models.SequenceStateEvent{
												Type: "my-old-event-type",
												ID:   "my-old-event-id",
											},
										},
									},
								},
							},
						}, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			expectUpdateStateToBeCalled: true,
			project:                     "my-project",
			service:                     "my-service",
			stage:                       "my-stage",
			keptnContext:                "my-context",
			eventType:                   "my-event-type",
			eventId:                     "my-id",
		},
		{
			name: "invalid event scope",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{},
			},
			expectUpdateStateToBeCalled: false,
			project:                     "",
			service:                     "my-service",
			stage:                       "my-stage",
			keptnContext:                "my-context",
			eventType:                   "my-event-type",
			eventId:                     "my-id",
		},
		{
			name: "find state returns error - do not call update",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return nil, errors.New("oops")
					},
				},
			},
			expectUpdateStateToBeCalled: false,
			project:                     "my-project",
			service:                     "my-service",
			stage:                       "my-stage",
			keptnContext:                "my-context",
			eventType:                   "my-event-type",
			eventId:                     "my-id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := models.KeptnContextExtendedCE{
				Data: keptnv2.EventData{
					Project: tt.project,
					Stage:   tt.stage,
					Service: tt.service,
				},
				ID:             tt.eventId,
				Shkeptncontext: tt.keptnContext,
				Type:           &tt.eventType,
			}

			if tt.expectImageToBeUpdated {
				event.Data = &keptnv2.DeploymentTriggeredEventData{
					EventData: keptnv2.EventData{
						Project: tt.project,
						Stage:   tt.stage,
						Service: tt.service,
						Result:  keptnv2.ResultPass,
					},
					ConfigurationChange: keptnv2.ConfigurationChange{Values: map[string]interface{}{"image": "my-image"}},
				}
			}

			smv := controller.NewSequenceStateMaterializedView(tt.fields.SequenceStateRepo)

			smv.OnSequenceTaskEvent(event)

			if tt.expectUpdateStateToBeCalled {
				require.Equal(t, 1, len(tt.fields.SequenceStateRepo.UpdateSequenceStateCalls()))
				call := tt.fields.SequenceStateRepo.UpdateSequenceStateCalls()[0]
				require.Equal(t, tt.project, call.State.Project)
				require.Equal(t, tt.service, call.State.Service)
				require.Equal(t, tt.keptnContext, call.State.Shkeptncontext)
				require.Equal(t, tt.eventType, call.State.Stages[0].LatestEvent.Type)
				require.Equal(t, tt.eventId, call.State.Stages[0].LatestEvent.ID)

				if tt.expectImageToBeUpdated {
					require.NotEmpty(t, 1, call.State.Stages[0].Image)
					deploymentData := event.Data.(*keptnv2.DeploymentTriggeredEventData)
					require.Equal(t, deploymentData.ConfigurationChange.Values["image"], call.State.Stages[0].Image)
				}

			} else {
				require.Equal(t, 0, len(tt.fields.SequenceStateRepo.UpdateSequenceStateCalls()))
			}
		})
	}
}

func TestSequenceStateMaterializedView_OnSequenceTriggered(t *testing.T) {

	tests := []struct {
		name                        string
		fields                      SequenceStateMVTestFields
		expectCreateStateToBeCalled bool
		expectUpdateStateToBeCalled bool
		project                     string
		service                     string
		stage                       string
		keptnContext                string
		sequenceName                string
		problemTitle                string
	}{
		{
			name: "create a new sequence state",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					CreateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			expectCreateStateToBeCalled: true,
			project:                     "my-project",
			service:                     "my-service",
			stage:                       "my-stage",
			keptnContext:                "my-context",
			sequenceName:                "my-sequence",
		},
		{
			name: "create a new remediation sequence",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					CreateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			expectCreateStateToBeCalled: true,
			project:                     "my-project",
			service:                     "my-service",
			stage:                       "my-stage",
			keptnContext:                "my-context",
			sequenceName:                "remediation",
			problemTitle:                "This is a very serious issue",
		},
		{
			name: "no project available",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{},
			},
			expectCreateStateToBeCalled: false,
			project:                     "",
			service:                     "my-service",
			stage:                       "my-stage",
			keptnContext:                "my-context",
			sequenceName:                "my-sequence",
		},
		{
			name: "wrong event type",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{},
			},
			expectCreateStateToBeCalled: false,
			project:                     "",
			service:                     "my-service",
			stage:                       "my-stage",
			keptnContext:                "my-context",
			sequenceName:                ".",
		},
		{
			name: "state already exists",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					CreateSequenceStateFunc: func(state models.SequenceState) error {
						return db.ErrStateAlreadyExists
					},
					FindSequenceStatesFunc: func(filter apimodels.StateFilter) (*apimodels.SequenceStates, error) {
						return &apimodels.SequenceStates{States: []apimodels.SequenceState{
							{
								Name:           "my-sequence",
								Service:        "my-service",
								Project:        "my-project",
								Shkeptncontext: "my-context",
								State:          "",
								Stages: []apimodels.SequenceStateStage{
									{
										Name: "my-other-stage",
									},
								},
							},
						}}, nil
					},
					UpdateSequenceStateFunc: func(state apimodels.SequenceState) error {
						return nil
					},
				},
			},
			expectCreateStateToBeCalled: true,
			project:                     "my-project",
			service:                     "my-service",
			stage:                       "my-stage",
			keptnContext:                "my-context",
			sequenceName:                "my-sequence",
		},
		{
			name: "state already exists - updating fails",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					CreateSequenceStateFunc: func(state models.SequenceState) error {
						return db.ErrStateAlreadyExists
					},
					FindSequenceStatesFunc: func(filter apimodels.StateFilter) (*apimodels.SequenceStates, error) {
						return &apimodels.SequenceStates{States: []apimodels.SequenceState{
							{
								Name:           "my-sequence",
								Service:        "my-service",
								Project:        "my-project",
								Shkeptncontext: "my-context",
								State:          "",
								Stages: []apimodels.SequenceStateStage{
									{
										Name: "my-other-stage",
									},
								},
							},
						}}, nil
					},
					UpdateSequenceStateFunc: func(state apimodels.SequenceState) error {
						return errors.New("oops")
					},
				},
			},
			expectCreateStateToBeCalled: true,
			project:                     "my-project",
			service:                     "my-service",
			stage:                       "my-stage",
			keptnContext:                "my-context",
			sequenceName:                "my-sequence",
		},
		{
			name: "create state returns an error",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					CreateSequenceStateFunc: func(state models.SequenceState) error {
						return errors.New("oops")
					},
				},
			},
			expectCreateStateToBeCalled: true,
			project:                     "my-project",
			service:                     "my-service",
			stage:                       "my-stage",
			keptnContext:                "my-context",
			sequenceName:                "my-sequence",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Data := keptnv2.EventData{
				Project: tt.project,
				Stage:   tt.stage,
				Service: tt.service,
			}
			var event models.KeptnContextExtendedCE
			//construct a remediation event
			if tt.problemTitle != "" {
				event = models.KeptnContextExtendedCE{
					Data: keptnv2.GetActionTriggeredEventData{
						EventData: Data,
						Problem: keptnv2.ProblemDetails{
							ProblemTitle: tt.problemTitle,
						}},
					Shkeptncontext: tt.keptnContext,
					Type:           common.Stringp("sh.keptn.event." + tt.stage + "." + tt.sequenceName + ".triggered"),
				}
			} else {
				//construct a simple event
				event = models.KeptnContextExtendedCE{
					Data:           Data,
					Shkeptncontext: tt.keptnContext,
					Type:           common.Stringp("sh.keptn.event." + tt.stage + "." + tt.sequenceName + ".triggered"),
				}
			}
			smv := controller.NewSequenceStateMaterializedView(tt.fields.SequenceStateRepo)

			smv.OnSequenceTriggered(event)

			if tt.expectCreateStateToBeCalled {
				require.Equal(t, 1, len(tt.fields.SequenceStateRepo.CreateSequenceStateCalls()))
				call := tt.fields.SequenceStateRepo.CreateSequenceStateCalls()[0]
				require.Equal(t, tt.project, call.State.Project)
				require.Equal(t, tt.service, call.State.Service)
				require.Equal(t, tt.sequenceName, call.State.Name)
				require.Equal(t, tt.keptnContext, call.State.Shkeptncontext)
				require.Equal(t, "triggered", call.State.State)
				require.Equal(t, tt.problemTitle, call.State.ProblemTitle)

			} else {
				require.Equal(t, 0, len(tt.fields.SequenceStateRepo.CreateSequenceStateCalls()))
			}

			if tt.expectUpdateStateToBeCalled {
				require.Equal(t, 1, len(tt.fields.SequenceStateRepo.UpdateSequenceStateCalls()))
				call := tt.fields.SequenceStateRepo.UpdateSequenceStateCalls()[0]
				require.Equal(t, tt.project, call.State.Project)
				require.Equal(t, tt.service, call.State.Service)
				require.Equal(t, tt.stage, call.State.Stages[0])
			}
		})
	}
}

func TestSequenceStateMaterializedView_OnSequencePaused(t *testing.T) {
	tests := []struct {
		name                        string
		fields                      SequenceStateMVTestFields
		expectUpdateStateToBeCalled bool
		sequencePause               scmodels.EventScope
	}{
		{
			name: "overall sequence paused",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{
							States: []models.SequenceState{
								{
									Name:           "my-sequence",
									Service:        "my-service",
									Project:        "my-project",
									Shkeptncontext: "my-context",
									State:          "triggered",
									Stages:         nil,
								},
							},
						}, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			sequencePause: scmodels.EventScope{
				KeptnContext: "my-context",
				EventData:    keptnv2.EventData{Project: "my-project"},
			},
		},
		{
			name: "stage of sequence paused",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{
							States: []models.SequenceState{
								{
									Name:           "my-sequence",
									Service:        "my-service",
									Project:        "my-project",
									Shkeptncontext: "my-context",
									State:          "triggered",
									Stages: []models.SequenceStateStage{
										{
											Name: "my-stage",
										},
									},
								},
							},
						}, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			sequencePause: scmodels.EventScope{
				KeptnContext: "my-context",
				EventData:    keptnv2.EventData{Project: "my-project", Stage: "my-stage"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smv := controller.NewSequenceStateMaterializedView(tt.fields.SequenceStateRepo)

			smv.OnSequencePaused(tt.sequencePause)

			require.NotEmpty(t, tt.fields.SequenceStateRepo.UpdateSequenceStateCalls())
			if tt.sequencePause.Stage == "" {
				require.Equal(t, models.SequencePaused, tt.fields.SequenceStateRepo.UpdateSequenceStateCalls()[0].State.State)
			} else {
				require.Equal(t, models.SequencePaused, tt.fields.SequenceStateRepo.UpdateSequenceStateCalls()[0].State.Stages[0].State)
			}
		})
	}
}

func TestSequenceStateMaterializedView_OnSubSequenceFinished(t *testing.T) {
	type args struct {
		event models.KeptnContextExtendedCE
	}
	tests := []struct {
		name   string
		fields SequenceStateMVTestFields
		args   args
	}{
		{
			name: "abort subsequence",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{
							States: []models.SequenceState{
								{
									Name:           "my-sequence",
									Service:        "my-service",
									Project:        "my-project",
									Shkeptncontext: "my-context",
									State:          "triggered",
									Stages: []models.SequenceStateStage{
										{
											Name:              "my-stage",
											LatestEvent:       &models.SequenceStateEvent{},
											LatestFailedEvent: &models.SequenceStateEvent{},
										},
									},
								},
							},
						}, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			args: args{
				event: models.KeptnContextExtendedCE{
					Data: keptnv2.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
						Status:  keptnv2.StatusAborted,
					},
					Shkeptncontext: "my-context",
					Type:           common.Stringp("sh.keptn.event.dev.sequence.finished"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smv := controller.NewSequenceStateMaterializedView(tt.fields.SequenceStateRepo)
			smv.OnSubSequenceFinished(tt.args.event)
			require.Equal(t, "aborted", tt.fields.SequenceStateRepo.UpdateSequenceStateCalls()[0].State.Stages[0].State)
		})
	}
}

func TestSequenceStateMaterializedView_UpdateLastEventOfSequence(t *testing.T) {
	testSequence := &models.SequenceStates{
		States: []models.SequenceState{
			{
				Name:           "my-sequence",
				Service:        "my-service",
				Project:        "my-project",
				Shkeptncontext: "my-context",
				State:          "triggered",
				Stages: []models.SequenceStateStage{
					{
						Name:              "my-stage",
						LatestEvent:       &models.SequenceStateEvent{},
						LatestFailedEvent: &models.SequenceStateEvent{},
					},
				},
			},
		},
	}

	tests := []struct {
		name    string
		fields  SequenceStateMVTestFields
		event   apimodels.KeptnContextExtendedCE
		want    apimodels.SequenceState
		wantErr error
	}{
		{
			name: "pass",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return testSequence, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			event: models.KeptnContextExtendedCE{
				Data: keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
					Status:  keptnv2.StatusSucceeded,
					Result:  keptnv2.ResultPass,
				},
				Shkeptncontext: "my-context",
				Type:           common.Stringp("sh.keptn.event.dev.sequence.finished"),
			},
			want: apimodels.SequenceState{
				Project:        "my-project",
				Service:        "my-service",
				State:          "triggered",
				Shkeptncontext: "my-context",
				Stages: []apimodels.SequenceStateStage{
					{
						Name:  "my-stage",
						State: string(keptnv2.StatusSucceeded),
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "strange",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return testSequence, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			event: models.KeptnContextExtendedCE{
				Data: keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
					Status:  "strange",
					Result:  "result",
				},
				Shkeptncontext: "my-context",
				Type:           common.Stringp("sh.keptn.event.dev.sequence.finished"),
			},
			want: apimodels.SequenceState{
				Project:        "my-project",
				Service:        "my-service",
				State:          "triggered",
				Shkeptncontext: "my-context",
				Stages: []apimodels.SequenceStateStage{
					{
						Name:  "my-stage",
						State: "strange",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "aborted",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return testSequence, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			event: models.KeptnContextExtendedCE{
				Data: keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
					Status:  keptnv2.StatusAborted,
					Result:  keptnv2.ResultPass,
				},
				Shkeptncontext: "my-context",
				Type:           common.Stringp("sh.keptn.event.dev.sequence.finished"),
			},
			want: apimodels.SequenceState{
				Project:        "my-project",
				Service:        "my-service",
				State:          "triggered",
				Shkeptncontext: "my-context",
				Stages: []apimodels.SequenceStateStage{
					{
						Name:  "my-stage",
						State: "aborted",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "failed",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return testSequence, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			event: models.KeptnContextExtendedCE{
				Data: keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
					Status:  keptnv2.StatusSucceeded,
					Result:  keptnv2.ResultFailed,
				},
				Shkeptncontext: "my-context",
				Type:           common.Stringp("sh.keptn.event.dev.sequence.finished"),
			},
			want: apimodels.SequenceState{
				Project:        "my-project",
				Service:        "my-service",
				State:          "triggered",
				Shkeptncontext: "my-context",
				Stages: []apimodels.SequenceStateStage{
					{
						Name:  "my-stage",
						State: string(keptnv2.StatusSucceeded),
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "failed #2",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return testSequence, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			event: models.KeptnContextExtendedCE{
				Data: keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
					Status:  keptnv2.StatusErrored,
					Result:  keptnv2.ResultPass,
				},
				Shkeptncontext: "my-context",
				Type:           common.Stringp("sh.keptn.event.dev.sequence.finished"),
			},
			want: apimodels.SequenceState{
				Project:        "my-project",
				Service:        "my-service",
				State:          "triggered",
				Shkeptncontext: "my-context",
				Stages: []apimodels.SequenceStateStage{
					{
						Name:  "my-stage",
						State: string(keptnv2.StatusErrored),
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "fail no sequence state",
			fields: SequenceStateMVTestFields{
				SequenceStateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{}, nil
					},
					UpdateSequenceStateFunc: func(state models.SequenceState) error {
						return nil
					},
				},
			},
			event: models.KeptnContextExtendedCE{
				Data: keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
					Status:  keptnv2.StatusSucceeded,
					Result:  keptnv2.ResultPass,
				},
				Shkeptncontext: "my-context",
				Type:           common.Stringp("sh.keptn.event.dev.sequence.finished"),
			},
			wantErr: fmt.Errorf("could not find sequence state for keptnContext my-context"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smv := controller.NewSequenceStateMaterializedView(tt.fields.SequenceStateRepo)
			out, err := smv.UpdateLastEventOfSequence(tt.event)
			require.Equal(t, tt.wantErr, err)
			if err == nil {
				require.Equal(t, tt.want.Project, out.Project)
				require.Equal(t, tt.want.Service, out.Service)
				require.Equal(t, tt.want.State, out.State)
				require.Equal(t, tt.want.Shkeptncontext, out.Shkeptncontext)
				require.Equal(t, tt.want.Stages[0].Name, out.Stages[0].Name)
				require.Equal(t, tt.want.Stages[0].State, out.Stages[0].State)
			}
		})
	}
}
