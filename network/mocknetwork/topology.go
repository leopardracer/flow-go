// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocknetwork

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"
)

// Topology is an autogenerated mock type for the Topology type
type Topology struct {
	mock.Mock
}

// Fanout provides a mock function with given fields: ids
func (_m *Topology) Fanout(ids flow.GenericIdentityList[flow.Identity]) flow.GenericIdentityList[flow.Identity] {
	ret := _m.Called(ids)

	if len(ret) == 0 {
		panic("no return value specified for Fanout")
	}

	var r0 flow.GenericIdentityList[flow.Identity]
	if rf, ok := ret.Get(0).(func(flow.GenericIdentityList[flow.Identity]) flow.GenericIdentityList[flow.Identity]); ok {
		r0 = rf(ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flow.GenericIdentityList[flow.Identity])
		}
	}

	return r0
}

// NewTopology creates a new instance of Topology. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTopology(t interface {
	mock.TestingT
	Cleanup(func())
}) *Topology {
	mock := &Topology{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
