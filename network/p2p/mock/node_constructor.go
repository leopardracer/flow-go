// Code generated by mockery v2.53.3. DO NOT EDIT.

package mockp2p

import (
	p2p "github.com/onflow/flow-go/network/p2p"
	mock "github.com/stretchr/testify/mock"
)

// NodeConstructor is an autogenerated mock type for the NodeConstructor type
type NodeConstructor struct {
	mock.Mock
}

// Execute provides a mock function with given fields: config
func (_m *NodeConstructor) Execute(config *p2p.NodeConfig) (p2p.LibP2PNode, error) {
	ret := _m.Called(config)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 p2p.LibP2PNode
	var r1 error
	if rf, ok := ret.Get(0).(func(*p2p.NodeConfig) (p2p.LibP2PNode, error)); ok {
		return rf(config)
	}
	if rf, ok := ret.Get(0).(func(*p2p.NodeConfig) p2p.LibP2PNode); ok {
		r0 = rf(config)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.LibP2PNode)
		}
	}

	if rf, ok := ret.Get(1).(func(*p2p.NodeConfig) error); ok {
		r1 = rf(config)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewNodeConstructor creates a new instance of NodeConstructor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNodeConstructor(t interface {
	mock.TestingT
	Cleanup(func())
}) *NodeConstructor {
	mock := &NodeConstructor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
