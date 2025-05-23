// Code generated by mockery v2.53.3. DO NOT EDIT.

package mock

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"

	protocol "github.com/onflow/flow-go/state/protocol"
)

// ProtocolState is an autogenerated mock type for the ProtocolState type
type ProtocolState struct {
	mock.Mock
}

// EpochStateAtBlockID provides a mock function with given fields: blockID
func (_m *ProtocolState) EpochStateAtBlockID(blockID flow.Identifier) (protocol.EpochProtocolState, error) {
	ret := _m.Called(blockID)

	if len(ret) == 0 {
		panic("no return value specified for EpochStateAtBlockID")
	}

	var r0 protocol.EpochProtocolState
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) (protocol.EpochProtocolState, error)); ok {
		return rf(blockID)
	}
	if rf, ok := ret.Get(0).(func(flow.Identifier) protocol.EpochProtocolState); ok {
		r0 = rf(blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(protocol.EpochProtocolState)
		}
	}

	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(blockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GlobalParams provides a mock function with no fields
func (_m *ProtocolState) GlobalParams() protocol.GlobalParams {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GlobalParams")
	}

	var r0 protocol.GlobalParams
	if rf, ok := ret.Get(0).(func() protocol.GlobalParams); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(protocol.GlobalParams)
		}
	}

	return r0
}

// KVStoreAtBlockID provides a mock function with given fields: blockID
func (_m *ProtocolState) KVStoreAtBlockID(blockID flow.Identifier) (protocol.KVStoreReader, error) {
	ret := _m.Called(blockID)

	if len(ret) == 0 {
		panic("no return value specified for KVStoreAtBlockID")
	}

	var r0 protocol.KVStoreReader
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) (protocol.KVStoreReader, error)); ok {
		return rf(blockID)
	}
	if rf, ok := ret.Get(0).(func(flow.Identifier) protocol.KVStoreReader); ok {
		r0 = rf(blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(protocol.KVStoreReader)
		}
	}

	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(blockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProtocolState creates a new instance of ProtocolState. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProtocolState(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProtocolState {
	mock := &ProtocolState{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
