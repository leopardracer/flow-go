// Code generated by mockery v2.53.3. DO NOT EDIT.

package mock

import (
	ledger "github.com/onflow/flow-go/ledger"
	mock "github.com/stretchr/testify/mock"
)

// Reporter is an autogenerated mock type for the Reporter type
type Reporter struct {
	mock.Mock
}

// Name provides a mock function with no fields
func (_m *Reporter) Name() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Report provides a mock function with given fields: payloads, statecommitment
func (_m *Reporter) Report(payloads []ledger.Payload, statecommitment ledger.State) error {
	ret := _m.Called(payloads, statecommitment)

	if len(ret) == 0 {
		panic("no return value specified for Report")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]ledger.Payload, ledger.State) error); ok {
		r0 = rf(payloads, statecommitment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewReporter creates a new instance of Reporter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReporter(t interface {
	mock.TestingT
	Cleanup(func())
}) *Reporter {
	mock := &Reporter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
