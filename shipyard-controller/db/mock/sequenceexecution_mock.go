// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package db_mock

import (
	"github.com/keptn/keptn/shipyard-controller/models"
	"sync"
)

// SequenceExecutionRepoMock is a mock implementation of db.SequenceExecutionRepo.
//
// 	func TestSomethingThatUsesSequenceExecutionRepo(t *testing.T) {
//
// 		// make and configure a mocked db.SequenceExecutionRepo
// 		mockedSequenceExecutionRepo := &SequenceExecutionRepoMock{
// 			AppendTaskEventFunc: func(taskSequence models.SequenceExecution, event models.TaskEvent) (*models.SequenceExecution, error) {
// 				panic("mock out the AppendTaskEvent method")
// 			},
// 			ClearFunc: func(projectName string) error {
// 				panic("mock out the Clear method")
// 			},
// 			GetFunc: func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
// 				panic("mock out the Get method")
// 			},
// 			GetByTriggeredIDFunc: func(project string, triggeredID string) (*models.SequenceExecution, error) {
// 				panic("mock out the GetByTriggeredID method")
// 			},
// 			IsContextPausedFunc: func(eventScope models.EventScope) bool {
// 				panic("mock out the IsContextPaused method")
// 			},
// 			PauseContextFunc: func(eventScope models.EventScope) error {
// 				panic("mock out the PauseContext method")
// 			},
// 			ResumeContextFunc: func(eventScope models.EventScope) error {
// 				panic("mock out the ResumeContext method")
// 			},
// 			UpdateStatusFunc: func(taskSequence models.SequenceExecution) (*models.SequenceExecution, error) {
// 				panic("mock out the UpdateStatus method")
// 			},
// 			UpsertFunc: func(item models.SequenceExecution, options *models.SequenceExecutionUpsertOptions) error {
// 				panic("mock out the Upsert method")
// 			},
// 		}
//
// 		// use mockedSequenceExecutionRepo in code that requires db.SequenceExecutionRepo
// 		// and then make assertions.
//
// 	}
type SequenceExecutionRepoMock struct {
	// AppendTaskEventFunc mocks the AppendTaskEvent method.
	AppendTaskEventFunc func(taskSequence models.SequenceExecution, event models.TaskEvent) (*models.SequenceExecution, error)

	// ClearFunc mocks the Clear method.
	ClearFunc func(projectName string) error

	// GetFunc mocks the Get method.
	GetFunc func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error)

	// GetByTriggeredIDFunc mocks the GetByTriggeredID method.
	GetByTriggeredIDFunc func(project string, triggeredID string) (*models.SequenceExecution, error)

	// IsContextPausedFunc mocks the IsContextPaused method.
	IsContextPausedFunc func(eventScope models.EventScope) bool

	// PauseContextFunc mocks the PauseContext method.
	PauseContextFunc func(eventScope models.EventScope) error

	// ResumeContextFunc mocks the ResumeContext method.
	ResumeContextFunc func(eventScope models.EventScope) error

	// UpdateStatusFunc mocks the UpdateStatus method.
	UpdateStatusFunc func(taskSequence models.SequenceExecution) (*models.SequenceExecution, error)

	// UpsertFunc mocks the Upsert method.
	UpsertFunc func(item models.SequenceExecution, options *models.SequenceExecutionUpsertOptions) error

	// calls tracks calls to the methods.
	calls struct {
		// AppendTaskEvent holds details about calls to the AppendTaskEvent method.
		AppendTaskEvent []struct {
			// TaskSequence is the taskSequence argument value.
			TaskSequence models.SequenceExecution
			//models.KeptnContextExtendedCEis the event argument value.
			Event models.TaskEvent
		}
		// Clear holds details about calls to the Clear method.
		Clear []struct {
			// ProjectName is the projectName argument value.
			ProjectName string
		}
		// Get holds details about calls to the Get method.
		Get []struct {
			// Filter is the filter argument value.
			Filter models.SequenceExecutionFilter
		}
		// GetByTriggeredID holds details about calls to the GetByTriggeredID method.
		GetByTriggeredID []struct {
			// Project is the project argument value.
			Project string
			// TriggeredID is the triggeredID argument value.
			TriggeredID string
		}
		// IsContextPaused holds details about calls to the IsContextPaused method.
		IsContextPaused []struct {
			// EventScope is the eventScope argument value.
			EventScope models.EventScope
		}
		// PauseContext holds details about calls to the PauseContext method.
		PauseContext []struct {
			// EventScope is the eventScope argument value.
			EventScope models.EventScope
		}
		// ResumeContext holds details about calls to the ResumeContext method.
		ResumeContext []struct {
			// EventScope is the eventScope argument value.
			EventScope models.EventScope
		}
		// UpdateStatus holds details about calls to the UpdateStatus method.
		UpdateStatus []struct {
			// TaskSequence is the taskSequence argument value.
			TaskSequence models.SequenceExecution
		}
		// Upsert holds details about calls to the Upsert method.
		Upsert []struct {
			// Item is the item argument value.
			Item models.SequenceExecution
			// Options is the options argument value.
			Options *models.SequenceExecutionUpsertOptions
		}
	}
	lockAppendTaskEvent  sync.RWMutex
	lockClear            sync.RWMutex
	lockGet              sync.RWMutex
	lockGetByTriggeredID sync.RWMutex
	lockIsContextPaused  sync.RWMutex
	lockPauseContext     sync.RWMutex
	lockResumeContext    sync.RWMutex
	lockUpdateStatus     sync.RWMutex
	lockUpsert           sync.RWMutex
}

// AppendTaskEvent calls AppendTaskEventFunc.
func (mock *SequenceExecutionRepoMock) AppendTaskEvent(taskSequence models.SequenceExecution, event models.TaskEvent) (*models.SequenceExecution, error) {
	if mock.AppendTaskEventFunc == nil {
		panic("SequenceExecutionRepoMock.AppendTaskEventFunc: method is nil but SequenceExecutionRepo.AppendTaskEvent was just called")
	}
	callInfo := struct {
		TaskSequence models.SequenceExecution
		Event        models.TaskEvent
	}{
		TaskSequence: taskSequence,
		Event:        event,
	}
	mock.lockAppendTaskEvent.Lock()
	mock.calls.AppendTaskEvent = append(mock.calls.AppendTaskEvent, callInfo)
	mock.lockAppendTaskEvent.Unlock()
	return mock.AppendTaskEventFunc(taskSequence, event)
}

// AppendTaskEventCalls gets all the calls that were made to AppendTaskEvent.
// Check the length with:
//     len(mockedSequenceExecutionRepo.AppendTaskEventCalls())
func (mock *SequenceExecutionRepoMock) AppendTaskEventCalls() []struct {
	TaskSequence models.SequenceExecution
	Event        models.TaskEvent
} {
	var calls []struct {
		TaskSequence models.SequenceExecution
		Event        models.TaskEvent
	}
	mock.lockAppendTaskEvent.RLock()
	calls = mock.calls.AppendTaskEvent
	mock.lockAppendTaskEvent.RUnlock()
	return calls
}

// Clear calls ClearFunc.
func (mock *SequenceExecutionRepoMock) Clear(projectName string) error {
	if mock.ClearFunc == nil {
		panic("SequenceExecutionRepoMock.ClearFunc: method is nil but SequenceExecutionRepo.Clear was just called")
	}
	callInfo := struct {
		ProjectName string
	}{
		ProjectName: projectName,
	}
	mock.lockClear.Lock()
	mock.calls.Clear = append(mock.calls.Clear, callInfo)
	mock.lockClear.Unlock()
	return mock.ClearFunc(projectName)
}

// ClearCalls gets all the calls that were made to Clear.
// Check the length with:
//     len(mockedSequenceExecutionRepo.ClearCalls())
func (mock *SequenceExecutionRepoMock) ClearCalls() []struct {
	ProjectName string
} {
	var calls []struct {
		ProjectName string
	}
	mock.lockClear.RLock()
	calls = mock.calls.Clear
	mock.lockClear.RUnlock()
	return calls
}

// Get calls GetFunc.
func (mock *SequenceExecutionRepoMock) Get(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
	if mock.GetFunc == nil {
		panic("SequenceExecutionRepoMock.GetFunc: method is nil but SequenceExecutionRepo.Get was just called")
	}
	callInfo := struct {
		Filter models.SequenceExecutionFilter
	}{
		Filter: filter,
	}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc(filter)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//     len(mockedSequenceExecutionRepo.GetCalls())
func (mock *SequenceExecutionRepoMock) GetCalls() []struct {
	Filter models.SequenceExecutionFilter
} {
	var calls []struct {
		Filter models.SequenceExecutionFilter
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}

// GetByTriggeredID calls GetByTriggeredIDFunc.
func (mock *SequenceExecutionRepoMock) GetByTriggeredID(project string, triggeredID string) (*models.SequenceExecution, error) {
	if mock.GetByTriggeredIDFunc == nil {
		panic("SequenceExecutionRepoMock.GetByTriggeredIDFunc: method is nil but SequenceExecutionRepo.GetByTriggeredID was just called")
	}
	callInfo := struct {
		Project     string
		TriggeredID string
	}{
		Project:     project,
		TriggeredID: triggeredID,
	}
	mock.lockGetByTriggeredID.Lock()
	mock.calls.GetByTriggeredID = append(mock.calls.GetByTriggeredID, callInfo)
	mock.lockGetByTriggeredID.Unlock()
	return mock.GetByTriggeredIDFunc(project, triggeredID)
}

// GetByTriggeredIDCalls gets all the calls that were made to GetByTriggeredID.
// Check the length with:
//     len(mockedSequenceExecutionRepo.GetByTriggeredIDCalls())
func (mock *SequenceExecutionRepoMock) GetByTriggeredIDCalls() []struct {
	Project     string
	TriggeredID string
} {
	var calls []struct {
		Project     string
		TriggeredID string
	}
	mock.lockGetByTriggeredID.RLock()
	calls = mock.calls.GetByTriggeredID
	mock.lockGetByTriggeredID.RUnlock()
	return calls
}

// IsContextPaused calls IsContextPausedFunc.
func (mock *SequenceExecutionRepoMock) IsContextPaused(eventScope models.EventScope) bool {
	if mock.IsContextPausedFunc == nil {
		panic("SequenceExecutionRepoMock.IsContextPausedFunc: method is nil but SequenceExecutionRepo.IsContextPaused was just called")
	}
	callInfo := struct {
		EventScope models.EventScope
	}{
		EventScope: eventScope,
	}
	mock.lockIsContextPaused.Lock()
	mock.calls.IsContextPaused = append(mock.calls.IsContextPaused, callInfo)
	mock.lockIsContextPaused.Unlock()
	return mock.IsContextPausedFunc(eventScope)
}

// IsContextPausedCalls gets all the calls that were made to IsContextPaused.
// Check the length with:
//     len(mockedSequenceExecutionRepo.IsContextPausedCalls())
func (mock *SequenceExecutionRepoMock) IsContextPausedCalls() []struct {
	EventScope models.EventScope
} {
	var calls []struct {
		EventScope models.EventScope
	}
	mock.lockIsContextPaused.RLock()
	calls = mock.calls.IsContextPaused
	mock.lockIsContextPaused.RUnlock()
	return calls
}

// PauseContext calls PauseContextFunc.
func (mock *SequenceExecutionRepoMock) PauseContext(eventScope models.EventScope) error {
	if mock.PauseContextFunc == nil {
		panic("SequenceExecutionRepoMock.PauseContextFunc: method is nil but SequenceExecutionRepo.PauseContext was just called")
	}
	callInfo := struct {
		EventScope models.EventScope
	}{
		EventScope: eventScope,
	}
	mock.lockPauseContext.Lock()
	mock.calls.PauseContext = append(mock.calls.PauseContext, callInfo)
	mock.lockPauseContext.Unlock()
	return mock.PauseContextFunc(eventScope)
}

// PauseContextCalls gets all the calls that were made to PauseContext.
// Check the length with:
//     len(mockedSequenceExecutionRepo.PauseContextCalls())
func (mock *SequenceExecutionRepoMock) PauseContextCalls() []struct {
	EventScope models.EventScope
} {
	var calls []struct {
		EventScope models.EventScope
	}
	mock.lockPauseContext.RLock()
	calls = mock.calls.PauseContext
	mock.lockPauseContext.RUnlock()
	return calls
}

// ResumeContext calls ResumeContextFunc.
func (mock *SequenceExecutionRepoMock) ResumeContext(eventScope models.EventScope) error {
	if mock.ResumeContextFunc == nil {
		panic("SequenceExecutionRepoMock.ResumeContextFunc: method is nil but SequenceExecutionRepo.ResumeContext was just called")
	}
	callInfo := struct {
		EventScope models.EventScope
	}{
		EventScope: eventScope,
	}
	mock.lockResumeContext.Lock()
	mock.calls.ResumeContext = append(mock.calls.ResumeContext, callInfo)
	mock.lockResumeContext.Unlock()
	return mock.ResumeContextFunc(eventScope)
}

// ResumeContextCalls gets all the calls that were made to ResumeContext.
// Check the length with:
//     len(mockedSequenceExecutionRepo.ResumeContextCalls())
func (mock *SequenceExecutionRepoMock) ResumeContextCalls() []struct {
	EventScope models.EventScope
} {
	var calls []struct {
		EventScope models.EventScope
	}
	mock.lockResumeContext.RLock()
	calls = mock.calls.ResumeContext
	mock.lockResumeContext.RUnlock()
	return calls
}

// UpdateStatus calls UpdateStatusFunc.
func (mock *SequenceExecutionRepoMock) UpdateStatus(taskSequence models.SequenceExecution) (*models.SequenceExecution, error) {
	if mock.UpdateStatusFunc == nil {
		panic("SequenceExecutionRepoMock.UpdateStatusFunc: method is nil but SequenceExecutionRepo.UpdateStatus was just called")
	}
	callInfo := struct {
		TaskSequence models.SequenceExecution
	}{
		TaskSequence: taskSequence,
	}
	mock.lockUpdateStatus.Lock()
	mock.calls.UpdateStatus = append(mock.calls.UpdateStatus, callInfo)
	mock.lockUpdateStatus.Unlock()
	return mock.UpdateStatusFunc(taskSequence)
}

// UpdateStatusCalls gets all the calls that were made to UpdateStatus.
// Check the length with:
//     len(mockedSequenceExecutionRepo.UpdateStatusCalls())
func (mock *SequenceExecutionRepoMock) UpdateStatusCalls() []struct {
	TaskSequence models.SequenceExecution
} {
	var calls []struct {
		TaskSequence models.SequenceExecution
	}
	mock.lockUpdateStatus.RLock()
	calls = mock.calls.UpdateStatus
	mock.lockUpdateStatus.RUnlock()
	return calls
}

// Upsert calls UpsertFunc.
func (mock *SequenceExecutionRepoMock) Upsert(item models.SequenceExecution, options *models.SequenceExecutionUpsertOptions) error {
	if mock.UpsertFunc == nil {
		panic("SequenceExecutionRepoMock.UpsertFunc: method is nil but SequenceExecutionRepo.Upsert was just called")
	}
	callInfo := struct {
		Item    models.SequenceExecution
		Options *models.SequenceExecutionUpsertOptions
	}{
		Item:    item,
		Options: options,
	}
	mock.lockUpsert.Lock()
	mock.calls.Upsert = append(mock.calls.Upsert, callInfo)
	mock.lockUpsert.Unlock()
	return mock.UpsertFunc(item, options)
}

// UpsertCalls gets all the calls that were made to Upsert.
// Check the length with:
//     len(mockedSequenceExecutionRepo.UpsertCalls())
func (mock *SequenceExecutionRepoMock) UpsertCalls() []struct {
	Item    models.SequenceExecution
	Options *models.SequenceExecutionUpsertOptions
} {
	var calls []struct {
		Item    models.SequenceExecution
		Options *models.SequenceExecutionUpsertOptions
	}
	mock.lockUpsert.RLock()
	calls = mock.calls.Upsert
	mock.lockUpsert.RUnlock()
	return calls
}
