// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	crypto "github.com/onflow/crypto"
	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// WeightedSignatureAggregator is an autogenerated mock type for the WeightedSignatureAggregator type
type WeightedSignatureAggregator struct {
	mock.Mock
}

// Aggregate provides a mock function with no fields
func (_m *WeightedSignatureAggregator) Aggregate() (flow.IdentifierList, []byte, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Aggregate")
	}

	var r0 flow.IdentifierList
	var r1 []byte
	var r2 error
	if rf, ok := ret.Get(0).(func() (flow.IdentifierList, []byte, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() flow.IdentifierList); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flow.IdentifierList)
		}
	}

	if rf, ok := ret.Get(1).(func() []byte); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]byte)
		}
	}

	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// TotalWeight provides a mock function with no fields
func (_m *WeightedSignatureAggregator) TotalWeight() uint64 {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for TotalWeight")
	}

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// TrustedAdd provides a mock function with given fields: signerID, sig
func (_m *WeightedSignatureAggregator) TrustedAdd(signerID flow.Identifier, sig crypto.Signature) (uint64, error) {
	ret := _m.Called(signerID, sig)

	if len(ret) == 0 {
		panic("no return value specified for TrustedAdd")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(flow.Identifier, crypto.Signature) (uint64, error)); ok {
		return rf(signerID, sig)
	}
	if rf, ok := ret.Get(0).(func(flow.Identifier, crypto.Signature) uint64); ok {
		r0 = rf(signerID, sig)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(flow.Identifier, crypto.Signature) error); ok {
		r1 = rf(signerID, sig)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Verify provides a mock function with given fields: signerID, sig
func (_m *WeightedSignatureAggregator) Verify(signerID flow.Identifier, sig crypto.Signature) error {
	ret := _m.Called(signerID, sig)

	if len(ret) == 0 {
		panic("no return value specified for Verify")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier, crypto.Signature) error); ok {
		r0 = rf(signerID, sig)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewWeightedSignatureAggregator creates a new instance of WeightedSignatureAggregator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWeightedSignatureAggregator(t interface {
	mock.TestingT
	Cleanup(func())
}) *WeightedSignatureAggregator {
	mock := &WeightedSignatureAggregator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
