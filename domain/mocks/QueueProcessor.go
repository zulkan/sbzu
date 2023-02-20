// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// QueueProcessor is an autogenerated mock type for the QueueProcessor type
type QueueProcessor struct {
	mock.Mock
}

// ProcessQueueMessage provides a mock function with given fields: rawData
func (_m *QueueProcessor) ProcessQueueMessage(rawData string) error {
	ret := _m.Called(rawData)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(rawData)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewQueueProcessor interface {
	mock.TestingT
	Cleanup(func())
}

// NewQueueProcessor creates a new instance of QueueProcessor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewQueueProcessor(t mockConstructorTestingTNewQueueProcessor) *QueueProcessor {
	mock := &QueueProcessor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
