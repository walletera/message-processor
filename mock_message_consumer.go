// Code generated by mockery v2.20.0. DO NOT EDIT.

package main

import mock "github.com/stretchr/testify/mock"

// MockMessageConsumer is an autogenerated mock type for the MessageConsumer type
type MockMessageConsumer struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *MockMessageConsumer) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Consume provides a mock function with given fields:
func (_m *MockMessageConsumer) Consume() (<-chan Message, error) {
	ret := _m.Called()

	var r0 <-chan Message
	var r1 error
	if rf, ok := ret.Get(0).(func() (<-chan Message, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() <-chan Message); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan Message)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockMessageConsumer interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockMessageConsumer creates a new instance of MockMessageConsumer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockMessageConsumer(t mockConstructorTestingTNewMockMessageConsumer) *MockMessageConsumer {
	mock := &MockMessageConsumer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
