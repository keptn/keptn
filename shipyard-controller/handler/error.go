package handler

import "errors"

var ErrProjectAlreadyExists = errors.New("project already exists")

var ErrServiceAlreadyExists = errors.New("service already exists")

var ErrServiceNotFound = errors.New("service not found")

var ErrProjectNotFound = errors.New("project not found")

var ErrInvalidStageChange = errors.New("stage name cannot be changed or removed")

var ErrStageNotFound = errors.New("stage not found")

var ErrChangesRollback = errors.New("failed to rollback changes")

var ErrOtherActiveSequencesRunning = errors.New("other sequences are currently running in the same stage for the same service")

var ErrSequencePaused = errors.New("sequence is paused")

var ErrSequenceBlocked = errors.New("sequence is currently blocked")

var ErrNoMatchingEvent = errors.New("no matching event found")

var ErrSequenceNotFound = errors.New("sequence not found")
