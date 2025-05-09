// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"

	model "github.com/onflow/flow-go/consensus/hotstuff/model"
)

// TimeoutAggregationConsumer is an autogenerated mock type for the TimeoutAggregationConsumer type
type TimeoutAggregationConsumer struct {
	mock.Mock
}

// OnDoubleTimeoutDetected provides a mock function with given fields: _a0, _a1
func (_m *TimeoutAggregationConsumer) OnDoubleTimeoutDetected(_a0 *model.TimeoutObject, _a1 *model.TimeoutObject) {
	_m.Called(_a0, _a1)
}

// OnInvalidTimeoutDetected provides a mock function with given fields: err
func (_m *TimeoutAggregationConsumer) OnInvalidTimeoutDetected(err model.InvalidTimeoutError) {
	_m.Called(err)
}

// OnNewQcDiscovered provides a mock function with given fields: certificate
func (_m *TimeoutAggregationConsumer) OnNewQcDiscovered(certificate *flow.QuorumCertificate) {
	_m.Called(certificate)
}

// OnNewTcDiscovered provides a mock function with given fields: certificate
func (_m *TimeoutAggregationConsumer) OnNewTcDiscovered(certificate *flow.TimeoutCertificate) {
	_m.Called(certificate)
}

// OnPartialTcCreated provides a mock function with given fields: view, newestQC, lastViewTC
func (_m *TimeoutAggregationConsumer) OnPartialTcCreated(view uint64, newestQC *flow.QuorumCertificate, lastViewTC *flow.TimeoutCertificate) {
	_m.Called(view, newestQC, lastViewTC)
}

// OnTcConstructedFromTimeouts provides a mock function with given fields: certificate
func (_m *TimeoutAggregationConsumer) OnTcConstructedFromTimeouts(certificate *flow.TimeoutCertificate) {
	_m.Called(certificate)
}

// OnTimeoutProcessed provides a mock function with given fields: timeout
func (_m *TimeoutAggregationConsumer) OnTimeoutProcessed(timeout *model.TimeoutObject) {
	_m.Called(timeout)
}

// NewTimeoutAggregationConsumer creates a new instance of TimeoutAggregationConsumer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTimeoutAggregationConsumer(t interface {
	mock.TestingT
	Cleanup(func())
}) *TimeoutAggregationConsumer {
	mock := &TimeoutAggregationConsumer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
