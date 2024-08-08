// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import (
	context "context"

	flow "github.com/onflow/flow-go/model/flow"
	execution_data "github.com/onflow/flow-go/module/executiondatasync/execution_data"

	mock "github.com/stretchr/testify/mock"
)

// Downloader is an autogenerated mock type for the Downloader type
type Downloader struct {
	mock.Mock
}

// AddHeightUpdatesConsumer provides a mock function with given fields: _a0
func (_m *Downloader) AddHeightUpdatesConsumer(_a0 execution_data.HeightUpdatesConsumer) {
	_m.Called(_a0)
}

// Done provides a mock function with given fields:
func (_m *Downloader) Done() <-chan struct{} {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Done")
	}

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// Get provides a mock function with given fields: ctx, rootID
func (_m *Downloader) Get(ctx context.Context, rootID flow.Identifier) (*execution_data.BlockExecutionData, error) {
	ret := _m.Called(ctx, rootID)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *execution_data.BlockExecutionData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) (*execution_data.BlockExecutionData, error)); ok {
		return rf(ctx, rootID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) *execution_data.BlockExecutionData); ok {
		r0 = rf(ctx, rootID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*execution_data.BlockExecutionData)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(ctx, rootID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HighestCompleteHeight provides a mock function with given fields:
func (_m *Downloader) HighestCompleteHeight() uint64 {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for HighestCompleteHeight")
	}

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// OnBlockProcessed provides a mock function with given fields: _a0
func (_m *Downloader) OnBlockProcessed(_a0 uint64) {
	_m.Called(_a0)
}

// Ready provides a mock function with given fields:
func (_m *Downloader) Ready() <-chan struct{} {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Ready")
	}

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// NewDownloader creates a new instance of Downloader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDownloader(t interface {
	mock.TestingT
	Cleanup(func())
}) *Downloader {
	mock := &Downloader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
