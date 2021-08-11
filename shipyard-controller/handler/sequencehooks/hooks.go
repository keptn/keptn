package sequencehooks

import (
	"github.com/keptn/keptn/shipyard-controller/models"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencetriggered.go . ISequenceTriggeredHook
type ISequenceTriggeredHook interface {
	OnSequenceTriggered(models.Event)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencestarted.go . ISequenceStartedHook
type ISequenceStartedHook interface {
	OnSequenceStarted(models.Event)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencetasktriggered.go . ISequenceTaskTriggeredHook
type ISequenceTaskTriggeredHook interface {
	OnSequenceTaskTriggered(models.Event)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencetaskstarted.go . ISequenceTaskStartedHook
type ISequenceTaskStartedHook interface {
	OnSequenceTaskStarted(models.Event)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencetaskfinished.go . ISequenceTaskFinishedHook
type ISequenceTaskFinishedHook interface {
	OnSequenceTaskFinished(models.Event)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/subsequencefinished.go . ISubSequenceFinishedHook
type ISubSequenceFinishedHook interface {
	OnSubSequenceFinished(event models.Event)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencefinished.go . ISequenceFinishedHook
type ISequenceFinishedHook interface {
	OnSequenceFinished(event models.Event)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencetimeout.go . ISequenceTimeoutHook
type ISequenceTimeoutHook interface {
	OnSequenceTimeout(event models.Event)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencepause.go . ISequencePausedHook
type ISequencePausedHook interface {
	OnSequencePaused(pause models.EventScope)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencepause.go . ISequenceResumedHook
type ISequenceResumedHook interface {
	OnSequenceResumed(resume models.EventScope)
}
