package db_test

import (
	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMongoDBLogRepo_InsertAndRetrieve(t *testing.T) {
	repo := db.NewMongoDBLogRepo(db.GetMongoDBConnectionInstance())

	mockClock := clock.NewMock()
	repo.TheClock = mockClock

	timeFrom := repo.TheClock.Now().UTC()
	timeTo := timeFrom.Add(10 * time.Second).UTC()

	// insert the first log entry
	err := repo.CreateLogEntries([]models.LogEntry{
		{
			IntegrationID: "my-integration-id",
			Message:       "my message",
		},
	})
	require.Nil(t, err)

	mockClock.Add(11 * time.Second)

	// insert the second log entry
	err = repo.CreateLogEntries([]models.LogEntry{
		{
			IntegrationID: "my-integration-id",
			Message:       "my second message",
		},
	})
	require.Nil(t, err)

	// retrieve entries without time filter - should return both
	entries, err := repo.GetLogEntries(models.GetLogParams{
		LogFilter: models.LogFilter{
			IntegrationID: "my-integration-id",
		},
	})

	require.Nil(t, err)
	require.Len(t, entries.Logs, 2)
	require.Equal(t, int64(2), entries.TotalCount)

	// retrieve with time filter - should only return the first entry
	entries, err = repo.GetLogEntries(models.GetLogParams{
		LogFilter: models.LogFilter{
			IntegrationID: "my-integration-id",
			FromTime:      timeutils.GetKeptnTimeStamp(timeFrom),
			BeforeTime:    timeutils.GetKeptnTimeStamp(timeTo),
		},
	})

	require.Nil(t, err)
	require.Len(t, entries.Logs, 1)
	require.Equal(t, int64(1), entries.TotalCount)
	require.Equal(t, "my message", entries.Logs[0].Message)

	// check if pagination works
	entries, err = repo.GetLogEntries(models.GetLogParams{
		PageSize: 1,
		LogFilter: models.LogFilter{
			IntegrationID: "my-integration-id",
		},
	})

	require.Nil(t, err)
	require.Len(t, entries.Logs, 1)
	require.Equal(t, int64(2), entries.TotalCount)
	require.Equal(t, int64(1), entries.NextPageKey)
	require.Equal(t, "my second message", entries.Logs[0].Message)

	entries, err = repo.GetLogEntries(models.GetLogParams{
		PageSize:    1,
		NextPageKey: 1,
		LogFilter: models.LogFilter{
			IntegrationID: "my-integration-id",
		},
	})

	require.Nil(t, err)
	require.Len(t, entries.Logs, 1)
	require.Equal(t, int64(2), entries.TotalCount)
	require.Equal(t, int64(0), entries.NextPageKey)
	require.Equal(t, "my message", entries.Logs[0].Message)

	// delete log entries for my-integration-id
	err = repo.DeleteLogEntries(models.DeleteLogParams{
		LogFilter: models.LogFilter{IntegrationID: "my-integration-id"},
	})

	require.Nil(t, err)

	entries, err = repo.GetLogEntries(models.GetLogParams{
		LogFilter: models.LogFilter{
			IntegrationID: "my-integration-id",
		},
	})

	require.Nil(t, err)
	require.Len(t, entries.Logs, 0)
	require.Equal(t, int64(0), entries.TotalCount)
}

func TestMongoDBLogRepo_SetupTTLIndex(t *testing.T) {

	mdbrepo := db.NewMongoDBLogRepo(db.GetMongoDBConnectionInstance())

	err := mdbrepo.SetupTTLIndex(10)
	require.Nil(t, err)
}
