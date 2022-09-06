package db

import (
	"fmt"
	"testing"

	"github.com/keptn/keptn/shipyard-controller/internal/common"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/require"
)

func TestMongoDBProjectsRepo_InsertAndRetrieve(t *testing.T) {
	r := NewMongoDBProjectsRepo(GetMongoDBConnectionInstance())

	err := r.CreateProject(&apimodels.ExpandedProject{
		ProjectName: "my-project",
	})

	require.Nil(t, err)

	prj, err := r.GetProject("my-project")

	require.Nil(t, err)
	require.NotNil(t, prj)

	err = r.DeleteProject("my-project")
	require.Nil(t, err)

	prj, err = r.GetProject("my-project")

	require.ErrorIs(t, err, common.ErrProjectNotFound)
	require.Nil(t, prj)

	CleanupDB(t, r)
}

func TestMongoDBProjectsRepo_InsertAndRetrieveMultiple(t *testing.T) {
	r := NewMongoDBProjectsRepo(GetMongoDBConnectionInstance())

	nrProjects := 5

	for i := 0; i < nrProjects; i++ {
		projectName := fmt.Sprintf("project-%d", i)
		err := r.CreateProject(&apimodels.ExpandedProject{
			ProjectName: projectName,
		})

		require.Nil(t, err)
	}

	prj, err := r.GetProjects()

	require.Nil(t, err)
	require.NotNil(t, prj)
	require.Len(t, prj, 5)

	CleanupDB(t, r)
}

func TestMongoDBProjectsRepo_UpdateProject(t *testing.T) {
	r := NewMongoDBProjectsRepo(GetMongoDBConnectionInstance())

	err := r.CreateProject(&apimodels.ExpandedProject{
		ProjectName: "my-project",
	})

	require.Nil(t, err)

	updatedProject := &apimodels.ExpandedProject{
		ProjectName: "my-project",
		Shipyard:    "shipyard-content",
	}
	err = r.UpdateProject(updatedProject)

	require.Nil(t, err)

	prj, err := r.GetProject("my-project")

	require.Nil(t, err)
	require.NotNil(t, prj)

	require.Equal(t, updatedProject, prj)

	CleanupDB(t, r)
}

func CleanupDB(t *testing.T, r *MongoDBProjectsRepo) {
	prjs, err := r.GetProjects()
	require.Nil(t, err)

	for _, p := range prjs {
		err = r.DeleteProject(p.ProjectName)
		require.Nil(t, err)
	}
}

func TestMongoDBProjectsRepo_UpdateProjectService(t *testing.T) {
	r := NewMongoDBProjectsRepo(GetMongoDBConnectionInstance())

	err := r.CreateProject(&apimodels.ExpandedProject{
		ProjectName: "my-project",
		Stages: []*apimodels.ExpandedStage{
			{
				Services: []*apimodels.ExpandedService{
					{
						ServiceName: "my-service",
					},
					{
						ServiceName: "my-other-service",
					},
				},
				StageName: "dev",
			},
			{
				Services: []*apimodels.ExpandedService{
					{
						ServiceName: "my-service",
					},
				},
				StageName: "staging",
			},
		},
	})

	require.Nil(t, err)

	eventContextInfo := apimodels.EventContextInfo{
		EventID:      "event-id",
		KeptnContext: "event-context",
		Time:         "event-timestamp",
	}

	encodedEventType := encodeKey("sh.keptn.event.test.triggered")
	err = r.UpdateProjectService("my-project", "dev", "my-service", map[string]interface{}{
		"deployedImage":                      "my-new-image",
		"lastEventTypes." + encodedEventType: eventContextInfo,
	})

	require.Nil(t, err)

	projectInfo, err := r.GetProject("my-project")
	require.Nil(t, err)

	require.Equal(t, "my-new-image", projectInfo.Stages[0].Services[0].DeployedImage)
	require.Equal(t, eventContextInfo, projectInfo.Stages[0].Services[0].LastEventTypes[encodedEventType])

	require.Empty(t, projectInfo.Stages[0].Services[1].DeployedImage)
	require.Empty(t, projectInfo.Stages[0].Services[1].LastEventTypes)
}

func TestMongoDBKeyEncodingProjectsRepo_UpdateProjectService(t *testing.T) {
	r := NewMongoDBKeyEncodingProjectsRepo(GetMongoDBConnectionInstance())

	err := r.CreateProject(&apimodels.ExpandedProject{
		ProjectName: "my-project",
		Stages: []*apimodels.ExpandedStage{
			{
				Services: []*apimodels.ExpandedService{
					{
						ServiceName: "my-service",
					},
					{
						ServiceName: "my-other-service",
					},
				},
				StageName: "dev",
			},
			{
				Services: []*apimodels.ExpandedService{
					{
						ServiceName: "my-service",
					},
				},
				StageName: "staging",
			},
		},
	})

	require.Nil(t, err)

	eventContextInfo := apimodels.EventContextInfo{
		EventID:      "event-id",
		KeptnContext: "event-context",
		Time:         "event-timestamp",
	}

	err = r.UpdateProjectService("my-project", "dev", "my-service", map[string]interface{}{
		"deployedImage": "my-new-image",
		"lastEventTypes.sh.keptn.event.test.triggered": eventContextInfo,
	})

	require.Nil(t, err)

	projectInfo, err := r.GetProject("my-project")
	require.Nil(t, err)

	require.Equal(t, "my-new-image", projectInfo.Stages[0].Services[0].DeployedImage)
	require.Equal(t, eventContextInfo, projectInfo.Stages[0].Services[0].LastEventTypes["sh.keptn.event.test.triggered"])

	require.Empty(t, projectInfo.Stages[0].Services[1].DeployedImage)
	require.Empty(t, projectInfo.Stages[0].Services[1].LastEventTypes)
}
