package sequencehooks

import "github.com/keptn/keptn/shipyard-controller/models"

type ISequenceTriggeredHook interface {
	OnSequenceTriggered(models.Event)
}

type ISequenceTaskTriggeredHook interface {
	OnSequenceTaskTriggered(models.Event)
}

type ISequenceTaskStartedHook interface {
	OnSequenceTaskStarted(models.Event)
}

type ISequenceTaskFinishedHook interface {
	OnSequenceTaskFinished(models.Event)
}

type ISequenceFinishedHook interface {
	OnSequenceFinished(event models.Event)
}
