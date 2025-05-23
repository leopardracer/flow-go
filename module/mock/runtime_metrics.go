// Code generated by mockery v2.53.3. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// RuntimeMetrics is an autogenerated mock type for the RuntimeMetrics type
type RuntimeMetrics struct {
	mock.Mock
}

// RuntimeSetNumberOfAccounts provides a mock function with given fields: count
func (_m *RuntimeMetrics) RuntimeSetNumberOfAccounts(count uint64) {
	_m.Called(count)
}

// RuntimeTransactionChecked provides a mock function with given fields: dur
func (_m *RuntimeMetrics) RuntimeTransactionChecked(dur time.Duration) {
	_m.Called(dur)
}

// RuntimeTransactionInterpreted provides a mock function with given fields: dur
func (_m *RuntimeMetrics) RuntimeTransactionInterpreted(dur time.Duration) {
	_m.Called(dur)
}

// RuntimeTransactionParsed provides a mock function with given fields: dur
func (_m *RuntimeMetrics) RuntimeTransactionParsed(dur time.Duration) {
	_m.Called(dur)
}

// RuntimeTransactionProgramsCacheHit provides a mock function with no fields
func (_m *RuntimeMetrics) RuntimeTransactionProgramsCacheHit() {
	_m.Called()
}

// RuntimeTransactionProgramsCacheMiss provides a mock function with no fields
func (_m *RuntimeMetrics) RuntimeTransactionProgramsCacheMiss() {
	_m.Called()
}

// NewRuntimeMetrics creates a new instance of RuntimeMetrics. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRuntimeMetrics(t interface {
	mock.TestingT
	Cleanup(func())
}) *RuntimeMetrics {
	mock := &RuntimeMetrics{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
