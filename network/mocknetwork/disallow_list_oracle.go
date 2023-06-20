// Code generated by mockery v2.21.4. DO NOT EDIT.

package mocknetwork

import (
	network "github.com/onflow/flow-go/network"
	mock "github.com/stretchr/testify/mock"

	peer "github.com/libp2p/go-libp2p/core/peer"
)

// DisallowListOracle is an autogenerated mock type for the DisallowListOracle type
type DisallowListOracle struct {
	mock.Mock
}

// GetAllDisallowedListCausesFor provides a mock function with given fields: _a0
func (_m *DisallowListOracle) GetAllDisallowedListCausesFor(_a0 peer.ID) []network.DisallowListedCause {
	ret := _m.Called(_a0)

	var r0 []network.DisallowListedCause
	if rf, ok := ret.Get(0).(func(peer.ID) []network.DisallowListedCause); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]network.DisallowListedCause)
		}
	}

	return r0
}

type mockConstructorTestingTNewDisallowListOracle interface {
	mock.TestingT
	Cleanup(func())
}

// NewDisallowListOracle creates a new instance of DisallowListOracle. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDisallowListOracle(t mockConstructorTestingTNewDisallowListOracle) *DisallowListOracle {
	mock := &DisallowListOracle{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
