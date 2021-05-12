package sequencehooks_test

import (
	"errors"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/handler/sequencehooks"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"testing"
)

type SequenceStateMVTestFields struct {
	SequenceStateRepo *db_mock.SequenceStateRepoMock
}

func TestSequenceStateMaterializedView_OnSequenceFinished(t *testing.T) {

	type args struct {
		event models.Event
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
				event: models.Event{
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
				event: models.Event{
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
				event: models.Event{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smv := &sequencehooks.SequenceStateMaterializedView{
				SequenceStateRepo: tt.fields.SequenceStateRepo,
			}
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
			event := models.Event{
				Data:           tt.eventData,
				ID:             tt.eventId,
				Shkeptncontext: tt.keptnContext,
				Type:           &tt.eventType,
			}

			smv := sequencehooks.NewSequenceStateMaterializedView(tt.fields.SequenceStateRepo)

			smv.OnSequenceTaskFinished(event)

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

func TestSequenceStateMaterializedView_OnSequenceTaskStarted(t *testing.T) {
	type args struct {
		event models.Event
	}
	tests := []struct {
		name   string
		fields SequenceStateMVTestFields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smv := &sequencehooks.SequenceStateMaterializedView{
				SequenceStateRepo: tt.fields.SequenceStateRepo,
			}
			smv.OnSequenceTaskStarted(tt.args.event)
		})
	}
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
			expectUpdateStateToBeCalled: true,
			project:                     "my-project",
			service:                     "my-service",
			stage:                       "my-stage",
			keptnContext:                "my-context",
			eventType:                   "my-event-type",
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
			event := models.Event{
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

			smv := sequencehooks.NewSequenceStateMaterializedView(tt.fields.SequenceStateRepo)

			smv.OnSequenceTaskTriggered(event)

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
		project                     string
		service                     string
		stage                       string
		keptnContext                string
		sequenceName                string
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
			event := models.Event{
				Data: keptnv2.EventData{
					Project: tt.project,
					Stage:   tt.stage,
					Service: tt.service,
				},
				Shkeptncontext: tt.keptnContext,
				Type:           common.Stringp("sh.keptn.event." + tt.stage + "." + tt.sequenceName + ".triggered"),
			}

			smv := sequencehooks.NewSequenceStateMaterializedView(tt.fields.SequenceStateRepo)

			smv.OnSequenceTriggered(event)
			if tt.expectCreateStateToBeCalled {
				require.Equal(t, 1, len(tt.fields.SequenceStateRepo.CreateSequenceStateCalls()))
				call := tt.fields.SequenceStateRepo.CreateSequenceStateCalls()[0]
				require.Equal(t, tt.project, call.State.Project)
				require.Equal(t, tt.service, call.State.Service)
				require.Equal(t, tt.sequenceName, call.State.Name)
				require.Equal(t, tt.keptnContext, call.State.Shkeptncontext)
				require.Equal(t, "triggered", call.State.State)

			} else {
				require.Equal(t, 0, len(tt.fields.SequenceStateRepo.CreateSequenceStateCalls()))
			}
		})
	}
}
