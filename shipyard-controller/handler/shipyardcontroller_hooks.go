package handler

import (
	"github.com/keptn/keptn/shipyard-controller/handler/sequencehooks"
	"github.com/keptn/keptn/shipyard-controller/models"
)

func (sc *shipyardController) AddSequenceTriggeredHook(hook sequencehooks.ISequenceTriggeredHook) {
	sc.sequenceTriggeredHooks = append(sc.sequenceTriggeredHooks, hook)
}

func (sc *shipyardController) AddSequenceStartedHook(hook sequencehooks.ISequenceStartedHook) {
	sc.sequenceStartedHooks = append(sc.sequenceStartedHooks, hook)
}

func (sc *shipyardController) AddSequenceTaskTriggeredHook(hook sequencehooks.ISequenceTaskTriggeredHook) {
	sc.sequenceTaskTriggeredHooks = append(sc.sequenceTaskTriggeredHooks, hook)
}

func (sc *shipyardController) AddSequenceTaskStartedHook(hook sequencehooks.ISequenceTaskStartedHook) {
	sc.sequenceTaskStartedHooks = append(sc.sequenceTaskStartedHooks, hook)
}

func (sc *shipyardController) AddSequenceTaskFinishedHook(hook sequencehooks.ISequenceTaskFinishedHook) {
	sc.sequenceTaskFinishedHooks = append(sc.sequenceTaskFinishedHooks, hook)
}

func (sc *shipyardController) AddSubSequenceFinishedHook(hook sequencehooks.ISubSequenceFinishedHook) {
	sc.subSequenceFinishedHooks = append(sc.subSequenceFinishedHooks, hook)
}

func (sc *shipyardController) AddSequenceFinishedHook(hook sequencehooks.ISequenceFinishedHook) {
	sc.sequenceFinishedHooks = append(sc.sequenceFinishedHooks, hook)
}

func (sc *shipyardController) AddSequenceTimeoutHook(hook sequencehooks.ISequenceTimeoutHook) {
	sc.sequenceTimoutHooks = append(sc.sequenceTimoutHooks, hook)
}

func (sc *shipyardController) AddSequencePausedHook(hook sequencehooks.ISequencePausedHook) {
	sc.sequencePausedHooks = append(sc.sequencePausedHooks, hook)
}

func (sc *shipyardController) AddSequenceResumedHook(hook sequencehooks.ISequenceResumedHook) {
	sc.sequenceResumedHooks = append(sc.sequenceResumedHooks, hook)
}

func (sc *shipyardController) onSequenceTriggered(event models.Event) {
	for _, hook := range sc.sequenceTriggeredHooks {
		hook.OnSequenceTriggered(event)
	}
}

func (sc *shipyardController) onSequenceStarted(event models.Event) {
	for _, hook := range sc.sequenceStartedHooks {
		hook.OnSequenceStarted(event)
	}
}

func (sc *shipyardController) onSequenceTaskStarted(event models.Event) {
	for _, hook := range sc.sequenceTaskStartedHooks {
		hook.OnSequenceTaskStarted(event)
	}
}

func (sc *shipyardController) onSequenceTaskTriggered(event models.Event) {
	for _, hook := range sc.sequenceTaskTriggeredHooks {
		hook.OnSequenceTaskTriggered(event)
	}
}

func (sc *shipyardController) onSequenceTaskFinished(event models.Event) {
	for _, hook := range sc.sequenceTaskFinishedHooks {
		hook.OnSequenceTaskFinished(event)
	}
}

func (sc *shipyardController) onSubSequenceFinished(event models.Event) {
	for _, hook := range sc.subSequenceFinishedHooks {
		hook.OnSubSequenceFinished(event)
	}
}

func (sc *shipyardController) onSequenceFinished(event models.Event) {
	for _, hook := range sc.sequenceFinishedHooks {
		hook.OnSequenceFinished(event)
	}
}

func (sc *shipyardController) onSequenceTimeout(event models.Event) {
	for _, hook := range sc.sequenceTimoutHooks {
		hook.OnSequenceTimeout(event)
	}
}

func (sc *shipyardController) onSequencePaused(pause models.EventScope) {
	for _, hook := range sc.sequencePausedHooks {
		hook.OnSequencePaused(pause)
	}
}

func (sc *shipyardController) onSequenceResumed(resume models.EventScope) {
	for _, hook := range sc.sequenceResumedHooks {
		hook.OnSequenceResumed(resume)
	}
}
