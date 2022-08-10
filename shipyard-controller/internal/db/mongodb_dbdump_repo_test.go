package db_test

import (
	"testing"

	"github.com/keptn/keptn/shipyard-controller/internal/db"
	"github.com/stretchr/testify/require"
)

func TestMongoDBDumpRepo_GetDump(t *testing.T) {
	repo := db.NewMongoDBDumpRepo(db.GetMongoDBConnectionInstance())

	projectNames, err := repo.ListAllCollections()
	require.NotEmpty(t, projectNames)

	_, err = repo.GetDump("test")
	require.Nil(t, err)
}
