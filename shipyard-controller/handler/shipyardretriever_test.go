package handler

import (
	"errors"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	common_mock "github.com/keptn/keptn/shipyard-controller/common/fake"
	"github.com/keptn/keptn/shipyard-controller/db"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

const validShipyardResourceContent = `apiVersion: spec.keptn.sh/0.2.0
kind: Shipyard
metadata:
  name: test-shipyard
spec:
  stages:
  - name: dev
    sequences:
    - name: artifact-delivery
      tasks:
      - name: deployment
        properties:  
          strategy: direct`

const shipyardWithInvalidVersion = `apiVersion: 0.1.0
kind: Shipyard
metadata:
  name: test-shipyard
spec:
  stages:
  - name: dev
    sequences:
    - name: artifact-delivery
      tasks:
      - name: deployment
        properties:  
          strategy: direct`

const invalidShipyardContent = "invalid"

func TestShipyardRetriever_GetShipyard(t *testing.T) {
	type fields struct {
		configurationStore common.ConfigurationStore
		projectRepo        db.ProjectMVRepo
	}
	type args struct {
		projectName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *keptnv2.Shipyard
		wantErr bool
	}{
		{
			name: "get shipyard from configuration service",
			fields: fields{
				configurationStore: &common_mock.ConfigurationStoreMock{
					GetProjectResourceFunc: func(projectName string, resourceURI string) (*models.Resource, error) {
						return &models.Resource{
							ResourceContent: validShipyardResourceContent,
							ResourceURI:     stringp("shipyard.yaml"),
						}, nil
					},
				},
				projectRepo: &db_mock.ProjectMVRepoMock{
					UpdateShipyardFunc: func(projectName string, shipyard string) error {
						return nil
					},
					GetProjectFunc: func(projectName string) (*models.ExpandedProject, error) {
						return &models.ExpandedProject{ProjectName: "my-project"}, nil
					},
				},
			},
			args: args{
				projectName: "my-project",
			},
			want:    getTestShipyard(),
			wantErr: false,
		},
		{
			name: "Updating shipyard content fails -> shipyard should still be returned",
			fields: fields{
				configurationStore: &common_mock.ConfigurationStoreMock{
					GetProjectResourceFunc: func(projectName string, resourceURI string) (*models.Resource, error) {
						return &models.Resource{
							ResourceContent: validShipyardResourceContent,
							ResourceURI:     stringp("shipyard.yaml"),
						}, nil
					},
				},
				projectRepo: &db_mock.ProjectMVRepoMock{
					UpdateShipyardFunc: func(projectName string, shipyard string) error {
						return errors.New("oops")
					},
					GetProjectFunc: func(projectName string) (*models.ExpandedProject, error) {
						return &models.ExpandedProject{ProjectName: "my-project"}, nil
					},
				},
			},
			args: args{
				projectName: "my-project",
			},
			want:    getTestShipyard(),
			wantErr: false,
		},
		{
			name: "invalid shipyard version",
			fields: fields{
				configurationStore: &common_mock.ConfigurationStoreMock{
					GetProjectResourceFunc: func(projectName string, resourceURI string) (*models.Resource, error) {
						return &models.Resource{
							ResourceContent: shipyardWithInvalidVersion,
							ResourceURI:     stringp("shipyard.yaml"),
						}, nil
					},
				},
				projectRepo: &db_mock.ProjectMVRepoMock{
					UpdateShipyardFunc: func(projectName string, shipyard string) error {
						return errors.New("oops")
					},
					GetProjectFunc: func(projectName string) (*models.ExpandedProject, error) {
						return &models.ExpandedProject{ProjectName: "my-project"}, nil
					},
				},
			},
			args: args{
				projectName: "my-project",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid shipyard content",
			fields: fields{
				configurationStore: &common_mock.ConfigurationStoreMock{
					GetProjectResourceFunc: func(projectName string, resourceURI string) (*models.Resource, error) {
						return &models.Resource{
							ResourceContent: invalidShipyardContent,
							ResourceURI:     stringp("shipyard.yaml"),
						}, nil
					},
				},
				projectRepo: &db_mock.ProjectMVRepoMock{
					UpdateShipyardFunc: func(projectName string, shipyard string) error {
						return errors.New("oops")
					},
					GetProjectFunc: func(projectName string) (*models.ExpandedProject, error) {
						return &models.ExpandedProject{ProjectName: "my-project"}, nil
					},
				},
			},
			args: args{
				projectName: "my-project",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "resource cannot be retrieved",
			fields: fields{
				configurationStore: &common_mock.ConfigurationStoreMock{
					GetProjectResourceFunc: func(projectName string, resourceURI string) (*models.Resource, error) {
						return nil, errors.New("oops")
					},
				},
				projectRepo: &db_mock.ProjectMVRepoMock{
					UpdateShipyardFunc: func(projectName string, shipyard string) error {
						return nil
					},
					GetProjectFunc: func(projectName string) (*models.ExpandedProject, error) {
						return &models.ExpandedProject{ProjectName: "my-project"}, nil
					},
				},
			},
			args: args{
				projectName: "my-project",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := NewShipyardRetriever(tt.fields.configurationStore, tt.fields.projectRepo)
			got, err := sr.GetShipyard(tt.args.projectName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCachedShipyard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestShipyardRetriever_GetCachedShipyard(t *testing.T) {
	type fields struct {
		configurationStore common.ConfigurationStore
		projectRepo        db.ProjectMVRepo
	}
	type args struct {
		projectName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *keptnv2.Shipyard
		wantErr bool
	}{
		{
			name: "get shipyard from projects MV",
			fields: fields{
				projectRepo: &db_mock.ProjectMVRepoMock{
					UpdateShipyardFunc: func(projectName string, shipyard string) error {
						return nil
					},
					GetProjectFunc: func(projectName string) (*models.ExpandedProject, error) {
						return &models.ExpandedProject{ProjectName: "my-project", Shipyard: validShipyardResourceContent}, nil
					},
				},
			},
			args: args{
				projectName: "my-project",
			},
			want:    getTestShipyard(),
			wantErr: false,
		},
		{
			name: "get shipyard from projects MV - invalid shipyard",
			fields: fields{
				projectRepo: &db_mock.ProjectMVRepoMock{
					UpdateShipyardFunc: func(projectName string, shipyard string) error {
						return nil
					},
					GetProjectFunc: func(projectName string) (*models.ExpandedProject, error) {
						return &models.ExpandedProject{ProjectName: "my-project", Shipyard: invalidShipyardContent}, nil
					},
				},
			},
			args: args{
				projectName: "my-project",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "get shipyard from projects MV - cannot retrieve project",
			fields: fields{
				projectRepo: &db_mock.ProjectMVRepoMock{
					UpdateShipyardFunc: func(projectName string, shipyard string) error {
						return nil
					},
					GetProjectFunc: func(projectName string) (*models.ExpandedProject, error) {
						return nil, errors.New("oops")
					},
				},
			},
			args: args{
				projectName: "my-project",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := NewShipyardRetriever(tt.fields.configurationStore, tt.fields.projectRepo)
			got, err := sr.GetCachedShipyard(tt.args.projectName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetShipyard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetShipyard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestShipyard() *keptnv2.Shipyard {
	return &keptnv2.Shipyard{
		ApiVersion: "spec.keptn.sh/0.2.0",
		Kind:       "Shipyard",
		Metadata: keptnv2.Metadata{
			Name: "test-shipyard",
		},
		Spec: keptnv2.ShipyardSpec{
			Stages: []keptnv2.Stage{
				{
					Name: "dev",
					Sequences: []keptnv2.Sequence{
						{
							Name:        "artifact-delivery",
							TriggeredOn: nil,
							Tasks: []keptnv2.Task{
								{
									Name:           "deployment",
									TriggeredAfter: "",
									Properties: map[string]interface{}{
										"strategy": "direct",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func TestShipyardRetriever_GetLatestCommitID(t *testing.T) {
	type fields struct {
		configurationStore *common_mock.ConfigurationStoreMock
	}
	type args struct {
		projectName string
		stageName   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "get latest git commit id",
			fields: fields{
				configurationStore: &common_mock.ConfigurationStoreMock{
					GetStageResourceFunc: func(projectName string, stageName string, resourceURI string) (*models.Resource, error) {
						return &models.Resource{Metadata: &models.Version{
							Version: "my-commit-id",
						}}, nil
					},
				},
			},
			args: args{
				projectName: "my-project",
				stageName:   "my-stage",
			},
			want:    "my-commit-id",
			wantErr: false,
		},
		{
			name: "error while retrieving resource",
			fields: fields{
				configurationStore: &common_mock.ConfigurationStoreMock{
					GetStageResourceFunc: func(projectName string, stageName string, resourceURI string) (*models.Resource, error) {
						return nil, errors.New("oops")
					},
				},
			},
			args: args{
				projectName: "my-project",
				stageName:   "my-stage",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "no commit ID set",
			fields: fields{
				configurationStore: &common_mock.ConfigurationStoreMock{
					GetStageResourceFunc: func(projectName string, stageName string, resourceURI string) (*models.Resource, error) {
						return &models.Resource{}, nil
					},
				},
			},
			args: args{
				projectName: "my-project",
				stageName:   "my-stage",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := &ShipyardRetriever{
				configurationStore: tt.fields.configurationStore,
			}
			got, err := sr.GetLatestCommitID(tt.args.projectName, tt.args.stageName)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}
