// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocknetwork

import (
	io "io"

	network "github.com/onflow/flow-go/network"
	mock "github.com/stretchr/testify/mock"
)

// Codec is an autogenerated mock type for the Codec type
type Codec struct {
	mock.Mock
}

// Decode provides a mock function with given fields: data
func (_m *Codec) Decode(data []byte) (interface{}, error) {
	ret := _m.Called(data)

	if len(ret) == 0 {
		panic("no return value specified for Decode")
	}

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (interface{}, error)); ok {
		return rf(data)
	}
	if rf, ok := ret.Get(0).(func([]byte) interface{}); ok {
		r0 = rf(data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Encode provides a mock function with given fields: v
func (_m *Codec) Encode(v interface{}) ([]byte, error) {
	ret := _m.Called(v)

	if len(ret) == 0 {
		panic("no return value specified for Encode")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}) ([]byte, error)); ok {
		return rf(v)
	}
	if rf, ok := ret.Get(0).(func(interface{}) []byte); ok {
		r0 = rf(v)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(v)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewDecoder provides a mock function with given fields: r
func (_m *Codec) NewDecoder(r io.Reader) network.Decoder {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for NewDecoder")
	}

	var r0 network.Decoder
	if rf, ok := ret.Get(0).(func(io.Reader) network.Decoder); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(network.Decoder)
		}
	}

	return r0
}

// NewEncoder provides a mock function with given fields: w
func (_m *Codec) NewEncoder(w io.Writer) network.Encoder {
	ret := _m.Called(w)

	if len(ret) == 0 {
		panic("no return value specified for NewEncoder")
	}

	var r0 network.Encoder
	if rf, ok := ret.Get(0).(func(io.Writer) network.Encoder); ok {
		r0 = rf(w)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(network.Encoder)
		}
	}

	return r0
}

// NewCodec creates a new instance of Codec. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCodec(t interface {
	mock.TestingT
	Cleanup(func())
}) *Codec {
	mock := &Codec{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
