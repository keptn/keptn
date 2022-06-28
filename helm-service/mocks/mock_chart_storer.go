// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/keptn/keptn/helm-service/pkg/types (interfaces: IChartStorer)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	"github.com/keptn/keptn/helm-service/pkg/common"
)

// MockIChartStorer is a mock of IChartStorer interface.
type MockIChartStorer struct {
	ctrl     *gomock.Controller
	recorder *MockIChartStorerMockRecorder
}

// MockIChartStorerMockRecorder is the mock recorder for MockIChartStorer.
type MockIChartStorerMockRecorder struct {
	mock *MockIChartStorer
}

// NewMockIChartStorer creates a new mock instance.
func NewMockIChartStorer(ctrl *gomock.Controller) *MockIChartStorer {
	mock := &MockIChartStorer{ctrl: ctrl}
	mock.recorder = &MockIChartStorerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIChartStorer) EXPECT() *MockIChartStorerMockRecorder {
	return m.recorder
}

// Store mocks base method.
func (m *MockIChartStorer) Store(arg0 common.StoreChartOptions) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Store indicates an expected call of Store.
func (mr *MockIChartStorerMockRecorder) Store(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockIChartStorer)(nil).Store), arg0)
}
