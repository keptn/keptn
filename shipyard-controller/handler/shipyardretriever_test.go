package handler

import (
	"errors"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	common_mock "github.com/keptn/keptn/shipyard-controller/common/fake"
	"github.com/keptn/keptn/shipyard-controller/db"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	scmodels "github.com/keptn/keptn/shipyard-controller/models"
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

const invalidShipyardResourceContent = `apiVersion: 0.1.0
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
					GetProjectFunc: func(projectName string) (*scmodels.ExpandedProject, error) {
						return &scmodels.ExpandedProject{ProjectName: "my-project"}, nil
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
					GetProjectFunc: func(projectName string) (*scmodels.ExpandedProject, error) {
						return &scmodels.ExpandedProject{ProjectName: "my-project"}, nil
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
			name: "invalid shipyard",
			fields: fields{
				configurationStore: &common_mock.ConfigurationStoreMock{
					GetProjectResourceFunc: func(projectName string, resourceURI string) (*models.Resource, error) {
						return &models.Resource{
							ResourceContent: invalidShipyardResourceContent,
							ResourceURI:     stringp("shipyard.yaml"),
						}, nil
					},
				},
				projectRepo: &db_mock.ProjectMVRepoMock{
					UpdateShipyardFunc: func(projectName string, shipyard string) error {
						return errors.New("oops")
					},
					GetProjectFunc: func(projectName string) (*scmodels.ExpandedProject, error) {
						return &scmodels.ExpandedProject{ProjectName: "my-project"}, nil
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := NewShipyardRetriever(tt.fields.configurationStore, tt.fields.projectRepo)
			got, err := sr.GetShipyard(tt.args.projectName)
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
