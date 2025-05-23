// Code generated by mockery v2.53.3. DO NOT EDIT.

package mempool

import (
	flow "github.com/onflow/flow-go/model/flow"
	execution_data "github.com/onflow/flow-go/module/executiondatasync/execution_data"

	mock "github.com/stretchr/testify/mock"
)

// ExecutionData is an autogenerated mock type for the ExecutionData type
type ExecutionData struct {
	mock.Mock
}

// Add provides a mock function with given fields: _a0
func (_m *ExecutionData) Add(_a0 *execution_data.BlockExecutionDataEntity) bool {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(*execution_data.BlockExecutionDataEntity) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// All provides a mock function with no fields
func (_m *ExecutionData) All() []*execution_data.BlockExecutionDataEntity {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for All")
	}

	var r0 []*execution_data.BlockExecutionDataEntity
	if rf, ok := ret.Get(0).(func() []*execution_data.BlockExecutionDataEntity); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*execution_data.BlockExecutionDataEntity)
		}
	}

	return r0
}

// ByID provides a mock function with given fields: _a0
func (_m *ExecutionData) ByID(_a0 flow.Identifier) (*execution_data.BlockExecutionDataEntity, bool) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for ByID")
	}

	var r0 *execution_data.BlockExecutionDataEntity
	var r1 bool
	if rf, ok := ret.Get(0).(func(flow.Identifier) (*execution_data.BlockExecutionDataEntity, bool)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(flow.Identifier) *execution_data.BlockExecutionDataEntity); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*execution_data.BlockExecutionDataEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(flow.Identifier) bool); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Clear provides a mock function with no fields
func (_m *ExecutionData) Clear() {
	_m.Called()
}

// Has provides a mock function with given fields: _a0
func (_m *ExecutionData) Has(_a0 flow.Identifier) bool {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Has")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(flow.Identifier) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Remove provides a mock function with given fields: _a0
func (_m *ExecutionData) Remove(_a0 flow.Identifier) bool {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Remove")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(flow.Identifier) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Size provides a mock function with no fields
func (_m *ExecutionData) Size() uint {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Size")
	}

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

// NewExecutionData creates a new instance of ExecutionData. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExecutionData(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExecutionData {
	mock := &ExecutionData{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
