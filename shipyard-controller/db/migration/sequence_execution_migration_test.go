package migration

import (
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSequenceExecutionMigrator_MigrateSequenceExecutions(t *testing.T) {
	defer setupLocalMongoDB()()
	sm := NewSequenceExecutionMigrator(db.GetMongoDBConnectionInstance())

	err := sm.MigrateSequenceExecutions()

	require.Nil(t, err)
}
