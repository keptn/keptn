package db

import (
	"encoding/json"
	"testing"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/require"
)

func TestMongoDBDumpRepo_ListAllCollections(t *testing.T) {
	dumpRepo := NewMongoDBDumpRepo(GetMongoDBConnectionInstance())
	projectsRepo := NewMongoDBProjectsRepo(GetMongoDBConnectionInstance())

	err := projectsRepo.CreateProject(&apimodels.ExpandedProject{
		ProjectName: "my-project",
	})
	require.Nil(t, err)

	projectNames, err := dumpRepo.ListAllCollections()

	require.Len(t, projectNames, 1)
	require.Nil(t, err)
	require.Contains(t, projectNames, "keptnProjectsMV")

	CleanupDB(t, projectsRepo)
}

func TestMongoDBDumpRepo_GetDump(t *testing.T) {
	dumpRepo := NewMongoDBDumpRepo(GetMongoDBConnectionInstance())
	projectsRepo := NewMongoDBProjectsRepo(GetMongoDBConnectionInstance())

	// get random collection dump
	_, err := dumpRepo.GetDump("test")
	require.Nil(t, err)

	// compare to project repo: 0 projects
	projects, err := projectsRepo.GetProjects()
	require.Nil(t, err)

	dump, err := dumpRepo.GetDump("keptnProjectsMV")
	require.Nil(t, err)

	var projects_dump []*apimodels.ExpandedProject
	bytes, err := json.Marshal(dump)
	err = json.Unmarshal(bytes, &projects_dump)

	require.Nil(t, err)
	require.Equal(t, projects_dump, projects)

	// compare to project repo: 1 project
	err = projectsRepo.CreateProject(&apimodels.ExpandedProject{
		ProjectName: "my-project",
	})
	require.Nil(t, err)

	projects, err = projectsRepo.GetProjects()
	require.Nil(t, err)

	dump, err = dumpRepo.GetDump("keptnProjectsMV")
	require.Nil(t, err)

	bytes, err = json.Marshal(dump)
	err = json.Unmarshal(bytes, &projects_dump)

	require.Nil(t, err)
	require.Equal(t, projects_dump, projects)

	// compare to project repo: 2 projects
	err = projectsRepo.CreateProject(&apimodels.ExpandedProject{
		ProjectName: "my-project2",
	})
	require.Nil(t, err)

	projects, err = projectsRepo.GetProjects()
	require.Nil(t, err)

	dump, err = dumpRepo.GetDump("keptnProjectsMV")
	require.Nil(t, err)

	bytes, err = json.Marshal(dump)
	err = json.Unmarshal(bytes, &projects_dump)

	require.Nil(t, err)
	require.Equal(t, projects_dump, projects)

	CleanupDB(t, projectsRepo)
}
