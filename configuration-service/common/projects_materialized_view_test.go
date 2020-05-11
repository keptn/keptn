package common

import (
	"errors"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/configuration-service/models"
	"reflect"
	"testing"
)

type CreateProjectMock func(project *models.ExpandedProject) error
type GetProjectMock func(projectName string) (*models.ExpandedProject, error)
type GetProjectsMock func() ([]*models.ExpandedProject, error)
type UpdateProjectMock func(project *models.ExpandedProject) error
type DeleteProjectMock func(projectName string) error

type mockProjectRepo struct {
	CreateProjectMock CreateProjectMock
	GetProjectMock    GetProjectMock
	UpdateProjectMock UpdateProjectMock
	DeleteProjectMock DeleteProjectMock
	GetProjectsMock   GetProjectsMock
}

func (m mockProjectRepo) CreateProject(project *models.ExpandedProject) error {
	return m.CreateProjectMock(project)
}

func (m mockProjectRepo) GetProject(projectName string) (*models.ExpandedProject, error) {
	return m.GetProjectMock(projectName)
}

func (m mockProjectRepo) GetProjects() ([]*models.ExpandedProject, error) {
	return m.GetProjectsMock()
}

func (m mockProjectRepo) UpdateProject(project *models.ExpandedProject) error {
	return m.UpdateProjectMock(project)
}

func (m mockProjectRepo) DeleteProject(projectName string) error {
	return m.DeleteProjectMock(projectName)
}

func TestGetProjectsMaterializedView(t *testing.T) {
	tests := []struct {
		name string
		want *projectsMaterializedView
	}{
		{
			name: "get MV instance",
			want: &projectsMaterializedView{
				ProjectRepo: &MongoDBProjectRepo{},
				Logger:      keptn.NewLogger("", "", "configuration-service"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetProjectsMaterializedView(); !reflect.DeepEqual(got, tt.want) {
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
		prj *models.Project
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
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: func(project *models.ExpandedProject) error {
						return nil
					},
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
						return nil, errors.New("")
					},
					UpdateProjectMock: nil,
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
				},
			},
			args: args{
				prj: &models.Project{
					ProjectName: "test-project",
				},
			},
			wantErr: false,
		},
		{
			name: "create project that did exist before",
			fields: fields{
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: func(project *models.ExpandedProject) error {
						return nil
					},
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
						}, nil
					},
					UpdateProjectMock: nil,
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
				},
			},
			args: args{
				prj: &models.Project{
					ProjectName: "test-project",
				},
			},
			wantErr: false,
		},
		{
			name: "return error if creating project failed",
			fields: fields{
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: func(project *models.ExpandedProject) error {
						return errors.New("")
					},
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
						return nil, nil
					},
					UpdateProjectMock: nil,
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
				},
			},
			args: args{
				prj: &models.Project{
					ProjectName: "test-project",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &projectsMaterializedView{
				ProjectRepo: tt.fields.ProjectRepo,
				Logger:      keptn.NewLogger("", "", "configuration-service"),
			}
			if err := mv.CreateProject(tt.args.prj); (err != nil) != tt.wantErr {
				t.Errorf("CreateProject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_projectsMaterializedView_UpdateShipyard(t *testing.T) {
	type fields struct {
		ProjectRepo ProjectRepo
	}
	type args struct {
		projectName     string
		shipyardContent string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Update shipyard",
			fields: fields{
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{ProjectName: "test-project"}, nil
					},
					UpdateProjectMock: func(project *models.ExpandedProject) error {
						if project.Shipyard == "test-content" {
							return nil
						}
						return errors.New("shipyard content was not updated properly")
					},
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
				},
			},
			args: args{
				projectName:     "test-project",
				shipyardContent: "test-content",
			},
			wantErr: false,
		},
		{
			name: "project does not exist",
			fields: fields{
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
						return nil, errors.New("")
					},
					UpdateProjectMock: nil,
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
				},
			},
			args: args{
				projectName:     "test-project",
				shipyardContent: "test-content",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &projectsMaterializedView{
				ProjectRepo: tt.fields.ProjectRepo,
			}
			if err := mv.UpdateShipyard(tt.args.projectName, tt.args.shipyardContent); (err != nil) != tt.wantErr {
				t.Errorf("UpdateShipyard() error = %v, wantErr %v", err, tt.wantErr)
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
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{ProjectName: "test-project"}, nil
					},
					UpdateProjectMock: func(project *models.ExpandedProject) error {
						if len(project.Stages) == 0 {
							return errors.New("unexpected length of stages array")
						}
						if project.Stages[0].StageName != "dev" {
							return errors.New("stage was not named properly")
						}
						return nil
					},
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
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
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
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
					UpdateProjectMock: func(project *models.ExpandedProject) error {
						// should not be called in this case
						return errors.New("update func should not be called in this case")
					},
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
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
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
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
					UpdateProjectMock: func(project *models.ExpandedProject) error {
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
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
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
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
						return nil, errors.New("")
					},
					UpdateProjectMock: nil,
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
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
			mv := &projectsMaterializedView{
				ProjectRepo: tt.fields.ProjectRepo,
				Logger:      keptn.NewLogger("", "", "configuration-service"),
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
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
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
					UpdateProjectMock: func(project *models.ExpandedProject) error {
						if len(project.Stages) != 0 {
							return errors.New("unexpected length of stages array")
						}
						return nil
					},
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
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
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
						}, nil
					},
					UpdateProjectMock: func(project *models.ExpandedProject) error {
						// should not be called in this case
						return errors.New("update func should not be called in this case")
					},
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
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
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
						return nil, errors.New("")
					},
					UpdateProjectMock: nil,
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
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
			mv := &projectsMaterializedView{
				ProjectRepo: tt.fields.ProjectRepo,
				Logger:      keptn.NewLogger("", "", "configuration-service"),
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
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
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
					UpdateProjectMock: func(project *models.ExpandedProject) error {
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
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
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
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
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
					UpdateProjectMock: func(project *models.ExpandedProject) error {
						// should not be called in this case
						return errors.New("")
					},
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
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
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
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
					UpdateProjectMock: func(project *models.ExpandedProject) error {
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
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
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
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
						return nil, errors.New("")
					},
					UpdateProjectMock: nil,
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
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
			mv := &projectsMaterializedView{
				ProjectRepo: tt.fields.ProjectRepo,
				Logger:      keptn.NewLogger("", "", "configuration-service"),
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
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
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
					UpdateProjectMock: func(project *models.ExpandedProject) error {
						return errors.New("should not be called in this case")
					},
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
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
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
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
					UpdateProjectMock: func(project *models.ExpandedProject) error {
						if len(project.Stages[0].Services) != 0 {
							return errors.New("service was not removed properly before update")
						}
						return nil
					},
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
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
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
						return nil, errors.New("")
					},
					UpdateProjectMock: nil,
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
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
			mv := &projectsMaterializedView{
				ProjectRepo: tt.fields.ProjectRepo,
				Logger:      keptn.NewLogger("", "", "configuration-service"),
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
		ProjectRepo ProjectRepo
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
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
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
					UpdateProjectMock: func(project *models.ExpandedProject) error {
						if project.Stages[0].Services[0].LastEventTypes[keptn.ConfigurationChangeEventType].KeptnContext == "test-context" {
							return nil
						}
						return errors.New("project was not updated correctly")
					},
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
				},
			},
			args: args{
				keptnBase:    &keptn.KeptnBase{Project: "test-project", Stage: "dev", Service: "test-service"},
				eventType:    keptn.ConfigurationChangeEventType,
				keptnContext: "test-context",
			},
			wantErr: false,
		},
		{
			name: "deployment-finished",
			fields: fields{
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
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
					UpdateProjectMock: func(project *models.ExpandedProject) error {
						if project.Stages[0].Services[0].LastEventTypes[keptn.DeploymentFinishedEventType].KeptnContext == "test-context" &&
							project.Stages[0].Services[0].DeployedImage == "test-image:0.1" {
							return nil
						}
						return errors.New("project was not updated correctly")
					},
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
				},
			},
			args: args{
				keptnBase: &keptn.KeptnBase{
					Project: "test-project",
					Stage:   "dev",
					Service: "test-service",
					Image:   stringp("test-image"),
					Tag:     stringp("0.1"),
				},
				eventType:    keptn.DeploymentFinishedEventType,
				keptnContext: "test-context",
			},
			wantErr: false,
		},
		{
			name: "approval.triggered",
			fields: fields{
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
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
					UpdateProjectMock: func(project *models.ExpandedProject) error {
						if len(project.Stages[0].Services[0].OpenApprovals) == 0 {
							return errors.New("project was not updated correctly - no approval finished event has been added")
						}

						if project.Stages[0].Services[0].OpenApprovals[0].EventID != "test-event-id" {
							return errors.New("project was not updated correctly: no event id in approval event")
						}

						return nil
					},
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
				},
			},
			args: args{
				keptnBase: &keptn.ApprovalTriggeredEventData{
					Project: "test-project",
					Stage:   "dev",
					Service: "test-service",
				},
				eventType:    keptn.ApprovalTriggeredEventType,
				keptnContext: "test-context",
			},
			wantErr: false,
		},
		{
			name: "approval.finished",
			fields: fields{
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
							Stages: []*models.ExpandedStage{
								{
									Services: []*models.ExpandedService{
										{
											ServiceName: "test-service",
											OpenApprovals: []*models.Approval{
												{
													EventID:      "test-event-id",
													KeptnContext: "test-context",
													Time:         "",
													Image:        "my-image",
													Tag:          "my-tag",
												},
											},
										},
									},
									StageName: "dev",
								},
							},
						}, nil
					},
					UpdateProjectMock: func(project *models.ExpandedProject) error {
						if len(project.Stages[0].Services[0].OpenApprovals) != 0 {
							return errors.New("project was not updated correctly - open approval was not removed")
						}

						return nil
					},
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
				},
			},
			args: args{
				keptnBase: &keptn.ApprovalFinishedEventData{
					Project: "test-project",
					Stage:   "dev",
					Service: "test-service",
					Approval: keptn.ApprovalData{
						TriggeredID: "test-event-id",
						Result:      "Pass",
						Status:      "Success",
					},
					Image: "my-image",
					Tag:   "my-tag",
				},
				eventType:    keptn.ApprovalFinishedEventType,
				keptnContext: "test-context",
			},
			wantErr: false,
		},
		{
			name: "approval.finished: no matching tags",
			fields: fields{
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
							Stages: []*models.ExpandedStage{
								{
									Services: []*models.ExpandedService{
										{
											ServiceName: "test-service",
											OpenApprovals: []*models.Approval{
												{
													EventID:      "test-event-id",
													KeptnContext: "test-context",
													Time:         "",
													Image:        "my-image",
													Tag:          "my-tag",
												},
											},
										},
									},
									StageName: "dev",
								},
							},
						}, nil
					},
					UpdateProjectMock: func(project *models.ExpandedProject) error {
						if len(project.Stages[0].Services[0].OpenApprovals) != 1 {
							return errors.New("project was not updated correctly - open approval was removed even though tags didn't match")
						}

						return nil
					},
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
				},
			},
			args: args{
				keptnBase: &keptn.ApprovalFinishedEventData{
					Project: "test-project",
					Stage:   "dev",
					Service: "test-service",
					Approval: keptn.ApprovalData{
						TriggeredID: "test-event-id",
						Result:      "Pass",
						Status:      "Success",
					},
					Image: "my-image",
					Tag:   "my-wrong-tag",
				},
				eventType:    keptn.ApprovalFinishedEventType,
				keptnContext: "test-context",
			},
			wantErr: false,
		},
		{
			name: "approval.finished: no matching images",
			fields: fields{
				ProjectRepo: &mockProjectRepo{
					CreateProjectMock: nil,
					GetProjectMock: func(projectName string) (project *models.ExpandedProject, err error) {
						return &models.ExpandedProject{
							ProjectName: "test-project",
							Stages: []*models.ExpandedStage{
								{
									Services: []*models.ExpandedService{
										{
											ServiceName: "test-service",
											OpenApprovals: []*models.Approval{
												{
													EventID:      "test-event-id",
													KeptnContext: "test-context",
													Time:         "",
													Image:        "my-image",
													Tag:          "my-tag",
												},
											},
										},
									},
									StageName: "dev",
								},
							},
						}, nil
					},
					UpdateProjectMock: func(project *models.ExpandedProject) error {
						if len(project.Stages[0].Services[0].OpenApprovals) != 1 {
							return errors.New("project was not updated correctly - open approval was removed even though images didn't match")
						}

						return nil
					},
					DeleteProjectMock: nil,
					GetProjectsMock:   nil,
				},
			},
			args: args{
				keptnBase: &keptn.ApprovalFinishedEventData{
					Project: "test-project",
					Stage:   "dev",
					Service: "test-service",
					Approval: keptn.ApprovalData{
						TriggeredID: "test-event-id",
						Result:      "Pass",
						Status:      "Success",
					},
					Image: "my-wrong-image",
					Tag:   "my-tag",
				},
				eventType:    keptn.ApprovalFinishedEventType,
				keptnContext: "test-context",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := &projectsMaterializedView{
				ProjectRepo: tt.fields.ProjectRepo,
				Logger:      keptn.NewLogger("", "", "configuration-service"),
			}
			if err := mv.UpdateEventOfService(tt.args.keptnBase, tt.args.eventType, tt.args.keptnContext, "test-event-id"); (err != nil) != tt.wantErr {
				t.Errorf("UpdateEventOfService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func stringp(s string) *string {
	return &s
}
