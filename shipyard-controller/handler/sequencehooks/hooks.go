package sequencehooks

import "github.com/keptn/keptn/shipyard-controller/models"

type ISequenceTriggeredHook interface {
	OnSequenceTriggered(models.Event) error
}

type ISequenceTaskTriggeredHook interface {
	OnSequenceTaskTriggered(models.Event) error
}

type ISequenceTaskStartedHook interface {
	OnSequenceTaskStarted(models.Event) error
}

type ISequenceTaskFinishedHook interface {
	OnSequenceTaskFinished(models.Event) error
}

type ISequenceFinishedHook interface {
	OnSequenceFinished(event models.Event) error
}
