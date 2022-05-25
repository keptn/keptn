package models

import (
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestSequenceExecution_GetNextTriggeredEventData(t *testing.T) {
	type fields struct {
		ID              string
		Sequence        keptnv2.Sequence
		Status          SequenceExecutionStatus
		Scope           EventScope
		InputProperties map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{
			name: "get initial triggered event - no input data",
			fields: fields{
				Sequence: keptnv2.Sequence{
					Name: "delivery",
					Tasks: []keptnv2.Task{
						{
							Name: "mytask",
							Properties: map[string]interface{}{
								"deploymentstrategy": "direct",
							},
						},
					},
				},
				Status: SequenceExecutionStatus{
					PreviousTasks: nil,
					CurrentTask:   TaskExecutionState{},
				},
				Scope: EventScope{
					EventData: keptnv2.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
					},
				},
				InputProperties: map[string]interface{}{
					"mytask": map[string]interface{}{
						"deploymentstrategy": "",
					},
				},
			},
			want: map[string]interface{}{
				"project": "my-project",
				"stage":   "my-stage",
				"service": "my-service",
				"mytask": map[string]interface{}{
					"deploymentstrategy": "direct",
				},
			},
		},
		{
			name: "get initial triggered event - with input data",
			fields: fields{
				Sequence: keptnv2.Sequence{
					Name: "delivery",
					Tasks: []keptnv2.Task{
						{
							Name: "mytask",
							Properties: map[string]interface{}{
								"foo": "bar",
							},
						},
						{
							Name: "my-second-task",
							Properties: map[string]interface{}{
								"foo": "bar",
							},
						},
					},
				},
				Status: SequenceExecutionStatus{
					PreviousTasks: nil,
					CurrentTask:   TaskExecutionState{},
				},
				Scope: EventScope{
					EventData: keptnv2.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
					},
				},
				InputProperties: map[string]interface{}{
					"configurationChange": map[string]interface{}{
						"image": "1.0",
					},
				},
			},
			want: map[string]interface{}{
				"project": "my-project",
				"stage":   "my-stage",
				"service": "my-service",
				"configurationChange": map[string]interface{}{
					"image": "1.0",
				},
				"mytask": map[string]interface{}{
					"foo": "bar",
				},
			},
		},
		{
			name: "get next triggered event - with input data and completed tasks",
			fields: fields{
				Sequence: keptnv2.Sequence{
					Name: "delivery",
					Tasks: []keptnv2.Task{
						{
							Name: "mytask",
							Properties: map[string]interface{}{
								"foo": "bar",
							},
						},
						{
							Name: "my-second-task",
							Properties: map[string]interface{}{
								"foo": "bar",
							},
						},
					},
				},
				Status: SequenceExecutionStatus{
					PreviousTasks: []TaskExecutionResult{
						{
							Name:   "mytask",
							Result: keptnv2.ResultPass,
							Status: keptnv2.StatusSucceeded,
							Properties: map[string]interface{}{
								"mytask": map[string]interface{}{
									"bar": "foo",
								},
							},
						},
					},
					CurrentTask: TaskExecutionState{},
				},
				Scope: EventScope{
					EventData: keptnv2.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
					},
				},
				InputProperties: map[string]interface{}{
					"configurationChange": map[string]interface{}{
						"image": "1.0",
					},
				},
			},
			want: map[string]interface{}{
				"project": "my-project",
				"stage":   "my-stage",
				"service": "my-service",
				"configurationChange": map[string]interface{}{
					"image": "1.0",
				},
				"mytask": map[string]interface{}{
					"bar": "foo",
				},
				"my-second-task": map[string]interface{}{
					"foo": "bar",
				},
				"result": keptnv2.ResultPass,
				"status": keptnv2.StatusSucceeded,
			},
		},
		{
			name: "get next triggered event - with input data and completed tasks with same properties",
			fields: fields{
				Sequence: keptnv2.Sequence{
					Name: "delivery",
					Tasks: []keptnv2.Task{
						{
							Name: "deployment",
							Properties: map[string]interface{}{
								"foo": "bar",
							},
						},
						{
							Name: "my-second-task",
						},
					},
				},
				Status: SequenceExecutionStatus{
					PreviousTasks: []TaskExecutionResult{
						{
							Name:   "deployment",
							Result: keptnv2.ResultPass,
							Status: keptnv2.StatusSucceeded,
							Properties: map[string]interface{}{
								"deployment": map[string]interface{}{
									"deploymentURIsLocal": []interface{}{
										"http://carts.sockshop-dev:80",
									},
									"deploymentURIsPublic": []interface{}{
										"http://carts.sockshop-staging.svc.cluster.local:80",
									},
								},
								"message": "task finished",
							},
						},
						{
							Name:   "test",
							Result: keptnv2.ResultPass,
							Status: keptnv2.StatusSucceeded,
							Properties: map[string]interface{}{
								"test": map[string]interface{}{
									"start": "3",
									"end":   "4",
								},
								"message": "task finished",
							},
						},
					},
					CurrentTask: TaskExecutionState{},
				},
				Scope: EventScope{
					EventData: keptnv2.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
					},
				},
				InputProperties: map[string]interface{}{
					"configurationChange": map[string]interface{}{
						"image": "1.0",
					},
					"deployment": map[string]interface{}{
						"deploymentURIsLocal": nil,
						"deploymentURIsPublic": []interface{}{
							"http://carts.sockshop-dev.svc.cluster.local:80",
						},
					},
					"test": map[string]interface{}{
						"start": "1",
						"end":   "2",
					},
				},
			},
			want: map[string]interface{}{
				"project": "my-project",
				"stage":   "my-stage",
				"service": "my-service",
				"configurationChange": map[string]interface{}{
					"image": "1.0",
				},
				"deployment": map[string]interface{}{
					"deploymentURIsLocal": []interface{}{
						"http://carts.sockshop-dev:80",
					},
					"deploymentURIsPublic": []interface{}{
						"http://carts.sockshop-staging.svc.cluster.local:80",
						"http://carts.sockshop-dev.svc.cluster.local:80",
					},
				},
				"test": map[string]interface{}{
					"start": "3",
					"end":   "4",
				},
				"result":  keptnv2.ResultPass,
				"status":  keptnv2.StatusSucceeded,
				"message": "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &SequenceExecution{
				ID:              tt.fields.ID,
				Sequence:        tt.fields.Sequence,
				Status:          tt.fields.Status,
				Scope:           tt.fields.Scope,
				InputProperties: tt.fields.InputProperties,
			}
			got := e.GetNextTriggeredEventData()

			require.Equal(t, tt.want, got)
		})
	}
}

func TestSequenceExecution_CompleteCurrentTask(t *testing.T) {
	type fields struct {
		Status SequenceExecutionStatus
	}
	tests := []struct {
		name              string
		fields            fields
		wantResult        keptnv2.ResultType
		wantStatus        keptnv2.StatusType
		wantPreviousTasks []TaskExecutionResult
	}{
		{
			name: "successful task with one executor",
			fields: fields{
				Status: SequenceExecutionStatus{
					CurrentTask: TaskExecutionState{
						Name:        "deployment",
						TriggeredID: "my-triggered-id",
						Events: []TaskEvent{
							{
								EventType: "deployment.started",
								Source:    "my-service",
								Result:    "",
								Status:    "",
							},
							{
								EventType: "deployment.finished",
								Source:    "my-service",
								Result:    keptnv2.ResultPass,
								Status:    keptnv2.StatusSucceeded,
								Properties: map[string]interface{}{
									"deploymentURI": "my-deployment-uri",
								},
							},
						},
					},
				},
			},
			wantResult: keptnv2.ResultPass,
			wantStatus: keptnv2.StatusSucceeded,
			wantPreviousTasks: []TaskExecutionResult{
				{
					Name:        "deployment",
					TriggeredID: "my-triggered-id",
					Result:      keptnv2.ResultPass,
					Status:      keptnv2.StatusSucceeded,
					Properties: map[string]interface{}{
						"deploymentURI": "my-deployment-uri",
					},
				},
			},
		},
		{
			name: "successful task with multiple executors",
			fields: fields{
				Status: SequenceExecutionStatus{
					CurrentTask: TaskExecutionState{
						Name:        "deployment",
						TriggeredID: "my-triggered-id",
						Events: []TaskEvent{
							{
								EventType: "deployment.started",
								Source:    "my-service",
								Result:    "",
								Status:    "",
							},
							{
								EventType: "deployment.finished",
								Source:    "my-service",
								Result:    keptnv2.ResultPass,
								Status:    keptnv2.StatusSucceeded,
								Properties: map[string]interface{}{
									"deploymentURI": "my-deployment-uri",
								},
							},
							{
								EventType: "deployment.started",
								Source:    "my-second-service",
								Result:    "",
								Status:    "",
							},
							{
								EventType: "deployment.finished",
								Source:    "my-second-service",
								Result:    keptnv2.ResultPass,
								Status:    keptnv2.StatusSucceeded,
								Properties: map[string]interface{}{
									"otherProperty": "otherValue",
								},
							},
						},
					},
				},
			},
			wantResult: keptnv2.ResultPass,
			wantStatus: keptnv2.StatusSucceeded,
			wantPreviousTasks: []TaskExecutionResult{
				{
					Name:        "deployment",
					TriggeredID: "my-triggered-id",
					Result:      keptnv2.ResultPass,
					Status:      keptnv2.StatusSucceeded,
					Properties: map[string]interface{}{
						"deploymentURI": "my-deployment-uri",
						"otherProperty": "otherValue",
					},
				},
			},
		},
		{
			name: "multiple executors - one of them failed",
			fields: fields{
				Status: SequenceExecutionStatus{
					CurrentTask: TaskExecutionState{
						Name:        "deployment",
						TriggeredID: "my-triggered-id",
						Events: []TaskEvent{
							{
								EventType: "deployment.started",
								Source:    "my-service",
								Result:    "",
								Status:    "",
							},
							{
								EventType: "deployment.finished",
								Source:    "my-service",
								Result:    keptnv2.ResultPass,
								Status:    keptnv2.StatusSucceeded,
								Properties: map[string]interface{}{
									"deploymentURI": "my-deployment-uri",
								},
							},
							{
								EventType: "deployment.started",
								Source:    "my-second-service",
								Result:    "",
								Status:    "",
							},
							{
								EventType: "deployment.finished",
								Source:    "my-second-service",
								Result:    keptnv2.ResultFailed,
								Status:    keptnv2.StatusSucceeded,
								Properties: map[string]interface{}{
									"otherProperty": "otherValue",
								},
							},
						},
					},
				},
			},
			wantResult: keptnv2.ResultFailed,
			wantStatus: keptnv2.StatusSucceeded,
			wantPreviousTasks: []TaskExecutionResult{
				{
					Name:        "deployment",
					TriggeredID: "my-triggered-id",
					Result:      keptnv2.ResultFailed,
					Status:      keptnv2.StatusSucceeded,
					Properties: map[string]interface{}{
						"deploymentURI": "my-deployment-uri",
						"otherProperty": "otherValue",
					},
				},
			},
		},
		{
			name: "multiple executors - one of them has 'warning' result",
			fields: fields{
				Status: SequenceExecutionStatus{
					CurrentTask: TaskExecutionState{
						Name:        "deployment",
						TriggeredID: "my-triggered-id",
						Events: []TaskEvent{
							{
								EventType: "deployment.started",
								Source:    "my-service",
								Result:    "",
								Status:    "",
							},
							{
								EventType: "deployment.finished",
								Source:    "my-service",
								Result:    keptnv2.ResultPass,
								Status:    keptnv2.StatusSucceeded,
								Properties: map[string]interface{}{
									"deploymentURI": "my-deployment-uri",
								},
							},
							{
								EventType: "deployment.started",
								Source:    "my-second-service",
								Result:    "",
								Status:    "",
							},
							{
								EventType: "deployment.finished",
								Source:    "my-second-service",
								Result:    keptnv2.ResultWarning,
								Status:    keptnv2.StatusSucceeded,
								Properties: map[string]interface{}{
									"otherProperty": "otherValue",
								},
							},
						},
					},
				},
			},
			wantResult: keptnv2.ResultWarning,
			wantStatus: keptnv2.StatusSucceeded,
			wantPreviousTasks: []TaskExecutionResult{
				{
					Name:        "deployment",
					TriggeredID: "my-triggered-id",
					Result:      keptnv2.ResultWarning,
					Status:      keptnv2.StatusSucceeded,
					Properties: map[string]interface{}{
						"deploymentURI": "my-deployment-uri",
						"otherProperty": "otherValue",
					},
				},
			},
		},
		{
			name: "multiple executors - one of them errored",
			fields: fields{
				Status: SequenceExecutionStatus{
					CurrentTask: TaskExecutionState{
						Name:        "deployment",
						TriggeredID: "my-triggered-id",
						Events: []TaskEvent{
							{
								EventType: "deployment.started",
								Source:    "my-service",
								Result:    "",
								Status:    "",
							},
							{
								EventType: "deployment.finished",
								Source:    "my-service",
								Result:    keptnv2.ResultPass,
								Status:    keptnv2.StatusSucceeded,
								Properties: map[string]interface{}{
									"deploymentURI": "my-deployment-uri",
								},
							},
							{
								EventType: "deployment.started",
								Source:    "my-second-service",
								Result:    "",
								Status:    "",
							},
							{
								EventType: "deployment.finished",
								Source:    "my-second-service",
								Result:    keptnv2.ResultFailed,
								Status:    keptnv2.StatusErrored,
								Properties: map[string]interface{}{
									"otherProperty": "otherValue",
								},
							},
						},
					},
				},
			},
			wantResult: keptnv2.ResultFailed,
			wantStatus: keptnv2.StatusErrored,
			wantPreviousTasks: []TaskExecutionResult{
				{
					Name:        "deployment",
					TriggeredID: "my-triggered-id",
					Result:      keptnv2.ResultFailed,
					Status:      keptnv2.StatusErrored,
					Properties: map[string]interface{}{
						"deploymentURI": "my-deployment-uri",
						"otherProperty": "otherValue",
					},
				},
			},
		},
		{
			name: "multiple executors - all properties nil",
			fields: fields{
				Status: SequenceExecutionStatus{
					CurrentTask: TaskExecutionState{
						Name:        "deployment",
						TriggeredID: "my-triggered-id",
						Events: []TaskEvent{
							{
								EventType: "deployment.started",
								Source:    "my-service",
								Result:    "",
								Status:    "",
							},
							{
								EventType:  "deployment.finished",
								Source:     "my-service",
								Result:     keptnv2.ResultPass,
								Status:     keptnv2.StatusSucceeded,
								Properties: nil,
							},
							{
								EventType: "deployment.started",
								Source:    "my-second-service",
								Result:    "",
								Status:    "",
							},
							{
								EventType:  "deployment.finished",
								Source:     "my-second-service",
								Result:     keptnv2.ResultFailed,
								Status:     keptnv2.StatusErrored,
								Properties: nil,
							},
						},
					},
				},
			},
			wantResult: keptnv2.ResultFailed,
			wantStatus: keptnv2.StatusErrored,
			wantPreviousTasks: []TaskExecutionResult{
				{
					Name:        "deployment",
					TriggeredID: "my-triggered-id",
					Result:      keptnv2.ResultFailed,
					Status:      keptnv2.StatusErrored,
					Properties:  nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &SequenceExecution{
				Status: tt.fields.Status,
			}
			result, status := e.CompleteCurrentTask()

			require.Equal(t, tt.wantResult, result)
			require.Equal(t, tt.wantStatus, status)

			require.Equal(t, tt.wantPreviousTasks, e.Status.PreviousTasks)
		})
	}
}

func TestSequenceExecution_GetNextTaskOfSequence(t *testing.T) {
	type fields struct {
		ID              string
		Sequence        keptnv2.Sequence
		Status          SequenceExecutionStatus
		Scope           EventScope
		InputProperties map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   *keptnv2.Task
	}{
		{
			name: "failed previous task - should return nil",
			fields: fields{
				Status: SequenceExecutionStatus{
					PreviousTasks: []TaskExecutionResult{
						{
							Name:        "deployment",
							TriggeredID: "my-triggered-id",
							Result:      keptnv2.ResultFailed,
							Status:      keptnv2.StatusSucceeded,
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "errored previous task - should return nil",
			fields: fields{
				Status: SequenceExecutionStatus{
					PreviousTasks: []TaskExecutionResult{
						{
							Name:        "deployment",
							TriggeredID: "my-triggered-id",
							Status:      keptnv2.StatusErrored,
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "previous task succeeded - get next task",
			fields: fields{
				Status: SequenceExecutionStatus{
					PreviousTasks: []TaskExecutionResult{
						{
							Name:        "deployment",
							TriggeredID: "my-triggered-id",
							Result:      keptnv2.ResultPass,
							Status:      keptnv2.StatusSucceeded,
						},
					},
				},
				Sequence: keptnv2.Sequence{
					Tasks: []keptnv2.Task{
						{
							Name: "deployment",
						},
						{
							Name: "evaluation",
						},
					},
				},
			},
			want: &keptnv2.Task{
				Name: "evaluation",
			},
		},
		{
			name: "no previous task - get first task",
			fields: fields{
				Status: SequenceExecutionStatus{},
				Sequence: keptnv2.Sequence{
					Tasks: []keptnv2.Task{
						{
							Name: "deployment",
						},
						{
							Name: "evaluation",
						},
					},
				},
			},
			want: &keptnv2.Task{
				Name: "deployment",
			},
		},
		{
			name: "all tasks finished - return nil",
			fields: fields{
				Status: SequenceExecutionStatus{
					PreviousTasks: []TaskExecutionResult{
						{
							Name:        "deployment",
							TriggeredID: "my-triggered-id",
							Result:      keptnv2.ResultPass,
							Status:      keptnv2.StatusSucceeded,
						},
						{
							Name:        "evaluation",
							TriggeredID: "my-triggered-id-2",
							Result:      keptnv2.ResultPass,
							Status:      keptnv2.StatusSucceeded,
						},
					},
				},
				Sequence: keptnv2.Sequence{
					Tasks: []keptnv2.Task{
						{
							Name: "deployment",
						},
						{
							Name: "evaluation",
						},
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &SequenceExecution{
				ID:              tt.fields.ID,
				Sequence:        tt.fields.Sequence,
				Status:          tt.fields.Status,
				Scope:           tt.fields.Scope,
				InputProperties: tt.fields.InputProperties,
			}
			if got := e.GetNextTaskOfSequence(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNextTaskOfSequence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSequenceExecution_IsPaused(t *testing.T) {
	type fields struct {
		ID              string
		Sequence        keptnv2.Sequence
		Status          SequenceExecutionStatus
		Scope           EventScope
		InputProperties map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "is paused",
			fields: fields{
				Status: SequenceExecutionStatus{
					State: models.SequencePaused,
				},
			},
			want: true,
		},
		{
			name: "is not paused",
			fields: fields{
				Status: SequenceExecutionStatus{
					State: models.SequenceStartedState,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &SequenceExecution{
				ID:              tt.fields.ID,
				Sequence:        tt.fields.Sequence,
				Status:          tt.fields.Status,
				Scope:           tt.fields.Scope,
				InputProperties: tt.fields.InputProperties,
			}
			if got := e.IsPaused(); got != tt.want {
				t.Errorf("IsPaused() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSequenceExecution_CanBePaused(t *testing.T) {
	type fields struct {
		ID              string
		Sequence        keptnv2.Sequence
		Status          SequenceExecutionStatus
		Scope           EventScope
		InputProperties map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "can be paused",
			fields: fields{
				Status: SequenceExecutionStatus{
					State: models.SequenceStartedState,
				},
			},
			want: true,
		},
		{
			name: "can not be paused",
			fields: fields{
				Status: SequenceExecutionStatus{
					State: models.SequenceFinished,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &SequenceExecution{
				ID:              tt.fields.ID,
				Sequence:        tt.fields.Sequence,
				Status:          tt.fields.Status,
				Scope:           tt.fields.Scope,
				InputProperties: tt.fields.InputProperties,
			}
			if got := e.CanBePaused(); got != tt.want {
				t.Errorf("IsPaused() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSequenceExecution_Pause(t *testing.T) {
	type fields struct {
		ID              string
		Sequence        keptnv2.Sequence
		Status          SequenceExecutionStatus
		Scope           EventScope
		InputProperties map[string]interface{}
	}
	tests := []struct {
		name              string
		fields            fields
		want              bool
		wantPreviousState string
	}{
		{
			name: "try to pause finished sequence",
			fields: fields{
				Status: SequenceExecutionStatus{
					State: models.SequenceFinished,
				},
			},
			want: false,
		},
		{
			name: "pause sequence - keep track of previous state",
			fields: fields{
				Status: SequenceExecutionStatus{
					State: models.SequenceStartedState,
				},
			},
			want:              true,
			wantPreviousState: models.SequenceStartedState,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &SequenceExecution{
				ID:              tt.fields.ID,
				Sequence:        tt.fields.Sequence,
				Status:          tt.fields.Status,
				Scope:           tt.fields.Scope,
				InputProperties: tt.fields.InputProperties,
			}
			if got := e.Pause(); got != tt.want {
				t.Errorf("Pause() = %v, want %v", got, tt.want)
			}
			if tt.want {
				require.Equal(t, tt.wantPreviousState, e.Status.StateBeforePause)
			}
		})
	}
}

func TestSequenceExecution_Resume(t *testing.T) {
	type fields struct {
		ID              string
		Sequence        keptnv2.Sequence
		Status          SequenceExecutionStatus
		Scope           EventScope
		InputProperties map[string]interface{}
	}
	tests := []struct {
		name      string
		fields    fields
		want      bool
		wantState string
	}{
		{
			name: "try to resume non-paused sequence",
			fields: fields{
				Status: SequenceExecutionStatus{
					State: models.SequenceStartedState,
				},
			},
			want: false,
		},
		{
			name: "resume sequence - set state back to previous state",
			fields: fields{
				Status: SequenceExecutionStatus{
					State:            models.SequencePaused,
					StateBeforePause: models.SequenceTriggeredState,
				},
			},
			want:      true,
			wantState: models.SequenceTriggeredState,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &SequenceExecution{
				ID:              tt.fields.ID,
				Sequence:        tt.fields.Sequence,
				Status:          tt.fields.Status,
				Scope:           tt.fields.Scope,
				InputProperties: tt.fields.InputProperties,
			}
			if got := e.Resume(); got != tt.want {
				t.Errorf("Pause() = %v, want %v", got, tt.want)
			}
			if tt.want {
				require.Equal(t, tt.wantState, e.Status.State)
			}
		})
	}
}

func TestTaskExecutionState_IsFinished(t *testing.T) {
	type fields struct {
		Name        string
		TriggeredID string
		Events      []TaskEvent
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "no events received yet",
			fields: fields{
				Events: nil,
			},
			want: false,
		},
		{
			name: "only received .started event",
			fields: fields{
				Events: []TaskEvent{
					{
						EventType: keptnv2.GetStartedEventType("task"),
					},
				},
			},
			want: false,
		},
		{
			name: "received two .started events, but only one .finished event",
			fields: fields{
				Events: []TaskEvent{
					{
						EventType: keptnv2.GetStartedEventType("task"),
					},
					{
						EventType: keptnv2.GetStartedEventType("task"),
					},
					{
						EventType: keptnv2.GetFinishedEventType("task"),
					},
				},
			},
			want: false,
		},
		{
			name: "received two .started events, and two .finished events",
			fields: fields{
				Events: []TaskEvent{
					{
						EventType: keptnv2.GetStartedEventType("task"),
					},
					{
						EventType: keptnv2.GetStartedEventType("task"),
					},
					{
						EventType: keptnv2.GetFinishedEventType("task"),
					},
					{
						EventType: keptnv2.GetFinishedEventType("task"),
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &TaskExecutionState{
				Name:        tt.fields.Name,
				TriggeredID: tt.fields.TriggeredID,
				Events:      tt.fields.Events,
			}
			if got := e.IsFinished(); got != tt.want {
				t.Errorf("IsFinished() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskExecutionState_IsPassed(t *testing.T) {
	type fields struct {
		Name        string
		TriggeredID string
		Events      []TaskEvent
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "failed task",
			fields: fields{
				Events: []TaskEvent{
					{
						EventType: keptnv2.GetStartedEventType("task"),
					},
					{
						EventType: keptnv2.GetFinishedEventType("task"),
						Result:    keptnv2.ResultPass,
					},
					{
						EventType: keptnv2.GetStartedEventType("task"),
					},
					{
						EventType: keptnv2.GetFinishedEventType("task"),
						Result:    keptnv2.ResultFailed,
					},
				},
			},
			want: false,
		},
		{
			name: "one task result is set to 'warning'",
			fields: fields{
				Events: []TaskEvent{
					{
						EventType: keptnv2.GetStartedEventType("task"),
					},
					{
						EventType: keptnv2.GetFinishedEventType("task"),
						Result:    keptnv2.ResultPass,
					},
					{
						EventType: keptnv2.GetStartedEventType("task"),
					},
					{
						EventType: keptnv2.GetFinishedEventType("task"),
						Result:    keptnv2.ResultWarning,
					},
				},
			},
			want: false,
		},
		{
			name: "successful task",
			fields: fields{
				Events: []TaskEvent{
					{
						EventType: keptnv2.GetStartedEventType("task"),
					},
					{
						EventType: keptnv2.GetFinishedEventType("task"),
						Result:    keptnv2.ResultPass,
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &TaskExecutionState{
				Name:        tt.fields.Name,
				TriggeredID: tt.fields.TriggeredID,
				Events:      tt.fields.Events,
			}
			if got := e.IsPassed(); got != tt.want {
				t.Errorf("IsPassed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSequenceExecution_SetNextCurrentTask(t *testing.T) {
	type fields struct {
		currentState     string
		stateBeforePause string
	}
	type args struct {
		taskName         string
		triggeredEventID string
	}
	tests := []struct {
		name                 string
		fields               fields
		args                 args
		wantCurrentState     string
		wantStateBeforePause string
	}{
		{
			name: "currently paused, next task = approval",
			fields: fields{
				stateBeforePause: models.SequenceStartedState,
				currentState:     models.SequencePaused,
			},
			args: args{
				taskName:         keptnv2.ApprovalTaskName,
				triggeredEventID: "1",
			},
			wantStateBeforePause: models.SequenceWaitingForApprovalState,
			wantCurrentState:     models.SequencePaused,
		},
		{
			name: "currently paused, next task = anything",
			fields: fields{
				stateBeforePause: models.SequenceStartedState,
				currentState:     models.SequencePaused,
			},
			args: args{
				taskName:         "anything",
				triggeredEventID: "1",
			},
			wantStateBeforePause: models.SequenceStartedState,
			wantCurrentState:     models.SequencePaused,
		},
		{
			name: "currently not paused, next task = approval",
			fields: fields{
				currentState: models.SequenceStartedState,
			},
			args: args{
				taskName:         keptnv2.ApprovalTaskName,
				triggeredEventID: "1",
			},
			wantCurrentState: models.SequenceWaitingForApprovalState,
		},
		{
			name: "currently not paused, next task = anything",
			fields: fields{
				currentState: models.SequenceStartedState,
			},
			args: args{
				taskName:         "anything",
				triggeredEventID: "1",
			},
			wantCurrentState: models.SequenceStartedState,
		},
		{
			name: "currently waiting for approval, next task = anything",
			fields: fields{
				currentState: models.SequenceWaitingForApprovalState,
			},
			args: args{
				taskName:         "anything",
				triggeredEventID: "1",
			},
			wantCurrentState: models.SequenceStartedState,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &SequenceExecution{
				Status: SequenceExecutionStatus{
					State:            tt.fields.currentState,
					StateBeforePause: tt.fields.stateBeforePause,
				},
			}
			e.SetNextCurrentTask(tt.args.taskName, tt.args.triggeredEventID)

			require.Equal(t, tt.wantCurrentState, e.Status.State)
			require.Equal(t, tt.wantStateBeforePause, e.Status.StateBeforePause)
		})
	}
}
