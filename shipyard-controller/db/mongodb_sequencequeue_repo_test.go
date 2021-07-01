package db

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_MongoDBSequenceRepoInsertAndRetrieve(t *testing.T) {

	nowTime := time.Now()

	queueItem1 := models.QueueItem{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			KeptnContext: "my-context-1",
			EventType:    keptnv2.GetTriggeredEventType("dev.delivery"),
		},
		EventID:   "my-id-1",
		Timestamp: nowTime,
	}

	queueItem2 := models.QueueItem{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			KeptnContext: "my-context-2",
			EventType:    keptnv2.GetTriggeredEventType("dev.delivery"),
		},
		EventID:   "my-id-2",
		Timestamp: nowTime.Add(2 * time.Second),
	}

	queueItem3 := models.QueueItem{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			KeptnContext: "my-context-2",
			EventType:    keptnv2.GetTriggeredEventType("dev.delivery"),
		},
		EventID:   "my-id-3",
		Timestamp: nowTime.Add(1 * time.Second),
	}

	mdbrepo := NewMongoDBSequenceQueueRepo(GetMongoDBConnectionInstance())

	// first, delete all sequences
	err := mdbrepo.DeleteQueuedSequences(models.QueueItem{})
	require.Nil(t, err)

	// insert queue items into collection
	err = mdbrepo.QueueSequence(queueItem1)
	require.Nil(t, err)

	err = mdbrepo.QueueSequence(queueItem2)
	require.Nil(t, err)

	err = mdbrepo.QueueSequence(queueItem3)
	require.Nil(t, err)

	// retrieve items -> should be sorted in an ascending order, not in the order they have been stored
	sequences, err := mdbrepo.GetQueuedSequences()

	require.Nil(t, err)
	require.Len(t, sequences, 3)
	require.Equal(t, queueItem1, sequences[0])
	require.Equal(t, queueItem3, sequences[1])
	require.Equal(t, queueItem2, sequences[2])

	// filter sequence with "my-id-1"
	err = mdbrepo.DeleteQueuedSequences(models.QueueItem{EventID: queueItem1.EventID})
	require.Nil(t, err)

	// retrieve sequences again - first one should not be included anymore
	sequences, err = mdbrepo.GetQueuedSequences()
	require.Nil(t, err)
	require.Len(t, sequences, 2)
	require.Equal(t, queueItem3, sequences[0])
	require.Equal(t, queueItem2, sequences[1])

	// delete all sequences for the project "my-project"
	err = mdbrepo.DeleteQueuedSequences(models.QueueItem{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: "my-project",
			},
		},
	})
	require.Nil(t, err)

	// retrieve sequence queue again - should now be empty
	sequences, err = mdbrepo.GetQueuedSequences()
	require.Equal(t, ErrNoEventFound, err)
}
