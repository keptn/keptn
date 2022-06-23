package migration

import (
	"fmt"
	"testing"

	"github.com/keptn/keptn/shipyard-controller/db"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
)

func TestProjectCredentialsMigration_Transform(t *testing.T) {
	tests := []struct {
		name        string
		projectRepo db.ProjectCredentialsRepo
		secretRepo  db.SecretCredentialsRepo
		want        error
	}{
		{
			name: "valid",
			secretRepo: db_mock.SecretCredentialsRepoMock{
				UpdateSecretFunc: func(project *models.ExpandedProjectOld) error {
					return nil
				},
			},
			projectRepo: db_mock.ProjectCredentialsRepoMock{
				GetOldCredentialsProjectsFunc: func() ([]*models.ExpandedProjectOld, error) {
					return []*models.ExpandedProjectOld{
						{
							CreationDate:     "date",
							LastEventContext: nil,
							ProjectName:      "project",
							Shipyard:         "shippy",
							ShipyardVersion:  "ship",
							Stages:           nil,
							GitRemoteURI:     "http://some-url",
							GitUser:          "user",
							InsecureSkipTLS:  false,
						},
						{
							CreationDate:     "date",
							LastEventContext: nil,
							ProjectName:      "project2",
							Shipyard:         "shippy",
							ShipyardVersion:  "ship",
							Stages:           nil,
							GitRemoteURI:     "http://some-url2",
							GitUser:          "user",
							InsecureSkipTLS:  false,
						},
					}, nil
				},
				UpdateProjectFunc: func(project *models.ExpandedProjectOld) error {
					return nil
				},
			},
			want: nil,
		},
		{
			name: "no projects",
			secretRepo: db_mock.SecretCredentialsRepoMock{
				UpdateSecretFunc: func(project *models.ExpandedProjectOld) error {
					return nil
				},
			},
			projectRepo: db_mock.ProjectCredentialsRepoMock{
				GetOldCredentialsProjectsFunc: func() ([]*models.ExpandedProjectOld, error) {
					return []*models.ExpandedProjectOld{}, nil
				},
				UpdateProjectFunc: func(project *models.ExpandedProjectOld) error {
					return nil
				},
			},
			want: nil,
		},
		{
			name: "nil projects",
			secretRepo: db_mock.SecretCredentialsRepoMock{
				UpdateSecretFunc: func(project *models.ExpandedProjectOld) error {
					return nil
				},
			},
			projectRepo: db_mock.ProjectCredentialsRepoMock{
				GetOldCredentialsProjectsFunc: func() ([]*models.ExpandedProjectOld, error) {
					return nil, nil
				},
				UpdateProjectFunc: func(project *models.ExpandedProjectOld) error {
					return nil
				},
			},
			want: nil,
		},
		{
			name: "get projects err",
			projectRepo: db_mock.ProjectCredentialsRepoMock{
				GetOldCredentialsProjectsFunc: func() ([]*models.ExpandedProjectOld, error) {
					return nil, fmt.Errorf("some err")
				},
			},
			want: fmt.Errorf("could not transform git credentials to new format: some err"),
		},
		{
			name: "update project err",
			secretRepo: db_mock.SecretCredentialsRepoMock{
				UpdateSecretFunc: func(project *models.ExpandedProjectOld) error {
					return nil
				},
			},
			projectRepo: db_mock.ProjectCredentialsRepoMock{
				GetOldCredentialsProjectsFunc: func() ([]*models.ExpandedProjectOld, error) {
					return []*models.ExpandedProjectOld{
						{
							CreationDate:     "date",
							LastEventContext: nil,
							ProjectName:      "project",
							Shipyard:         "shippy",
							ShipyardVersion:  "ship",
							Stages:           nil,
							GitRemoteURI:     "http://some-url",
							GitUser:          "user",
							InsecureSkipTLS:  false,
						},
						{
							CreationDate:     "date",
							LastEventContext: nil,
							ProjectName:      "project2",
							Shipyard:         "shippy",
							ShipyardVersion:  "ship",
							Stages:           nil,
							GitRemoteURI:     "http://some-url2",
							GitUser:          "user",
							InsecureSkipTLS:  false,
						},
					}, nil
				},
				UpdateProjectFunc: func(project *models.ExpandedProjectOld) error {
					return fmt.Errorf("some err")
				},
			},
			want: fmt.Errorf("could not transform git credentials for project project: some err"),
		},
		{
			name: "update secret err",
			secretRepo: db_mock.SecretCredentialsRepoMock{
				UpdateSecretFunc: func(project *models.ExpandedProjectOld) error {
					return fmt.Errorf("some err")
				},
			},
			projectRepo: db_mock.ProjectCredentialsRepoMock{
				GetOldCredentialsProjectsFunc: func() ([]*models.ExpandedProjectOld, error) {
					return []*models.ExpandedProjectOld{
						{
							CreationDate:     "date",
							LastEventContext: nil,
							ProjectName:      "project",
							Shipyard:         "shippy",
							ShipyardVersion:  "ship",
							Stages:           nil,
							GitRemoteURI:     "http://some-url",
							GitUser:          "user",
							InsecureSkipTLS:  false,
						},
						{
							CreationDate:     "date",
							LastEventContext: nil,
							ProjectName:      "project2",
							Shipyard:         "shippy",
							ShipyardVersion:  "ship",
							Stages:           nil,
							GitRemoteURI:     "http://some-url2",
							GitUser:          "user",
							InsecureSkipTLS:  false,
						},
					}, nil
				},
			},
			want: fmt.Errorf("could not transform git credentials for project project: some err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := ProjectCredentialsMigrator{
				secretRepo:  tt.secretRepo,
				projectRepo: tt.projectRepo,
			}
			err := repo.Transform()
			require.Equal(t, tt.want, err)
		})
	}
}
