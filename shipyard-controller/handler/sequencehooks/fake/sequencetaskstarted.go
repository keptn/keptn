// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package fake

import (
	"github.com/keptn/keptn/shipyard-controller/models"
	"sync"
)

// ISequenceTaskStartedHookMock is a mock implementation of sequencehooks.ISequenceTaskStartedHook.
//
// 	func TestSomethingThatUsesISequenceTaskStartedHook(t *testing.T) {
//
// 		// make and configure a mocked sequencehooks.ISequenceTaskStartedHook
// 		mockedISequenceTaskStartedHook := &ISequenceTaskStartedHookMock{
// 			OnSequenceTaskStartedFunc: func(event models.Event)  {
// 				panic("mock out the OnSequenceTaskStarted method")
// 			},
// 		}
//
// 		// use mockedISequenceTaskStartedHook in code that requires sequencehooks.ISequenceTaskStartedHook
// 		// and then make assertions.
//
// 	}
type ISequenceTaskStartedHookMock struct {
	// OnSequenceTaskStartedFunc mocks the OnSequenceTaskStarted method.
	OnSequenceTaskStartedFunc func(event models.Event)

	// calls tracks calls to the methods.
	calls struct {
		// OnSequenceTaskStarted holds details about calls to the OnSequenceTaskStarted method.
		OnSequenceTaskStarted []struct {
			// Event is the event argument value.
			Event models.Event
		}
	}
	lockOnSequenceTaskStarted sync.RWMutex
}

// OnSequenceTaskStarted calls OnSequenceTaskStartedFunc.
func (mock *ISequenceTaskStartedHookMock) OnSequenceTaskStarted(event models.Event) {
	if mock.OnSequenceTaskStartedFunc == nil {
		panic("ISequenceTaskStartedHookMock.OnSequenceTaskStartedFunc: method is nil but ISequenceTaskStartedHook.OnSequenceTaskStarted was just called")
	}
	callInfo := struct {
		Event models.Event
	}{
		Event: event,
	}
	mock.lockOnSequenceTaskStarted.Lock()
	mock.calls.OnSequenceTaskStarted = append(mock.calls.OnSequenceTaskStarted, callInfo)
	mock.lockOnSequenceTaskStarted.Unlock()
	mock.OnSequenceTaskStartedFunc(event)
}

// OnSequenceTaskStartedCalls gets all the calls that were made to OnSequenceTaskStarted.
// Check the length with:
//     len(mockedISequenceTaskStartedHook.OnSequenceTaskStartedCalls())
func (mock *ISequenceTaskStartedHookMock) OnSequenceTaskStartedCalls() []struct {
	Event models.Event
} {
	var calls []struct {
		Event models.Event
	}
	mock.lockOnSequenceTaskStarted.RLock()
	calls = mock.calls.OnSequenceTaskStarted
	mock.lockOnSequenceTaskStarted.RUnlock()
	return calls
}
