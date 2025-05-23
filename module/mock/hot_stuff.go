// Code generated by mockery v2.53.3. DO NOT EDIT.

package mock

import (
	irrecoverable "github.com/onflow/flow-go/module/irrecoverable"
	mock "github.com/stretchr/testify/mock"

	model "github.com/onflow/flow-go/consensus/hotstuff/model"
)

// HotStuff is an autogenerated mock type for the HotStuff type
type HotStuff struct {
	mock.Mock
}

// Done provides a mock function with no fields
func (_m *HotStuff) Done() <-chan struct{} {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Done")
	}

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// Ready provides a mock function with no fields
func (_m *HotStuff) Ready() <-chan struct{} {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Ready")
	}

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// Start provides a mock function with given fields: _a0
func (_m *HotStuff) Start(_a0 irrecoverable.SignalerContext) {
	_m.Called(_a0)
}

// SubmitProposal provides a mock function with given fields: proposal
func (_m *HotStuff) SubmitProposal(proposal *model.SignedProposal) {
	_m.Called(proposal)
}

// NewHotStuff creates a new instance of HotStuff. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHotStuff(t interface {
	mock.TestingT
	Cleanup(func())
}) *HotStuff {
	mock := &HotStuff{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
