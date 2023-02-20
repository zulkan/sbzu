// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	domain "gozu/domain"

	mock "github.com/stretchr/testify/mock"
)

// StockUseCase is an autogenerated mock type for the StockUseCase type
type StockUseCase struct {
	mock.Mock
}

// GetStockSummary provides a mock function with given fields: stockCode
func (_m *StockUseCase) GetStockSummary(stockCode string) (*domain.StockSummary, error) {
	ret := _m.Called(stockCode)

	var r0 *domain.StockSummary
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.StockSummary, error)); ok {
		return rf(stockCode)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.StockSummary); ok {
		r0 = rf(stockCode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.StockSummary)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(stockCode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProcessFileData provides a mock function with given fields: rawData
func (_m *StockUseCase) ProcessFileData(rawData string) error {
	ret := _m.Called(rawData)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(rawData)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ProcessQueueMessage provides a mock function with given fields: rawData
func (_m *StockUseCase) ProcessQueueMessage(rawData string) error {
	ret := _m.Called(rawData)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(rawData)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteStockSummary provides a mock function with given fields: record
func (_m *StockUseCase) WriteStockSummary(record *domain.StockRecord) error {
	ret := _m.Called(record)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.StockRecord) error); ok {
		r0 = rf(record)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewStockUseCase interface {
	mock.TestingT
	Cleanup(func())
}

// NewStockUseCase creates a new instance of StockUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStockUseCase(t mockConstructorTestingTNewStockUseCase) *StockUseCase {
	mock := &StockUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
