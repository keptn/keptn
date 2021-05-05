package db

import (
	"errors"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetProjectsMaterializedView(t *testing.T) {
	tests := []struct {
		name string
		want *ProjectsMaterializedView
	}{
		{
			name: "get MV instance",
			want: GetProjectsMaterializedView(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetProjectsMaterializedView(); got != tt.want {
				t.Errorf("GetProjectsMaterializedView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_projectsMaterializedView_CreateProject(t *testing.T) {
	type fields struct {
		ProjectRepo ProjectRepo
	}
	type args struct {
		prj *models.ExpandedProject
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create project that did not exist before",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: func(project *models.ExpandedProject) error {
						return nil
					},

					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return nil, errors.New("")
					},
				},
			},
			args: args{
				prj: &models.ExpandedProject{
					ProjectName: "test-project",
				},
			},
			wantErr: false,
		},
		{
			name: "create project that did exist before",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: func(project *models.ExpandedProject) error {
						return nil
					},
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
						}, nil
					},
					UpdateProjectFunc: nil,
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				prj: &models.ExpandedProject{
					ProjectName: "test-project",
				},
			},
			wantErr: false,
		},
		{
			name: "return error if creating project failed",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: func(project *models.ExpandedProject) error {
						return errors.New("")
					},
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return nil, nil
					},
					UpdateProjectFunc: nil,
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				prj: &models.ExpandedProject{
					ProjectName: "test-project",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &ProjectsMaterializedView{
				ProjectRepo: tt.fields.ProjectRepo,
			}
			if err := mv.CreateProject(tt.args.prj); (err != nil) != tt.wantErr {
				t.Errorf("CreateProject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

const testShipyardContent = `apiVersion: "spec.keptn.sh/0.2.0"
kind: "Shipyard"
metadata:
  name: "shipyard-sockshop"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "delivery"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "test"
              properties:
                teststrategy: "functional"
            - name: "evaluation"
            - name: "release"
        - name: "delivery-direct"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "release"

    - name: "staging"
      sequences:
        - name: "delivery"
          triggeredOn:
            - event: "dev.delivery.finished"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "blue_green_service"
            - name: "test"
              properties:
                teststrategy: "performance"
            - name: "evaluation"
            - name: "release"
        - name: "rollback"
          triggeredOn:
            - event: "staging.delivery.finished"
              selector:
                match:
                  result: "fail"
          tasks:
            - name: "rollback"`

func Test_projectsMaterializedView_UpdateShipyard(t *testing.T) {
	type fields struct {
		ProjectRepo ProjectRepo
	}
	type args struct {
		projectName     string
		shipyardContent string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		expectProject *models.ExpandedProject
		wantErr       bool
	}{
		{
			name: "Update shipyard",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
							Stages: []*models.ExpandedStage{
								{
									StageName: "dev",
								},
								{
									StageName: "staging",
								},
							}}, nil
					},
					UpdateProjectFunc: func(project *models.ExpandedProject) error {
						return nil
					},
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				projectName:     "test-project",
				shipyardContent: testShipyardContent,
			},
			expectProject: &models.ExpandedProject{
				ProjectName:     "test-project",
				Shipyard:        testShipyardContent,
				ShipyardVersion: "spec.keptn.sh/0.2.0",
				Stages: []*models.ExpandedStage{
					{
						StageName:    "dev",
						ParentStages: []string{},
					},
					{
						StageName:    "staging",
						ParentStages: []string{"dev"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "project does not exist",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return nil, errors.New("")
					},
					UpdateProjectFunc: nil,
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				projectName:     "test-project",
				shipyardContent: testShipyardContent,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &ProjectsMaterializedView{
				ProjectRepo: tt.fields.ProjectRepo,
			}
			if err := mv.UpdateShipyard(tt.args.projectName, tt.args.shipyardContent); (err != nil) != tt.wantErr {
				t.Errorf("UpdateShipyard() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockRepo := mv.ProjectRepo.(*db_mock.ProjectRepoMock)

			if tt.expectProject != nil {
				require.Equal(t, 1, len(mockRepo.UpdateProjectCalls()))
				call := mockRepo.UpdateProjectCalls()[0]

				require.Equal(t, tt.expectProject.ShipyardVersion, call.Project.ShipyardVersion)
				require.Equal(t, tt.expectProject.Shipyard, call.Project.Shipyard)
				require.Equal(t, tt.expectProject.Stages, call.Project.Stages)
				mockRepo.UpdateProjectCalls()
			} else {
				require.Empty(t, mockRepo.UpdateProjectCalls())
			}
		})
	}
}

func Test_projectsMaterializedView_CreateStage(t *testing.T) {
	type fields struct {
		ProjectRepo ProjectRepo
	}
	type args struct {
		project string
		stage   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create stage that did not exist before",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{ProjectName: "test-project"}, nil
					},
					UpdateProjectFunc: func(project *models.ExpandedProject) error {
						if len(project.Stages) == 0 {
							return errors.New("unexpected length of stages array")
						}
						if project.Stages[0].StageName != "dev" {
							return errors.New("stage was not named properly")
						}
						return nil
					},
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				project: "test-project",
				stage:   "dev",
			},
			wantErr: false,
		},
		{
			name: "Create stage that did exist before",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
							Stages: []*models.ExpandedStage{
								{
									Services:  nil,
									StageName: "dev",
								},
							},
						}, nil
					},
					UpdateProjectFunc: func(project *models.ExpandedProject) error {
						// should not be called in this case
						return errors.New("update func should not be called in this case")
					},
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				project: "test-project",
				stage:   "dev",
			},
			wantErr: false,
		},
		{
			name: "Create stage that did not exist before",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
							Stages: []*models.ExpandedStage{
								{
									Services:  nil,
									StageName: "dev",
								},
							},
						}, nil
					},
					UpdateProjectFunc: func(project *models.ExpandedProject) error {
						if len(project.Stages) != 2 {
							return errors.New("unexpected length of stages array")
						}
						if project.Stages[0].StageName != "dev" {
							return errors.New("stages was not named properly")
						}
						if project.Stages[1].StageName != "staging" {
							return errors.New("stage was not named properly")
						}
						return nil
					},
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				project: "test-project",
				stage:   "staging",
			},
			wantErr: false,
		},
		{
			name: "Create stage to non-existing project",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return nil, errors.New("")
					},
					UpdateProjectFunc: nil,
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				project: "test-project",
				stage:   "staging",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &ProjectsMaterializedView{
				ProjectRepo: tt.fields.ProjectRepo,
			}
			if err := mv.CreateStage(tt.args.project, tt.args.stage); (err != nil) != tt.wantErr {
				t.Errorf("CreateStage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_projectsMaterializedView_DeleteStage(t *testing.T) {
	type fields struct {
		ProjectRepo ProjectRepo
	}
	type args struct {
		project string
		stage   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Delete stage that did exist before",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
							Stages: []*models.ExpandedStage{
								{
									Services:  nil,
									StageName: "dev",
								},
							},
						}, nil
					},
					UpdateProjectFunc: func(project *models.ExpandedProject) error {
						if len(project.Stages) != 0 {
							return errors.New("unexpected length of stages array")
						}
						return nil
					},
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				project: "test-project",
				stage:   "dev",
			},
			wantErr: false,
		},
		{
			name: "Delete stage that did not exist before",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
						}, nil
					},
					UpdateProjectFunc: func(project *models.ExpandedProject) error {
						// should not be called in this case
						return errors.New("update func should not be called in this case")
					},
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				project: "test-project",
				stage:   "dev",
			},
			wantErr: false,
		},
		{
			name: "Delete stage from non-existing project",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return nil, errors.New("")
					},
					UpdateProjectFunc: nil,
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				project: "test-project",
				stage:   "staging",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &ProjectsMaterializedView{
				ProjectRepo: tt.fields.ProjectRepo,
			}
			if err := mv.DeleteStage(tt.args.project, tt.args.stage); (err != nil) != tt.wantErr {
				t.Errorf("DeleteStage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_projectsMaterializedView_CreateService(t *testing.T) {
	type fields struct {
		ProjectRepo ProjectRepo
	}
	type args struct {
		project string
		stage   string
		service string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create service that did not exist before",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
							Stages: []*models.ExpandedStage{
								{
									Services:  nil,
									StageName: "dev",
								},
							},
						}, nil
					},
					UpdateProjectFunc: func(project *models.ExpandedProject) error {
						if len(project.Stages) != 1 {
							return errors.New("unexpected length of stages array")
						}
						if project.Stages[0].StageName != "dev" {
							return errors.New("stage was not named properly")
						}
						if len(project.Stages[0].Services) != 1 {
							return errors.New("unexpected length of services array")
						}
						if project.Stages[0].Services[0].ServiceName != "test-service" {
							return errors.New("service was not named properly")
						}
						return nil
					},
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				project: "test-project",
				stage:   "dev",
				service: "test-service",
			},
			wantErr: false,
		},
		{
			name: "Create service that did exist before",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
							Stages: []*models.ExpandedStage{
								{
									Services: []*models.ExpandedService{
										{
											ServiceName: "test-service",
										},
									},
									StageName: "dev",
								},
							},
						}, nil
					},
					UpdateProjectFunc: func(project *models.ExpandedProject) error {
						// should not be called in this case
						return errors.New("")
					},
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				project: "test-project",
				stage:   "dev",
				service: "test-service",
			},
			wantErr: false,
		},
		{
			name: "Create service that did not exist before",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
							Stages: []*models.ExpandedStage{
								{
									Services: []*models.ExpandedService{
										{
											ServiceName: "test-service",
										},
									},
									StageName: "dev",
								},
							},
						}, nil
					},
					UpdateProjectFunc: func(project *models.ExpandedProject) error {
						if len(project.Stages) != 1 {
							return errors.New("unexpected length of stages array")
						}
						if project.Stages[0].StageName != "dev" {
							return errors.New("stage was not named properly")
						}
						if len(project.Stages[0].Services) != 2 {
							return errors.New("unexpected length of services array")
						}
						if project.Stages[0].Services[0].ServiceName != "test-service" {
							return errors.New("service was not named properly")
						}
						if project.Stages[0].Services[1].ServiceName != "test-service2" {
							return errors.New("service was not named properly")
						}

						return nil
					},
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				project: "test-project",
				stage:   "dev",
				service: "test-service2",
			},
			wantErr: false,
		},
		{
			name: "Create service to non-existing project",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return nil, errors.New("")
					},
					UpdateProjectFunc: nil,
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				project: "test-project",
				stage:   "dev",
				service: "test-service",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &ProjectsMaterializedView{
				ProjectRepo: tt.fields.ProjectRepo,
			}
			if err := mv.CreateService(tt.args.project, tt.args.stage, tt.args.service); (err != nil) != tt.wantErr {
				t.Errorf("CreateService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_projectsMaterializedView_DeleteService(t *testing.T) {
	type fields struct {
		ProjectRepo ProjectRepo
	}
	type args struct {
		project string
		stage   string
		service string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Delete service that did not exist before",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
							Stages: []*models.ExpandedStage{
								{
									Services:  nil,
									StageName: "dev",
								},
							},
						}, nil
					},
					UpdateProjectFunc: func(project *models.ExpandedProject) error {
						return errors.New("should not be called in this case")
					},
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				project: "test-project",
				stage:   "dev",
				service: "test-service",
			},
			wantErr: false,
		},
		{
			name: "Delete service that did exist before",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
							Stages: []*models.ExpandedStage{
								{
									Services: []*models.ExpandedService{
										{
											ServiceName: "test-service",
										},
									},
									StageName: "dev",
								},
							},
						}, nil
					},
					UpdateProjectFunc: func(project *models.ExpandedProject) error {
						if len(project.Stages[0].Services) != 0 {
							return errors.New("service was not removed properly before update")
						}
						return nil
					},
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				project: "test-project",
				stage:   "dev",
				service: "test-service",
			},
			wantErr: false,
		},
		{
			name: "Delete service from non-existing project",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return nil, errors.New("")
					},
					UpdateProjectFunc: nil,
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				project: "test-project",
				stage:   "dev",
				service: "test-service",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &ProjectsMaterializedView{
				ProjectRepo: tt.fields.ProjectRepo,
			}
			if err := mv.DeleteService(tt.args.project, tt.args.stage, tt.args.service); (err != nil) != tt.wantErr {
				t.Errorf("DeleteService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_updateServiceInStage(t *testing.T) {
	type args struct {
		project *models.ExpandedProject
		stage   string
		service string
		fn      serviceUpdateFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Expect function to be called",
			args: args{
				project: &models.ExpandedProject{
					ProjectName: "test-project",
					Stages: []*models.ExpandedStage{
						{
							Services: []*models.ExpandedService{
								{
									ServiceName: "test-service",
								},
							},
							StageName: "dev",
						},
					},
				},
				stage:   "dev",
				service: "test-service",
				fn: func(service *models.ExpandedService) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "Expect function to not be called",
			args: args{
				project: &models.ExpandedProject{
					ProjectName: "test-project",
					Stages: []*models.ExpandedStage{
						{
							Services:  []*models.ExpandedService{},
							StageName: "dev",
						},
					},
				},
				stage:   "dev",
				service: "test-service",
				fn: func(service *models.ExpandedService) error {
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "Expect error if nil project is provided",
			args: args{
				project: nil,
				stage:   "dev",
				service: "test-service",
				fn: func(service *models.ExpandedService) error {
					return nil
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := updateServiceInStage(tt.args.project, tt.args.stage, tt.args.service, tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("updateServiceInStage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_projectsMaterializedView_UpdateEventOfService(t *testing.T) {
	type fields struct {
		ProjectRepo    ProjectRepo
		EventRetriever EventRepo
	}
	type args struct {
		keptnBase    interface{}
		eventType    string
		keptnContext string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "configuration-change",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
							Stages: []*models.ExpandedStage{
								{
									Services: []*models.ExpandedService{
										{
											ServiceName: "test-service",
										},
									},
									StageName: "dev",
								},
							},
						}, nil
					},
					UpdateProjectFunc: func(project *models.ExpandedProject) error {
						if project.Stages[0].Services[0].LastEventTypes["keptn.sh.some.event"].KeptnContext == "test-context" {
							return nil
						}
						return errors.New("project was not updated correctly")
					},
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				keptnBase:    &keptnv2.EventData{Project: "test-project", Stage: "dev", Service: "test-service"},
				eventType:    "keptn.sh.some.event",
				keptnContext: "test-context",
			},
			wantErr: false,
		},
		{
			name: "deployment-finished",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
							Stages: []*models.ExpandedStage{
								{
									Services: []*models.ExpandedService{
										{
											ServiceName: "test-service",
										},
									},
									StageName: "dev",
								},
							},
						}, nil
					},
					UpdateProjectFunc: func(project *models.ExpandedProject) error {
						if project.Stages[0].Services[0].DeployedImage == "the-service-image:latest" {
							return nil
						}
						return errors.New("project was not updated correctly")
					},
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
				EventRetriever: &db_mock.EventRepoMock{
					GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
						//e1 := models.Event{Triggeredid: "a-triggered-id"}
						e2 := models.Event{
							Data: keptnv2.DeploymentTriggeredEventData{
								EventData: keptnv2.EventData{
									Project: "test-project",
									Stage:   "dev",
									Service: "test-service",
								},
								ConfigurationChange: keptnv2.ConfigurationChange{
									Values: map[string]interface{}{"image": "the-service-image:latest"},
								},
								Deployment: keptnv2.DeploymentTriggeredData{},
							},
							ID: "the-triggered-id",
						}
						return []models.Event{e2}, nil

					},
				},
			},
			args: args{
				keptnBase: &keptnv2.EventData{
					Project: "test-project",
					Stage:   "dev",
					Service: "test-service",
				},
				eventType:    keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName),
				keptnContext: "test-context",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &ProjectsMaterializedView{
				ProjectRepo:     tt.fields.ProjectRepo,
				EventsRetriever: tt.fields.EventRetriever,
			}
			if err := mv.UpdateEventOfService(tt.args.keptnBase, tt.args.eventType, tt.args.keptnContext, "test-event-id", "the-triggered-id"); (err != nil) != tt.wantErr {
				t.Errorf("UpdateEventOfService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_projectsMaterializedView_CreateRemediation(t *testing.T) {
	type fields struct {
		ProjectRepo ProjectRepo
	}
	type args struct {
		project     string
		stage       string
		service     string
		remediation *models.Remediation
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
							Stages: []*models.ExpandedStage{
								{
									Services: []*models.ExpandedService{
										{
											ServiceName: "test-service",
										},
									},
									StageName: "dev",
								},
							},
						}, nil
					},
					UpdateProjectFunc: func(project *models.ExpandedProject) error {
						if len(project.Stages[0].Services[0].OpenRemediations) == 0 {
							return errors.New("project was not updated correctly - no approval has been added")
						}

						if project.Stages[0].Services[0].OpenRemediations[0].EventID != "test-event-id" {
							return errors.New("project was not updated correctly: no event id in approval event")
						}

						if project.Stages[0].Services[0].OpenRemediations[0].KeptnContext != "test-context" {
							return errors.New("project was not updated correctly: keptnContext was not set correctly")
						}

						if project.Stages[0].Services[0].OpenRemediations[0].Type != "remediation.status.changed" {
							return errors.New("project was not updated correctly: type was not set correctly")
						}

						if project.Stages[0].Services[0].OpenRemediations[0].Action != "scale" {
							return errors.New("project was not updated correctly: action was not set correctly")
						}

						return nil
					},
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				project: "test-project",
				stage:   "dev",
				service: "test-service",
				remediation: &models.Remediation{
					EventID:      "test-event-id",
					KeptnContext: "test-context",
					Time:         "1",
					Type:         "remediation.status.changed",
					Action:       "scale",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &ProjectsMaterializedView{
				ProjectRepo: tt.fields.ProjectRepo,
			}
			if err := mv.CreateRemediation(tt.args.project, tt.args.stage, tt.args.service, tt.args.remediation); (err != nil) != tt.wantErr {
				t.Errorf("CreateRemediation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_projectsMaterializedView_CloseOpenRemediations(t *testing.T) {
	type fields struct {
		ProjectRepo ProjectRepo
	}
	type args struct {
		project      string
		stage        string
		service      string
		keptnContext string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "close open approval",
			fields: fields{
				ProjectRepo: &db_mock.ProjectRepoMock{
					CreateProjectFunc: nil,
					GetProjectFunc: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
							Stages: []*models.ExpandedStage{
								{
									Services: []*models.ExpandedService{
										{
											ServiceName: "test-service",
											OpenRemediations: []*models.Remediation{
												{
													EventID:      "test-event-id",
													KeptnContext: "test-context",
													Time:         "1",
													Type:         "remediation.triggered",
												},
												{
													EventID:      "test-event-id",
													KeptnContext: "test-context",
													Time:         "1",
													Type:         "remediation.progressed",
												},
											},
										},
									},
									StageName: "dev",
								},
							},
						}, nil
					},
					UpdateProjectFunc: func(project *models.ExpandedProject) error {
						if len(project.Stages[0].Services[0].OpenRemediations) != 0 {
							return errors.New("project was not updated correctly - open approval was not removed")
						}

						return nil
					},
					DeleteProjectFunc: nil,
					GetProjectsFunc:   nil,
				},
			},
			args: args{
				project:      "test-project",
				stage:        "dev",
				service:      "test-service",
				keptnContext: "test-context",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &ProjectsMaterializedView{
				ProjectRepo: tt.fields.ProjectRepo,
			}
			if err := mv.CloseOpenRemediations(tt.args.project, tt.args.stage, tt.args.service, tt.args.keptnContext); (err != nil) != tt.wantErr {
				t.Errorf("CloseOpenRemediations() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
