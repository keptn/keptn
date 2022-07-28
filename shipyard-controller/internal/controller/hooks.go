package controller

import (
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/models"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencetriggered.go . ISequenceTriggeredHook
type ISequenceTriggeredHook interface {
	OnSequenceTriggered(apimodels.KeptnContextExtendedCE)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencestarted.go . ISequenceStartedHook
type ISequenceStartedHook interface {
	OnSequenceStarted(apimodels.KeptnContextExtendedCE)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencewaiting.go . ISequenceWaitingHook
type ISequenceWaitingHook interface {
	OnSequenceWaiting(apimodels.KeptnContextExtendedCE)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencetasktriggered.go . ISequenceTaskTriggeredHook
type ISequenceTaskTriggeredHook interface {
	OnSequenceTaskTriggered(apimodels.KeptnContextExtendedCE)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencetaskstarted.go . ISequenceTaskStartedHook
type ISequenceTaskStartedHook interface {
	OnSequenceTaskStarted(apimodels.KeptnContextExtendedCE)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencetaskfinished.go . ISequenceTaskFinishedHook
type ISequenceTaskFinishedHook interface {
	OnSequenceTaskFinished(apimodels.KeptnContextExtendedCE)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/subsequencefinished.go . ISubSequenceFinishedHook
type ISubSequenceFinishedHook interface {
	OnSubSequenceFinished(event apimodels.KeptnContextExtendedCE)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencefinished.go . ISequenceFinishedHook
type ISequenceFinishedHook interface {
	OnSequenceFinished(event apimodels.KeptnContextExtendedCE)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequenceaborted.go . ISequenceAbortedHook
type ISequenceAbortedHook interface {
	OnSequenceAborted(event models.EventScope)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencetimeout.go . ISequenceTimeoutHook
type ISequenceTimeoutHook interface {
	OnSequenceTimeout(event apimodels.KeptnContextExtendedCE)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencepause.go . ISequencePausedHook
type ISequencePausedHook interface {
	OnSequencePaused(pause models.EventScope)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencepause.go . ISequenceResumedHook
type ISequenceResumedHook interface {
	OnSequenceResumed(resume models.EventScope)
}
