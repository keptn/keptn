package db_test

import (
	"encoding/json"
	"testing"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/internal/db"
	"github.com/stretchr/testify/require"
)

func TestMongoDBDumpRepo_ListAllCollections(t *testing.T) {
	repo := db.NewMongoDBDumpRepo(db.GetMongoDBConnectionInstance())

	projectNames, err := repo.ListAllCollections()
	require.NotEmpty(t, projectNames)
	require.Nil(t, err)
}

func TestMongoDBDumpRepo_GetDump(t *testing.T) {
	dumpRepo := db.NewMongoDBDumpRepo(db.GetMongoDBConnectionInstance())
	projectsRepo := db.NewMongoDBProjectsRepo(db.GetMongoDBConnectionInstance())

	// get random collection dump
	_, err := dumpRepo.GetDump("test")
	require.Nil(t, err)

	// compare project collection to projectsrepo collection
	projects_repo, err := projectsRepo.GetProjects()
	require.Nil(t, err)

	dump, err := dumpRepo.GetDump("keptnProjectsMV")
	require.Nil(t, err)

	var projects_dump []*apimodels.ExpandedProject
	bytes, err := json.Marshal(dump)
	err = json.Unmarshal(bytes, &projects_dump)

	require.Nil(t, err)
	require.Equal(t, projects_dump, projects_repo)
}
