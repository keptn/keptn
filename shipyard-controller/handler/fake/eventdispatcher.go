// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package fake

import (
	"context"
	"github.com/keptn/keptn/shipyard-controller/models"
	"sync"
)

// IEventDispatcherMock is a mock implementation of handler.IEventDispatcher.
//
// 	func TestSomethingThatUsesIEventDispatcher(t *testing.T) {
//
// 		// make and configure a mocked handler.IEventDispatcher
// 		mockedIEventDispatcher := &IEventDispatcherMock{
// 			AddFunc: func(event models.DispatcherEvent, skipQueue bool) error {
// 				panic("mock out the Add method")
// 			},
// 			RunFunc: func(ctx context.Context)  {
// 				panic("mock out the Run method")
// 			},
// 			StopFunc: func()  {
// 				panic("mock out the Stop method")
// 			},
// 		}
//
// 		// use mockedIEventDispatcher in code that requires handler.IEventDispatcher
// 		// and then make assertions.
//
// 	}
type IEventDispatcherMock struct {
	// AddFunc mocks the Add method.
	AddFunc func(event models.DispatcherEvent, skipQueue bool) error

	// RunFunc mocks the Run method.
	RunFunc func(ctx context.Context)

	// StopFunc mocks the Stop method.
	StopFunc func()

	// calls tracks calls to the methods.
	calls struct {
		// Add holds details about calls to the Add method.
		Add []struct {
			//models.KeptnContextExtendedCEis the event argument value.
			Event models.DispatcherEvent
			// SkipQueue is the skipQueue argument value.
			SkipQueue bool
		}
		// Run holds details about calls to the Run method.
		Run []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// Stop holds details about calls to the Stop method.
		Stop []struct {
		}
	}
	lockAdd  sync.RWMutex
	lockRun  sync.RWMutex
	lockStop sync.RWMutex
}

// Add calls AddFunc.
func (mock *IEventDispatcherMock) Add(event models.DispatcherEvent, skipQueue bool) error {
	if mock.AddFunc == nil {
		panic("IEventDispatcherMock.AddFunc: method is nil but IEventDispatcher.Add was just called")
	}
	callInfo := struct {
		Event     models.DispatcherEvent
		SkipQueue bool
	}{
		Event:     event,
		SkipQueue: skipQueue,
	}
	mock.lockAdd.Lock()
	mock.calls.Add = append(mock.calls.Add, callInfo)
	mock.lockAdd.Unlock()
	return mock.AddFunc(event, skipQueue)
}

// AddCalls gets all the calls that were made to Add.
// Check the length with:
//     len(mockedIEventDispatcher.AddCalls())
func (mock *IEventDispatcherMock) AddCalls() []struct {
	Event     models.DispatcherEvent
	SkipQueue bool
} {
	var calls []struct {
		Event     models.DispatcherEvent
		SkipQueue bool
	}
	mock.lockAdd.RLock()
	calls = mock.calls.Add
	mock.lockAdd.RUnlock()
	return calls
}

// Run calls RunFunc.
func (mock *IEventDispatcherMock) Run(ctx context.Context) {
	if mock.RunFunc == nil {
		panic("IEventDispatcherMock.RunFunc: method is nil but IEventDispatcher.Run was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockRun.Lock()
	mock.calls.Run = append(mock.calls.Run, callInfo)
	mock.lockRun.Unlock()
	mock.RunFunc(ctx)
}

// RunCalls gets all the calls that were made to Run.
// Check the length with:
//     len(mockedIEventDispatcher.RunCalls())
func (mock *IEventDispatcherMock) RunCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockRun.RLock()
	calls = mock.calls.Run
	mock.lockRun.RUnlock()
	return calls
}

// Stop calls StopFunc.
func (mock *IEventDispatcherMock) Stop() {
	if mock.StopFunc == nil {
		panic("IEventDispatcherMock.StopFunc: method is nil but IEventDispatcher.Stop was just called")
	}
	callInfo := struct {
	}{}
	mock.lockStop.Lock()
	mock.calls.Stop = append(mock.calls.Stop, callInfo)
	mock.lockStop.Unlock()
	mock.StopFunc()
}

// StopCalls gets all the calls that were made to Stop.
// Check the length with:
//     len(mockedIEventDispatcher.StopCalls())
func (mock *IEventDispatcherMock) StopCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockStop.RLock()
	calls = mock.calls.Stop
	mock.lockStop.RUnlock()
	return calls
}
