// Code generated by mockery v2.53.3. DO NOT EDIT.

package mockp2p

import (
	control "github.com/libp2p/go-libp2p/core/control"
	mock "github.com/stretchr/testify/mock"

	multiaddr "github.com/multiformats/go-multiaddr"

	network "github.com/libp2p/go-libp2p/core/network"

	p2p "github.com/onflow/flow-go/network/p2p"

	peer "github.com/libp2p/go-libp2p/core/peer"
)

// ConnectionGater is an autogenerated mock type for the ConnectionGater type
type ConnectionGater struct {
	mock.Mock
}

// InterceptAccept provides a mock function with given fields: _a0
func (_m *ConnectionGater) InterceptAccept(_a0 network.ConnMultiaddrs) bool {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for InterceptAccept")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(network.ConnMultiaddrs) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// InterceptAddrDial provides a mock function with given fields: _a0, _a1
func (_m *ConnectionGater) InterceptAddrDial(_a0 peer.ID, _a1 multiaddr.Multiaddr) bool {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for InterceptAddrDial")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(peer.ID, multiaddr.Multiaddr) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// InterceptPeerDial provides a mock function with given fields: p
func (_m *ConnectionGater) InterceptPeerDial(p peer.ID) bool {
	ret := _m.Called(p)

	if len(ret) == 0 {
		panic("no return value specified for InterceptPeerDial")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(peer.ID) bool); ok {
		r0 = rf(p)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// InterceptSecured provides a mock function with given fields: _a0, _a1, _a2
func (_m *ConnectionGater) InterceptSecured(_a0 network.Direction, _a1 peer.ID, _a2 network.ConnMultiaddrs) bool {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for InterceptSecured")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(network.Direction, peer.ID, network.ConnMultiaddrs) bool); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// InterceptUpgraded provides a mock function with given fields: _a0
func (_m *ConnectionGater) InterceptUpgraded(_a0 network.Conn) (bool, control.DisconnectReason) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for InterceptUpgraded")
	}

	var r0 bool
	var r1 control.DisconnectReason
	if rf, ok := ret.Get(0).(func(network.Conn) (bool, control.DisconnectReason)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(network.Conn) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(network.Conn) control.DisconnectReason); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Get(1).(control.DisconnectReason)
	}

	return r0, r1
}

// SetDisallowListOracle provides a mock function with given fields: oracle
func (_m *ConnectionGater) SetDisallowListOracle(oracle p2p.DisallowListOracle) {
	_m.Called(oracle)
}

// NewConnectionGater creates a new instance of ConnectionGater. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConnectionGater(t interface {
	mock.TestingT
	Cleanup(func())
}) *ConnectionGater {
	mock := &ConnectionGater{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
