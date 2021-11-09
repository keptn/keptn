package handler

import (
	"github.com/go-test/deep"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"reflect"
	"testing"
)

func Test_shipyardController_getTaskSequenceInStage(t *testing.T) {
	type fields struct {
		projectRepo      db.ProjectMVRepo
		eventRepo        db.EventRepo
		taskSequenceRepo db.TaskSequenceRepo
	}
	type args struct {
		stageName        string
		taskSequenceName string
		shipyard         *keptnv2.Shipyard
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *keptnv2.Sequence
		wantErr bool
	}{
		{
			name: "get built-in evaluation task sequence",
			fields: fields{
				projectRepo:      nil,
				eventRepo:        nil,
				taskSequenceRepo: nil,
			},
			args: args{
				stageName:        "dev",
				taskSequenceName: "evaluation",
				shipyard: &keptnv2.Shipyard{
					ApiVersion: "0.2.0",
					Kind:       "shipyard",
					Metadata:   keptnv2.Metadata{},
					Spec: keptnv2.ShipyardSpec{
						Stages: []keptnv2.Stage{
							{
								Name:      "dev",
								Sequences: []keptnv2.Sequence{},
							},
						},
					},
				},
			},
			want: &keptnv2.Sequence{
				Name:        "evaluation",
				TriggeredOn: nil,
				Tasks: []keptnv2.Task{
					{
						Name:       "evaluation",
						Properties: nil,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "get user-defined evaluation task sequence",
			fields: fields{
				projectRepo:      nil,
				eventRepo:        nil,
				taskSequenceRepo: nil,
			},
			args: args{
				stageName:        "dev",
				taskSequenceName: "evaluation",
				shipyard: &keptnv2.Shipyard{
					ApiVersion: "0.2.0",
					Kind:       "shipyard",
					Metadata:   keptnv2.Metadata{},
					Spec: keptnv2.ShipyardSpec{
						Stages: []keptnv2.Stage{
							{
								Name: "dev",
								Sequences: []keptnv2.Sequence{
									{
										Name:        "evaluation",
										TriggeredOn: nil,
										Tasks: []keptnv2.Task{
											{
												Name:       "evaluation",
												Properties: nil,
											},
											{
												Name:       "notify",
												Properties: nil,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: &keptnv2.Sequence{
				Name:        "evaluation",
				TriggeredOn: nil,
				Tasks: []keptnv2.Task{
					{
						Name:       "evaluation",
						Properties: nil,
					},
					{
						Name:       "notify",
						Properties: nil,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTaskSequenceInStage(tt.args.stageName, tt.args.taskSequenceName, tt.args.shipyard)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTaskSequenceInStage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("GetTaskSequenceInStage() got = %v, want %v", got, tt.want)
				for _, d := range diff {
					t.Log(d)
				}
			}
		})
	}
}

func Test_GetTaskSequencesByTrigger(t *testing.T) {
	type args struct {
		eventScope            models.EventScope
		completedTaskSequence string
		shipyard              *keptnv2.Shipyard
		previousTask          string
	}
	tests := []struct {
		name string
		args args
		want []NextTaskSequence
	}{
		{
			name: "default behavior - get sequence triggered by result=pass,warning",
			args: args{
				eventScope: models.EventScope{EventData: keptnv2.EventData{
					Result: keptnv2.ResultPass,
					Stage:  "dev",
				}},
				completedTaskSequence: "artifact-delivery",
				shipyard: &keptnv2.Shipyard{
					ApiVersion: shipyardVersion,
					Kind:       "shipyard",
					Metadata:   keptnv2.Metadata{},
					Spec: keptnv2.ShipyardSpec{
						Stages: []keptnv2.Stage{
							{
								Name: "dev",
								Sequences: []keptnv2.Sequence{
									{
										Name:        "artifact-delivery",
										TriggeredOn: nil,
										Tasks:       nil,
									},
								},
							},
							{
								Name: "hardening",
								Sequences: []keptnv2.Sequence{
									{
										Name: "artifact-delivery",
										TriggeredOn: []keptnv2.Trigger{
											{
												Event:    "dev.artifact-delivery.finished",
												Selector: keptnv2.Selector{},
											},
										},
										Tasks: nil,
									},
									{
										Name:        "artifact-delivery-2",
										TriggeredOn: nil,
										Tasks:       nil,
									},
								},
							},
						},
					},
				},
			},
			want: []NextTaskSequence{
				{
					Sequence: keptnv2.Sequence{
						Name: "artifact-delivery",
						TriggeredOn: []keptnv2.Trigger{
							{
								Event:    "dev.artifact-delivery.finished",
								Selector: keptnv2.Selector{},
							},
						},
						Tasks: nil,
					},
					StageName: "hardening",
				},
			},
		},
		{
			name: "get sequence triggered by result=fail",
			args: args{
				eventScope: models.EventScope{EventData: keptnv2.EventData{
					Result: keptnv2.ResultFailed,
					Stage:  "dev",
				}},
				completedTaskSequence: "artifact-delivery",
				shipyard: &keptnv2.Shipyard{
					ApiVersion: shipyardVersion,
					Kind:       "shipyard",
					Metadata:   keptnv2.Metadata{},
					Spec: keptnv2.ShipyardSpec{
						Stages: []keptnv2.Stage{
							{
								Name: "dev",
								Sequences: []keptnv2.Sequence{
									{
										Name:        "artifact-delivery",
										TriggeredOn: nil,
										Tasks:       nil,
									},
								},
							},
							{
								Name: "hardening",
								Sequences: []keptnv2.Sequence{
									{
										Name: "artifact-delivery",
										TriggeredOn: []keptnv2.Trigger{
											{
												Event:    "dev.artifact-delivery.finished",
												Selector: keptnv2.Selector{},
											},
										},
										Tasks: nil,
									},
									{
										Name: "artifact-delivery-2",
										TriggeredOn: []keptnv2.Trigger{
											{
												Event: "dev.artifact-delivery.finished",
												Selector: keptnv2.Selector{
													Match: map[string]string{
														"result": string(keptnv2.ResultFailed),
													},
												},
											},
										},
										Tasks: nil,
									},
								},
							},
							{
								Name: "production",
								Sequences: []keptnv2.Sequence{
									{
										Name: "artifact-delivery",
										TriggeredOn: []keptnv2.Trigger{
											{
												Event:    "dev.artifact-delivery.finished",
												Selector: keptnv2.Selector{},
											},
										},
										Tasks: nil,
									},
									{
										Name: "artifact-delivery-2",
										TriggeredOn: []keptnv2.Trigger{
											{
												Event: "dev.artifact-delivery.finished",
												Selector: keptnv2.Selector{
													Match: map[string]string{
														"result": string(keptnv2.ResultFailed),
													},
												},
											},
										},
										Tasks: nil,
									},
								},
							},
						},
					},
				},
			},
			want: []NextTaskSequence{
				{
					Sequence: keptnv2.Sequence{
						Name: "artifact-delivery-2",
						TriggeredOn: []keptnv2.Trigger{
							{
								Event: "dev.artifact-delivery.finished",
								Selector: keptnv2.Selector{
									Match: map[string]string{
										"result": string(keptnv2.ResultFailed),
									},
								},
							},
						},
						Tasks: nil,
					},
					StageName: "hardening",
				},
				{
					Sequence: keptnv2.Sequence{
						Name: "artifact-delivery-2",
						TriggeredOn: []keptnv2.Trigger{
							{
								Event: "dev.artifact-delivery.finished",
								Selector: keptnv2.Selector{
									Match: map[string]string{
										"result": string(keptnv2.ResultFailed),
									},
								},
							},
						},
						Tasks: nil,
					},
					StageName: "production",
				},
			},
		},
		{
			name: "get sequence triggered by result=fail of specific task",
			args: args{
				eventScope: models.EventScope{EventData: keptnv2.EventData{
					Result: keptnv2.ResultFailed,
					Stage:  "dev",
				}},
				completedTaskSequence: "artifact-delivery",
				previousTask:          "evaluation",
				shipyard: &keptnv2.Shipyard{
					ApiVersion: shipyardVersion,
					Kind:       "shipyard",
					Metadata:   keptnv2.Metadata{},
					Spec: keptnv2.ShipyardSpec{
						Stages: []keptnv2.Stage{
							{
								Name: "dev",
								Sequences: []keptnv2.Sequence{
									{
										Name:        "artifact-delivery",
										TriggeredOn: nil,
										Tasks:       nil,
									},
								},
							},
							{
								Name: "hardening",
								Sequences: []keptnv2.Sequence{
									{
										Name: "artifact-delivery",
										TriggeredOn: []keptnv2.Trigger{
											{
												Event:    "dev.artifact-delivery.finished",
												Selector: keptnv2.Selector{},
											},
										},
										Tasks: nil,
									},
									{
										Name: "artifact-delivery-2",
										TriggeredOn: []keptnv2.Trigger{
											{
												Event: "dev.artifact-delivery.finished",
												Selector: keptnv2.Selector{
													Match: map[string]string{
														"evaluation.result": string(keptnv2.ResultFailed),
													},
												},
											},
										},
										Tasks: nil,
									},
								},
							},
							{
								Name: "production",
								Sequences: []keptnv2.Sequence{
									{
										Name: "artifact-delivery",
										TriggeredOn: []keptnv2.Trigger{
											{
												Event:    "dev.artifact-delivery.finished",
												Selector: keptnv2.Selector{},
											},
										},
										Tasks: nil,
									},
									{
										Name: "artifact-delivery-2",
										TriggeredOn: []keptnv2.Trigger{
											{
												Event: "dev.artifact-delivery.finished",
												Selector: keptnv2.Selector{
													Match: map[string]string{
														"deployment.result": string(keptnv2.ResultFailed),
													},
												},
											},
										},
										Tasks: nil,
									},
								},
							},
						},
					},
				},
			},
			want: []NextTaskSequence{
				{
					Sequence: keptnv2.Sequence{
						Name: "artifact-delivery-2",
						TriggeredOn: []keptnv2.Trigger{
							{
								Event: "dev.artifact-delivery.finished",
								Selector: keptnv2.Selector{
									Match: map[string]string{
										"evaluation.result": string(keptnv2.ResultFailed),
									},
								},
							},
						},
						Tasks: nil,
					},
					StageName: "hardening",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTaskSequencesByTrigger(tt.args.eventScope, tt.args.completedTaskSequence, tt.args.shipyard, tt.args.previousTask); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTaskSequencesByTrigger() = %v, want %v", got, tt.want)
			}
		})
	}
}
