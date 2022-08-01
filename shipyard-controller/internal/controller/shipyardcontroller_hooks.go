package controller

import (
	"github.com/keptn/go-utils/pkg/api/models"
	scmodels "github.com/keptn/keptn/shipyard-controller/models"
)

func (sc *ShipyardController) AddSequenceTriggeredHook(hook ISequenceTriggeredHook) {
	sc.sequenceTriggeredHooks = append(sc.sequenceTriggeredHooks, hook)
}

func (sc *ShipyardController) AddSequenceStartedHook(hook ISequenceStartedHook) {
	sc.sequenceStartedHooks = append(sc.sequenceStartedHooks, hook)
}

func (sc *ShipyardController) AddSequenceWaitingHook(hook ISequenceWaitingHook) {
	sc.sequenceWaitingHooks = append(sc.sequenceWaitingHooks, hook)
}

func (sc *ShipyardController) AddSequenceTaskTriggeredHook(hook ISequenceTaskTriggeredHook) {
	sc.sequenceTaskTriggeredHooks = append(sc.sequenceTaskTriggeredHooks, hook)
}

func (sc *ShipyardController) AddSequenceTaskStartedHook(hook ISequenceTaskStartedHook) {
	sc.sequenceTaskStartedHooks = append(sc.sequenceTaskStartedHooks, hook)
}

func (sc *ShipyardController) AddSequenceTaskFinishedHook(hook ISequenceTaskFinishedHook) {
	sc.sequenceTaskFinishedHooks = append(sc.sequenceTaskFinishedHooks, hook)
}

func (sc *ShipyardController) AddSubSequenceFinishedHook(hook ISubSequenceFinishedHook) {
	sc.subSequenceFinishedHooks = append(sc.subSequenceFinishedHooks, hook)
}

func (sc *ShipyardController) AddSequenceFinishedHook(hook ISequenceFinishedHook) {
	sc.sequenceFinishedHooks = append(sc.sequenceFinishedHooks, hook)
}

func (sc *ShipyardController) AddSequenceTimeoutHook(hook ISequenceTimeoutHook) {
	sc.sequenceTimoutHooks = append(sc.sequenceTimoutHooks, hook)
}

func (sc *ShipyardController) AddSequencePausedHook(hook ISequencePausedHook) {
	sc.sequencePausedHooks = append(sc.sequencePausedHooks, hook)
}

func (sc *ShipyardController) AddSequenceResumedHook(hook ISequenceResumedHook) {
	sc.sequenceResumedHooks = append(sc.sequenceResumedHooks, hook)
}

func (sc *ShipyardController) AddSequenceAbortedHook(hook ISequenceAbortedHook) {
	sc.sequenceAbortedHooks = append(sc.sequenceAbortedHooks, hook)
}

func (sc *ShipyardController) onSequenceTriggered(event models.KeptnContextExtendedCE) {
	for _, hook := range sc.sequenceTriggeredHooks {
		hook.OnSequenceTriggered(event)
	}
}

func (sc *ShipyardController) onSequenceStarted(event models.KeptnContextExtendedCE) {
	for _, hook := range sc.sequenceStartedHooks {
		hook.OnSequenceStarted(event)
	}
}

func (sc *ShipyardController) onSequenceWaiting(event models.KeptnContextExtendedCE) {
	for _, hook := range sc.sequenceWaitingHooks {
		hook.OnSequenceWaiting(event)
	}
}

func (sc *ShipyardController) onSequenceTaskStarted(event models.KeptnContextExtendedCE) {
	for _, hook := range sc.sequenceTaskStartedHooks {
		hook.OnSequenceTaskStarted(event)
	}
}

func (sc *ShipyardController) onSequenceTaskTriggered(event models.KeptnContextExtendedCE) {
	for _, hook := range sc.sequenceTaskTriggeredHooks {
		hook.OnSequenceTaskTriggered(event)
	}
}

func (sc *ShipyardController) onSequenceTaskFinished(event models.KeptnContextExtendedCE) {
	for _, hook := range sc.sequenceTaskFinishedHooks {
		hook.OnSequenceTaskFinished(event)
	}
}

func (sc *ShipyardController) onSubSequenceFinished(event models.KeptnContextExtendedCE) {
	for _, hook := range sc.subSequenceFinishedHooks {
		hook.OnSubSequenceFinished(event)
	}
}

func (sc *ShipyardController) onSequenceFinished(event models.KeptnContextExtendedCE) {
	for _, hook := range sc.sequenceFinishedHooks {
		hook.OnSequenceFinished(event)
	}
}

func (sc *ShipyardController) onSequenceAborted(eventScope scmodels.EventScope) {
	for _, hook := range sc.sequenceAbortedHooks {
		hook.OnSequenceAborted(eventScope)
	}
}

func (sc *ShipyardController) onSequenceTimeout(event models.KeptnContextExtendedCE) {
	for _, hook := range sc.sequenceTimoutHooks {
		hook.OnSequenceTimeout(event)
	}
}

func (sc *ShipyardController) onSequencePaused(pause scmodels.EventScope) {
	for _, hook := range sc.sequencePausedHooks {
		hook.OnSequencePaused(pause)
	}
}

func (sc *ShipyardController) onSequenceResumed(resume scmodels.EventScope) {
	for _, hook := range sc.sequenceResumedHooks {
		hook.OnSequenceResumed(resume)
	}
}
