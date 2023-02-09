// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"

	storage "github.com/onflow/flow-go/storage"
)

// Headers is an autogenerated mock type for the Headers type
type Headers struct {
	mock.Mock
}

// BatchIndexByChunkID provides a mock function with given fields: headerID, chunkID, batch
func (_m *Headers) BatchIndexByChunkID(headerID flow.Identifier, chunkID flow.Identifier, batch storage.BatchStorage) error {
	ret := _m.Called(headerID, chunkID, batch)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier, flow.Identifier, storage.BatchStorage) error); ok {
		r0 = rf(headerID, chunkID, batch)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BatchRemoveChunkBlockIndexByChunkID provides a mock function with given fields: chunkID, batch
func (_m *Headers) BatchRemoveChunkBlockIndexByChunkID(chunkID flow.Identifier, batch storage.BatchStorage) error {
	ret := _m.Called(chunkID, batch)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier, storage.BatchStorage) error); ok {
		r0 = rf(chunkID, batch)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BlockIDByHeight provides a mock function with given fields: height
func (_m *Headers) BlockIDByHeight(height uint64) (flow.Identifier, error) {
	ret := _m.Called(height)

	var r0 flow.Identifier
	if rf, ok := ret.Get(0).(func(uint64) flow.Identifier); ok {
		r0 = rf(height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flow.Identifier)
		}
	}

	var r1 error
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

	var r0 *flow.Header
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.Header); ok {
		r0 = rf(blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Header)
		}
	}

	var r1 error
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

	var r0 *flow.Header
	if rf, ok := ret.Get(0).(func(uint64) *flow.Header); ok {
		r0 = rf(height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Header)
		}
	}

	var r1 error
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

	var r0 []*flow.Header
	if rf, ok := ret.Get(0).(func(flow.Identifier) []*flow.Header); ok {
		r0 = rf(parentID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*flow.Header)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(parentID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IDByChunkID provides a mock function with given fields: chunkID
func (_m *Headers) IDByChunkID(chunkID flow.Identifier) (flow.Identifier, error) {
	ret := _m.Called(chunkID)

	var r0 flow.Identifier
	if rf, ok := ret.Get(0).(func(flow.Identifier) flow.Identifier); ok {
		r0 = rf(chunkID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flow.Identifier)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(chunkID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IndexByChunkID provides a mock function with given fields: headerID, chunkID
func (_m *Headers) IndexByChunkID(headerID flow.Identifier, chunkID flow.Identifier) error {
	ret := _m.Called(headerID, chunkID)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier, flow.Identifier) error); ok {
		r0 = rf(headerID, chunkID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Store provides a mock function with given fields: header
func (_m *Headers) Store(header *flow.Header) error {
	ret := _m.Called(header)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.Header) error); ok {
		r0 = rf(header)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewHeaders interface {
	mock.TestingT
	Cleanup(func())
}

// NewHeaders creates a new instance of Headers. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewHeaders(t mockConstructorTestingTNewHeaders) *Headers {
	mock := &Headers{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
