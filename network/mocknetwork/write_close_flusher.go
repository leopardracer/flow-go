// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocknetwork

import mock "github.com/stretchr/testify/mock"

// WriteCloseFlusher is an autogenerated mock type for the WriteCloseFlusher type
type WriteCloseFlusher struct {
	mock.Mock
}

// Close provides a mock function with no fields
func (_m *WriteCloseFlusher) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Flush provides a mock function with no fields
func (_m *WriteCloseFlusher) Flush() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Flush")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Write provides a mock function with given fields: p
func (_m *WriteCloseFlusher) Write(p []byte) (int, error) {
	ret := _m.Called(p)

	if len(ret) == 0 {
		panic("no return value specified for Write")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (int, error)); ok {
		return rf(p)
	}
	if rf, ok := ret.Get(0).(func([]byte) int); ok {
		r0 = rf(p)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewWriteCloseFlusher creates a new instance of WriteCloseFlusher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWriteCloseFlusher(t interface {
	mock.TestingT
	Cleanup(func())
}) *WriteCloseFlusher {
	mock := &WriteCloseFlusher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
