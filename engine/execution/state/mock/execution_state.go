// Code generated by mockery v2.21.4. DO NOT EDIT.

package mock

import (
	context "context"

	execution "github.com/onflow/flow-go/engine/execution"
	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"

	snapshot "github.com/onflow/flow-go/fvm/storage/snapshot"
)

// ExecutionState is an autogenerated mock type for the ExecutionState type
type ExecutionState struct {
	mock.Mock
}

// ChunkDataPackByChunkID provides a mock function with given fields: _a0
func (_m *ExecutionState) ChunkDataPackByChunkID(_a0 flow.Identifier) (*flow.ChunkDataPack, error) {
	ret := _m.Called(_a0)

	var r0 *flow.ChunkDataPack
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) (*flow.ChunkDataPack, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.ChunkDataPack); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.ChunkDataPack)
		}
	}

	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateStorageSnapshot provides a mock function with given fields: blockID
func (_m *ExecutionState) CreateStorageSnapshot(blockID flow.Identifier) (snapshot.StorageSnapshot, *flow.Header, error) {
	ret := _m.Called(blockID)

	var r0 snapshot.StorageSnapshot
	var r1 *flow.Header
	var r2 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) (snapshot.StorageSnapshot, *flow.Header, error)); ok {
		return rf(blockID)
	}
	if rf, ok := ret.Get(0).(func(flow.Identifier) snapshot.StorageSnapshot); ok {
		r0 = rf(blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(snapshot.StorageSnapshot)
		}
	}

	if rf, ok := ret.Get(1).(func(flow.Identifier) *flow.Header); ok {
		r1 = rf(blockID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*flow.Header)
		}
	}

	if rf, ok := ret.Get(2).(func(flow.Identifier) error); ok {
		r2 = rf(blockID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetExecutionResultID provides a mock function with given fields: _a0, _a1
func (_m *ExecutionState) GetExecutionResultID(_a0 context.Context, _a1 flow.Identifier) (flow.Identifier, error) {
	ret := _m.Called(_a0, _a1)

	var r0 flow.Identifier
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) (flow.Identifier, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) flow.Identifier); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flow.Identifier)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHighestExecutedBlockID provides a mock function with given fields: _a0
func (_m *ExecutionState) GetHighestExecutedBlockID(_a0 context.Context) (uint64, flow.Identifier, error) {
	ret := _m.Called(_a0)

	var r0 uint64
	var r1 flow.Identifier
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context) (uint64, flow.Identifier, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) uint64); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context) flow.Identifier); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(flow.Identifier)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context) error); ok {
		r2 = rf(_a0)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// HasState provides a mock function with given fields: _a0
func (_m *ExecutionState) HasState(_a0 flow.StateCommitment) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(flow.StateCommitment) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewStorageSnapshot provides a mock function with given fields: commit, blockID, height
func (_m *ExecutionState) NewStorageSnapshot(commit flow.StateCommitment, blockID flow.Identifier, height uint64) snapshot.StorageSnapshot {
	ret := _m.Called(commit, blockID, height)

	var r0 snapshot.StorageSnapshot
	if rf, ok := ret.Get(0).(func(flow.StateCommitment, flow.Identifier, uint64) snapshot.StorageSnapshot); ok {
		r0 = rf(commit, blockID, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(snapshot.StorageSnapshot)
		}
	}

	return r0
}

// SaveExecutionResults provides a mock function with given fields: ctx, result
func (_m *ExecutionState) SaveExecutionResults(ctx context.Context, result *execution.ComputationResult) error {
	ret := _m.Called(ctx, result)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *execution.ComputationResult) error); ok {
		r0 = rf(ctx, result)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StateCommitmentByBlockID provides a mock function with given fields: _a0, _a1
func (_m *ExecutionState) StateCommitmentByBlockID(_a0 context.Context, _a1 flow.Identifier) (flow.StateCommitment, error) {
	ret := _m.Called(_a0, _a1)

	var r0 flow.StateCommitment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) (flow.StateCommitment, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) flow.StateCommitment); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flow.StateCommitment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateHighestExecutedBlockIfHigher provides a mock function with given fields: _a0, _a1
func (_m *ExecutionState) UpdateHighestExecutedBlockIfHigher(_a0 context.Context, _a1 *flow.Header) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *flow.Header) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewExecutionState interface {
	mock.TestingT
	Cleanup(func())
}

// NewExecutionState creates a new instance of ExecutionState. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewExecutionState(t mockConstructorTestingTNewExecutionState) *ExecutionState {
	mock := &ExecutionState{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
