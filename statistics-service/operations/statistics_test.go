package operations

import (
	"github.com/go-test/deep"
	"testing"
	"time"
)

func TestStatistics_ensureProjectAndServiceExist(t *testing.T) {
	type fields struct {
		From     time.Time
		To       time.Time
		Projects map[string]*Project
	}
	type args struct {
		projectName string
		serviceName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Create project/service in empty Statistics object",
			fields: fields{
				From:     time.Time{},
				To:       time.Time{},
				Projects: nil,
			},
			args: args{
				projectName: "my-project",
				serviceName: "my-service",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Statistics{
				From:     tt.fields.From,
				To:       tt.fields.To,
				Projects: tt.fields.Projects,
			}

			s.ensureProjectAndServiceExist(tt.args.projectName, tt.args.serviceName)

			if s.Projects == nil {
				t.Error("Statistics.ensureProjectAndServiceExist(): Projects map has not been initialized")
				return
			}
			if s.Projects[tt.args.projectName] == nil {
				t.Error("Statistics.ensureProjectAndServiceExist(): Projects has not been initialized")
				return
			}
			if s.Projects[tt.args.projectName].Services == nil {
				t.Error("Statistics.ensureProjectAndServiceExist(): Project Service map has not been initialized")
				return
			}
			if s.Projects[tt.args.projectName].Services[tt.args.serviceName] == nil {
				t.Error("Statistics.ensureProjectAndServiceExist(): Project Service has not been initialized")
				return
			}
			if s.Projects[tt.args.projectName].Services[tt.args.serviceName].Events == nil {
				t.Error("Statistics.ensureProjectAndServiceExist(): Project Service Events map has not been initialized")
				return
			}
		})
	}
}

func TestStatistics_IncreaseExecutedSequencesCount(t *testing.T) {
	type fields struct {
		From     time.Time
		To       time.Time
		Projects map[string]*Project
	}
	type args struct {
		projectName string
		serviceName string
		increment   int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult int
	}{
		{
			name: "increase by one - initially nil",
			fields: fields{
				From:     time.Time{},
				To:       time.Time{},
				Projects: nil,
			},
			args: args{
				projectName: "my-project",
				serviceName: "my-service",
				increment:   1,
			},
			wantResult: 1,
		},
		{
			name: "increase by one - previous value = 1",
			fields: fields{
				From: time.Time{},
				To:   time.Time{},
				Projects: map[string]*Project{
					"my-project": &Project{
						Name: "my-project",
						Services: map[string]*Service{
							"my-service": &Service{
								Name:              "my-service",
								ExecutedSequences: 1,
								Events:            nil,
							},
						},
					},
				},
			},
			args: args{
				projectName: "my-project",
				serviceName: "my-service",
				increment:   1,
			},
			wantResult: 2,
		},
		{
			name: "increase by two - previous value = 1",
			fields: fields{
				From: time.Time{},
				To:   time.Time{},
				Projects: map[string]*Project{
					"my-project": &Project{
						Name: "my-project",
						Services: map[string]*Service{
							"my-service": &Service{
								Name:              "my-service",
								ExecutedSequences: 1,
								Events:            nil,
							},
						},
					},
				},
			},
			args: args{
				projectName: "my-project",
				serviceName: "my-service",
				increment:   2,
			},
			wantResult: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Statistics{
				From:     tt.fields.From,
				To:       tt.fields.To,
				Projects: tt.fields.Projects,
			}

			s.IncreaseExecutedSequencesCount(tt.args.projectName, tt.args.serviceName, tt.args.increment)
			result := s.Projects[tt.args.projectName].Services[tt.args.serviceName].ExecutedSequences
			if result != tt.wantResult {
				t.Errorf("Statistics.IncreaseExecutedSequencesCount(): want %d, got %d", tt.wantResult, result)
			}
		})
	}
}

func TestStatistics_IncreaseEventTypeCount(t *testing.T) {
	type fields struct {
		From     time.Time
		To       time.Time
		Projects map[string]*Project
	}
	type args struct {
		projectName string
		serviceName string
		eventType   string
		increment   int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult int
	}{
		{
			fields: fields{
				From:     time.Time{},
				To:       time.Time{},
				Projects: nil,
			},
			args: args{
				projectName: "my-project",
				serviceName: "my-service",
				eventType:   "my-event",
				increment:   1,
			},
			wantResult: 1,
		},
		{
			name: "increase by one - previous value = 1",
			fields: fields{
				From: time.Time{},
				To:   time.Time{},
				Projects: map[string]*Project{
					"my-project": &Project{
						Name: "my-project",
						Services: map[string]*Service{
							"my-service": &Service{
								Name:              "my-service",
								ExecutedSequences: 1,
								Events: map[string]int{
									"my-event": 1,
								},
							},
						},
					},
				},
			},
			args: args{
				projectName: "my-project",
				serviceName: "my-service",
				eventType:   "my-event",
				increment:   1,
			},
			wantResult: 2,
		},
		{
			name: "increase by two - previous value = 1",
			fields: fields{
				From: time.Time{},
				To:   time.Time{},
				Projects: map[string]*Project{
					"my-project": &Project{
						Name: "my-project",
						Services: map[string]*Service{
							"my-service": &Service{
								Name:              "my-service",
								ExecutedSequences: 1,
								Events: map[string]int{
									"my-event": 1,
								},
							},
						},
					},
				},
			},
			args: args{
				projectName: "my-project",
				serviceName: "my-service",
				eventType:   "my-event",
				increment:   2,
			},
			wantResult: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Statistics{
				From:     tt.fields.From,
				To:       tt.fields.To,
				Projects: tt.fields.Projects,
			}

			s.IncreaseEventTypeCount(tt.args.projectName, tt.args.serviceName, tt.args.eventType, tt.args.increment)
			result := s.Projects[tt.args.projectName].Services[tt.args.serviceName].Events[tt.args.eventType]
			if result != tt.wantResult {
				t.Errorf("Statistics.IncreaseEventTypeCount(): want %d, got %d", tt.wantResult, result)
			}
		})
	}
}

func TestMergeStatistics(t *testing.T) {
	type args struct {
		target     Statistics
		statistics []Statistics
	}
	tests := []struct {
		name string
		args args
		want Statistics
	}{
		{
			name: "mege statistics",
			args: args{
				target: Statistics{
					From:     time.Time{},
					To:       time.Time{},
					Projects: nil,
				},
				statistics: []Statistics{
					{
						From: time.Time{},
						To:   time.Time{},
						Projects: map[string]*Project{
							"my-project": {
								Name: "my-project",
								Services: map[string]*Service{
									"my-service": {
										Name:              "my-service",
										ExecutedSequences: 2,
										Events: map[string]int{
											"my-type": 1,
										},
										KeptnServiceExecutions: map[string]*KeptnService{
											"my-keptn-service": {
												Name: "my-keptn-service",
												Executions: map[string]int{
													"my-type": 1,
												},
											},
										},
									},
								},
							},
						},
					},
					{
						From: time.Time{},
						To:   time.Time{},
						Projects: map[string]*Project{
							"my-project": {
								Name: "my-project",
								Services: map[string]*Service{
									"my-service": {
										Name:              "my-service",
										ExecutedSequences: 2,
										Events: map[string]int{
											"my-type":   1,
											"my-type-2": 1,
										},
										KeptnServiceExecutions: map[string]*KeptnService{
											"my-keptn-service": {
												Name: "my-keptn-service",
												Executions: map[string]int{
													"my-type": 1,
												},
											},
										},
									},
								},
							},
						},
					},
					{
						From: time.Time{},
						To:   time.Time{},
						Projects: map[string]*Project{
							"my-project-2": {
								Name: "my-project-2",
								Services: map[string]*Service{
									"my-service-2": {
										Name:              "my-service-2",
										ExecutedSequences: 2,
										Events: map[string]int{
											"my-type":   2,
											"my-type-2": 1,
										},
									},
								},
							},
						},
					},
				},
			},
			want: Statistics{
				From: time.Time{},
				To:   time.Time{},
				Projects: map[string]*Project{
					"my-project": {
						Name: "my-project",
						Services: map[string]*Service{
							"my-service": {
								Name:              "my-service",
								ExecutedSequences: 4,
								Events: map[string]int{
									"my-type":   2,
									"my-type-2": 1,
								},
								KeptnServiceExecutions: map[string]*KeptnService{
									"my-keptn-service": {
										Name: "my-keptn-service",
										Executions: map[string]int{
											"my-type": 2,
										},
									},
								},
								ExecutedSequencesPerType: map[string]int{},
							},
						},
					},
					"my-project-2": {
						Name: "my-project-2",
						Services: map[string]*Service{
							"my-service-2": {
								Name:              "my-service-2",
								ExecutedSequences: 2,
								Events: map[string]int{
									"my-type":   2,
									"my-type-2": 1,
								},
								KeptnServiceExecutions:   map[string]*KeptnService{},
								ExecutedSequencesPerType: map[string]int{},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeStatistics(tt.args.target, tt.args.statistics)

			diff := deep.Equal(got, tt.want)
			if len(diff) > 0 {
				t.Error("MergeStatistics(): did not return expected value")
				for _, d := range diff {
					t.Log(d)
				}
			}
		})
	}
}

func TestStatistics_IncreaseKeptnServiceExecutionCount(t *testing.T) {
	type fields struct {
		From     time.Time
		To       time.Time
		Projects map[string]*Project
	}
	type args struct {
		projectName      string
		serviceName      string
		keptnServiceName string
		eventType        string
		increment        int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult int
	}{

		{
			name: "increase by one - previous value = 1",
			fields: fields{
				From: time.Time{},
				To:   time.Time{},
				Projects: map[string]*Project{
					"my-project": &Project{
						Name: "my-project",
						Services: map[string]*Service{
							"my-service": &Service{
								Name:              "my-service",
								ExecutedSequences: 1,
								Events: map[string]int{
									"my-event": 1,
								},
								KeptnServiceExecutions: map[string]*KeptnService{
									"my-keptn-service": {
										Name: "my-keptn-service",
										Executions: map[string]int{
											"my-event": 1,
										},
									},
								},
							},
						},
					},
				},
			},
			args: args{
				projectName:      "my-project",
				serviceName:      "my-service",
				keptnServiceName: "my-keptn-service",
				eventType:        "my-event",
				increment:        1,
			},
			wantResult: 2,
		},
		{
			name: "increase by two - previous value = 1",
			fields: fields{
				From: time.Time{},
				To:   time.Time{},
				Projects: map[string]*Project{
					"my-project": &Project{
						Name: "my-project",
						Services: map[string]*Service{
							"my-service": &Service{
								Name:              "my-service",
								ExecutedSequences: 1,
								Events: map[string]int{
									"my-event": 1,
								},
								KeptnServiceExecutions: map[string]*KeptnService{
									"my-keptn-service": {
										Name: "my-keptn-service",
										Executions: map[string]int{
											"my-event": 1,
										},
									},
								},
							},
						},
					},
				},
			},
			args: args{
				projectName:      "my-project",
				serviceName:      "my-service",
				keptnServiceName: "my-keptn-service",
				eventType:        "my-event",
				increment:        2,
			},
			wantResult: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Statistics{
				From:     tt.fields.From,
				To:       tt.fields.To,
				Projects: tt.fields.Projects,
			}

			s.IncreaseKeptnServiceExecutionCount(tt.args.projectName, tt.args.serviceName, tt.args.keptnServiceName, tt.args.eventType, tt.args.increment)
			result := s.Projects[tt.args.projectName].Services[tt.args.serviceName].KeptnServiceExecutions[tt.args.keptnServiceName].Executions[tt.args.eventType]
			if result != tt.wantResult {
				t.Errorf("Statistics.IncreaseEventTypeCount(): want %d, got %d", tt.wantResult, result)
			}
		})
	}
}
