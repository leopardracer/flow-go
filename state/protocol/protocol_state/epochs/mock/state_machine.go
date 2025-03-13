// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"
)

// StateMachine is an autogenerated mock type for the StateMachine type
type StateMachine struct {
	mock.Mock
}

// Build provides a mock function with given fields:
func (_m *StateMachine) Build() (*flow.EpochStateEntry, flow.Identifier, bool) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Build")
	}

	var r0 *flow.EpochStateEntry
	var r1 flow.Identifier
	var r2 bool
	if rf, ok := ret.Get(0).(func() (*flow.EpochStateEntry, flow.Identifier, bool)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *flow.EpochStateEntry); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.EpochStateEntry)
		}
	}

	if rf, ok := ret.Get(1).(func() flow.Identifier); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(flow.Identifier)
		}
	}

	if rf, ok := ret.Get(2).(func() bool); ok {
		r2 = rf()
	} else {
		r2 = ret.Get(2).(bool)
	}

	return r0, r1, r2
}

// EjectIdentity provides a mock function with given fields: ejectionEvent
func (_m *StateMachine) EjectIdentity(ejectionEvent *flow.EjectNode) bool {
	ret := _m.Called(ejectionEvent)

	if len(ret) == 0 {
		panic("no return value specified for EjectIdentity")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(*flow.EjectNode) bool); ok {
		r0 = rf(ejectionEvent)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ParentState provides a mock function with given fields:
func (_m *StateMachine) ParentState() *flow.RichEpochStateEntry {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ParentState")
	}

	var r0 *flow.RichEpochStateEntry
	if rf, ok := ret.Get(0).(func() *flow.RichEpochStateEntry); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.RichEpochStateEntry)
		}
	}

	return r0
}

// ProcessEpochCommit provides a mock function with given fields: epochCommit
func (_m *StateMachine) ProcessEpochCommit(epochCommit *flow.EpochCommit) (bool, error) {
	ret := _m.Called(epochCommit)

	if len(ret) == 0 {
		panic("no return value specified for ProcessEpochCommit")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(*flow.EpochCommit) (bool, error)); ok {
		return rf(epochCommit)
	}
	if rf, ok := ret.Get(0).(func(*flow.EpochCommit) bool); ok {
		r0 = rf(epochCommit)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*flow.EpochCommit) error); ok {
		r1 = rf(epochCommit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProcessEpochRecover provides a mock function with given fields: epochRecover
func (_m *StateMachine) ProcessEpochRecover(epochRecover *flow.EpochRecover) (bool, error) {
	ret := _m.Called(epochRecover)

	if len(ret) == 0 {
		panic("no return value specified for ProcessEpochRecover")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(*flow.EpochRecover) (bool, error)); ok {
		return rf(epochRecover)
	}
	if rf, ok := ret.Get(0).(func(*flow.EpochRecover) bool); ok {
		r0 = rf(epochRecover)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*flow.EpochRecover) error); ok {
		r1 = rf(epochRecover)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProcessEpochSetup provides a mock function with given fields: epochSetup
func (_m *StateMachine) ProcessEpochSetup(epochSetup *flow.EpochSetup) (bool, error) {
	ret := _m.Called(epochSetup)

	if len(ret) == 0 {
		panic("no return value specified for ProcessEpochSetup")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(*flow.EpochSetup) (bool, error)); ok {
		return rf(epochSetup)
	}
	if rf, ok := ret.Get(0).(func(*flow.EpochSetup) bool); ok {
		r0 = rf(epochSetup)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*flow.EpochSetup) error); ok {
		r1 = rf(epochSetup)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransitionToNextEpoch provides a mock function with given fields:
func (_m *StateMachine) TransitionToNextEpoch() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for TransitionToNextEpoch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// View provides a mock function with given fields:
func (_m *StateMachine) View() uint64 {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for View")
	}

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// NewStateMachine creates a new instance of StateMachine. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStateMachine(t interface {
	mock.TestingT
	Cleanup(func())
}) *StateMachine {
	mock := &StateMachine{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
