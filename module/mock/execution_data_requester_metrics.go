// Code generated by mockery v2.53.3. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// ExecutionDataRequesterMetrics is an autogenerated mock type for the ExecutionDataRequesterMetrics type
type ExecutionDataRequesterMetrics struct {
	mock.Mock
}

// ExecutionDataFetchFinished provides a mock function with given fields: duration, success, height
func (_m *ExecutionDataRequesterMetrics) ExecutionDataFetchFinished(duration time.Duration, success bool, height uint64) {
	_m.Called(duration, success, height)
}

// ExecutionDataFetchStarted provides a mock function with no fields
func (_m *ExecutionDataRequesterMetrics) ExecutionDataFetchStarted() {
	_m.Called()
}

// FetchRetried provides a mock function with no fields
func (_m *ExecutionDataRequesterMetrics) FetchRetried() {
	_m.Called()
}

// NotificationSent provides a mock function with given fields: height
func (_m *ExecutionDataRequesterMetrics) NotificationSent(height uint64) {
	_m.Called(height)
}

// NewExecutionDataRequesterMetrics creates a new instance of ExecutionDataRequesterMetrics. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExecutionDataRequesterMetrics(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExecutionDataRequesterMetrics {
	mock := &ExecutionDataRequesterMetrics{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
