// Code generated by mockery v2.53.3. DO NOT EDIT.

package mock

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"
)

// Headers is an autogenerated mock type for the Headers type
type Headers struct {
	mock.Mock
}

// BlockIDByHeight provides a mock function with given fields: height
func (_m *Headers) BlockIDByHeight(height uint64) (flow.Identifier, error) {
	ret := _m.Called(height)

	if len(ret) == 0 {
		panic("no return value specified for BlockIDByHeight")
	}

	var r0 flow.Identifier
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64) (flow.Identifier, error)); ok {
		return rf(height)
	}
	if rf, ok := ret.Get(0).(func(uint64) flow.Identifier); ok {
		r0 = rf(height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flow.Identifier)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ByBlockID provides a mock function with given fields: blockID
func (_m *Headers) ByBlockID(blockID flow.Identifier) (*flow.Header, error) {
	ret := _m.Called(blockID)

	if len(ret) == 0 {
		panic("no return value specified for ByBlockID")
	}

	var r0 *flow.Header
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) (*flow.Header, error)); ok {
		return rf(blockID)
	}
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.Header); ok {
		r0 = rf(blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Header)
		}
	}

	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(blockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ByHeight provides a mock function with given fields: height
func (_m *Headers) ByHeight(height uint64) (*flow.Header, error) {
	ret := _m.Called(height)

	if len(ret) == 0 {
		panic("no return value specified for ByHeight")
	}

	var r0 *flow.Header
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64) (*flow.Header, error)); ok {
		return rf(height)
	}
	if rf, ok := ret.Get(0).(func(uint64) *flow.Header); ok {
		r0 = rf(height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Header)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ByParentID provides a mock function with given fields: parentID
func (_m *Headers) ByParentID(parentID flow.Identifier) ([]*flow.Header, error) {
	ret := _m.Called(parentID)

	if len(ret) == 0 {
		panic("no return value specified for ByParentID")
	}

	var r0 []*flow.Header
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) ([]*flow.Header, error)); ok {
		return rf(parentID)
	}
	if rf, ok := ret.Get(0).(func(flow.Identifier) []*flow.Header); ok {
		r0 = rf(parentID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*flow.Header)
		}
	}

	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(parentID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Exists provides a mock function with given fields: blockID
func (_m *Headers) Exists(blockID flow.Identifier) (bool, error) {
	ret := _m.Called(blockID)

	if len(ret) == 0 {
		panic("no return value specified for Exists")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) (bool, error)); ok {
		return rf(blockID)
	}
	if rf, ok := ret.Get(0).(func(flow.Identifier) bool); ok {
		r0 = rf(blockID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(blockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: header
func (_m *Headers) Store(header *flow.Header) error {
	ret := _m.Called(header)

	if len(ret) == 0 {
		panic("no return value specified for Store")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.Header) error); ok {
		r0 = rf(header)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewHeaders creates a new instance of Headers. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHeaders(t interface {
	mock.TestingT
	Cleanup(func())
}) *Headers {
	mock := &Headers{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
