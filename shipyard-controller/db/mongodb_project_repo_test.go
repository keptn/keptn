package db

import (
	"fmt"
	"testing"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/common"
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
}
