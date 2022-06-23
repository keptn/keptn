package migration

import (
	"fmt"
	"testing"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

var projectNew = &apimodels.ExpandedProject{
	ProjectName: "test-project1",
	Stages: []*apimodels.ExpandedStage{
		{
			Services: []*apimodels.ExpandedService{
				{
					ServiceName: "test-service1",
				},
			},
		},
	},
	GitCredentials: &apimodels.GitAuthCredentialsSecure{
		RemoteURL: "http://some-url1",
		User:      "user1",
		HttpsAuth: &apimodels.HttpsGitAuthSecure{
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuthSecure{
				User:   "proxy-user1",
				Scheme: "http1",
				URL:    "url1",
			},
		},
	},
}

var projectOld = &models.ExpandedProjectOld{
	ProjectName: "test-project2",
	Stages: []*apimodels.ExpandedStage{
		{
			Services: []*apimodels.ExpandedService{
				{
					ServiceName: "test-service2",
				},
			},
		},
	},
	GitRemoteURI:    "http://some-url2",
	GitUser:         "user2",
	InsecureSkipTLS: true,
	GitProxyURL:     "url2",
	GitProxyScheme:  "http2",
	GitProxyUser:    "proxy-user2",
}

var projectOldToNew = &apimodels.ExpandedProject{
	ProjectName: "test-project2",
	Stages: []*apimodels.ExpandedStage{
		{
			Services: []*apimodels.ExpandedService{
				{
					ServiceName: "test-service2",
				},
			},
		},
	},
	GitCredentials: &apimodels.GitAuthCredentialsSecure{
		RemoteURL: "http://some-url2",
		User:      "user2",
		HttpsAuth: &apimodels.HttpsGitAuthSecure{
			InsecureSkipTLS: true,
			Proxy: &apimodels.ProxyGitAuthSecure{
				User:   "proxy-user2",
				Scheme: "http2",
				URL:    "url2",
			},
		},
	},
}

func TestProjectCredentialsMigration_TransformNewModel(t *testing.T) {
	config, err := rest.InClusterConfig()
	require.Nil(t, err)

	kubeAPI, err := kubernetes.NewForConfig(config)
	require.Nil(t, err)

	projectRepo := db.NewMongoDBProjectsRepo(db.GetMongoDBConnectionInstance())
	err = projectRepo.CreateProject(projectNew)
	require.Nil(t, err)

	migrator := NewProjectCredentialsMigrator(db.GetMongoDBConnectionInstance(), common.NewK8sSecretStore(kubeAPI))
	migrator.Transform()

	migratedProject, err := projectRepo.GetProject(projectNew.ProjectName)
	require.Nil(t, err)
	require.Equal(t, projectNew, migratedProject)

	err = projectRepo.DeleteProject(projectNew.ProjectName)
	require.Nil(t, err)
}

func TestProjectCredentialsMigration_TransformOldModel(t *testing.T) {
	config, err := rest.InClusterConfig()
	require.Nil(t, err)

	kubeAPI, err := kubernetes.NewForConfig(config)
	require.Nil(t, err)

	projectRepo := db.NewProjectCredentialsRepo(db.GetMongoDBConnectionInstance())
	err = projectRepo.CreateOldCredentialsProject(projectOld)
	require.Nil(t, err)

	migrator := NewProjectCredentialsMigrator(db.GetMongoDBConnectionInstance(), common.NewK8sSecretStore(kubeAPI))
	migrator.Transform()

	migratedProject, err := projectRepo.ProjectRepo.GetProject(projectOld.ProjectName)
	require.Nil(t, err)
	require.Equal(t, projectOldToNew, migratedProject)

	err = projectRepo.ProjectRepo.DeleteProject(projectOld.ProjectName)
	require.Nil(t, err)
}

func TestProjectCredentialsMigration_TransformMixed(t *testing.T) {
	config, err := rest.InClusterConfig()
	require.Nil(t, err)

	kubeAPI, err := kubernetes.NewForConfig(config)
	require.Nil(t, err)

	projectRepo := db.NewProjectCredentialsRepo(db.GetMongoDBConnectionInstance())

	err = projectRepo.CreateOldCredentialsProject(projectOld)
	require.Nil(t, err)

	err = projectRepo.ProjectRepo.CreateProject(projectNew)
	require.Nil(t, err)

	migrator := NewProjectCredentialsMigrator(db.GetMongoDBConnectionInstance(), common.NewK8sSecretStore(kubeAPI))
	migrator.Transform()

	migratedProject, err := projectRepo.ProjectRepo.GetProject(projectOld.ProjectName)
	require.Nil(t, err)
	require.Equal(t, projectOldToNew, migratedProject)

	migratedProject, err = projectRepo.ProjectRepo.GetProject(projectNew.ProjectName)
	require.Nil(t, err)
	require.Equal(t, projectNew, migratedProject)

	err = projectRepo.ProjectRepo.DeleteProject(projectOld.ProjectName)
	require.Nil(t, err)

	err = projectRepo.ProjectRepo.DeleteProject(projectNew.ProjectName)
	require.Nil(t, err)
}
